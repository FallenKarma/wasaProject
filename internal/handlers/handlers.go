package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/fallenkarma/wasatext/internal/models"
	"github.com/fallenkarma/wasatext/internal/service"
	"github.com/gorilla/mux"
)

// Handler defines the HTTP handlers for the API
type Handler struct {
	service *service.Service
}

// New creates a new Handler
func New(svc *service.Service) *Handler {
	return &Handler{
		service: svc,
	}
}

// extractToken extracts bearer token from Authorization header
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if bearerToken == "" {
		return ""
	}
	
	// Remove "Bearer " prefix if present
	return strings.TrimPrefix(bearerToken, "Bearer ")
}

// AuthMiddleware handles authentication for protected endpoints
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from request
		token := extractToken(r)
		log.Println("Extracted token: " + token)

		if token == "" {
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		// Verify token by checking if the user exists
		// In a real system, you would validate against JWT or other token mechanism
		user, err := h.service.GetUser(r.Context(), token)
		if err != nil || user == nil {
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context for use in handlers
		ctx := context.WithValue(r.Context(), "userID", token)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// getUserIDFromContext extracts the user ID from the request context
func getUserIDFromContext(r *http.Request) string {
	if userID, ok := r.Context().Value("userID").(string); ok {
		return userID
	}
	return ""
}

// respondWithJSON writes a JSON response
func respondWithJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
		}
	}
}

// respondWithError writes an error response
func respondWithError(w http.ResponseWriter, status int, message string) {
	respondWithJSON(w, status, map[string]string{"error": message})
}

// Login handles user login/creation
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	response, err := h.service.Login(r.Context(), req.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, response)
}

// SetMyUserName handles updating the user's name
func (h *Handler) SetMyUserName(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req models.UpdateUsernameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.UpdateUsername(r.Context(), userID, req.Name); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

// SetMyPhoto handles setting the user's profile photo
func (h *Handler) SetMyPhoto(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		respondWithError(w, http.StatusBadRequest, "Could not parse multipart form")
		return
	}

	// Get file from form
	file, _, err := r.FormFile("photo")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid file")
		return
	}
	defer file.Close()

	// Save photo
	photoURL, err := h.service.SetUserPhoto(r.Context(), userID, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"photo": photoURL})
}

// GetMyConversations returns all conversations for the authenticated user
func (h *Handler) GetMyConversations(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	conversations, err := h.service.GetConversations(r.Context(), userID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, conversations)
}

// GetConversation returns a specific conversation
func (h *Handler) GetConversation(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	conversationID := vars["id"]

	conversation, err := h.service.GetConversation(r.Context(), conversationID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if user is a participant in this conversation
	isParticipant := false
	for _, id := range conversation.Participants {
		if id == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		respondWithError(w, http.StatusForbidden, "Not a participant in this conversation")
		return
	}

	respondWithJSON(w, http.StatusOK, conversation)
}

// SendMessage handles sending a new message
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Extract conversation ID and message content
	conversationID := msg.ID // Assuming the ID field contains the conversation ID
	content := msg.Content
	replyToID := msg.ReplyTo

	// Check message type and process accordingly
	var newMsg *models.Message
	var err error

	if msg.Type == models.TextMessage {
		newMsg, err = h.service.SendTextMessage(r.Context(), userID, conversationID, content, replyToID)
	} else {
		respondWithError(w, http.StatusBadRequest, "Invalid message type for this endpoint")
		return
	}

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, newMsg)
}

// ForwardMessage handles forwarding a message to another conversation
func (h *Handler) ForwardMessage(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	var req models.ForwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.ForwardMessage(r.Context(), userID, req.MessageID, req.TargetConversationID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

// CommentMessage adds a reaction to a message
func (h *Handler) CommentMessage(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	messageID := vars["id"]

	var reaction models.Reaction
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if err := h.service.AddReaction(r.Context(), userID, messageID, reaction.Emoji); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

// UncommentMessage removes a reaction from a message
func (h *Handler) UncommentMessage(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	messageID := vars["id"]

	if err := h.service.RemoveReaction(r.Context(), userID, messageID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

// DeleteMessage deletes a message
func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	messageID := vars["id"]

	if err := h.service.DeleteMessage(r.Context(), userID, messageID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}

// AddToGroup adds a user to a group conversation
func (h *Handler) AddToGroup(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	groupID := vars["id"]

	var req models.AddToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check if the user is in the group (to verify they have permission)
	group, err := h.service.GetConversation(r.Context(), groupID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Verify user is in the group
	isInGroup := false
	for _, id := range group.Participants {
		if id == userID {
			isInGroup = true
			break
		}
	}

	if !isInGroup {
		respondWithError(w, http.StatusForbidden, "You must be in the group to add members")
		return
	}

	if err := h.service.AddToGroup(r.Context(), groupID, req.UserID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

// LeaveGroup handles a user leaving a group conversation
func (h *Handler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	groupID := vars["id"]

	if err := h.service.LeaveGroup(r.Context(), groupID, userID); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

// SetGroupName updates a group's name
func (h *Handler) SetGroupName(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	groupID := vars["id"]

	var req models.SetGroupNameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Check if the user is in the group (to verify they have permission)
	group, err := h.service.GetConversation(r.Context(), groupID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Verify user is in the group
	isInGroup := false
	for _, id := range group.Participants {
		if id == userID {
			isInGroup = true
			break
		}
	}

	if !isInGroup {
		respondWithError(w, http.StatusForbidden, "You must be in the group to change its name")
		return
	}

	if err := h.service.SetGroupName(r.Context(), groupID, req.Name); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, nil)
}

// SetGroupPhoto updates a group's photo
func (h *Handler) SetGroupPhoto(w http.ResponseWriter, r *http.Request) {
	userID := getUserIDFromContext(r)
	if userID == "" {
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	groupID := vars["id"]

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		respondWithError(w, http.StatusBadRequest, "Could not parse multipart form")
		return
	}

	// Get file from form
	file, _, err := r.FormFile("photo")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid file")
		return
	}
	defer file.Close()

	// Check if the user is in the group (to verify they have permission)
	group, err := h.service.GetConversation(r.Context(), groupID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Verify user is in the group
	isInGroup := false
	for _, id := range group.Participants {
		if id == userID {
			isInGroup = true
			break
		}
	}

	if !isInGroup {
		respondWithError(w, http.StatusForbidden, "You must be in the group to change its photo")
		return
	}

	// Save photo
	photoURL, err := h.service.SetGroupPhoto(r.Context(), groupID, file)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"photo": photoURL})
}