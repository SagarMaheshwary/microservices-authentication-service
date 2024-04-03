# MICROSERVICES - USER SERVICE

This service is a part of the Microservices project built for handling authentication and authorization stuff.

### TECHNOLOGIES

- Golang
- gRPC
- JWT

### SETUP

cd into the project directory and copy **.env.example** to **.env** and update the required variables.

Create executable and start the server:

```bash
go build cmd/server/main.go && ./main
```

### APIs (RPC)

| SERVICE     | RPC                    | METADATA                               | DESCRIPTION                       | USECASE                                                                                           |
| ----------- | ---------------------- | -------------------------------------- | --------------------------------- | ------------------------------------------------------------------------------------------------- |
| AuthService | Register               | -                                      | User registration                 | User registration using "user" microservice's UserService.Store RPC                               |
| AuthService | Login                  | -                                      | User authentication via JWT token | User login using "user" microservice's UserService.FindByCredential RPC                           |
| AuthService | VerifyToken            | Bearer token in "authorization" header | Authorization Request             | Token Verification and returning user's data using "user" microservice's UserService.FindById RPC |
| AuthService | Logout (unimplemented) | -                                      | Add jwt token to blacklist        | User logout                                                                                       |
