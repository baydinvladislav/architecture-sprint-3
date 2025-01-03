# Table of Contents:
- [Description](#description)
- [Project Structure](#project-structure)
- [Auth logic](#auth-logic)
- [C4 User Service component diagram](#c4-user-service-component-diagram)
- [C4 User Service code diagram](#c4-user-service-code-diagram)

### Description:
The User Service is a component of the Clever Village System. It handles user registration, authentication, and profile management. 
The service includes JWT-based authentication, enabling secure and scalable user session management. 
It also manages interactions between users and connected houses, ensuring proper data handling and integration with other system components.

### Project Structure:
```
user-service/
├── business/                                     # Business logic configurations
│   └── feature_settings.go                       # Feature flag and settings management
├── docs/                                         # Documentation for the service
│   ├── docs.go                                   # GoDocs for the service
│   ├── swagger.json                              # API specification in JSON format
│   └── swagger.yaml                              # API specification in YAML format
├── middleware/                                   # Middleware for request handling
│   └── auth_middleware.go                        # Middleware for authentication
├── persistence/                                  # Database models
│   ├── house_model.go                            # House entity model
│   └── user_model.go                             # User entity model
├── presentation/                                 # API controllers for handling requests
│   ├── auth_controllers.go                       # Controllers for authentication endpoints
│   ├── house_controllers.go                      # Controllers for house-related endpoints
│   └── user_controllers.go                       # Controllers for user-related endpoints
├── repository/                                   # Data access layer
│   ├── house_repository.go                       # Repository for house data
│   └── user_repository.go                        # Repository for user data
├── schemas/                                      # Project data structures
│   ├── dto/                                      # Data transfer objects
│   ├── events/                                   # Event schemas
│   └── web/                                      # Web-related schemas
├── service/                                      # Core business logic
│   ├── auth_service.go                           # Service for authentication logic
│   ├── house_service.go                          # Service for house-related logic
│   ├── user_service.go                           # Service for user-related logic
│   ├── verify_connection_service.go              # Service for verifying connections
│   └── verify_connection_service_test.go         # Unit tests for verify connection service
├── shared/                                       # Shared utilities and configurations
│   ├── container.go                              # Dependency injection container
│   └── settings.go                               # Configuration settings
├── suppliers/                                    # External integrations
│   └── kafka_supplier.go                         # Kafka producer logic
├── Code_CleverVillageSystem_UserService.puml     # PlantUML diagrams
├── Dockerfile                                    # Docker image configuration
├── go.mod                                        # Go module dependencies
└── README.md                                     # Project documentation
```

### Auth logic:
During registration and login, the user receives a JWT access token and a JWT refresh token. 
When the access token expires (15 minutes), the user can send a request to `/refresh-token` to obtain a new pair 
of tokens using the refresh token. The refresh token is valid for one week. If the user does not use the application 
for a period exceeding the refresh token's validity, the application will require them to log in again.

### C4 User Service component diagram:
![System Architecture](./Component_CleverVillageSystem_UserService.svg)

### C4 User Service code diagram:
![System Architecture](./Code_CleverVillageSystem_UserService.svg)
