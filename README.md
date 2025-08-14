# Chat Server

A real-time chat server built in Go with robust connection management and messaging features.

## Features

### Core Functionality

- **Real-time messaging** via WebSocket connections
- **User authentication** with JWT-based sessions
- **Chat rooms** supporting both direct (private) and group conversations
- **Message persistence** with PostgreSQL database
- **RESTful API** with OpenAPI 3.0 specification

### Technical Features

- **Connection management** with graceful establishment, maintenance, and termination
- **Scalable WebSocket hub** for handling multiple concurrent connections
- **Comprehensive logging** and error handling
- **Docker containerization** for easy deployment

## 🛠️ Tech Stack

- **Language**: Go 1.24.2
- **Web Framework**: Chi Router
- **Real-time**: Gorilla WebSocket
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Authentication**: JWT with bcrypt
- **API Documentation**: OpenAPI 3.0 + Swagger UI
- **Containerization**: Docker & Docker Compose

## ️ Architecture

The project follows clean architecture principles with clear separation of concerns:

```
internal/app/
├── handlers/     # HTTP request handlers
├── services/     # Business logic layer
├── repositories/ # Data access layer
├── models/       # Data models
├── middlewares/  # HTTP middlewares
├── websocket/    # WebSocket handling
└── validations/  # Input validation
```

<code_block_to_apply_changes_from>
chatserver/
├── api/ # OpenAPI specifications
├── cmd/ # Application entry points
├── configs/ # Configuration management
├── internal/ # Private application code
├── docker-compose.yml
└── Dockerfile

````

### Running Locally
```bash
# Install dependencies
go mod download

# Set up database connection
# Configure your .env file

# Run the server
go run cmd/main.go
````

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]

## 🚀 Quick Start

### Prerequisites

- Go 1.24.2+
- Docker & Docker Compose
- PostgreSQL 15

### Running with Docker

```bash
# Clone the repository
git clone <repository-url>
cd chatserver

# Set up environment variables
cp .env.example .env
# Edit .env with your configuration

# Start the services
docker-compose up -d
```

The server will be available at `http://localhost:8080`

### API Documentation

- **OpenAPI Spec**: `http://localhost:8080/openapi.yaml`
- **Swagger UI**: `http://localhost:8080/docs`

## API Endpoints

### Authentication

- `POST /auth/signup` - User registration
- `POST /auth/signin` - User login

### User Management

- `GET /user/me` - Get current user info

### Chat Rooms

- `GET /rooms` - List available rooms
- `POST /rooms` - Create new room
- `POST /rooms/direct` - Send direct message

### WebSocket

- `GET /ws/room/{roomID}` - Join chat room

