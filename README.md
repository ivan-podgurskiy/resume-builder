# AI Resume Builder

A modern web application that enables users to build professional, ATS-optimized resumes with AI assistance.

## Tech Stack

### Backend
- **Language:** Go 1.22+
- **Framework:** Fiber v2
- **Database:** PostgreSQL 16
- **ORM:** GORM
- **Queue:** Asynq (Redis-based background jobs)
- **Cache:** Redis
- **AI:** Anthropic Claude Sonnet 4

### Frontend
- **Framework:** SvelteKit
- **Language:** TypeScript 5+
- **Styling:** Tailwind CSS
- **UI Components:** Custom shadcn-style components
- **Icons:** Lucide

### Infrastructure
- **Hosting:** Fly.io
- **File Storage:** Cloudflare R2
- **CDN:** Cloudflare

## Project Structure

```
resume-builder/
├── resume-backend/          # Go backend API
│   ├── cmd/
│   │   ├── api/            # HTTP server entry point
│   │   ├── worker/         # Background worker
│   │   └── migrate/        # Database migrations
│   ├── internal/
│   │   ├── api/            # HTTP handlers, middleware, router
│   │   ├── service/        # Business logic
│   │   ├── repository/     # Data access layer
│   │   ├── models/         # Domain models
│   │   └── config/         # Configuration
│   ├── Dockerfile
│   └── fly.toml
│
├── resume-frontend/         # SvelteKit frontend
│   ├── src/
│   │   ├── routes/         # Page routes
│   │   │   ├── (marketing)/    # Public pages
│   │   │   ├── (auth)/         # Authentication
│   │   │   └── (app)/          # Authenticated app
│   │   └── lib/
│   │       ├── components/     # UI components
│   │       ├── stores/         # Svelte stores
│   │       ├── api/            # API client
│   │       ├── types/          # TypeScript types
│   │       └── utils/          # Utility functions
│   ├── Dockerfile
│   └── fly.toml
│
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.22+
- Node.js 20+
- Docker & Docker Compose (for backend stack)
- PostgreSQL 16
- Redis (for background jobs)

### Development (Recommended)

Run backend and frontend separately so you can change the frontend without rebuilding Docker:

**Terminal 1 – Backend (Docker):**
```bash
make dev-backend
# or: docker compose up postgres redis backend
```

**Terminal 2 – Frontend (local, hot reload):**
```bash
make dev-frontend
# or: cd resume-frontend && npm run dev
```

- Backend: `http://localhost:8080`
- Frontend: `http://localhost:5173` (proxies API requests to backend)

### Backend Setup (without Docker)

1. Navigate to the backend directory:
   ```bash
   cd resume-backend
   ```

2. Copy the environment file:
   ```bash
   cp .env.example .env
   ```

3. Update the `.env` file with your configuration

4. Install dependencies:
   ```bash
   go mod download
   ```

5. Start the server:
   ```bash
   go run cmd/api/main.go
   ```

The API will be available at `http://localhost:8080`

### Frontend Setup

1. Navigate to the frontend directory:
   ```bash
   cd resume-frontend
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

The app will be available at `http://localhost:5173`

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/logout` - Logout
- `POST /api/v1/auth/refresh` - Refresh token
- `GET /api/v1/auth/me` - Get current user

### Resumes
- `GET /api/v1/resumes` - List resumes
- `POST /api/v1/resumes` - Create resume
- `GET /api/v1/resumes/:id` - Get resume
- `PUT /api/v1/resumes/:id` - Update resume
- `DELETE /api/v1/resumes/:id` - Delete resume
- `POST /api/v1/resumes/:id/duplicate` - Duplicate resume

### Templates
- `GET /api/v1/templates` - List templates
- `GET /api/v1/templates/:id` - Get template

### Public
- `GET /api/v1/share/:slug` - View public resume

## Features

### Phase 1 (MVP)
- [x] User authentication (JWT)
- [x] Resume CRUD operations
- [x] Basic resume editor with live preview
- [x] 3 professional templates
- [ ] AI-powered data extraction
- [ ] AI content improvement
- [ ] PDF export
- [ ] DOCX export

### Phase 2 (Growth)
- [ ] LinkedIn import
- [ ] More templates
- [ ] Cover letter generator
- [ ] ATS compatibility checker
- [ ] Version history

### Phase 3 (Scale)
- [ ] Stripe integration
- [ ] Team collaboration
- [ ] Custom branding

## Deployment

Both applications are configured for deployment on Fly.io:

```bash
# Deploy backend
cd resume-backend
fly deploy

# Deploy frontend
cd resume-frontend
fly deploy
```

## Environment Variables

### Backend
- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string
- `JWT_SECRET` - JWT signing secret
- `ANTHROPIC_API_KEY` - Anthropic API key

### Frontend
- `PUBLIC_API_URL` - Backend API URL

## License

MIT
