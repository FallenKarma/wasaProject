package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

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
	log.Println("Initializing API handlers")
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

// logRequest logs information about an incoming request
func logRequest(handler string, r *http.Request, userID string) {
	log.Printf("[%s] %s %s | UserID: %s | IP: %s | User-Agent: %s", 
		handler, 
		r.Method, 
		r.URL.Path, 
		userID,
		r.RemoteAddr,
		r.UserAgent(),
	)
}

// logRequestWithDuration logs information about a request and its duration
func logRequestWithDuration(handler string, r *http.Request, userID string, start time.Time, statusCode int) {
	duration := time.Since(start)
	log.Printf("[%s] %s %s | UserID: %s | Status: %d | Duration: %s | IP: %s", 
		handler, 
		r.Method, 
		r.URL.Path, 
		userID,
		statusCode,
		duration,
		r.RemoteAddr,
	)
}

// logError logs an error with request context
func logError(handler string, r *http.Request, userID string, err error, msg string) {
	log.Printf("[ERROR][%s] %s %s | UserID: %s | %s: %v",
		handler,
		r.Method,
		r.URL.Path,
		userID,
		msg,
		err,
	)
}

// AuthMiddleware handles authentication for protected endpoints
func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract token from request
		token := extractToken(r)
		if token == "" {
			log.Printf("[AuthMiddleware] %s %s | No token provided | IP: %s", r.Method, r.URL.Path, r.RemoteAddr)
			http.Error(w, "Unauthorized: No token provided", http.StatusUnauthorized)
			return
		}

		// Verify token by checking if the user exists
		// In a real system, you would validate against JWT or other token mechanism
		start := time.Now()
		user, err := h.service.GetUser(r.Context(), token)
		if err != nil || user == nil {
			log.Printf("[AuthMiddleware] %s %s | Invalid token | Token: %s | IP: %s | Error: %v", 
				r.Method, r.URL.Path, token, r.RemoteAddr, err)
			http.Error(w, "Unauthorized: Invalid token", http.StatusUnauthorized)
			return
		}
		
		log.Printf("[AuthMiddleware] %s %s | User authenticated | UserID: %s | Duration: %s", 
			r.Method, r.URL.Path, token, time.Since(start))

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
	handlerName := "Login"
	start := time.Now()
	
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logError(handlerName, r, "", err, "Invalid request payload")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Printf("[%s] Login attempt | Username: %s | IP: %s", handlerName, req.Name, r.RemoteAddr)
	
	response, err := h.service.Login(r.Context(), req.Name)
	if err != nil {
		logError(handlerName, r, "", err, "Login failed")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[%s] Login successful | UserID: %s | Username: %s | Duration: %s", 
		handlerName, response.Identifier, req.Name, time.Since(start))
	
	respondWithJSON(w, http.StatusCreated, response)
}

func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
	handlerName := "GetUsers"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	logRequest(handlerName, r, userID)

	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		logError(handlerName, r, userID, err, "Failed to get users")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[%s] Retrieved users | UserID: %s | Count: %d | Duration: %s", 
		handlerName, userID, len(users), time.Since(start))
	
	respondWithJSON(w, http.StatusOK, users)
}


