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

| SERVICE     | RPC         | METADATA                               | DESCRIPTION                                                                                                  |
| ----------- | ----------- | -------------------------------------- | ------------------------------------------------------------------------------------------------------------ |
| AuthService | Register    | -                                      | User registration via [user microservice](https://github.com/SagarMaheshwary/microservices-user-service) RPC |
| AuthService | Login       | -                                      | User login via user microservice                                                                             |
| AuthService | VerifyToken | Bearer token in "authorization" header | Token verification and getting user data via user microservice                                               |
| AuthService | Logout      | Bearer token in "authorization" header | User logout by adding token to redis blacklist.                                                              |
