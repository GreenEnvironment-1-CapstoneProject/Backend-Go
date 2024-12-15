# EcoMate 🌱

<p align="center">
  <img src="./assets/Logo Ecomate.png" width="350" alt="EcoMate Logo" />
</p>

<div align="center">
  
  [![API Docs](https://img.shields.io/badge/API_Docs-Open-green.svg)](https://greenenvironment.my.id/swagger/index.html#/)
  [![Made with Go](https://img.shields.io/badge/Go-00ADD8?style=flat&logo=go&logoColor=white)](https://go.dev/)
  [![Echo Framework](https://img.shields.io/badge/Echo-000000?style=flat&logo=go&logoColor=white)](https://echo.labstack.com/)
  [![MySQL](https://img.shields.io/badge/MySQL-4479A1?style=flat&logo=mysql&logoColor=white)](https://www.mysql.com/)
  
</div>

## Table of Contents

- [About](#-about-ecomate)
- [Features](#-key-features)
- [Architecture](#-architecture)
- [Project Structure](#-project-structure)
- [Technology Stack](#️-technology-stack)
- [Getting Started](#-getting-started)
- [API Documentation](#-api-documentation)
- [Contributing](#-contributing)
- [License](#-license)

## 📖 About EcoMate

EcoMate is a comprehensive environmental platform designed to empower users in their journey towards sustainable living. Our mission is to make eco-friendly choices accessible and rewarding for everyone.

## 🌟 Key Features

### Core Features

- **Eco-friendly Marketplace**: Browse and purchase sustainable products
- **Community Challenges**: Participate in eco-challenges
- **Discussion Forums**: Connect with like-minded individuals
- **Impact Tracking**: Monitor environmental contributions
- **Leaderboard System**: Compete and earn recognition

### Technical Features

- **AI-Powered Chatbot**: Environmental information assistant
- **Secure Authentication**: JWT and Google OAuth integration
- **Cloud Storage**: Efficient media handling with Cloudinary
- **Payment Processing**: Secure transactions via Midtrans
- **Email Notifications**: Automated communications via SMTP

## 📊 Architecture

### Clean Architecture Implementation

The project follows clean architecture principles with clear separation of concerns:

1. **Entity Layer** (`entity.go`)

   - Core business rules and entities
   - Independent of external frameworks

2. **Use Case Layer** (`service/`)

   - Application-specific business rules
   - Orchestrates data flow between entities

3. **Interface Adapters** (`controller/`, `repository/`)

   - Controllers handle HTTP requests
   - Repositories manage data persistence

4. **External Interfaces** (`utils/`, `helper/`)
   - Frameworks and tools
   - External service integrations

### Visual Architecture

<div align="center">
  <img src="./assets/HLA Capastone Project.png" alt="High-Level Architecture" width="800"/>
</div>

### Database Schema

<div align="center">
  <img src="./assets/Capstone-_Kelompok1-ERD.drawio (2).png" alt="Database Schema" width="800"/>
</div>

## 📂 Project Structure

```
├── assets/                    # Static assets and images
├── configs/                   # Application configuration
├── constant/                  # Constants and enums
│   ├── route/                # Route constants
│   └── ...                   # Other constants
├── docs/                      # API documentation
├── features/                  # Core business features
│   ├── users/                # User management
│   ├── products/             # Product management
│   ├── challenges/           # Challenge system
│   ├── forum/                # Discussion forums
│   ├── transactions/         # Payment processing
│   └── ...                   # Other features
├── helper/                    # Utility helpers
├── routes/                    # Route definitions
├── utils/                     # Shared utilities
│   ├── databases/            # Database utilities
│   ├── google/               # Google OAuth
│   ├── midtrans/             # Payment gateway
│   ├── openai/               # AI integration
│   └── storages/             # File storage
└── main.go                    # Application entry point
```

Each feature follows a consistent structure:

```
feature/
├── controller/               # HTTP handlers
│   ├── handler.go           # Request handling
│   ├── request.go           # Request DTOs
│   └── response.go          # Response DTOs
├── repository/              # Data access
│   ├── model.go             # Database models
│   └── query.go             # Database queries
├── service/                 # Business logic
│   ├── service.go           # Core logic
│   └── service_test.go      # Unit tests
└── entity.go                # Domain entities
```

## 🛠️ Technology Stack

### Core Technologies

- **Backend**: Go 1.23+
- **Framework**: Echo
- **Database**: MySQL 7.0+ with GORM
- **Documentation**: Swagger

### External Services

- **Authentication**: JWT, Google OAuth
- **Storage**: Cloudinary
- **Payments**: Midtrans
- **AI**: OpenAI
- **Email**: SMTP

## 🚀 Getting Started

### Prerequisites

- Go 1.23+
- MySQL 7.0+
- Git

### Installation

1. **Clone Repository**

```bash
git clone https://github.com/GreenEnvironment-1-CapstoneProject/Backend-Go.git
cd backend-capstone
```

2. **Install Dependencies**

```bash
go mod tidy
```

3. **Configure Environment**
   Create `.env` file:

```env
# Application
APP_PORT=your_app_port

# Database
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=your_db_name

# Authentication
JWT_SECRET=your_jwt_secret
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret

# External Services
OPENAI_API_KEY=your_openai_key
MIDTRANS_CLIENT_KEY=your_midtrans_client_key
MIDTRANS_SERVER_KEY=your_midtrans_server_key

# Email
SMTP_USER=your_smtp_user
SMTP_PASS=your_smtp_password
SMTP_HOST=your_smtp_host
SMTP_PORT=your_smtp_port
```

4. **Run Application**

```bash
go run main.go
```

## 👥 Contributors

<table>
  <tr>
    <td align="center">
      <a href="https://github.com/reinhardprs">
        <img src="https://github.com/reinhardprs.png" width="100px;" alt="Reinhard Prasetya"/><br />
        <sub><b>Reinhard Prasetya</b></sub>
      </a><br />
      <a href="https://www.linkedin.com/in/reinhardprasetya/">
        <img src="https://img.shields.io/badge/LinkedIn-blue?style=flat&logo=linkedin" alt="LinkedIn"/>
      </a>
    </td>
    <td align="center">
      <a href="https://github.com/MHafidafandi">
        <img src="https://github.com/MHafidafandi.png" width="100px;" alt="Muhammad Hafid Afandi"/><br />
        <sub><b>Muhammad Hafid Afandi</b></sub>
      </a><br />
      <a href="https://www.linkedin.com/in/m-hafid-afandi-23b725245/">
        <img src="https://img.shields.io/badge/LinkedIn-blue?style=flat&logo=linkedin" alt="LinkedIn"/>
      </a>
    </td>
  </tr>
</table>

## 📚 Documentation

- [API Documentation](https://greenenvironment.my.id/swagger/index.html#/)
- [GitHub Repository](https://github.com/GreenEnvironment-1-CapstoneProject/Backend-Go.git)

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
