# AI Resume Builder

A modern web application that enables users to build professional, ATS-optimized resumes with AI assistance.

## Tech Stack

### Backend
- **Language:** Go 1.24
- **Framework:** Fiber v2
- **Database:** PostgreSQL (via GORM, schema applied with `AutoMigrate` on startup)
- **AI:** Anthropic Claude (Sonnet) for extraction, content improvement, summaries, and job analysis
- **PDF export:** Headless Chrome via `chromedp`
- **File parsing:** PDF / DOCX / DOC / TXT (text extraction for AI import)
- **Auth:** JWT access + refresh tokens (`golang-jwt`), bcrypt password hashing

### Frontend
- **Framework:** SvelteKit (Svelte 4)
- **Language:** TypeScript 5+
- **Styling:** Tailwind CSS
- **UI Components:** Custom shadcn-style components
- **Icons:** Lucide
- **Tests:** Vitest

### Infrastructure
- **Hosting:** Fly.io (`fly.toml` in each app) — `docker-compose` for local Postgres
- **File Storage (optional):** Cloudflare R2 (config present; wire up as needed)

## Project Structure

```
resume-builder/
├── resume-backend/          # Go backend API
│   ├── cmd/
│   │   └── api/             # HTTP server entry point
│   ├── internal/
│   │   ├── api/             # HTTP handlers, middleware, router
│   │   ├── service/         # Business logic (auth, resume, ai, pdf, fileparser)
│   │   ├── repository/      # Data access layer (GORM)
│   │   ├── models/          # Domain models
│   │   └── config/          # Configuration
│   ├── Dockerfile
│   └── fly.toml
│
├── resume-frontend/         # SvelteKit frontend
│   ├── src/
│   │   ├── routes/          # Page routes
│   │   │   ├── (marketing)  # Public pages (home, about, privacy, terms, templates)
│   │   │   ├── (auth)/      # Login / signup
│   │   │   └── (app)/       # Dashboard, resume editor, settings
│   │   └── lib/
│   │       ├── components/  # UI components
│   │       ├── stores/      # Svelte stores
│   │       ├── api/         # API client
│   │       ├── types/       # TypeScript types
│   │       └── utils/       # Utility functions
│   ├── Dockerfile
│   └── fly.toml
│
└── README.md
```

## Getting Started

### Prerequisites

- Go 1.24+
- Node.js 20+
- Docker & Docker Compose (for the local Postgres stack)

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
- Frontend: `http://localhost:5173` (proxies API requests to the backend)

### Backend Setup (without Docker)

1. `cd resume-backend`
2. `cp .env.example .env` and fill in your configuration
3. `go mod download`
4. `go run cmd/api/main.go`

The API will be available at `http://localhost:8080`. The database schema is applied
automatically via GORM `AutoMigrate`, and the default templates are seeded on startup.

### Frontend Setup

1. `cd resume-frontend`
2. `npm install`
3. `npm run dev`

The app will be available at `http://localhost:5173`.

## Testing

### Backend
```bash
cd resume-backend
go test ./...
```

### Frontend
```bash
cd resume-frontend
npm test          # run once
npm run test:watch
npm run check     # svelte-check (type + a11y)
```

## API Endpoints

### Authentication
- `POST /api/v1/auth/register` - Register new user
- `POST /api/v1/auth/login` - Login
- `POST /api/v1/auth/logout` - Logout
- `POST /api/v1/auth/refresh` - Refresh token
- `POST /api/v1/auth/forgot-password` - Request password reset
- `POST /api/v1/auth/reset-password` - Reset password with token
- `POST /api/v1/auth/verify-email` - Verify email with token
- `GET /api/v1/auth/me` - Get current user (protected)

### Resumes (protected)
- `GET /api/v1/resumes` - List resumes
- `POST /api/v1/resumes` - Create resume
- `GET /api/v1/resumes/:id` - Get resume
- `PUT /api/v1/resumes/:id` - Update resume
- `PATCH /api/v1/resumes/:id` - Partial update (auto-save)
- `DELETE /api/v1/resumes/:id` - Delete resume
- `POST /api/v1/resumes/:id/duplicate` - Duplicate resume
- `GET /api/v1/resumes/:id/versions` - List version history
- `POST /api/v1/resumes/:id/versions/:versionId/restore` - Restore a version
- `PATCH /api/v1/resumes/:id/visibility` - Toggle public sharing

### Templates
- `GET /api/v1/templates` - List templates
- `GET /api/v1/templates/:id` - Get template
- `GET /api/v1/templates/:id/preview` - Get rendered preview config
- `GET /api/v1/templates/:id/preview-image` - Get preview image

### AI (protected)
- `POST /api/v1/ai/extract` - Extract structured resume data from text
- `POST /api/v1/ai/extract-file` - Extract from an uploaded file (PDF/DOCX/DOC/TXT)
- `GET /api/v1/ai/supported-formats` - List supported import formats
- `POST /api/v1/ai/improve` - Improve a snippet of text
- `POST /api/v1/ai/generate-summary` - Generate a professional summary
- `POST /api/v1/ai/analyze-job` - Analyze a resume against a job description

### Export (protected)
- `POST /api/v1/export/pdf` - Export as PDF
- `POST /api/v1/export/txt` - Export as plain text
- `POST /api/v1/export/json` - Export as JSON

### Public
- `GET /api/v1/share/:id` - View a publicly shared resume

## Features

### Phase 1 — MVP (done)
- [x] User authentication (JWT access + refresh, email verification, password reset)
- [x] Resume CRUD operations + duplicate
- [x] Resume editor with live preview and auto-save
- [x] Professional templates (16 seeded across modern/professional/creative/academic categories)
- [x] AI-powered data extraction (from text and uploaded files)
- [x] AI content improvement and summary generation
- [x] PDF / TXT / JSON export
- [x] Version history (with restore)
- [x] Public resume sharing (backend endpoint + visibility toggle)

### Phase 2 — Growth (in progress)
- [x] Public share page (frontend `/share/[slug]` viewer, template-aware, print-friendly)
- [ ] Tailor-to-job / ATS keyword match UI — `analyze-job` backend exists
- [ ] Cover letter generator
- [ ] LinkedIn import
- [ ] Expanded test coverage

### Phase 3 — Scale (not started)
- [ ] Stripe billing / subscriptions
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
- `ENVIRONMENT` - `development` or `production`
- `DATABASE_URL` - PostgreSQL connection string
- `REDIS_URL` - Redis connection string (reserved for future use)
- `JWT_SECRET` - JWT signing secret
- `JWT_EXPIRATION_HOURS` - Access token lifetime (default: 24)
- `REFRESH_TOKEN_DAYS` - Refresh token lifetime (default: 30)
- `ANTHROPIC_API_KEY` - Anthropic API key
- `ANTHROPIC_MODEL` - Claude model id
- `CORS_ORIGINS` - Comma-separated allowed origins

### Frontend
- `PUBLIC_API_URL` - Backend API URL

## License

MIT
