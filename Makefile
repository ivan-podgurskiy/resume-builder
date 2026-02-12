.PHONY: dev-backend dev-frontend

# Run backend stack (postgres, redis, backend) in Docker. Leave running while developing.
# Frontend changes won't require any rebuild.
dev-backend:
	docker compose up postgres redis backend

# Run frontend dev server with hot reload. Use alongside dev-backend.
dev-frontend:
	cd resume-frontend && npm run dev
