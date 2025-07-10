# Hardcover Book Embed

An embeddable HTML component that displays currently reading or last read books from any Hardcover user, with a Go backend for API proxying and caching.

## Features

- 📚 Displays currently reading or last read books from any Hardcover user
- 🚀 Go backend with caching to respect API rate limits
- 🎨 Responsive, embeddable HTML component
- 🔒 Secure API token handling (server-side only)
- ⚡ 30-minute caching by default (configurable)
- 🌐 CORS support for cross-domain embedding
- 👥 Support for multiple users on the same page

## Quick Start

### 1. Get Your Hardcover API Token

1. Go to your [Hardcover account settings](https://hardcover.app/account/api)
2. Generate an API token
3. Keep this token secure - it will be used in your environment variables

### 2. Set Up Environment Variables

Copy the example environment file:
```bash
cp .env.example .env
```

Edit `.env` and add your Hardcover API token:
```bash
HARDCOVER_API_TOKEN=your_hardcover_api_token_here
PORT=8080
CACHE_TTL_MINUTES=30
ALLOWED_ORIGINS=*
```

### 3. Run the Server

```bash
# Using Make (recommended)
make setup   # First time setup - creates .env from .env.example
make run     # Build and run the server

# Or manually
go mod tidy
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`.

### 4. Embed the Component

You can embed the component in any website. See [EMBEDDING.md](EMBEDDING.md) for detailed instructions.

Quick example:

```html
<!-- Currently reading books -->
<div data-hardcover-widget data-api-url="http://localhost:8080" data-username="your-username"></div>

<!-- Last read books -->
<div data-hardcover-widget data-api-url="http://localhost:8080" data-username="your-username" data-book-type="last-read"></div>

<script src="http://localhost:8080/widget.js"></script>
```

Or use the API directly:
```html
<script>
    // Currently reading books
    fetch('http://localhost:8080/api/books/currently-reading/your-username')
        .then(response => response.json())
        .then(data => {
            console.log('Currently reading:', data.books);
        });
    
    // Last read books
    fetch('http://localhost:8080/api/books/last-read/your-username')
        .then(response => response.json())
        .then(data => {
            console.log('Last read:', data.books);
        });
</script>
```

## API Endpoints

- `GET /api/books/currently-reading/:username` - Returns currently reading books for a user
- `GET /api/books/last-read/:username` - Returns last read books for a user  
- `GET /api/health` - Health check endpoint
- `GET /embed.html` - Embeddable HTML component
- `GET /widget.js` - JavaScript widget for embedding

## Configuration

Environment variables:

- `HARDCOVER_API_TOKEN` (required) - Your Hardcover API token
- `PORT` (optional) - Server port (default: 8080)
- `CACHE_TTL_MINUTES` (optional) - Cache duration in minutes (default: 30)
- `ALLOWED_ORIGINS` (optional) - CORS allowed origins (default: *)

## Development

### Project Structure

```
├── cmd/server/          # Server entry point
├── internal/
│   ├── api/            # HTTP handlers
│   ├── hardcover/      # Hardcover API client
│   └── cache/          # Caching layer
├── test/               # Test scripts
├── web/                # Static files
├── Makefile            # Build automation
├── Dockerfile          # Container definition
├── .air.toml           # Auto-reload config
└── .env.example        # Environment template
```

### Building

```bash
# Using Make
make build       # Build the application
make build-all   # Build for multiple platforms

# Or manually
go build -o hardcover-embed cmd/server/main.go
```

### Testing

```bash
# Run API tests (provide a username)
./test/test.sh your-username

# Or manually test the API
curl http://localhost:8080/api/health
curl http://localhost:8080/api/books/currently-reading/your-username
```

### Development Commands

```bash
make dev         # Run with auto-reload (requires air)
make fmt         # Format code
make lint        # Run linter
make test        # Run tests
make clean       # Clean build artifacts
make help        # Show all available commands
```

## Deployment

### Docker

```bash
# Build Docker image
make docker-build

# Run Docker container
make docker-run

# Or manually
docker build -t hardcover-embed:latest .
docker run -p 8080:8080 --env-file .env hardcover-embed:latest
```

### Environment Variables for Production

Ensure these are set in your production environment:
- `HARDCOVER_API_TOKEN` - Your Hardcover API token
- `PORT` - Port for the server to listen on
- `ALLOWED_ORIGINS` - Comma-separated list of allowed origins for CORS

## Customization

The HTML component includes CSS custom properties for easy theming:

```css
:root {
    --primary-color: #2563eb;
    --text-color: #374151;
    --bg-color: #ffffff;
    /* ... and more */
}
```

## Rate Limiting

The Hardcover API has a rate limit of 60 requests per minute. This server implements caching with a default TTL of 30 minutes to ensure you stay well within these limits.

## License

MIT License - see LICENSE file for details.