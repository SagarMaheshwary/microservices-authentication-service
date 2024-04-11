# MICROSERVICES - AUTHENTICATION SERVICE

This service is a part of the Microservices project built for handling authentication and authorization stuff.

### TECHNOLOGIES

- Golang (1.22.2)
- Redis (7.2)
- gRPC
- JWT

### SETUP

cd into the project directory and copy **.env.example** to **.env** and update the required variables.

Create executable and start the server:

```bash
go build cmd/server/main.go && ./main
```

Or install "[air](https://github.com/cosmtrek/air)" and run it to autoreload when making file changes:

```bash
air -c .air-toml
```

### APIs (RPC)

| SERVICE     | RPC         | METADATA                               | DESCRIPTION                                                                                   |
| ----------- | ----------- | -------------------------------------- | --------------------------------------------------------------------------------------------- |
| AuthService | Register    | -                                      | User registration using "user" microservice's UserService.Store RPC                           |
| AuthService | Login       | -                                      | User login using "user" microservice's UserService.FindByCredential RPC                       |
| AuthService | VerifyToken | Bearer token in "authorization" header | Token verification and getting user data using "user" microservice's UserService.FindById RPC |
| AuthService | Logout      | Bearer token in "authorization" header | User logout by adding token to redis blacklist.                                               |
