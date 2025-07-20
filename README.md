# WE-SPARK-CANVAS

A full-stack image sharing platform with decentralized storage using IPFS, built with Next.js and Go.

## Architecture

- **Frontend**: Next.js with React, Chakra UI, and Framer Motion
- **Backend**: Go with Gin framework
- **Database**: PostgreSQL
- **Storage**: IPFS for decentralized image storage

## Prerequisites

- Node.js (v18 or higher)
- Go (v1.23 or higher)
- PostgreSQL
- IPFS Desktop or daemon

## Installation Guide

### Installing Node.js and npm

#### Windows
1. Download Node.js from https://nodejs.org/
2. Run the installer and follow the setup wizard
3. Verify installation:
```bash
node --version
npm --version
```

#### macOS
Using Homebrew:
```bash
brew install node
```

Or download from https://nodejs.org/

#### Linux (Ubuntu/Debian)
```bash
# Using NodeSource repository
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Verify installation
node --version
npm --version
```

### Installing Go

#### Windows
1. Download Go from https://golang.org/dl/
2. Run the MSI installer
3. Add Go to your PATH (usually done automatically)
4. Verify installation:
```bash
go version
```

#### macOS
Using Homebrew:
```bash
brew install go
```

Or download from https://golang.org/dl/

#### Linux
```bash
# Download and extract
wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz

# Add to PATH (add to ~/.bashrc or ~/.profile)
export PATH=$PATH:/usr/local/go/bin

# Reload shell or source the file
source ~/.bashrc

# Verify installation
go version
```

### Installing PostgreSQL

#### Windows
1. Download PostgreSQL from https://www.postgresql.org/download/windows/
2. Run the installer
3. Remember the password you set for the postgres user
4. PostgreSQL service should start automatically

#### macOS
Using Homebrew:
```bash
brew install postgresql
brew services start postgresql

# Create a database user (optional, for development)
createuser -s postgres
```

#### Linux (Ubuntu/Debian)
```bash
# Install PostgreSQL
sudo apt update
sudo apt install postgresql postgresql-contrib

# Start PostgreSQL service
sudo systemctl start postgresql
sudo systemctl enable postgresql

# Switch to postgres user and create database
sudo -i -u postgres
createuser --interactive
createdb your_database_name
exit
```

### Installing IPFS

#### Windows
1. Download IPFS Desktop from https://github.com/ipfs-shipyard/ipfs-desktop/releases
2. Run the installer
3. Start IPFS Desktop
4. The daemon will run automatically

#### macOS
Using Homebrew:
```bash
brew install ipfs
```

Or download IPFS Desktop from https://github.com/ipfs-shipyard/ipfs-desktop/releases

#### Linux
```bash
# Download IPFS
wget https://dist.ipfs.tech/kubo/v0.24.0/kubo_v0.24.0_linux-amd64.tar.gz
tar -xvzf kubo_v0.24.0_linux-amd64.tar.gz
cd kubo
sudo bash install.sh

# Initialize IPFS
ipfs init

# Start daemon
ipfs daemon
```

Or install IPFS Desktop:
```bash
# Download AppImage
wget https://github.com/ipfs-shipyard/ipfs-desktop/releases/download/v0.28.0/ipfs-desktop-0.28.0-linux-x86_64.AppImage
chmod +x ipfs-desktop-0.28.0-linux-x86_64.AppImage
./ipfs-desktop-0.28.0-linux-x86_64.AppImage
```

## Setup Instructions

### 1. Clone the Repository

```bash
git clone <repository-url>
cd we-spark-canvas
```

### 2. Database Setup

1. Install PostgreSQL if not already installed
2. Create a database and user:

```sql
CREATE DATABASE weheartit_db;
CREATE USER weheartit_user WITH ENCRYPTED PASSWORD 'secretpassword';
GRANT ALL PRIVILEGES ON DATABASE weheartit_db TO weheartit_user;
```

3. Run migrations:

```bash
cd server
psql -h localhost -U weheartit_user -d weheartit_db -f migrations/001_auth_schema.sql
```

### 3. IPFS Setup

1. Install IPFS Desktop from https://github.com/ipfs-shipyard/ipfs-desktop/releases
2. Start IPFS daemon (usually runs on port 5001)
3. Verify IPFS is running:

```bash
curl http://127.0.0.1:5001/api/v0/version
```

### 4. Backend Setup

1. Navigate to server directory:

```bash
cd server
```

2. Install Go dependencies:

```bash
go mod download
```

3. Verify environment variables in `.env`:

```
PORT=8080
DATABASE_URL=postgres://weheartit_user:secretpassword@localhost:5432/weheartit_db?sslmode=disable
IPFS_API_URL=http://127.0.0.1:5001
```

4. Run the server:

```bash
go run main.go
```

The server will start on http://localhost:8080

### 5. Frontend Setup

1. Navigate to frontend directory:

```bash
cd frontend
```

2. Install dependencies:

```bash
npm install
```

3. Start development server:

```bash
npm run dev
```

The frontend will start on http://localhost:3000

## Development Commands

### Frontend Commands
- `npm run dev` - Start development server
- `npm run build` - Build for production  
- `npm run start` - Start production server
- `npm run lint` - Run ESLint

### Backend Commands
- `go run main.go` - Run development server
- `go test ./...` - Run all tests
- `go test ./handlers` - Run handler tests
- `go build` - Build binary

## Project Structure

```
we-spark-canvas/
├── frontend/                 # Next.js frontend
│   ├── components/          # React components
│   ├── pages/              # Next.js pages
│   ├── hooks/              # Custom React hooks
│   ├── lib/                # Utilities and configurations
│   └── styles/             # CSS styles
├── server/                  # Go backend
│   ├── handlers/           # HTTP handlers
│   ├── models/             # Data models
│   ├── storage/            # Database layer
│   ├── migrations/         # SQL migrations
│   └── main.go            # Server entry point
└── CLAUDE.md              # Development guide
```

## Features

- User authentication with JWT
- Image upload to IPFS
- Collection management
- Image hearting/favoriting
- Responsive design with Chakra UI
- Tag-based image filtering

## Design System

- **Colors**: Pastel palette (#FFE5F7, #E0F7E9, #F3E8FF) with neutrals
- **Typography**: Montserrat for headings, Inter for body text
- **Components**: Rounded corners (8-12px), generous whitespace
- **Interactions**: Subtle hover effects, card overlays, minimal animations

## Environment Variables

### Backend (.env)
- `PORT` - Server port (default: 8080)
- `DATABASE_URL` - PostgreSQL connection string
- `IPFS_API_URL` - IPFS API endpoint (default: http://127.0.0.1:5001)

## API Endpoints

- `GET /api/images` - Get all images
- `POST /api/images` - Upload new image
- `GET /api/collections` - Get user collections
- `POST /api/collections` - Create new collection
- `POST /api/auth/login` - User login
- `POST /api/auth/register` - User registration

## Troubleshooting

### Common Issues

1. **IPFS Connection Error**: Ensure IPFS daemon is running on port 5001
2. **Database Connection Error**: Verify PostgreSQL is running and credentials are correct
3. **Frontend Build Errors**: Run `npm install` to ensure all dependencies are installed
4. **CORS Issues**: Backend includes CORS middleware for development

### Logs

- Backend logs are output to console
- Frontend logs available in browser console
- IPFS logs available in IPFS Desktop

## Contributing

1. Create a feature branch from `development`
2. Make your changes following the code style guidelines in CLAUDE.md
3. Test your changes locally
4. Submit a pull request to `main` branch