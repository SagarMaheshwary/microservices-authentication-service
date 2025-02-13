# MICROSERVICES - AUTHENTICATION SERVICE

Authentication Service for the [Microservices](https://github.com/SagarMaheshwary/microservices) project.

### OVERVIEW

- Golang
- ZeroLog
- gRPC – Acts as both the main server and client for the User service
- JWT (JSON Web Tokens) – Used for authentication and secure communication
- Redis - Maintains the token blacklist
- Prometheus Client – Exports default and custom metrics for Prometheus server monitoring

### SETUP

Follow the instructions in the [README](https://github.com/SagarMaheshwary/microservices?tab=readme-ov-file#setup) of the main microservices repository to run this service along with others using Docker Compose.

### APIs (gRPC)

Proto files are located in the **internal/proto** directory.

| SERVICE                                                        | RPC         | METADATA                            | DESCRIPTION                                    |
| -------------------------------------------------------------- | ----------- | ----------------------------------- | ---------------------------------------------- |
| AuthService                                                    | Register    | -                                   | User registration                              |
| AuthService                                                    | Login       | -                                   | User login                                     |
| AuthService                                                    | VerifyToken | Bearer token in "authorization" key | Token verification and getting user data       |
| AuthService                                                    | Logout      | Bearer token in "authorization" key | User logout by adding token to redis blacklist |
| [Health](https://google.golang.org/grpc/health/grpc_health_v1) | Check       | -                                   | Service health check                           |

### APIs (REST)

| API      | METHOD | BODY | Headers | Description                 |
| -------- | ------ | ---- | ------- | --------------------------- |
| /metrics | GET    | -    | -       | Prometheus metrics endpoint |
