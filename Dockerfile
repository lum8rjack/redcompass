# Stage 1: Build Node.js app
FROM node:25.2 AS node-builder

# Set working directory
WORKDIR /app

# Copy package.json and package-lock.json first (for better caching)
COPY frontend/package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application files
COPY frontend .

# Build the Node.js application (adjust as needed)
RUN npm run build

# Stage 2: Build Go app
FROM golang:1.25 AS go-builder

# Set working directory
WORKDIR /goapp

# Copy the Go application source code
COPY . .

# Copy the Node.js files
COPY --from=node-builder /app/dist ./frontend/dist

# Set environment variable for a static Go build
ENV CGO_ENABLED=0

# Build the Go application
RUN go build -ldflags "-s -w" -trimpath -o redcompass.bin

# Final Stage: Use a distroless image to run the compiled Go app
FROM gcr.io/distroless/base

# Set working directory
WORKDIR /

# Copy the compiled Go binary
COPY --from=go-builder /goapp/redcompass.bin .

# Expose necessary port
EXPOSE 8090

# Run the application
CMD ["./redcompass.bin", "serve", "--http", "0.0.0.0:8090"]
