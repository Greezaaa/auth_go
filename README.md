[![Stand With Ukraine](https://raw.githubusercontent.com/vshymanskyy/StandWithUkraine/main/banner-direct-single.svg)](https://stand-with-ukraine.pp.ua)

# 🚀 Auth Service (Go Backend)

A modular, high-performance authentication service built with Go and PostgreSQL.

## 🏗 Project Architecture

Moved away from a "God-file" approach to a Layered Design to ensure maintainability:

*   **Handler**: Manages HTTP logic, JSON decoding/encoding, and routing using `chi`.
*   **Service**: The "Brain" of the app. Handles business logic (bcrypt hashing, secure verification code generation).
*   **Repository**: Data persistence layer using `pgxpool` for high-concurrency connection management.
*   **Migrations**: Database versioning via `goose` and Go's `embed` package to keep the schema "baked" into the binary.

## 🛠 Key Features & Decisions

*   **UUID v4**: Used non-sequential unique identifiers (via `google/uuid` and Postgres `uuid-ossp`) instead of typical integers for security and distributed safety.
*   **Connection Pooling**: Implemented `pgxpool` to allow hundreds of concurrent users to connect simultaneously without bottlenecks.
*   **Security**: One-way password hashing using `bcrypt` and cryptographically secure 6-digit verification codes.
*   **Environment Isolation**: Used a `.env` file and Docker Compose to manage secrets and infrastructure independently.

## 🚦 Getting Started

### 1. Environment Setup

Ensure your `.env` file matches the following structure:

```env
# Networking & Server
PORT=8080
DOMAIN=localhost
API_PREFIX=api
API_VERSION=v1

# Database Credentials
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=auth_db
DB_PORT=5432

# Docker Specific
DB_HOST=db
APP_CONTAINER_NAME=auth_app
DB_CONTAINER_NAME=auth_db
COMPOSE_PROJECT_NAME=auth_service
```

### 2. Launch with Docker

Rebuild the Go binary and start the infrastructure containers:

```bash
docker-compose up --build
```

### 3. Database Access (TablePlus)

To retrieve the `email_code` for manual verification testing, use the following connection details:

*   **Host**: `localhost`
*   **Port**: `5433` (as mapped in `docker-compose.yml`)
*   **User/Pass**: From your `.env` file

## 📡 API Endpoints

| Method | Endpoint | Description | Request Body / Params |
| :--- | :--- | :--- | :--- |
| **POST** | `/api/v1/users/register` | Register a new user | **JSON**<br>`email` (string, required)<br>`password` (string, required)<br>`name` (string, required)<br>`allowed_apps` (array of strings) |
| **GET** | `/api/v1/users` | List all registered users | None |
| **GET** | `/api/v1/users/{id}` | Retrieve a specific user | `id` (UUID) in URL path |
| **PATCH** | `/api/v1/users/verify` | Confirm email | **JSON**<br>`email` (string, required)<br>`code` (string, required) |

## 📝 Developer Notes

> [!NOTE]
> **Strict Paths**: We use the versioned path `/api/v1/users` to allow for future API iterations without breaking existing clients.

> [!TIP]
> **Migrations**: If you need to change the table structure, create a new `.sql` file in the `/migrations` folder instead of editing existing ones.

## 🗺 Roadmap

### 🛡️ Security & Anti-Fraud
- [ ] **IP-Based Rate Limiting**: Implement a 5-attempt-per-minute limit per IP on `/login` and `/register`.
- [ ] **Automated IP Ban**: Temporary blacklisting of IPs that exceed brute-force thresholds.
- [ ] **Account Lockout**: Secure user accounts after multiple failed password attempts.

### 📊 Metadata & Compliance
- [ ] **Fingerprinting**: Track User OS, Browser, and Language for every login attempt.
- [ ] **GDPR Compliance**: Implementation of IP anonymization (xxx masking) for stored logs.
- [ ] **Audit Logs**: Record of sensitive actions (password changes, email updates) with masked IP data.

### 📧 Communication & UX
- [ ] **Automated Verification**: Immediate 6-digit code delivery via SMTP (SendGrid/Mailtrap).
- [ ] **Welcome Sequence**: Trigger a "Welcome" email automatically after a user confirms their email.
- [ ] **Password Recovery**: Secure "Forgot Password" flow with time-limited tokens.

### 🔑 Identity & Access
- [ ] **JWT & Refresh Tokens**: Full session management with short-lived access tokens and secure rotation.
- [ ] **RBAC (Role-Based Access Control)**: Permission layers for Owner, Admin, and SuperAdmin.
- [ ] **Profile Management**: Allow data owners to update their own info while restricting admin-only fields.

### 📈 Operations
- [ ] **Structured Logging**: Integration of `slog` for JSON-based observability.
- [ ] **Health Monitoring**: Real-time tracking of DB connection pool health and latency.

## 🤝 Contributing

Contributions are what make the open-source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

### THANK YOU!
