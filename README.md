# Chirpy HTTP Server
A social media backend for the Chirpy application built with Go.

## Features
- User registration and authentication with JWT.
- Refresh token support for secure sessions.
- Chirp creation, retrieval, and deletion.
- Webhook integration for user upgrades.
- Admin metrics and database reset capabilities.

## Setup

### Prerequisites
- Go 1.21+
- PostgreSQL database

### Installation
1. Clone the repository.
2. Install dependencies: go mod download.
3. Set up environment variables in a .env file:
   - DB_URL: PostgreSQL connection string.
   - JWT_SECRET: Secret key for token signing.
   - PLATFORM: Set to "dev" to enable reset endpoints.
   - POLKA_KEY: API key for Polka webhooks.

### Running
Build and run the server:
go build -o out/server && ./out/server

The server starts on port 8080.

## Usage
- Access static files at /app/.
- API endpoints are under /api/.
- Admin metrics are available at /admin/metrics.
