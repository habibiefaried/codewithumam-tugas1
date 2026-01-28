# Category CRUD API

A REST API for managing categories with PostgreSQL database backend built with Go.

![CI Tests](https://github.com/habibiefaried/codewithumam-tugas1/workflows/CI%20Tests/badge.svg)

## Features

- ✅ Full CRUD operations for categories
- ✅ PostgreSQL database backend
- ✅ Configuration from `secrets.yml` or environment variables
- ✅ Database migrations with `CREATE IF NOT EXISTS`
- ✅ Comprehensive unit tests for database queries
- ✅ RESTful API design with proper separation of concerns
- ✅ JSON request/response
- ✅ Automated CI tests
- ✅ Deployed on Railway

## Live API

🚀 **Production URL:** https://codewithumam-tugas1-production.up.railway.app/

## Getting Started

### Prerequisites

- Go 1.22 or higher
- PostgreSQL 12 or higher

### Installation

1. Clone the repository
```bash
git clone https://github.com/habibiefaried/codewithumam-tugas1
cd codewithumam-tugas1
```

2. Download dependencies
```bash
go mod download
```

3. Setup configuration - Copy the example file and update with your database credentials:
```bash
cp secrets.yml.example secrets.yml
```

Edit `secrets.yml` with your database details:
```yaml
db_url: localhost
db_port: 5432
db_name: your_database
db_user: your_user
db_password: your_password
port: 8080
```

4. Run the server
```bash
go run main.go
```

The server will automatically create the `category` table if it doesn't exist.

### Build

```bash
go build -o api-server main.go
./api-server
```

## API Endpoints

### Base URLs

**Production:** `https://codewithumam-tugas1-production.up.railway.app`  
**Local Development:** `http://localhost:8080`

### Health Check
```bash
# Production
curl https://codewithumam-tugas1-production.up.railway.app/health

# Local
curl http://localhost:8080/health
```

**Response:**
```
OK
```

### Version
```bash
# Production
curl https://codewithumam-tugas1-production.up.railway.app/version

# Local
curl http://localhost:8080/version
```

**Response:**
```
Commit: unknown
```

---

## Category Endpoints

### 1. Get All Categories

**Endpoint:** `GET /categories`

**Request:**
```bash
# Production
curl https://codewithumam-tugas1-production.up.railway.app/categories

# Local
curl http://localhost:8080/categories
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Electronics",
    "description": "Electronic devices and gadgets"
  },
  {
    "id": 2,
    "name": "Books",
    "description": "Physical and digital books"
  }
]
```

---

### 2. Get Category by ID

**Endpoint:** `GET /categories/{id}`

**Request:**
```bash
# Production
curl https://codewithumam-tugas1-production.up.railway.app/categories/1

# Local
curl http://localhost:8080/categories/1
```

**Response (Success - 200):**
```json
{
  "id": 1,
  "name": "Electronics",
  "description": "Electronic devices and gadgets"
}
```

**Response (Not Found - 404):**
```
Category not found
```

---

### 3. Create Category

**Endpoint:** `POST /categories`

**Request:**
```bash
# Production
curl -X POST https://codewithumam-tugas1-production.up.railway.app/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electronics",
    "description": "Electronic devices and gadgets"
  }'

# Local
curl -X POST http://localhost:8080/categories \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electronics",
    "description": "Electronic devices and gadgets"
  }'
```

**Response (Success - 201):**
```json
{
  "id": 1,
  "name": "Electronics",
  "description": "Electronic devices and gadgets"
}
```

**Response (Bad Request - 400):**
```
Name is required
```

---

### 4. Update Category

**Endpoint:** `PUT /categories/{id}`

**Request:**
```bash
# Production
curl -X PUT https://codewithumam-tugas1-production.up.railway.app/categories/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electronics and Tech",
    "description": "Electronic devices, gadgets, and technology products"
  }'

# Local
curl -X PUT http://localhost:8080/categories/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Electronics and Tech",
    "description": "Electronic devices, gadgets, and technology products"
  }'
```

**Response (Success - 200):**
```json
{
  "id": 1,
  "name": "Electronics and Tech",
  "description": "Electronic devices, gadgets, and technology products"
}
```

**Response (Not Found - 404):**
```
Category not found
```

---

### 5. Delete Category

**Endpoint:** `DELETE /categories/{id}`

**Request:**
```bash
# Production
curl -X DELETE https://codewithumam-tugas1-production.up.railway.app/categories/1

# Local
curl -X DELETE http://localhost:8080/categories/1
```

**Response (Success - 204):**
```
(No content)
```

**Response (Not Found - 404):**
```
Category not found
```

---

## Quick Testing Examples

### Complete Workflow (Production)

```bash
# 1. Create a category
curl -X POST https://codewithumam-tugas1-production.up.railway.app/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Electronics","description":"Electronic devices"}'

# 2. Create another category
curl -X POST https://codewithumam-tugas1-production.up.railway.app/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Books","description":"Physical and digital books"}'

# 3. List all categories
curl https://codewithumam-tugas1-production.up.railway.app/categories

# 4. Get specific category
curl https://codewithumam-tugas1-production.up.railway.app/categories/1

# 5. Update category
curl -X PUT https://codewithumam-tugas1-production.up.railway.app/categories/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Electronics and Tech","description":"Updated description"}'

# 6. Delete category
curl -X DELETE https://codewithumam-tugas1-production.up.railway.app/categories/2

# 7. Verify deletion
curl https://codewithumam-tugas1-production.up.railway.app/categories
```

### Using jq for Pretty Output

If you have `jq` installed, you can format the JSON output:

```bash
curl https://codewithumam-tugas1-production.up.railway.app/categories | jq
```

---

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `204 No Content` - Resource deleted successfully
- `400 Bad Request` - Invalid request (missing required fields, invalid ID format)
- `404 Not Found` - Resource not found
- `500 Internal Server Error` - Server error

---

## Data Model

### Category

| Field       | Type   | Required | Description           |
|-------------|--------|----------|-----------------------|
| id          | int    | Auto     | Unique identifier     |
| name        | string | Yes      | Category name         |
| description | string | No       | Category description  |

---

## Configuration

### Environment Variables & secrets.yml

The application supports two methods of configuration in this priority order:

1. **`secrets.yml`** (local file, for development) - Loaded first if it exists
2. **Environment Variables** - Used as fallback if not in `secrets.yml`
3. **Default Values** - Used if neither file nor environment variable is set

#### Supported Configuration Keys

| Key | YAML Key | Env Var | Default | Description |
|-----|----------|---------|---------|-------------|
| Database URL | `db_url` | `DB_URL` | `localhost` | PostgreSQL server hostname/IP |
| Database Port | `db_port` | `DB_PORT` | `5432` | PostgreSQL server port |
| Database Name | `db_name` | `DB_NAME` | `postgres` | Database name |
| Database User | `db_user` | `DB_USER` | `postgres` | Database user |
| Database Password | `db_password` | `DB_PASSWORD` | `postgres` | Database password |
| Server Port | `port` | `PORT` | `8080` | HTTP server port |

#### Example: Using Environment Variables

```bash
DB_URL=db.example.com \
DB_PORT=5432 \
DB_NAME=mydb \
DB_USER=admin \
DB_PASSWORD=secret123 \
PORT=8080 \
go run main.go
```

#### Example: Using secrets.yml

```yaml
db_url: db.example.com
db_port: 5432
db_name: mydb
db_user: admin
db_password: secret123
port: 8080
```

**Note:** `secrets.yml` is in `.gitignore` to prevent committing sensitive data. Always use `secrets.yml.example` as a template.

---

## Development

### Running Tests

Run unit tests for the database queries:

```bash
# First, ensure secrets.yml exists with database credentials
# (copy from secrets.yml.example if it doesn't exist)
cp secrets.yml.example secrets.yml
# Then edit secrets.yml with your actual database credentials

# Run tests for database package
go test ./database -v

# Run all tests
go test ./... -v

# Run tests with coverage
go test ./... -v -cover
```

**Important:** Tests require `secrets.yml` to exist in the project root directory with valid database credentials. If `secrets.yml` is missing, tests will fail with an error message.

**Note:** Tests use a separate `category_test` table that is automatically created and cleaned up.

### Project Structure

```
.
├── main.go                 # Application entry point
├── config/
│   └── config.go          # Configuration management (secrets.yml + env vars)
├── api/
│   └── categories.go      # HTTP route handlers for category endpoints
├── database/
│   ├── model.go           # Category data model
│   ├── migrations.go      # Database schema creation
│   ├── queries.go         # Category CRUD operations
│   └── database_test.go   # Database query unit tests
├── secrets.yml.example    # Configuration template
└── README.md              # This file
```

#### Folder Responsibilities

**`api/`** - HTTP layer
- Handles incoming HTTP requests
- Validates request data
- Calls database functions
- Returns HTTP responses

**`database/`** - Data layer
- Defines data models
- Manages database schema (migrations)
- Implements CRUD queries
- Includes unit tests

---

## Technical Details

- **Backend:** Go 1.22+
- **Database:** PostgreSQL
- **Storage:** Persistent PostgreSQL database
- **Migrations:** Automatic table creation with `CREATE IF NOT EXISTS`
- **Testing:** Go `testing` package with database test fixtures
- **Configuration:** YAML-based (`secrets.yml`) with environment variable fallback
- **Time Complexity:** O(1) for all database operations (indexed by primary key)
- **Deployment:** Railway (https://railway.app)

---

## License

This project is open source and available under the MIT License.

---

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## Contact

For questions or issues, please open an issue on GitHub.

---

**Happy Coding! 🚀**
