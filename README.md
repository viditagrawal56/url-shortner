# SecureLink - Secure URL Shortener

SecureLink is a modern, secure URL shortening service built with React and Go. It provides features like authentication-protected links and email notifications.

## Features

- **Secure Short Links**: Create shortened URLs with optional authentication protection
- **User Authentication**: Secure user registration and login system
- **Access Control**: Control who can access your shortened links
- **Email Notifications**: Get notified when protected links are accessed
- **Dashboard**: User-friendly dashboard to manage all your links

## Technology Stack

### Frontend
- React 19.0 with TypeScript
- Vite 6.2 for development and building
- TailwindCSS for styling
- React Router DOM for navigation
- React Toastify for notifications
- Lucide React for icons

### Backend
- Go 1.24
- Chi router for HTTP routing
- GORM with PostgreSQL for database
- JWT for authentication
- SMTP for email notifications
- Godotenv for environment configuration

## Prerequisites

- Node.js 18+ 
- Go 1.24+
- PostgreSQL 15+
- SMTP server access for email notifications

## Setup & Installation

### Backend Setup

1. Navigate to the server directory:
```bash
cd server
```

2. Start PostgreSQL database using Docker Compose:
```bash
docker-compose -f docker-compose.dev.yaml up -d
```

This will start a PostgreSQL container with the following default credentials:
- Database: urlshortnerdb
- Username: postgres
- Password: postgres
- Port: 5432

3. Create a `.env` file based on provided template:
```bash
# Server Configuration
SERVER_PORT=8080
SERVER_BASE_URL=http://localhost:8080

# Database Configuration
DB_CONNECTION_STRING=postgresql://postgres:postgres@localhost:5432/urlshortnerdb
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=25
DB_CONN_MAX_LIFETIME=5m

# Email Configuration
EMAIL_SMTP_HOST=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_SMTP_USERNAME=your-email@gmail.com
EMAIL_SMTP_PASSWORD=your-app-password
EMAIL_FROM=your-email@gmail.com
```

4. Install dependencies and run:
```bash
go mod tidy
go run cmd/main.go
```

Note: To stop the database container, run:
```bash
docker-compose -f docker-compose.dev.yaml down
```
### Frontend Setup

1. Navigate to the client directory:
```bash
cd client
```

2. Install dependencies:
```bash
npm install
```

3. Start development server:
```bash
npm run dev
```

## Project Structure

### Frontend Structure
```
client/
├── src/
│   ├── components/       # Reusable React components
│   ├── context/         # React context providers
│   ├── pages/           # Page components
│   ├── App.tsx          # Main app component
│   └── main.tsx         # Entry point
├── public/              # Static assets
└── package.json         # Project configuration
```

### Backend Structure
```
server/
├── cmd/
│   └── main.go          # Application entry point
├── internal/
│   ├── auth/            # Authentication logic
│   ├── config/          # Configuration handling
│   ├── db/              # Database operations
│   ├── email/           # Email service
│   ├── handlers/        # HTTP handlers
│   ├── models/          # Data models
│   └── urlShortner/     # Core URL shortening logic
└── go.mod               # Go module file
```
