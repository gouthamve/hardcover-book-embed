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
- **web/static/embed.html** - Self-contained embeddable component with styling and JavaScript

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

## Frontend Testing and Debugging

### Using Playwright MCP Server

When troubleshooting HTML, CSS, and JavaScript issues in the application frontend, use the Playwright MCP server tools:

1. **Starting the application**:
   ```bash
   go run cmd/server/main.go
   ```

2. **Accessing frontend pages**:
   - Navigate to `http://localhost:8080/` or the configured port
   - Use `mcp__playwright__browser_navigate` to open pages
   - Use `mcp__playwright__browser_snapshot` to capture the current page state
   - Use `mcp__playwright__browser_console_messages` to view JavaScript console output

3. **Common debugging workflows**:
   - **HTML/CSS issues**: Use `mcp__playwright__browser_snapshot` to inspect the DOM structure
   - **JavaScript errors**: Check `mcp__playwright__browser_console_messages` for error logs
   - **Interactive debugging**: Use `mcp__playwright__browser_evaluate` to run JavaScript in the page context
   - **Network issues**: Monitor API calls with `mcp__playwright__browser_network_requests`

4. **Testing the embed component**:
   - Test the embed at `/web/static/embed.html`
   - Verify CORS headers are working properly
   - Check that the component loads and displays book data correctly