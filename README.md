# Go Clean Code DDD
This is a Go project that follows the Clean Code Architecture and Domain-Driven Design (DDD) principles. The project is still under development, and I'm still learning about these topics, so this is just my approach to implementing them.

### Project Structure
The project has the following directory structure:
```shell
├── README.md
├── cmd
│   └── main.go
├── go.mod
├── go.sum
├── internal
│   ├── application
│   │   └── usecase
│   │       ├── auth_usecase.go
│   │       └── user_usecase.go
│   ├── config
│   │   ├── config.go
│   │   └── config.json
│   ├── domain
│   │   ├── entity
│   │   │   ├── cache.go
│   │   │   ├── token.go
│   │   │   └── user.go
│   │   ├── repository
│   │   │   ├── cache_repository.go
│   │   │   └── user_repository.go
│   │   └── service
│   │       ├── auth_service.go
│   │       ├── cache_service.go
│   │       └── user_service.go
│   ├── infrastructure
│   │   ├── cache.go
│   │   ├── database.go
│   │   ├── memory
│   │   │   ├── cache_repository.go
│   │   │   └── cache_repository_test.go
│   │   └── persistence
│   │       └── user_repository.go
│   ├── interface
│   │   ├── api
│   │   │   ├── auth_handler.go
│   │   │   ├── routes.go
│   │   │   └── user_handler.go
│   │   ├── middleware
│   │   │   └── auth.go
│   │   └── validator
│   │       └── validator.go
│   └── presentation
│       └── model
│           ├── auth_request.go
│           ├── auth_response.go
│           ├── user_request.go
│           └── user_response.go
└── pkg
    ├── password
    │   └── password.go
    └── uuid
        └── uuid.go

```

### Usage
To run the project, first you need to create config file
```shell
cp config-example.json config.json
```
Then execute the following command:
```shell
go run cmd/main.go
```
This will start the application and listen for HTTP requests on the defined port in `config.json`.

### Contributing
As this is a learning project, I welcomes any feedback, suggestions, or contributions that can help improve the project's design and implementation.

### License
This project is licensed under the [MIT License](LICENSE).