// SetMyUserName handles updating the user's name
func (h *Handler) SetMyUserName(w http.ResponseWriter, r *http.Request) {
	handlerName := "SetMyUserName"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	logRequest(handlerName, r, userID)
	
	var req models.UpdateUsernameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logError(handlerName, r, userID, err, "Invalid request payload")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Printf("[%s] Username update | UserID: %s | New name: %s", handlerName, userID, req.Name)
	
	if err := h.service.UpdateUsername(r.Context(), userID, req.Name); err != nil {
		logError(handlerName, r, userID, err, "Failed to update username")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	logRequestWithDuration(handlerName, r, userID, start, http.StatusOK)
	respondWithJSON(w, http.StatusOK, nil)
}

// SetMyPhoto handles setting the user's profile photo
func (h *Handler) SetMyPhoto(w http.ResponseWriter, r *http.Request) {
	handlerName := "SetMyPhoto"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	logRequest(handlerName, r, userID)

	// Parse multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10 MB max
		logError(handlerName, r, userID, err, "Could not parse multipart form")
		respondWithError(w, http.StatusBadRequest, "Could not parse multipart form")
		return
	}

	// Get file from form
	file, fileHeader, err := r.FormFile("photo")
	if err != nil {
		logError(handlerName, r, userID, err, "Invalid file")
		respondWithError(w, http.StatusBadRequest, "Invalid file")
		return
	}
	defer file.Close()

	log.Printf("[%s] Processing profile photo | UserID: %s | File size: %d bytes | File name: %s", 
		handlerName, userID, fileHeader.Size, fileHeader.Filename)

	// Save photo
	photoURL, err := h.service.SetUserPhoto(r.Context(), userID, file)
	if err != nil {
		logError(handlerName, r, userID, err, "Failed to save profile photo")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[%s] Profile photo updated | UserID: %s | Photo URL: %s | Duration: %s", 
		handlerName, userID, photoURL, time.Since(start))
		
	respondWithJSON(w, http.StatusOK, map[string]string{"photo": photoURL})
}

// CreateConversation handles creating a new conversation
func (h *Handler) CreateConversation(w http.ResponseWriter, r *http.Request) {
	handlerName := "CreateConversation"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	logRequest(handlerName, r, userID)
	
	var req models.CreateConversationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logError(handlerName, r, userID, err, "Invalid request payload")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	log.Printf("type: %s", req.Type)

	// Validate request
	if len(req.Participants) == 0 {
		log.Printf("[%s] Invalid request: no participants | UserID: %s", handlerName, userID)
		respondWithError(w, http.StatusBadRequest, "At least one participant is required")
		return
	}

	log.Printf("[%s] Creating conversation | Type: %s | UserID: %s | Participants: %d | Name: %s", 
		handlerName, req.Type, userID, len(req.Participants), req.Name)

	// Create the conversation
	conversation, err := h.service.CreateConversation(r.Context(), userID, req.Participants, req.Type, req.Name)
	if err != nil {
		logError(handlerName, r, userID, err, "Failed to create conversation")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[%s] Conversation created | ConversationID: %s | UserID: %s | Type: %s | Participants: %d | Duration: %s", 
		handlerName, conversation.ID, userID, req.Type, len(conversation.Participants), time.Since(start))
	
	respondWithJSON(w, http.StatusCreated, conversation)
}

// GetMyConversations returns all conversations for the authenticated user
func (h *Handler) GetMyConversations(w http.ResponseWriter, r *http.Request) {
	handlerName := "GetMyConversations"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	logRequest(handlerName, r, userID)

	conversations, err := h.service.GetConversations(r.Context(), userID)
	if err != nil {
		logError(handlerName, r, userID, err, "Failed to get conversations")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[%s] Retrieved conversations | UserID: %s | Count: %d | Duration: %s", 
		handlerName, userID, len(conversations), time.Since(start))
	
	respondWithJSON(w, http.StatusOK, conversations)
}

