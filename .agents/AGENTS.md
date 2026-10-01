# AGENTS.md

## Aturan Utama (WAJIB BACA)
- Project ini adalah project latihan. **JANGAN memodifikasi code di project** (file di `services/`, `pkg/`, `proto/`, `sql/`, `postman/`, `go.mod`, `go.sum`, dll.) KECUALI user secara eksplisit mengatakan **"BUATKAN CODE"**.
- File yang BOLEH dimodifikasi secara bebas: file-file agent/config saja, seperti `AGENTS.md`, `opencode.json`, dan file sejenis.
- Jika user meminta penjelasan atau panduan, berikan instruksi langkah-demi-langkah dan jelaskan konsepnya — jangan langsung menulis/mengubah codenya.

## Build & Run
- Build everything: `go build ./...`
- Run a single service: `go run ./services/<name>/cmd` — **must be run from the repo root** because each service loads its `.env` via a relative path (e.g., `godotenv.Load("services/user/.env")`).
- There are no tests, no CI workflows, no linter config, and no Makefile/Taskfile in this repo. Do not invent commands like `make test` or `golangci-lint run`.

## Service Ports & Protocols (from each service's `.env`)
| Service | HTTP port | gRPC port | Notes |
|---------|-----------|-----------|-------|
| order   | 8080      | —         | Orchestrator: gRPC client to payment & product; RabbitMQ consumer |
| product | 8081      | 50052     | Dual HTTP + gRPC server |
| user    | 8082      | —         | HTTP only |
| payment | —         | 50051     | gRPC only; publishes to RabbitMQ; Midtrans gateway |

## Prerequisites & Startup Order
- **PostgreSQL** must be running at `localhost:5432` (DSN: `postgres://postgres:indocyber@localhost:5432/golang?sslmode=disable`). Apply schemas from `sql/*.sql`.
- **RabbitMQ** must be running at `amqp://guest:guest@localhost:5672/`. The order service declares exchange `payment.events` (topic) and queue `order.payment-status.queue` (routing key `payment.updated`); the payment service declares the same exchange.
- **Startup order matters:** start `product` and `payment` before `order` (order's `main` fatally exits if gRPC connections fail). RabbitMQ must be up before `order` and `payment`.
- The payment service requires `MIDTRANS_SERVER_KEY` / `MIDTRANS_CLIENT_KEY` / `MIDTRANS_PRODUCTION` in its `.env`.

## Configuration
- Each service has a gitignored `.env` and a committed `.env.example`. To set up a new service, copy `.env.example` → `.env` and fill in values.
- gRPC connections use **insecure** credentials (`insecure.NewCredentials()`) — no TLS.

## Architecture Conventions
- Each service follows clean architecture: `cmd/main.go` (wiring) → `internal/delivery/{http,grpc}` → `internal/usecase` → `internal/repository` → `internal/domain`.
- Shared helpers live in `pkg/utils` (DB connect, gRPC connect, RabbitMQ, logging, crypt, pagination, validator) and `pkg/request` / `pkg/response`.
- Proto definitions live in `proto/product` and `proto/payment`; the generated `*.pb.go` files are committed to the repo. After editing a `.proto`, regenerate with `protoc` (do not hand-edit `*.pb.go`).
- The order service's checkout flow: gRPC `GetProductDetail` → gRPC `DeductStock` → create order → background goroutine calls gRPC `ProcessPayment` → payment service publishes result to RabbitMQ → order service consumes and updates order status.
