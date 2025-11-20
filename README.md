# WalletServer — Production-Grade Digital Wallet API in Pure Go

![Go](https://img.shields.io/badge/go-1.22%2B-00ADD8?style=for-the-badge&logo=go)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql)
![net/http](https://img.shields.io/badge/net%2Fhttp-only-success?style=for-the-badge)
![License](https://img.shields.io/badge/license-MIT-blue?style=for-the-badge)

A fully functional, secure, and atomic digital wallet backend built **entirely in pure Go** — **zero web frameworks** (no Gin, no Echo, no Fiber).

This is how real payment systems work under the hood.

### Live Postman Demo
![Postman Collection Demo](https://raw.githubusercontent.com/niloy104/WalletServer/main/assets/postman-demo.gif)
> *(Full collection link below)*

### Core Features

- User registration & login with JWT authentication
- Atomic P2P money transfers (with row-level locking)
- Top-up wallet balance
- Real-time balance inquiry
- Complete transaction history per user
- Full input validation & error handling
- Rate limiting, CORS, structured logging
- Database migrations
- Clean, scalable, layered architecture


### Tech Stack

| Layer             | Technology                          |
|-------------------|-------------------------------------|
| Language          | Go (pure net/http)                  |
| Database          | PostgreSQL + sqlx                   |
| Auth              | JWT (self-implemented)              |
| Architecture      | Domain-Driven Design (DDD) inspired |
| Concurrency       | Database transactions + `FOR UPDATE`|
| No dependencies   | Only standard library + sqlx        |

### REST API Endpoints

| Method | Endpoint                  | Description                      | Auth  |
|--------|---------------------------|----------------------------------|-------|
| POST   | `/users`                  | Create account                   | Public|
| POST   | `/users/login`            | Login → JWT                      | Public|
| GET    | `/wallets/my-balance`     | Get current balance              | JWT   |
| POST   | `/wallets/topup`          | Add money                        | JWT   |
| POST   | `/wallets/transfer`       | Send money to another user       | JWT   |
| GET    | `/transactions`           | Full transaction history         | JWT   |

### Key Engineering Concepts Implemented

- Atomic transactions using `WithTx` + `FOR UPDATE` row locking
- Preventing race conditions in financial operations
- Clean separation of concerns (Domain → Service → Repository → Handler)
- Dependency inversion (interfaces in domain layer)
- Secure JWT implementation from scratch
- Structured error handling with custom domain errors
- Middleware chain: Logger → CORS → Rate Limit → Auth
- Production-ready response utilities

### Project Structure (Clean & Scalable)

```text
wallet/
├── cmd/ 
├── domain/           # Pure entities & business rules
├── migrations/
├── wallet/           # Application use cases 
├── user/ 
├── repo/             # Database implementations
├── rest/handlers/walletB/  # HTTP delivery layer
├── rest/middlewares/ # Auth, CORS, Logger, RateLimiter
├── infra/db/         # Connection + migrations
└── util/             # Response helpers, JWT, etc.
├── main.go           # main file

### Run Locally
Bash: 
git clone https://github.com/niloy104/VaultGo.git
cd wallet
cp .env.example .env
go run main.go

Server: http://localhost:4000


### Postman Collection
Download Full Postman Collection
View in Postman (link after upload)
Future Roadmap


### Future Roadmap
Admin dashboard
Transaction search & filters
Docker + PostgreSQL compose
Unit & integration tests
WebSocket real-time balance updates

License
MIT © niloy104
