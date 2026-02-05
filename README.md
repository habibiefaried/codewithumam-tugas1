# Category & Product CRUD API

A REST API for managing categories and products with PostgreSQL database backend built with Go.

![CI Tests](https://github.com/habibiefaried/codewithumam-tugas1/workflows/CI%20Tests/badge.svg)

## Features

- ✅ Full CRUD operations for categories and products
- ✅ Product-Category relationship with foreign keys
- ✅ PostgreSQL database backend
- ✅ Configuration from `secrets.yml` or environment variables
- ✅ Database migrations with `CREATE IF NOT EXISTS`
- ✅ Indexes on id and name columns for performance
- ✅ Comprehensive unit tests for database queries
- ✅ RESTful API design with proper separation of concerns
- ✅ JSON request/response
- ✅ Checkout endpoint with transactional stock updates
- ✅ Automated CI/CD tests (30+ test scenarios)
- ✅ Deployed on Railway

## Live API

🚀 **Production URL:** https://codewithumam-tugas-production.up.railway.app/

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

The server will automatically create the `category` and `product` tables if they don't exist.

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
curl https://codewithumam-tugas-production.up.railway.app/health

# Local
curl http://localhost:8080/health
```

**Response:**
```
OK
```

---

## Category Endpoints

### 1. Get All Categories

**Endpoint:** `GET /categories`

**Request:**
```bash
# Production
curl https://codewithumam-tugas-production.up.railway.app/categories

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
curl https://codewithumam-tugas-production.up.railway.app/categories/1

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
curl -X POST https://codewithumam-tugas-production.up.railway.app/categories \
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
curl -X PUT https://codewithumam-tugas-production.up.railway.app/categories/1 \
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
curl -X DELETE https://codewithumam-tugas-production.up.railway.app/categories/1

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

## Product Endpoints

### Products: Get All

**Endpoint:** `GET /products`

**Request:**
```bash
# Production
curl https://codewithumam-tugas-production.up.railway.app/products

# Local
curl http://localhost:8080/products
```

**Response:**
```json
[
  {
    "id": 1,
    "name": "Laptop",
    "price": 1299,
    "stock": 8,
    "category_id": 1,
    "category_name": "Electronics",
    "category_description": "Electronic devices and gadgets"
  },
  {
    "id": 2,
    "name": "Phone",
    "price": 299,
    "stock": 10,
    "category_id": 1,
    "category_name": "Electronics",
    "category_description": "Electronic devices and gadgets"
  }
]
```

---

### Products: Get by ID

**Endpoint:** `GET /products/{id}`

**Request:**
```bash
# Production
curl https://codewithumam-tugas-production.up.railway.app/products/1

# Local
curl http://localhost:8080/products/1
```

**Response (Success - 200):**
```json
{
  "id": 1,
  "name": "Laptop",
  "price": 1299,
  "stock": 8,
  "category_id": 1,
  "category_name": "Electronics",
  "category_description": "Electronic devices and gadgets"
}
```

**Response (Not Found - 404):**
```
Product not found
```

---

### Products: Create

**Endpoint:** `POST /products`

**Request:**
```bash
# Production
curl -X POST https://codewithumam-tugas-production.up.railway.app/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "price": 1299,
    "stock": 8,
    "category_id": 1
  }'

# Local
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop",
    "price": 1299,
    "stock": 8,
    "category_id": 1
  }'
```

**Response (Success - 201):**
```json
{
  "id": 1,
  "name": "Laptop",
  "price": 1299,
  "stock": 8,
  "category_id": 1,
  "category_name": "Electronics",
  "category_description": "Electronic devices and gadgets"
}
```

**Response (Bad Request - 400):**
```
Name is required
```

---

### Products: Update

**Endpoint:** `PUT /products/{id}`

**Request:**
```bash
# Production
curl -X PUT https://codewithumam-tugas-production.up.railway.app/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Pro",
    "price": 1599,
    "stock": 5,
    "category_id": 1
  }'

# Local
curl -X PUT http://localhost:8080/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Laptop Pro",
    "price": 1599,
    "stock": 5,
    "category_id": 1
  }'
```

**Response (Success - 200):**
```json
{
  "id": 1,
  "name": "Laptop Pro",
  "price": 1599,
  "stock": 5,
  "category_id": 1,
  "category_name": "Electronics",
  "category_description": "Electronic devices and gadgets"
}
```

**Response (Not Found - 404):**
```
Product not found
```

---

### Products: Delete

**Endpoint:** `DELETE /products/{id}`

**Request:**
```bash
# Production
curl -X DELETE https://codewithumam-tugas-production.up.railway.app/products/1

# Local
curl -X DELETE http://localhost:8080/products/1
```

**Response (Success - 204):**
```
(No content)
```

**Response (Not Found - 404):**
```
Product not found
```

---

## Checkout Endpoint

### Checkout: Create Transaction

**Endpoint:** `POST /checkout`

**Request:**
```bash
# Production
curl -X POST https://codewithumam-tugas-production.up.railway.app/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"product_id": 1, "quantity": 2},
      {"product_id": 2, "quantity": 1}
    ]
  }'

# Local
curl -X POST http://localhost:8080/checkout \
  -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"product_id": 1, "quantity": 2},
      {"product_id": 2, "quantity": 1}
    ]
  }'
```

**Response (Success - 201):**
```json
{
  "id": 1,
  "total_amount": 2897,
  "created_at": "2026-02-05T08:15:30Z",
  "details": [
    {
      "id": 1,
      "transaction_id": 1,
      "product_id": 1,
      "product_name": "Laptop",
      "product_description": "Electronic devices and gadgets",
      "unit_price": 1299,
      "quantity": 2,
      "subtotal": 2598
    },
    {
      "id": 2,
      "transaction_id": 1,
      "product_id": 2,
      "product_name": "Phone",
      "product_description": "Electronic devices and gadgets",
      "unit_price": 299,
      "quantity": 1,
      "subtotal": 299
    }
  ]
}
```

**Response (Insufficient Stock - 400):**
```
Insufficient stock
```

**Response (Product Not Found - 404):**
```
Product not found
```

---

## Report Endpoints

### Report: Hari Ini

**Endpoint:** `GET /report/hari-ini`

**Request:**
```bash
# Production
curl https://codewithumam-tugas-production.up.railway.app/report/hari-ini

# Local
curl http://localhost:8080/report/hari-ini
```

**Response (Success - 200):**
```json
{
  "total_revenue": 45000,
  "total_transaksi": 5,
  "produk_terlaris": { "nama": "Indomie Goreng", "qty_terjual": 12 }
}
```

---

### Report: Range

**Endpoint:** `GET /report?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD`

**Request:**
```bash
# Production
curl "https://codewithumam-tugas-production.up.railway.app/report?start_date=2026-01-01&end_date=2026-02-01"

# Local
curl "http://localhost:8080/report?start_date=2026-01-01&end_date=2026-02-01"
```

**Response (Success - 200):**
```json
{
  "total_revenue": 45000,
  "total_transaksi": 5,
  "produk_terlaris": { "nama": "Indomie Goreng", "qty_terjual": 12 }
}
```

**Response (Bad Request - 400):**
```
start_date and end_date are required
```

**Response (Bad Request - 400):**
```
Invalid start_date
```

**Response (Bad Request - 400):**
```
Invalid end_date
```

**Response (Bad Request - 400):**
```
end_date must be on or after start_date
```

---

## Quick Testing Examples

### Complete Category Workflow (Production)

```bash
# 1. Create categories
curl -X POST https://codewithumam-tugas-production.up.railway.app/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Electronics","description":"Electronic devices"}'

curl -X POST https://codewithumam-tugas-production.up.railway.app/categories \
  -H "Content-Type: application/json" \
  -d '{"name":"Books","description":"Physical and digital books"}'

# 2. List all categories
curl https://codewithumam-tugas-production.up.railway.app/categories

# 3. Get specific category
curl https://codewithumam-tugas-production.up.railway.app/categories/1

# 4. Update category
curl -X PUT https://codewithumam-tugas-production.up.railway.app/categories/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Electronics and Tech","description":"Updated description"}'

# 5. Delete category
curl -X DELETE https://codewithumam-tugas-production.up.railway.app/categories/2
```

### Complete Product Workflow (Production)

```bash
# 1. Create products (category_id must exist)
curl -X POST https://codewithumam-tugas-production.up.railway.app/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop","price":1299,"stock":8,"category_id":1}'

curl -X POST https://codewithumam-tugas-production.up.railway.app/products \
  -H "Content-Type: application/json" \
  -d '{"name":"Phone","price":299,"stock":10,"category_id":1}'

# 2. List all products (includes category information)
curl https://codewithumam-tugas-production.up.railway.app/products

# 3. Get specific product
curl https://codewithumam-tugas-production.up.railway.app/products/1

# 4. Update product
curl -X PUT https://codewithumam-tugas-production.up.railway.app/products/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"Laptop Pro","price":1599,"stock":5,"category_id":1}'

# 5. Delete product
curl -X DELETE https://codewithumam-tugas-production.up.railway.app/products/2
```

### Using jq for Pretty Output

If you have `jq` installed, you can format the JSON output:

```bash
curl https://codewithumam-tugas-production.up.railway.app/categories | jq
```

---

## Error Handling

The API returns appropriate HTTP status codes:

- `200 OK` - Request successful
- `201 Created` - Resource created successfully
- `204 No Content` - Resource deleted successfully
- `400 Bad Request` - Invalid request (missing required fields, invalid ID format)
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource conflict (e.g., deleting category with existing products)
- `500 Internal Server Error` - Server error

---

## Edge Cases & Constraints

### Foreign Key Constraint Protection
**Issue:** Attempting to delete a category that has products referencing it.

**Example:**
```bash
# Category 1 has products
curl -X DELETE http://localhost:8080/categories/1
# Response: 409 Conflict
# Body: Cannot delete category that has products
```

**Why:** PostgreSQL enforces referential integrity. Products have a foreign key constraint pointing to categories. You must delete all products referencing a category before deleting the category itself.

**Solution:**
1. Delete all products with that category_id first
2. Then delete the category

```bash
# 1. Find products for category 1
curl http://localhost:8080/products | jq '.[] | select(.category_id == 1) | .id'

# 2. Delete each product
curl -X DELETE http://localhost:8080/products/1
curl -X DELETE http://localhost:8080/products/2

# 3. Now delete the category
curl -X DELETE http://localhost:8080/categories/1  # 204 No Content
```

### Input Validation
- **Empty name:** Returns 400 Bad Request
- **Invalid ID format:** Returns 400 Bad Request (e.g., `/categories/abc`)
- **Non-existent ID:** Returns 404 Not Found
- **Missing required fields:** Returns 400 Bad Request

### JSON Encoding
Special characters in names/descriptions are JSON-encoded. For example:
- `&` becomes `\u0026`
- `<` becomes `\u003c`
- `>` becomes `\u003e`

This is normal JSON encoding and doesn't affect functionality.

---

## Data Model

### Category

| Field       | Type   | Required | Description           |
|-------------|--------|----------|-----------------------|
| id          | int    | Auto     | Unique identifier     |
| name        | string | Yes      | Category name         |
| description | string | No       | Category description  |

### Product

| Field                   | Type   | Required | Description                      |
|-------------------------|--------|----------|-------------------------------------|
| id                      | int    | Auto     | Unique identifier                  |
| name                    | string | Yes      | Product name                       |
| price                   | int    | Yes      | Product price (in cents)           |
| stock                   | int    | Yes      | Available stock quantity           |
| category_id             | int    | Yes      | Foreign key to category table      |
| category_name           | string | Read     | Category name (from join)          |
| category_description    | string | Read     | Category description (from join)   |

### Transaction

| Field        | Type      | Required | Description                     |
|--------------|-----------|----------|---------------------------------|
| id           | int       | Auto     | Unique identifier               |
| total_amount | int       | Auto     | Total transaction amount        |
| created_at   | timestamp | Auto     | Checkout timestamp (UTC)        |
| details      | array     | Read     | List of transaction details     |

### TransactionDetail

| Field           | Type   | Required | Description                         |
|-----------------|--------|----------|-------------------------------------|
| id              | int    | Auto     | Unique identifier                   |
| transaction_id  | int    | Yes      | Foreign key to transaction table    |
| product_id      | int    | Yes      | Product ID snapshot (no FK)         |
| product_name    | string | Yes      | Product name snapshot               |
| product_description | string | Yes   | Product description snapshot        |
| unit_price      | int    | Yes      | Unit price at purchase time         |
| quantity        | int    | Yes      | Quantity purchased                  |
| subtotal        | int    | Yes      | price × quantity                    |

### CheckoutRequest

| Field | Type  | Required | Description                   |
|-------|-------|----------|-------------------------------|
| items | array | Yes      | List of items to purchase     |

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

**Note:** Tests use separate `category_test` and `product_test` tables that are automatically created and cleaned up.

### Reset Database (Drop & Recreate)

Run the reset tool to drop all tables and recreate them using migrations:

```bash
go run ./cmd/reset-db
```

### Project Structure

```
.
├── main.go                 # Application entry point
├── config/
│   └── config.go          # Configuration management (secrets.yml + env vars)
├── api/
│   ├── categories.go      # HTTP route handlers for category endpoints
│   └── products.go        # HTTP route handlers for product endpoints
├── database/
│   ├── categories.go      # Category data model
│   ├── product.go         # Product data model
│   ├── migrations.go      # Database schema creation (tables + indexes)
│   ├── queries.go         # CRUD operations for categories and products
│   ├── categories_test.go # Category unit tests
│   └── product_test.go    # Product unit tests
├── .github/workflows/
│   └── ci.yml             # GitHub Actions CI/CD pipeline (30+ tests)
├── secrets.yml.example    # Configuration template
└── README.md              # This file
```

#### Folder Responsibilities

**`api/`** - HTTP layer
- Handles incoming HTTP requests for /categories and /products
- Validates request data (required fields, ID format)
- Calls database functions
- Returns JSON responses with appropriate HTTP status codes

**`database/`** - Data layer
- Defines data models (Category, Product)
- Manages database schema (CREATE IF NOT EXISTS, indexes)
- Implements CRUD queries for both entities
- Includes comprehensive unit tests with test table isolation

**`config/`** - Configuration layer
- Loads configuration from `secrets.yml` (priority)
- Falls back to environment variables
- Falls back to hardcoded defaults
- Supports multi-location path resolution (current directory and executable directory)

---

## Technical Details

- **Backend:** Go 1.22+
- **Database:** PostgreSQL 12+
- **Storage:** Persistent PostgreSQL database
- **Migrations:** Automatic table creation with `CREATE IF NOT EXISTS` and indexes
- **Indexes:** Optimized indexes on `id` and `name` columns for fast lookups
- **Relationships:** Foreign key constraints between products and categories
- **Testing:** Go `testing` package with database test fixtures and isolated test tables
- **Configuration:** YAML-based (`secrets.yml`) with environment variable fallback
- **CI/CD:** GitHub Actions with PostgreSQL service and 30+ automated tests
- **Time Complexity:** O(1) for indexed queries, O(n) for full table scans
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
