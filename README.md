# Todo List API

A Go backend API for a simple to-do list application.

## 🚧 Work in Progress

This project is currently under development, aiming to implement:

- **User Authentication** - Secure user registration and login
- **Database Integration** - Persistent storage for users and tasks
- **JWT Token Usage** - Stateless authentication with JSON Web Tokens

## Project Structure

```
todo-list/
├── main.go          # Application entry point
├── task/
│   └── task.go      # Task models and repository interface
└── user/
    └── user.go      # User models and repository interface
```

## Features

### Tasks
- Create, read, update, and delete tasks
- Task properties: name, description, status (high/moderate/low), label
- Timestamps for tracking creation time

### Users
- User registration and management
- Secure password storage (hashed)
- CRUD operations for user accounts

## Getting Started

### Prerequisites
- Go 1.25+

### Installation

```bash
git clone <repository-url>
cd todo-list
go mod download
```

### Running the Application

```bash
go run main.go
```

## License

MIT
