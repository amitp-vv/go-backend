# Go Backend Project

## Overview
This project is a backend application built using Go, designed to handle various operations related to real estate and blockchain. It includes features such as user authentication, property management, and event handling.

## Project Structure
```
go-backend
├── cmd
│   └── main.go                # Application entry point
├── config
│   └── config.go              # Configuration settings
├── docs
│   └── Backend-API.md         # API documentation
├── internal
│   ├── listeners
│   │   └── events.go          # Event listeners
│   ├── middleware
│   │   ├── auth.go            # Authentication middleware
│   │   └── rate_limit.go       # Rate limiting middleware
│   ├── models
│   │   ├── claim.go           # Claim model
│   │   ├── distribution.go     # Distribution model
│   │   ├── property.go        # Property model
│   │   ├── user.go            # User model
│   │   ├── wallet.go          # Wallet model
│   │   └── index.go           # Model aggregator
│   ├── routes
│   │   ├── admin.go           # Admin routes
│   │   ├── auth.go            # Authentication routes
│   │   ├── chain.go           # Blockchain routes
│   │   └── read.go            # Read operations routes
│   └── services
│       └── chain.go           # Blockchain services
├── go.mod                      # Module definition
├── go.sum                      # Dependency checksums
└── README.md                   # Project documentation
```

## Setup Instructions
1. **Clone the repository:**
   ```
   git clone <repository-url>
   cd go-backend
   ```

2. **Install dependencies:**
   ```
   go mod tidy
   ```

3. **Run the application:**
   ```
   go run cmd/main.go
   ```

## Usage
- The API provides various endpoints for managing users, properties, and blockchain interactions. Refer to the `docs/Backend-API.md` for detailed information on available endpoints and their usage.

## Contributing
Contributions are welcome! Please submit a pull request or open an issue for any enhancements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for more details.
