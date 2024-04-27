.PHONY: frontend backend

# Frontend target
frontend:
	@echo "Building and running frontend..."
	cd ui && npm install && npm start

# Backend target
backend:
	@echo "Building and running backend..."
	cd server && go run ./cmd/api

# Help target
help:
	@echo "Available targets:"
	@echo "  frontend    : Build and run the frontend"
	@echo "  backend     : Build and run the backend"
	@echo "  all         : Build and run both frontend and backend"
	@echo "  help        : Show this help message"

# If no target is provided, show help
.DEFAULT_GOAL := help
