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

## 📖 About EcoMate

EcoMate is a comprehensive environmental platform designed to empower users in their journey towards sustainable living. Our mission is to make eco-friendly choices accessible and rewarding for everyone.

### 🌟 Key Features

- **Eco-friendly Marketplace**: Browse and purchase sustainable products
- **Community Challenges**: Participate in eco-challenges
- **Discussion Forums**: Connect with like-minded individuals
- **Reward System**: Earn points for sustainable actions
- **Chatbot**: ask AI about the green environment

## 🛠️ Technology Stack

- **Backend**: Go (Golang)
- **Framework**: Echo
- **Database**: MySQL with GORM
- **Authentication**: JWT
- **Cloud Storage**: Cloudinary
- **Payment Gateway**: Midtrans
- **ChatBot**: OpenAI
- **Email Service**: SMTP

## 🚀 Getting Started

### Prerequisites

- Go 1.23 or higher
- MySQL 7.0 or higher
- Git

### Installation

1. Clone the repository

```bash
git clone https://github.com/GreenEnvironment-1-CapstoneProject/Backend-Go.git
cd backend-capstone
```

2. Install dependencies

```bash
go mod tidy
```

3. Configure environment variables

```env
# APP
APP_PORT=your_app_port

# Database Configuration
DB_HOST=your_db_host
DB_PORT=your_db_port
DB_USER=your_db_user
DB_PASS=your_db_password
DB_NAME=your_db_name

# Email Configuration
SMTP_USER=your_smtp_user
SMTP_PASS=your_smtp_password
SMTP_HOST=your_smtp_host
SMTP_PORT=your_smtp_port

# Authentication
JWT_SECRET=your_jwt_secret

# Google Cloud Storage
PROJECT_ID=your_project_id
BUCKET_NAME=your_bucket_name
GOOGLE_APPLICATION_CREDENTIALS=path_to_credentials

# Payment Gateway
MIDTRANS_CLIENT_KEY=your_client_key
MIDTRANS_SERVER_KEY=your_server_key

# Open AI
OPENAI_API_KEY=


# Google
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
```

4. Run the application

```bash
go run main.go
```

## 📊 System Architecture

### Database Schema

<div align="center">
  <img src="./assets/Capstone-_Kelompok1-ERD.drawio (2).png" alt="Database Schema" width="800"/>
</div>

### High-Level Architecture

<div align="center">
  <img src="./assets/HLA Capastone Project.png" alt="High-Level Architecture" width="800"/>
</div>

## 🔐 API Features

### User Management

- Registration and Authentication
- Profile Management
- Password Reset
- Role-based Access Control

### Product Management

- Product Catalog
- Shopping Cart
- Order Management
- Payment Processing

### Community Features

- Challenge Participation
- Leaderboards
- Discussion Forums
- Impacts

### Administrative Tools

- User Management
- Product Management
- Challenge Administration
- Analytics Dashboard

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
