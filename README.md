# Online Learning Platform

This is a monorepo for a full-stack online learning platform.

## Table of Contents

- [Project Structure](#project-structure)
- [Getting Started](#getting-started)
- [Running the Apps](#running-the-apps)
- [Running the Services](#running-the-services)
- [Testing](#testing)
- [Contributing](#contributing)

## Project Structure

This repository is a monorepo containing several applications and services.

- `apps/`: Contains the front-end applications.
  - `web-app/`: The main Next.js web application for students.
  - `mobile-app/`: A React Native application for students.
  - `moderation-panel/`: A React application for moderators.
- `services/`: Contains the back-end microservices.
  - `api-gateway/`: The main entry point for all API requests.
  - `user-service/`: Manages user accounts and authentication.
  - `content-service/`: Manages courses and other educational content.
  - `...and many more.`
- `docs/`: Contains important documentation about the platform.
- `infrastructure/`: Contains Terraform scripts for deploying the platform.
- `scripts/`: Contains useful scripts for development.

## Getting Started

To get started, you'll need to install the dependencies for each application and service.

### Prerequisites

- Node.js (v18 or higher)
- Go (v1.20 or higher)
- Python (v3.9 or higher)
- Docker

### Installation

Clone the repository:

```bash
git clone <repository-url>
cd <repository-name>
```

Then, navigate to each service and application directory and install the dependencies. For example:

```bash
# For Node.js services/apps
cd apps/web-app
npm install

# For Go services
cd services/user-service
go mod tidy

# For Python services
cd services/qna-service
pip install -r requirements.txt
```

## Running the Apps

Each application has its own start script. For example, to run the web app:

```bash
cd apps/web-app
npm run dev
```

## Running the Services

The services can be run individually. For example, to run the user service:

```bash
cd services/user-service
go run main.go
```

It is recommended to use Docker Compose to run all the services at once. A `docker-compose.yml` file will be added soon.

## Testing

To run the tests for a specific service or application, navigate to its directory and run the test command.

```bash
# For Node.js services/apps with a test script
cd apps/mobile-app
npm test

# For Go services
cd services/user-service
go test ./...
```

## Contributing

Please read `docs/CONTENT_STANDARDS.md` and `docs/MODERATION_HANDBOOK.md` for guidelines on contributing to the platform.
