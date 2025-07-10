# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based web service that provides an embeddable HTML component for displaying currently reading books from the Hardcover API. The project consists of a Go backend that proxies and caches Hardcover API requests, and a standalone HTML component for embedding.

## Development Commands

### Running the Application
```bash
# Set up environment variables first
cp .env.example .env
# Edit .env with your Hardcover API token

# Install dependencies and run
go mod tidy
go run cmd/server/main.go
```

### Building
```bash
go build -o hardcover-embed cmd/server/main.go
```

### Testing API Endpoints
```bash
curl http://localhost:8080/api/health
curl http://localhost:8080/api/books/currently-reading
```

## Architecture

### Backend Components
- **cmd/server/main.go** - Entry point, configuration, HTTP server setup
- **internal/hardcover/** - GraphQL client for Hardcover API with authentication
- **internal/cache/** - In-memory caching with TTL to respect API rate limits (60 req/min)
- **internal/api/** - HTTP handlers with CORS support for embedding

### Frontend
- **web/embed.html** - Self-contained embeddable component with styling and JavaScript

### Key Design Decisions
- Caching with 30-minute TTL to stay within Hardcover's 60 requests/minute limit
- CORS enabled for cross-domain embedding
- Environment-based configuration for security
- GraphQL queries target status_id: 2 (currently reading)
- Graceful error handling and loading states

## Configuration

Required environment variables:
- `HARDCOVER_API_TOKEN` - API token from Hardcover account settings
- Optional: `PORT`, `CACHE_TTL_MINUTES`, `ALLOWED_ORIGINS`