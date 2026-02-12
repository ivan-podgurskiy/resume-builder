# VPS Deployment Guide: resume.ivanpodgurskiy.name

Deploy the AI Resume Builder to your VPS at **resume.ivanpodgurskiy.name**.

## Prerequisites

- VPS: `92.205.19.23` (user: `ivanpod`)
- SSH access to the VPS
- Domain `ivanpodgurskiy.name` with DNS access

## One-Time Setup

### 1. Create `.env` on VPS

SSH into the VPS and create the backend environment file:

```bash
ssh ivanpod@92.205.19.23
mkdir -p ~/docker/projects/resume-builder
cd ~/docker/projects/resume-builder
nano resume-backend/.env
```

Copy from `resume-backend/.env.example` and set:

| Variable | Required | Notes |
|----------|----------|-------|
| `DATABASE_URL` | Auto (in compose) | Overridden by docker-compose |
| `REDIS_URL` | Auto | Overridden by docker-compose |
| `JWT_SECRET` | **Yes** | Generate: `openssl rand -base64 32` |
| `POSTGRES_PASSWORD` | **Yes** | Strong password for DB |
| `ANTHROPIC_API_KEY` | For AI features | Get from console.anthropic.com |
| `R2_*` | Optional | For file storage; can skip initially |

### 2. DNS: Add Subdomain

Add an A record in your DNS provider:

```
Type: A
Name: resume
Value: 92.205.19.23
TTL: 3600
```

### 3. Nginx: Add Resume Subdomain

The nginx config in **reddit-ai-monitor** (`nginx.ivanpodgurskiy.name.conf`) has been updated to include `resume.ivanpodgurskiy.name`. Deploy it:

```bash
cd ../reddit-ai-monitor
chmod +x deploy_nginx_config.sh
./deploy_nginx_config.sh
```

If you prefer to add manually, see `nginx.resume.snippet.conf` in this project.

### 4. SSL Certificate

After DNS propagates, add the subdomain to your certificate:

```bash
ssh ivanpod@92.205.19.23
sudo certbot certonly --nginx \
  -d ivanpodgurskiy.name \
  -d n8n.ivanpodgurskiy.name \
  -d metabase.ivanpodgurskiy.name \
  -d resume.ivanpodgurskiy.name
```

If you already have a wildcard cert for `*.ivanpodgurskiy.name`, no change needed.

### 5. Reload Nginx

```bash
sudo nginx -t
sudo systemctl reload nginx
```

## Deploy

From the `resume-builder` project directory:

```bash
chmod +x deploy.sh
./deploy.sh
```

This will:
1. Build Go backend locally
2. Build frontend locally (`npm run build`)
3. Show file count and total size to be copied
4. Rsync the project to the VPS (including pre-built artifacts, excluding node_modules, .git)
5. On VPS: run `docker compose` to build runtime images and start containers

## Verify

- **App:** https://resume.ivanpodgurskiy.name
- **Health:** `curl https://resume.ivanpodgurskiy.name` (should return HTML)

## If SSH Drops During Build

The backend image installs Chromium (~200MB), so the build can take 5–10 minutes. If your SSH session drops:

```bash
# SSH in and run the build directly on the VPS
ssh ivanpod@92.205.19.23
cd ~/docker/projects/resume-builder
docker compose -f docker-compose.prod.yml build --no-cache
docker compose -f docker-compose.prod.yml up -d
```

Or run in `tmux` so it survives disconnects:

```bash
ssh ivanpod@92.205.19.23
tmux new -s build
cd ~/docker/projects/resume-builder
docker compose -f docker-compose.prod.yml build --no-cache && docker compose -f docker-compose.prod.yml up -d
# Ctrl+B, D to detach; tmux attach -t build to reattach
```

## Troubleshooting

### Check container status
```bash
ssh ivanpod@92.205.19.23 "cd ~/docker/projects/resume-builder && docker compose -f docker-compose.prod.yml ps"
```

### View logs
```bash
ssh ivanpod@92.205.19.23 "cd ~/docker/projects/resume-builder && docker compose -f docker-compose.prod.yml logs -f"
```

### Restart services
```bash
ssh ivanpod@92.205.19.23 "cd ~/docker/projects/resume-builder && docker compose -f docker-compose.prod.yml restart"
```

### Nginx logs
```bash
ssh ivanpod@92.205.19.23 "sudo tail -f /var/log/nginx/resume_error.log"
```

## Port Summary

| Service | Internal Port | Exposed (localhost) |
|---------|---------------|---------------------|
| Frontend | 3000 | 127.0.0.1:4000 |
| Backend | 8080 | (internal only) |
| PostgreSQL | 5432 | (internal only) |
| Redis | 6379 | (internal only) |

Nginx proxies `resume.ivanpodgurskiy.name` → `127.0.0.1:4000`.
