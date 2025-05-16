# M-RPG Game Server

A Go-based game server template with PostgreSQL ORM, featuring user authentication, daily rewards, and WebSocket support.

## Features

- User registration and authentication using JWT
- Daily reward system with streak bonuses
- WebSocket support for real-time communication
- PostgreSQL database with GORM ORM
- RESTful API endpoints
- Authorization using JWT in headers

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Make (optional, for using Makefile commands)

## Setup

1. Clone the repository:
```bash
git clone https://github.com/m-rpg/2505-server.git
cd 2505-server
```

2. Create a `.env` file based on `.env.example`:
```bash
cp .env.example .env
```

3. Update the `.env` file with your database credentials and other settings.

4. Install dependencies:
```bash
go mod download
```

5. Run database migrations:
```bash
go run main.go migrate
```

6. Start the server:
```bash
go run main.go
```

## API Endpoints

### Authentication

- `POST /api/auth/register` - Register a new user
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```

- `POST /api/auth/login` - Login user
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```

### Protected Routes (Requires Authorization Header)

- `GET /api/profile` - Get user profile
- `GET /api/daily-reward` - Check daily reward status
- `POST /api/daily-reward/claim` - Claim daily reward

### WebSocket

- `GET /api/ws` - WebSocket endpoint for real-time communication
  - Requires Authorization header with JWT token

## WebSocket Message Format

```json
{
  "type": "string",
  "payload": {}
}
```

## Development

### Project Structure

```
.
├── main.go           # Application entry point
├── models/           # Database models
├── handlers/         # HTTP and WebSocket handlers
├── middleware/       # Custom middleware
├── config/          # Configuration files
└── migrations/      # Database migrations
```

### Running Tests

```bash
go test ./...
```

### Building

```bash
go build -o bin/server
```

## License

MIT License 
