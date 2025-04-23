# ChatterG-go 🚀  

ChatterG-go is a Go-based web application that provides a chat platform with WebSocket support, user management, and authentication.  

## Features ✨  

- **User Management**: Create, retrieve, and manage users. 👤  
- **WebSocket Support**: Real-time communication using WebSocket. 🔗  
- **Authentication**: JWT-based authentication for secure access. 🔒  
- **Database Integration**: Supports CockroachDB for relational data and MongoDB for chat messages. 🗄️  
- **Health Check**: Endpoint to verify the server's health. ✅  

## Project Structure 🗂️  

```  
.  
├── main.go  
├── internal/  
│   ├── database/  
│   │   └── databases.go  
│   ├── handlers/  
│   │   ├── handlers.go  
│   │   ├── user.go  
│   │   └── websockethandler.go  
│   ├── middleware/  
│   │   └── auth.go  
│   ├── models/  
│   │   └── user.go  
│   ├── routes/  
│   │   └── routes.go  
│   └── services/  
│       └── websocket.go  
├── pkg/  
│   └── utils/  
│       └── go-jwt.go  
├── go.mod  
├── go.sum  
└── README.md  
```  

## Installation 🛠️  

1. Clone the repository:  
    ```bash  
    git clone https://github.com/Piyu-Pika/ChatterG-go.git
    cd ChatterG-go
    ```  

2. Install dependencies:  
    ```bash  
    go mod tidy  
    ```  

3. Run the application:  
    ```bash  
    go run main.go  
    ```  

## Endpoints 🌐  

### Health Check ✅  
- **GET** `/health`  
  - Response: `{ "status": "ok" }`  

### User Management 👤  
- **POST** `/save-data`  
  - Save a new user.  
- **GET** `/users`  
  - Retrieve all users.  
- **GET** `/users/:id`  
  - Retrieve a user by ID.  

### WebSocket 🔗  
- **GET** `/ws`  
  - Establish a WebSocket connection.  

## Environment Variables 🌍  

- `COCKROACHDB_DSN`: Connection string for CockroachDB.  
- `MONGODB_URI`: Connection string for MongoDB.  

## Dependencies 📦  

- [Gin](https://github.com/gin-gonic/gin) - Web framework.  
- [GORM](https://gorm.io/) - ORM for relational databases.  
- [MongoDB Driver](https://github.com/mongodb/mongo-go-driver) - MongoDB integration.  
- [JWT](https://github.com/dgrijalva/jwt-go) - JSON Web Token for authentication.  

## License 📜  

This project is licensed under the MIT License.  

## Author ✍️  

[Piyu-Pika](https://github.com/Piyu-Pika)  