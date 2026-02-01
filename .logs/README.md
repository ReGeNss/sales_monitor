# Sales Monitor - Price Tracking System

A comprehensive price monitoring system for tracking product prices across multiple marketplaces (Fora, ATB, Silpo). Built as a monorepo with Go scraping services and Node.js REST API.

## Architecture

This is a polyglot monorepo with separated ecosystems:

```
sales-monitor/
├── scheduler/              # Go: Cron scheduler for scraping jobs
├── scraper_app/           # Go: Web scraper with Playwright
├── node/                  # Node.js workspace
│   ├── apps/
│   │   └── api-monitor/   # NestJS REST API
│   └── packages/
│       ├── database/      # Shared MikroORM entities
│       └── common/        # Shared types & utilities
├── internal/              # Go shared models
├── migrations/            # Atlas DB migrations (shared)
└── docker-compose.yaml    # Full stack orchestration
```

## Tech Stack

### Go Services
- **Scheduler**: Cron job orchestration with hot-reload
- **Scraper**: Playwright-based web scraping
- **Database**: GORM (MySQL/MariaDB)

### Node.js API
- **Framework**: NestJS 10+
- **ORM**: MikroORM 6+
- **Auth**: JWT with Passport.js
- **Validation**: class-validator
- **Documentation**: Swagger/OpenAPI
- **Database**: Shared MariaDB with Go services

### Infrastructure
- **Database**: MariaDB
- **Migrations**: Atlas
- **Containers**: Docker & Docker Compose
- **Monorepo**: npm workspaces

## Getting Started

### Prerequisites

- Docker & Docker Compose
- Node.js 20+ (for local development)
- Go 1.25+ (for local development)
- npm 7+ (for workspaces)

### Quick Start with Docker

1. **Clone and configure**:
```bash
git clone <repository>
cd "sales monitor"
cp .env.example .env
# Edit .env with your configuration
```

2. **Start all services**:
```bash
docker-compose up
```

Services will be available at:
- API: http://localhost:3000
- API Docs: http://localhost:3000/api/docs
- phpMyAdmin: http://localhost:8080
- MariaDB: localhost:3306

### Local Development

#### Node.js API Development

1. **Install dependencies**:
```bash
npm install
```

2. **Build shared packages**:
```bash
npm run build --workspace=@sales-monitor/database
npm run build --workspace=@sales-monitor/common
```

3. **Start API in dev mode**:
```bash
npm run dev:api
```

The API will start with hot-reload at http://localhost:3000

#### Go Services Development

```bash
# Run scheduler
cd scheduler
go run main.go

# Run scraper
cd scraper_app
go run main.go
```

## API Documentation

### Authentication

All protected endpoints require JWT Bearer token:
```bash
# Register
curl -X POST http://localhost:3000/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"login": "user", "password": "password123"}'

# Login (returns JWT token)
curl -X POST http://localhost:3000/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"login": "user", "password": "password123"}'

# Use token in subsequent requests
curl -X GET http://localhost:3000/api/users/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Main Endpoints

**Public Endpoints:**
- `GET /api/products` - List products with filters
- `GET /api/products/:id` - Product details
- `GET /api/products/:id/prices` - Price history
- `GET /api/categories` - List categories
- `GET /api/brands` - List brands
- `GET /api/marketplaces` - List marketplaces
- `GET /api/prices/latest` - Latest prices
- `GET /api/prices/trends` - Price trends

**Protected Endpoints (require JWT):**
- `GET /api/users/profile` - User profile
- `GET /api/favorites/products` - Favorite products
- `POST /api/favorites/products/:id` - Add to favorites
- `DELETE /api/favorites/products/:id` - Remove from favorites
- `GET /api/favorites/brands` - Favorite brands

Full API documentation with interactive testing: http://localhost:3000/api/docs

## Environment Variables

```env
# Database
DB_HOST=mariadb
DB_PORT=3306
DATABASE_NAME=sales_monitor
DB_USER_NAME=sales_user
DB_USER_PASSWORD=sales_password
DB_ROOT_PASSWORD=root_password

# API
API_PORT=3000
API_HOST=0.0.0.0
API_CORS_ORIGIN=http://localhost:3000,http://localhost:5173
API_PREFIX=api
JWT_SECRET=your-super-secret-jwt-key-change-in-production
NODE_ENV=development

# Scraper
SCRAPED_DATA_FOLDER=./scraped_data
SCRAPER_CONFIG_PATH=./scraper_config.yaml
```

## Database Migrations

Migrations are managed by Atlas:

```bash
# Apply migrations
docker-compose up migrate

# Create new migration (manual)
atlas migrate new --dir file://migrations
```

**Important**: MikroORM in the Node.js API is used only for data access. All schema changes must be done through Atlas migrations.

## Development Workflow

### Adding a new feature to API

1. Create/update entities in `node/packages/database/src/entities/`
2. Rebuild database package: `npm run build -w @sales-monitor/database`
3. Create service, controller, DTOs in appropriate module
4. Add Swagger decorators for documentation
5. Write tests
6. Update README if needed

### Running Tests

```bash
# Unit tests
npm run test -w @sales-monitor/api-monitor

# E2E tests
npm run test:e2e -w @sales-monitor/api-monitor

# Coverage
npm run test:cov -w @sales-monitor/api-monitor
```

## Project Structure Details

### Shared Database Package (`@sales-monitor/database`)
- Contains all MikroORM entities
- Shared between multiple Node.js apps
- Reflects MariaDB schema managed by Atlas

### Shared Common Package (`@sales-monitor/common`)
- TypeScript types
- Utility functions
- Constants
- Shared across all Node.js applications

### Workspaces Benefits
- Single `npm install` for all packages
- Shared dependencies
- Easy cross-package development
- Consistent tooling

## Troubleshooting

### Port already in use
```bash
# Change ports in .env
API_PORT=3001
DB_PORT=3307
```

### Database connection issues
```bash
# Check if MariaDB is healthy
docker-compose ps

# View logs
docker-compose logs mariadb
docker-compose logs api-monitor
```

### Build errors
```bash
# Clean and rebuild
rm -rf node_modules node/*/node_modules node/*/*/node_modules
npm install
npm run build:all
```

## Contributing

1. Follow existing code structure
2. Write tests for new features
3. Update documentation
4. Use TypeScript strictly
5. Add Swagger annotations
6. Keep shared packages DRY

## License

[Your License]

## Support

For issues and questions, please open an issue on GitHub.
