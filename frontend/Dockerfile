# Stage 1: Build the Vue app
FROM node:20 AS builder

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm install

# Copy the rest of the app and build it
COPY . .
RUN npm run build

# Stage 2: Serve the app using a lightweight static server
FROM node:20-slim

WORKDIR /app

# Install the static file server
RUN npm install -g serve

# Copy the built files from the builder stage
COPY --from=builder /app/dist ./dist

# Expose the port the server will run on
EXPOSE 4173

# Serve the static files
CMD ["serve", "-s", "dist", "-l", "4173"]