// GetConversation returns a specific conversation
func (h *Handler) GetConversation(w http.ResponseWriter, r *http.Request) {
	handlerName := "GetConversation"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	conversationID := vars["id"]
	
	logRequest(handlerName, r, userID)
	log.Printf("[%s] Retrieving conversation | UserID: %s | ConversationID: %s", 
		handlerName, userID, conversationID)

	conversation, err := h.service.GetConversation(r.Context(), conversationID)
	if err != nil {
		logError(handlerName, r, userID, err, fmt.Sprintf("Failed to get conversation ID: %s", conversationID))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Check if user is a participant in this conversation
	isParticipant := false
	for _, participant := range conversation.Participants {
		if participant.ID == userID {
			isParticipant = true
			break
		}
	}

	if !isParticipant {
		log.Printf("[%s] Access denied | UserID: %s | ConversationID: %s | Not a participant", 
			handlerName, userID, conversationID)
		respondWithError(w, http.StatusForbidden, "Not a participant in this conversation")
		return
	}

	log.Printf("[%s] Conversation retrieved | UserID: %s | ConversationID: %s | Participants: %d | Messages: %d | Duration: %s", 
		handlerName, userID, conversationID, len(conversation.Participants), len(conversation.Messages), time.Since(start))
	
	respondWithJSON(w, http.StatusOK, conversation)
}

// SendMessage handles sending a new message
func (h *Handler) SendMessage(w http.ResponseWriter, r *http.Request) {
	handlerName := "SendMessage"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	logRequest(handlerName, r, userID)
	
	var msg models.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		logError(handlerName, r, userID, err, "Invalid request payload")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	// Extract conversation ID and message content
	conversationID := msg.ID // Assuming the ID field contains the conversation ID
	content := msg.Content
	replyToID := msg.ReplyTo

	log.Printf("[%s] Processing message | UserID: %s | ConvID: %s | Type: %s | ReplyToID: %s | ContentLen: %d", 
		handlerName, userID, conversationID, msg.Type, replyToID, len(content))

	// Check message type and process accordingly
	var newMsg *models.Message
	var err error

	if msg.Type == models.TextMessage {
		newMsg, err = h.service.SendTextMessage(r.Context(), userID, conversationID, content, replyToID)
	} else {
		log.Printf("[%s] Invalid message type | UserID: %s | Type: %s", handlerName, userID, msg.Type)
		respondWithError(w, http.StatusBadRequest, "Invalid message type for this endpoint")
		return
	}

	if err != nil {
		logError(handlerName, r, userID, err, fmt.Sprintf("Failed to send message to conversation: %s", conversationID))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[%s] Message sent | UserID: %s | ConvID: %s | MessageID: %s | Duration: %s", 
		handlerName, userID, conversationID, newMsg.ID, time.Since(start))
	
	respondWithJSON(w, http.StatusCreated, newMsg)
}

// ForwardMessage handles forwarding a message to another conversation
func (h *Handler) ForwardMessage(w http.ResponseWriter, r *http.Request) {
	handlerName := "ForwardMessage"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	logRequest(handlerName, r, userID)
	
	var req models.ForwardMessageRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logError(handlerName, r, userID, err, "Invalid request payload")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Printf("[%s] Forwarding message | UserID: %s | MessageID: %s | TargetConvID: %s", 
		handlerName, userID, req.MessageID, req.TargetConversationID)

	if err := h.service.ForwardMessage(r.Context(), userID, req.MessageID, req.TargetConversationID); err != nil {
		logError(handlerName, r, userID, err, "Failed to forward message")
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[%s] Message forwarded | UserID: %s | MessageID: %s | TargetConvID: %s | Duration: %s", 
		handlerName, userID, req.MessageID, req.TargetConversationID, time.Since(start))
	
	respondWithJSON(w, http.StatusOK, nil)
}

// CommentMessage adds a reaction to a message
func (h *Handler) CommentMessage(w http.ResponseWriter, r *http.Request) {
	handlerName := "CommentMessage"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	messageID := vars["id"]
	
	logRequest(handlerName, r, userID)

	var reaction models.Reaction
	if err := json.NewDecoder(r.Body).Decode(&reaction); err != nil {
		logError(handlerName, r, userID, err, "Invalid request payload")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Printf("[%s] Adding reaction | UserID: %s | MessageID: %s | Emoji: %s", 
		handlerName, userID, messageID, reaction.Emoji)

	if err := h.service.AddReaction(r.Context(), userID, messageID, reaction.Emoji); err != nil {
		logError(handlerName, r, userID, err, fmt.Sprintf("Failed to add reaction to message: %s", messageID))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[%s] Reaction added | UserID: %s | MessageID: %s | Emoji: %s | Duration: %s", 
		handlerName, userID, messageID, reaction.Emoji, time.Since(start))
	
	respondWithJSON(w, http.StatusOK, nil)
}

// UncommentMessage removes a reaction from a message
func (h *Handler) UncommentMessage(w http.ResponseWriter, r *http.Request) {
	handlerName := "UncommentMessage"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	messageID := vars["id"]
	
	logRequest(handlerName, r, userID)
	log.Printf("[%s] Removing reaction | UserID: %s | MessageID: %s", handlerName, userID, messageID)

	if err := h.service.RemoveReaction(r.Context(), userID, messageID); err != nil {
		logError(handlerName, r, userID, err, fmt.Sprintf("Failed to remove reaction from message: %s", messageID))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[%s] Reaction removed | UserID: %s | MessageID: %s | Duration: %s", 
		handlerName, userID, messageID, time.Since(start))
	
	respondWithJSON(w, http.StatusOK, nil)
}

// DeleteMessage deletes a message
func (h *Handler) DeleteMessage(w http.ResponseWriter, r *http.Request) {
	handlerName := "DeleteMessage"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	messageID := vars["id"]
	
	logRequest(handlerName, r, userID)
	log.Printf("[%s] Deleting message | UserID: %s | MessageID: %s", handlerName, userID, messageID)

	if err := h.service.DeleteMessage(r.Context(), userID, messageID); err != nil {
		logError(handlerName, r, userID, err, fmt.Sprintf("Failed to delete message: %s", messageID))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[%s] Message deleted | UserID: %s | MessageID: %s | Duration: %s", 
		handlerName, userID, messageID, time.Since(start))
	
	respondWithJSON(w, http.StatusNoContent, nil)
}

// AddToGroup adds a user to a group conversation
func (h *Handler) AddToGroup(w http.ResponseWriter, r *http.Request) {
	handlerName := "AddToGroup"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	groupID := vars["id"]
	
	logRequest(handlerName, r, userID)

	var req models.AddToGroupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logError(handlerName, r, userID, err, "Invalid request payload")
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	log.Printf("[%s] Adding user to group | RequestedBy: %s | GroupID: %s | NewUserID: %s", 
		handlerName, userID, groupID, req.UserID)

	// Check if the user is in the group (to verify they have permission)
	group, err := h.service.GetConversation(r.Context(), groupID)
	if err != nil {
		logError(handlerName, r, userID, err, fmt.Sprintf("Failed to get group: %s", groupID))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Verify user is in the group
	isInGroup := false
	for _, participant := range group.Participants {
		if participant.ID == userID {
			isInGroup = true
			break
		}
	}

	if !isInGroup {
		log.Printf("[%s] Permission denied | UserID: %s | GroupID: %s | Not in group", 
			handlerName, userID, groupID)
		respondWithError(w, http.StatusForbidden, "You must be in the group to add members")
		return
	}

	if err := h.service.AddToGroup(r.Context(), groupID, req.UserID); err != nil {
		logError(handlerName, r, userID, err, fmt.Sprintf("Failed to add user %s to group %s", req.UserID, groupID))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[%s] User added to group | RequestedBy: %s | GroupID: %s | AddedUserID: %s | Duration: %s", 
		handlerName, userID, groupID, req.UserID, time.Since(start))
	
	respondWithJSON(w, http.StatusOK, nil)
}

// LeaveGroup handles a user leaving a group conversation
func (h *Handler) LeaveGroup(w http.ResponseWriter, r *http.Request) {
	handlerName := "LeaveGroup"
	start := time.Now()
	
	userID := getUserIDFromContext(r)
	if userID == "" {
		log.Printf("[%s] %s %s | Not authenticated | IP: %s", handlerName, r.Method, r.URL.Path, r.RemoteAddr)
		respondWithError(w, http.StatusUnauthorized, "Not authenticated")
		return
	}

	vars := mux.Vars(r)
	groupID := vars["id"]
	
	logRequest(handlerName, r, userID)
	log.Printf("[%s] User leaving group | UserID: %s | GroupID: %s", handlerName, userID, groupID)

	if err := h.service.LeaveGroup(r.Context(), groupID, userID); err != nil {
		logError(handlerName, r, userID, err, fmt.Sprintf("Failed to leave group: %s", groupID))
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("[%s] User left group | UserID: %s | GroupID: %s | Duration: %s", 
		handlerName, userID, groupID, time.Since(start))
	
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
	for _, participant := range group.Participants {
		if participant.ID == userID {
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
	for _, participant := range group.Participants {
		if participant.ID == userID {
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