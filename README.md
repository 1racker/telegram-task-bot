# Telegram Task Manager Bot (Go)

A simple PET-project in Go: Telegram bot helps users plan their day, received daily reminders, track which tasks are completed and generate weekly performance reports.

## Features
-Add tasks with start time and duration
-Inline-buttons for marking statuses ("Started", "Completed", "Postponed", "Declined")
-Daily reminders via cron (the user sets the notification time)
-Automatic reminders (via internal cron scheduler)
-Webhook integration with Telegram Bot API
-Persistent storage using PostgresQL via GORM
-GitHub Actions CI/CD (lint, test, build, deploy)
-Auto-deploy on [Render](https://render.com)

## Tech Stack
| Layer | Technology |
|-------|-------------|
| **Language** | Go 1.25 |
| **Framework** | `telebot.v3` |
| **Database** | PostgreSQL (via GORM) |
| **Scheduler** | Built-in `cron` package |
| **CI/CD** | GitHub Actions + Render |
| **Containerization** | Docker (multi-stage build) |

## Project Structure

.
├── main.go # Entry point (webhook setup + bot launch)
├── config/ # Config loader (env variables)
├── handlers/ # Telegram commands & interactions
├── storage/ # PostgreSQL models & GORM setup
├── cron/ # Task scheduler
├── Dockerfile # Multi-stage build
├── .github/workflows/ci.yml # CI/CD pipeline
├── go.mod / go.sum # Dependencies
└── README.md

## Environment Variables

| Variable | Description | Example |
|-----------|--------------|----------|
| `TELEGRAM_BOT_TOKEN` | Telegram Bot API Token | `123456789:ABC...` |
| `WEBHOOK_URL` | Public Render URL | `https://telegram-task-bot-eqhx.onrender.com` |
| `PORT` | Webhook listening port | `8080` |
| `DB_URL` | PostgreSQL connection string | `postgres://user:pass@host:5432/dbname` |
| `TZ` | Timezone | `Europe/Bucharest` |
| `WEEKLY_REPORT_DAY` | Day of weekly summary | `Sunday` |

---

##  Local Setup

You can run the bot locally using either Docker or native Go.

Option A

```bash
docker build -t telegram-task-bot .
docker run -p 8080:8080 \
  -e TELEGRAM_BOT_TOKEN=your_token \
  -e WEBHOOK_URL=http://localhost:8080 \
  -e DB_URL=postgres://user:password@localhost:5432/tasksdb \
  -e TZ=Europe/Bucharest \
  -e WEEKLY_REPORT_DAY=Sunday \
  telegram-task-bot
  ```

Option B

```bash
go mod tidy
export TELEGRAM_BOT_TOKEN=your_token
export DB_URL=postgres://user:password@localhost:5432/tasksdb
export WEBHOOK_URL=https://<your_ngrok_url>
go run main.go
```

## Deployment on Render
This project uses Render Web Service + PostgreSQL Database.
Steps:
1. Create a Web Service on Render → connect your GitHub repo
2. Select environment: Docker
3. Add Environment Variables:
```ini
TELEGRAM_BOT_TOKEN=...
DB_URL=postgres://... (from Render PostgreSQL dashboard)
WEBHOOK_URL=https://<your-render-app>.onrender.com
TZ=Europe/Bucharest
WEEKLY_REPORT_DAY=Sunday
```
4. Deploy — Render will build using your Dockerfile automatically.
5. Check logs → you should see:
Connected to PostgreSQL successfully
Bot launched on port 8080 via webhook ...

## CI/CD Overflow
The bot is built and tested automatically via GitHub Actions:
-Linting (Dockerfile + Go code)
-Unit tests (if added)
-Build check
-Automatic deploy to Render after successful checks

Workflow file:
.github/workflows/ci.yml
```yaml

name: CI/CD
on:
  push:
    branches: [ "main" ]
  pull_request:

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: 1.25
      - run: go mod tidy
      - run: go build -v ./...
      - name: Lint Dockerfile
        uses: hadolint/hadolint-action@v3.1.0

  deploy:
    needs: build
    runs-on: ubuntu-latest
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: actions/checkout@v4
      - name: Trigger Render Deploy
        run: |
          curl -X POST "$RENDER_DEPLOY_HOOK"
```

## For Recruiters
This project demonstrates my skills in:
- Backend development with Go (webhooks, concurrency, database design)
- Docker-based deployment and CI/CD with GitHub Actions
- Using PostgreSQL with GORM ORM
- Cloud deployment with Render

The codebase is fully functional and can be deployed to any cloud provider supporting Docker.
