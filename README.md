# RedCompass

RedCompass is a web application designed to help you manage and track domains across your various projects. It provides a centralized interface for monitoring domain health, categorization status, and project assignments.

## Features

- **Domain Management**
  - Track domain registration and expiration dates
  - Monitor domain health status
  - Manage domain locks and auto-renewal settings
  - Tag domains for better organization
  - Track DNS records for each domain

- **Project Organization**
  - Assign domains to specific projects
  - Track project progress and completion
  - Manage team member assignments
  - View project history and previously assigned domains

- **Domain Categorization**
  - Track domain categorization across multiple services
  - Submit categorization requests
  - Monitor categorization status
  - Follow best practices for domain categorization

## Quick Start with Docker

The easiest way to get started is using Docker:

```bash
# Clone the repository
git clone https://github.com/lum8rjack/redcompass
cd redcompass

# Build and run with Docker Compose
docker compose up -d

# Or build and run manually
docker build -t redcompass .
docker run -d -p 8090:8090 --name redcompass redcompass
```

The application will be available at http://localhost:8090

## Manual Setup

If you prefer to run the application without Docker, you'll need:

### Prerequisites
- Go (for backend)
- Node.js (for frontend)
  - Vue.js
  - TailwindCSS
  - PocketBase

### Installation Steps

1. Clone the repository:
```bash
git clone https://github.com/lum8rjack/redcompass
cd redcompass
```

2. Install frontend dependencies:
```bash
cd frontend
npm install
```

3. Build the frontend:
```bash
npm run build
```

4. Build and run the backend:
```bash
cd ..
go build
./redcompass serve --http "0.0.0.0:8090"
```

## Development

To work on the frontend in development mode:

```bash
cd frontend
npm run dev
```

For backend development, you can use:

```bash
go run . serve --http "0.0.0.0:8090"
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

