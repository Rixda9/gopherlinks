# 🔗 GopherLinks

A high-performance URL shortening service built with Go, featuring Redis caching and PostgreSQL storage.

![Demo](/docs/demo.gif)

**Try it live:** [https://gopherlinks.fly.dev](https://gopherlinks.fly.dev)

---

## ✨ Features

- ⚡ **Lightning Fast** - Redis caching for sub-millisecond redirects
- 🔒 **Secure** - Input validation and SQL injection prevention
- 🎯 **Smart Deduplication** - Same URL always gets the same short code
- 🐳 **Docker Ready** - One-command deployment
- 📚 **API Documentation** - Interactive Swagger docs

---

## 🚀 Quick Start

### Try the Web Interface

Visit [gopherlinks.fly.dev](https://gopherlinks.fly.dev) and paste your URL!

### Try the API
```bash
# Shorten a URL
curl -X POST https://gopherlinks.fly.dev/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"url": "https://github.com/yourusername"}'

# Response
{
  "short_url": "https://gopherlinks.fly.dev/aB3xK9"
}

# Use the short URL - automatically redirects!
```

---

## 🛠️ Tech Stack

- **Language:** Go 1.24
- **Database:** PostgreSQL 15
- **Cache:** Redis 7
- **Web Framework:** Chi Router
- **Deployment:** Fly.io + Neon + Upstash

---

## 📁 Project Structure
```
gopherlinks/
├── cmd/api/              # Application entry point
├── internal/
│   ├── database/         # Database connection & migrations
│   ├── models/           # Data structures
│   ├── repository/       # Data access layer (Postgres & Redis)
│   └── server/           # HTTP handlers & routing
├── web/
│   └── index.html        # Frontend UI
├── migrations/           # SQL schema migrations
├── Dockerfile            # Container configuration
└── docker-compose.yml    # Local development setup
```

---

## 🏗️ Architecture

### Caching Strategy
1. User requests short URL
2. Check Redis cache (sub-millisecond)
3. If miss, query PostgreSQL
4. Cache result for 24 hours
5. Redirect user

### Database Schema
```sql
CREATE TABLE urls (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    short_code VARCHAR(10) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

### Deduplication
- Same URL always returns the same short code
- Prevents database bloat
- Maintains consistent user experience

---

## 💻 Local Development

### Prerequisites
- Go 1.24+
- Docker & Docker Compose

### Setup
```bash
# Clone the repository
git clone https://github.com/Rixda9/gopherlinks.git
cd gopherlinks

# Start services (Postgres + Redis + App)
docker-compose up --build

# Visit http://localhost:8080
```

### Environment Variables
```env
DATABASE_URL=postgres://user:pass@localhost:5432/urlshortener
REDIS_ADDR=localhost:6379
PORT=8080
BASE_URL=http://localhost:8080
```

---

## 📡 API Endpoints

### Shorten URL
```http
POST /api/shorten
Content-Type: application/json

{
  "url": "https://example.com/very/long/url"
}
```

**Response:**
```json
{
  "short_url": "https://gopherlinks.fly.dev/aB3xK9"
}
```

### Redirect
```http
GET /{shortCode}
```
Redirects to the original URL.

### Health Check
```http
GET /health
```

---

## 🚀 Deployment

Deployed on **Fly.io** with:
- **Neon** - Serverless Postgres (free tier)
- **Upstash** - Serverless Redis (free tier)
- **Fly.io** - Application hosting (free tier)

**Total cost: $0/month** 💰

### Deploy Your Own
```bash
# Install Fly CLI
curl -L https://fly.io/install.sh | sh

# Login
flyctl auth login

# Deploy
flyctl launch
flyctl deploy
```

---

## 🎯 Key Features Explained

### Smart URL Deduplication
If you shorten the same URL twice, you get the same short code back. This:
- Saves database space
- Provides consistent links
- Improves cache hit rates

### Redis Caching
- Frequently accessed URLs cached for 24 hours
- Reduces database load by ~95%
- Average redirect time: < 10ms

---

## 🧪 Performance

- **Redirect Speed:** < 10ms (cached), < 50ms (uncached)
- **Concurrent Requests:** 1000+ req/s
- **Uptime:** 99.9% (Fly.io)

---

## 🤝 Contributing

Pull requests welcome! For major changes, please open an issue first.

---

## 📝 License

MIT License - see [LICENSE](LICENSE) file for details

---

## 👤 Author

**Rixda9**
- GitHub: [@Rixda9](https://github.com/Rixda9)
- Live Demo: [gopherlinks.fly.dev](https://gopherlinks.fly.dev)

---

Built using Go and deployed on Fly.io
