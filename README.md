# 🚀 GoRoute

A lightweight, fast, and minimal web framework for Go inspired by Express.js and Gin.

GoRoute is designed for simplicity, performance, and developer productivity while keeping the API clean and easy to use.

---

## ✨ Features

- Simple and expressive routing
- GET / POST / PUT / DELETE support
- JSON response helpers
- Middleware-ready architecture (coming soon)
- Lightweight and fast
- Minimal boilerplate
- Easy to learn API

---

## 📦 Installation

```bash
go get github.com/Sahadat20/goroute
```

---

## 🚀 Quick Start

```go
package main

import (
	"net/http"

	goroute "github.com/Sahadat20/goroute"
)

func main() {
	app := goroute.New()

	app.GET("/", func(c *goroute.Context) {
		c.String(http.StatusOK, "Welcome to GoRoute 🚀")
	})

	app.Run(":8080")
}
```

---

## 📚 Routing Examples

### ➤ Basic GET Route

```go
app.GET("/hello", func(c *goroute.Context) {
	c.String(200, "Hello World")
})
```

---

### ➤ JSON Response

```go
app.GET("/json", func(c *goroute.Context) {
	c.JSON(200, map[string]interface{}{
		"framework": "GoRoute",
		"version":   "1.0.0",
	})
})
```

---

### ➤ POST Request (JSON Body)

```go
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

app.POST("/user", func(c *goroute.Context) {
	var user User

	if err := c.BindJSON(&user); err != nil {
		c.String(400, "Invalid JSON")
		return
	}

	c.JSON(201, map[string]string{
		"message": "User created",
	})
})
```

---

## 🧪 Example Project

Run the included example:

```bash
go run examples/basic/main.go
```

Then open:

```
http://localhost:8081
```

---

## 📁 Project Structure

```
goroute/
├── context/
├── router/
├── examples/
│   └── basic/
├── go.mod
├── README.md
└── LICENSE
```

---

## 🛠 Roadmap

- [x] Basic routing (GET, POST, PUT, DELETE)
- [x] JSON response helpers
- [ ] Middleware system
- [ ] Route parameters (/user/:id)
- [ ] Group routing (/api/v1)
- [ ] Logger middleware
- [ ] Recovery middleware
- [ ] Static file serving
- [ ] WebSocket support

---

## ⚡ Why GoRoute?

GoRoute is built for developers who want:

- A minimal alternative to large frameworks
- Clean and readable API design
- Fast backend development in Go
- Express-like simplicity in Go ecosystem

---

## 📖 Philosophy

> “Simple APIs, fast performance, zero complexity.”

---

## 📜 License

MIT License

---

## 👨‍💻 Author

Sahadat Hossain  
https://github.com/Sahadat20

---

## ⭐ Support

If you like this project, consider giving it a star ⭐ on GitHub.
