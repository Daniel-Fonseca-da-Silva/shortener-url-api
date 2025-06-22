# 🔗 URL Shortener API

A simple and efficient URL shortening API developed in Go, with AES encryption, Redis-based rate limiting, and structured logging.

## 📋 Table of Contents

- [Features](#-features)
- [Architecture](#-architecture)
- [Technologies](#-technologies)
- [Installation](#-installation)
- [Usage](#-usage)
- [API Endpoints](#-api-endpoints)
- [Docker](#-docker)
- [Development](#-development)
- [Security](#-security)
- [Limitations](#-limitations)
- [Contributing](#-contributing)

## ✨ Features

- **URL Shortening**: Converts long URLs into short links
- **AES Encryption**: URLs are encrypted before storage
- **Automatic Redirection**: Direct access to original URLs
- **Structured Logging**: Detailed logs with Zap
- **Input Validation**: Valid URL verification
- **Thread-safe**: Concurrent access with mutex
- **Rate Limiting**: Performed using Redis
- **Containerization**: Docker and Docker Compose ready

## 🏗️ Architecture

### Workflow

```
1. Original URL → Validation → AES Encryption → Short ID → Storage
2. http://localhost:8080/ID → Search → Decryption → Redirection
```

### Components

- **HTTP Server**: Web server on port 8080
- **URL Store**: In-memory map for storage
- **Crypto Engine**: AES-CTR encryption/decryption
- **ID Generator**: Random short ID generator
- **Logger**: Structured logging system
- **Rate Limiter**: Uses Redis for request limiting
- **Redis**: Used for rate limiting (service runs as a container)

## 🛠️ Technologies

- **Go 1.24.1**: Main language
- **Zap**: Structured logging
- **AES-CTR**: Encryption
- **Redis**: Rate limiting backend
- **Docker**: Containerization
- **Alpine Linux**: Optimized base image

## 🚀 Installation

### Prerequisites

- Go 1.24.1 or higher
- Docker (optional)

### Local Installation

```bash
# Clone the repository
git clone https://github.com/Daniel-Fonseca-da-Silva/shortener-url-api.git
cd shortener-url-api

# Install dependencies
go mod download

# Run the application (requires a local Redis instance running on port 6379)
go run main.go
```

### Docker Installation

```bash
# Build image
docker build -t url-shortener .

# Run (requires a local Redis instance running on port 6379)
docker run -p 8080:8080 url-shortener
```

### Docker Compose Installation

```bash
# Build and run all services (including Redis)
docker compose up --build

# Run in background
docker compose up -d
```

## 📖 Usage

### URL Shortening

```bash
# Via curl
curl "http://localhost:8080/shorten?url=https://google.com"

# Response
The shortened url is: http://localhost:8080/AbC123
```

### Accessing Shortened URL

```bash
# Via curl (with redirect)
curl -L "http://localhost:8080/AbC123"

# Via browser
http://localhost:8080/AbC123
```

## 🔌 API Endpoints

### GET /shorten

Shortens a URL.

**Parameters:**
- `url` (query): URL to be shortened (required)

**Example:**
```bash
GET /shorten?url=https://example.com
```

**Success Response:**
```
HTTP/1.1 200 OK
Content-Type: text/plain

The shortened url is: http://localhost:8080/AbC123
```

**Error Response:**
```
HTTP/1.1 400 Bad Request
Content-Type: text/plain

URL parameter in query is required
```

### GET /

Redirects to the original URL.

**Parameters:**
- `{shortId}` (path): Short ID of the URL

**Example:**
```bash
GET /AbC123
```

**Response:**
```
HTTP/1.1 302 Found
Location: https://example.com
```

## 🐳 Docker

### Build Image

```bash
docker build -t url-shortener .
```

### Run

```bash
# Simple execution (requires a local Redis instance)
docker run -p 8080:8080 url-shortener
```

### Docker Compose

```bash
# Complete execution (runs both the app and Redis)
docker compose up --build

# Run in background
docker compose up -d

# View logs
docker compose logs -f

# Stop services
docker compose down
```

#### Example docker-compose.yml

```yaml
version: '3.8'

services:
  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    restart: unless-stopped
    networks:
      - url-shortener-network
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s

  url-shortener:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    restart: unless-stopped
    depends_on:
      redis:
        condition: service_healthy
    healthcheck:
      test: ["CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8080/"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 40s
    networks:
      - url-shortener-network

networks:
  url-shortener-network:
    driver: bridge
```

**Note:**
- The Go application is configured to connect to Redis using the hostname `redis:6379` (the service name in Docker Compose).
- No environment variables are required for Redis configuration.

## 🔧 Development

### Project Structure

```
shortener-url-api/
├── main.go              # Main file
├── go.mod               # Go dependencies
├── go.sum               # Dependency checksums
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker orchestration
├── .dockerignore        # Docker ignored files
└── README.md            # Documentation
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |

### Logs

The system uses structured logging with different levels:

- **INFO**: Main operations (server startup, shortened URLs)
- **WARN**: Warnings (invalid URLs, IDs not found)
- **DEBUG**: Detailed information (encryption, ID generation)
- **FATAL**: Critical errors (initialization failure)

## 🔒 Security

### Security Features

- **AES-CTR Encryption**: URLs are encrypted before storage
- **Input Validation**: Valid URL verification
- **Non-root User**: Container runs with non-privileged user
- **Health Checks**: Application health monitoring
- **Rate Limiting**: Performed using Redis

### Security Limitations

⚠️ **Warning**: This is an educational project with limitations:

- **Hardcoded Key**: Encryption key is in the code
- **In-memory Storage**: URLs are lost on restart
- **No Authentication**: No access control
- **Hardcoded Host**: Always returns `localhost:8080`

## ⚠️ Limitations

### Current Limitations

1. **Persistence**: URLs are stored only in memory
2. **Scalability**: Limited by memory size
3. **Collisions**: Possibility of duplicate IDs (very low)
4. **Configuration**: Host and port hardcoded
5. **Validation**: Basic URL verification

### Future Improvements

- [ ] Database persistence
- [ ] Environment variable configuration
- [ ] More robust URL validation
- [ ] URL expiration system
- [ ] Metrics and monitoring
- [ ] Authentication and authorization
- [ ] Rate limiting
- [ ] Distributed cache

## 🤝 Contributing

### How to Contribute

1. Fork the project
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### Code Standards

- Follow Clean Code principles
- Apply SOLID concepts
- Use Object Calisthenics
- Keep tests updated
- Document important changes

## 📄 License

This project is under the MIT license. See the [LICENSE](LICENSE) file for more details.

## 👨‍💻 Author

**Daniel Fonseca da Silva**

- GitHub: [@Daniel-Fonseca-da-Silva](https://github.com/Daniel-Fonseca-da-Silva)

## 🙏 Acknowledgments

- Go Community
- Uber Zap for logging
- Alpine Linux for optimized images
- Docker for containerization

---

⭐ If this project was useful to you, consider giving it a star!
