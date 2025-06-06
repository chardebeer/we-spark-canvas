# WE-SPARK-CANVAS DEV GUIDE

## Commands
- Frontend:
  - `cd frontend && npm run dev` - Start Next.js development server
  - `cd frontend && npm run build` - Build for production
  - `cd frontend && npm run lint` - Run ESLint
- Backend:
  - `cd server && go run main.go` - Run Go server
  - `cd server && go test ./...` - Run all tests
  - `cd server && go test ./handlers` - Run specific package tests

## Code Style
- Frontend:
  - Use ES modules with named exports
  - Follow Next.js patterns for pages and components
  - Use React hooks for state management
  - Implement Chakra UI components with project color palette
  - Use Framer Motion for animations and transitions
- Backend:
  - Follow Go standard layout (handlers, models, storage)
  - Use proper error handling with context
  - Keep functions small and focused
  - Write docstrings for exported functions

## Design System
- Colors: Use pastel palette (#FFE5F7, #E0F7E9, #F3E8FF) with neutrals (#FAFAFA, #F0F0F0, #333333)
- Typography: Montserrat for headings, Inter for body text
- Components: Rounded corners (8-12px), generous whitespace (24-32px margins)
- Interactions: Subtle hover effects, card overlays, minimal animations

## Architecture
- Next.js frontend with React and Chakra UI
- Go backend with Gin framework
- PostgreSQL database for metadata
- IPFS integration for decentralized image storage