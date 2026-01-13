# Fired Calendar - Countdown to Freedom

A simple web application that allows users to track their working days until they get fired. Users mark their workdays and see how many days remain until their expected termination date.

## Features

- User registration/login with 6-word mnemonic phrases
- Workday tracking with visual indicators
- "Mark Today" button for quick marking
- Settings page with fired date modification
- Account deletion with 7-day restoration window
- Responsive UI with Bootstrap
- Russian/English language switching
- Backend logging for each request

## Build and Deploy on Ubuntu from macOS

### Method 0: Static Binary Compilation (No Docker Required)

To build a static binary with SQLite support that runs on Ubuntu:

```bash
# Install GCC for Linux cross-compilation (if available)
brew install FiloSottile/musl-cross/musl-cross

# Set up the environment and build a static binary:
CC=x86_64-linux-musl-gcc CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags "-static" -w -s' -o fired-calendar-linux .
```

Then transfer the `fired-calendar-linux` binary to your Ubuntu server and run it.

### Method 1: Using Docker (Ubuntu-based, Recommended)

1. Install Docker on your macOS system
2. Navigate to the project directory
3. Build and run the application:

```bash
docker-compose up --build
```

The application will be accessible at `http://localhost:8080`

### Method 2: Direct Binary Deployment

1. Build the Docker image on macOS:

```bash
docker build -t fired-calendar .
```

2. Export the binary from the container:

```bash
# Run the container temporarily
docker run --name temp-container fired-calendar

# Copy the binary from the container
docker cp temp-container:/app/fired-calendar ./fired-calendar-linux

# Stop and remove the temporary container
docker rm temp-container
```

3. Transfer the `fired-calendar-linux` binary to your Ubuntu server

4. On the Ubuntu server, make sure you have the required dependencies:

```bash
sudo apt-get update
sudo apt-get install ca-certificates
```

5. Run the application on Ubuntu:

```bash
chmod +x fired-calendar-linux
./fired-calendar-linux
```

### Method 3: Building with Cross-compilation (Alternative)

If you have the necessary build tools installed on macOS:

```bash
# Install GCC for Linux cross-compilation (if available)
brew install FiloSottile/musl-cross/musl-cross

# Then set up the environment and build:
CC=x86_64-linux-musl-gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o fired-calendar-linux .
```

## Configuration

The application uses the following environment variables:

- `PORT`: Port to run the server on (default: 8080)
- `DATABASE_PATH`: Path to the SQLite database file (default: ./fired_calendar.db)
- `SESSION_KEY`: Secret key for session encryption (default: your-secret-key-change-this-in-production)

## Architecture

- Backend: Go with SQLite database
- Frontend: HTML/CSS/JavaScript with Bootstrap and Alpine.js
- Deployment: Docker container with static file serving
- Database: SQLite with encrypted user data and calendar entries

## Security

- Sessions stored securely with encrypted cookies
- User passwords replaced with 6-word mnemonic phrases
- Account restoration possible within 7 days after deletion
- Input validation and sanitization

## Technologies Used

- Go (Golang) for backend
- SQLite for database
- Bootstrap for responsive UI
- Alpine.js for frontend interactivity
- Docker for containerization
- Nginx for reverse proxy (configuration included)
