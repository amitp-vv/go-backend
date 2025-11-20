# Backend API Documentation

## Overview
This document provides an overview of the backend API, detailing the available endpoints, request and response formats, and usage examples.

## Base URL
The base URL for all API requests is:
```
http://localhost:8080/api
```

## Endpoints

### 1. Authentication

#### POST /auth/login
- **Description**: Authenticates a user and returns a token.
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string"
  }
  ```
- **Response**:
  - **200 OK**:
    ```json
    {
      "token": "string"
    }
    ```
  - **401 Unauthorized**:
    ```json
    {
      "error": "Invalid credentials"
    }
    ```

#### POST /auth/register
- **Description**: Registers a new user.
- **Request Body**:
  ```json
  {
    "username": "string",
    "password": "string",
    "email": "string"
  }
  ```
- **Response**:
  - **201 Created**:
    ```json
    {
      "message": "User registered successfully"
    }
    ```
  - **400 Bad Request**:
    ```json
    {
      "error": "User already exists"
    }
    ```

### 2. Admin Operations

#### GET /admin/users
- **Description**: Retrieves a list of all users.
- **Response**:
  - **200 OK**:
    ```json
    [
      {
        "id": "string",
        "username": "string",
        "email": "string"
      }
    ]
    ```

### 3. Property Management

#### GET /properties
- **Description**: Retrieves a list of properties.
- **Response**:
  - **200 OK**:
    ```json
    [
      {
        "id": "string",
        "name": "string",
        "location": "string"
      }
    ]
    ```

### 4. Claims

#### POST /claims
- **Description**: Submits a new claim.
- **Request Body**:
  ```json
  {
    "propertyId": "string",
    "description": "string"
  }
  ```
- **Response**:
  - **201 Created**:
    ```json
    {
      "message": "Claim submitted successfully"
    }
    ```

## Error Handling
All error responses will follow the format:
```json
{
  "error": "Error message"
}
```

## Usage Examples
### Example: User Login
```bash
curl -X POST http://localhost:8080/api/auth/login -H "Content-Type: application/json" -d '{"username":"user","password":"pass"}'
```

### Example: Get Properties
```bash
curl -X GET http://localhost:8080/api/properties -H "Authorization: Bearer your_token"
```

## Conclusion
This API provides a robust interface for managing users, properties, and claims within the application. For further details, refer to the specific endpoint documentation.