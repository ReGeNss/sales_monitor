# Sales Monitor - Node.js Workspace

This workspace contains Node.js applications and shared packages for the Sales Monitor project.

## Structure

```
node/
├── apps/
│   └── api-monitor/          # NestJS REST API
└── packages/
    ├── database/             # MikroORM entities (shared)
    └── common/               # Common types and utilities (shared)
```

## Getting Started

### Prerequisites

- Node.js 20+
- npm 7+ (for workspaces support)
- MariaDB (via Docker)

### Installation

From the repository root:

```bash
npm install
```

This will install dependencies for all workspaces.

### Development

```bash
# Run API in development mode
npm run dev:api

# Build all packages
npm run build:all

# Run tests
npm run test:all
```

### Environment Variables

Create a `.env` file in the root directory with:

```env
DB_HOST=localhost
DB_PORT=3306
DATABASE_NAME=sales_monitor
DB_USER_NAME=sales_user
DB_USER_PASSWORD=sales_password

JWT_SECRET=your-secret-key
API_PORT=3000
API_CORS_ORIGIN=http://localhost:3000
NODE_ENV=development
```

## API Documentation

Once the API is running, access Swagger documentation at:
- http://localhost:3000/api/docs

## Available Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login
- `GET /api/auth/me` - Get current user (protected)

### Users (Protected)
- `GET /api/users/profile` - Get user profile
- `PUT /api/users/notification-token` - Update notification token
- `DELETE /api/users/account` - Delete account

### Favorites (Protected)
- `GET /api/favorites/products` - Get favorite products
- `POST /api/favorites/products/:id` - Add product to favorites
- `DELETE /api/favorites/products/:id` - Remove from favorites
- `GET /api/favorites/brands` - Get favorite brands
- `POST /api/favorites/brands/:id` - Add brand to favorites
- `DELETE /api/favorites/brands/:id` - Remove from favorites

### Products (Public)
- `GET /api/products` - Get all products (with filters)
- `GET /api/products/:id` - Get product by ID
- `GET /api/products/:id/prices` - Get product price history

### Prices (Public)
- `GET /api/prices/latest` - Get latest prices
- `GET /api/prices/trends` - Get price trends

### Categories (Public)
- `GET /api/categories` - Get all categories
- `GET /api/categories/:id` - Get category by ID
- `GET /api/categories/:id/products` - Get products in category

### Brands (Public)
- `GET /api/brands` - Get all brands
- `GET /api/brands/:id` - Get brand by ID
- `GET /api/brands/:id/products` - Get products by brand

### Marketplaces (Public)
- `GET /api/marketplaces` - Get all marketplaces
- `GET /api/marketplaces/:id` - Get marketplace by ID

## Architecture

### Workspaces
- **apps/api-monitor**: NestJS REST API application
- **packages/database**: Shared MikroORM entities and configuration
- **packages/common**: Shared TypeScript types, utilities, and constants

### Technology Stack
- **Framework**: NestJS 10+
- **ORM**: MikroORM 6+ (MySQL/MariaDB)
- **Authentication**: JWT with Passport.js
- **Validation**: class-validator + class-transformer
- **Documentation**: Swagger/OpenAPI
- **Testing**: Jest

### Database
- Uses MikroORM with existing MariaDB schema
- Schema is managed by Atlas migrations (in Go workspace)
- MikroORM is used only for data access, not migrations

## Docker

Run with Docker Compose:

```bash
docker-compose up api-monitor
```

This will start:
- MariaDB database
- Atlas migrations
- API Monitor service

## Development Tips

### Adding a new endpoint
1. Create DTOs in `dto/` folder
2. Implement service logic in `.service.ts`
3. Add controller endpoints in `.controller.ts`
4. Add Swagger documentation with decorators

### Adding a new entity
1. Create entity in `packages/database/src/entities/`
2. Export from `packages/database/src/entities/index.ts`
3. Rebuild database package: `npm run build -w @sales-monitor/database`

### Testing
```bash
# Unit tests
npm run test -w @sales-monitor/api-monitor

# E2E tests
npm run test:e2e -w @sales-monitor/api-monitor

# Coverage
npm run test:cov -w @sales-monitor/api-monitor
```

## Contributing

1. Follow existing code structure and naming conventions
2. Add proper TypeScript types
3. Write tests for new features
4. Update Swagger documentation
5. Keep shared packages DRY
