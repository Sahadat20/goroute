package main

import (
	"fmt"
	"net/http"

	goroute "github.com/Sahadat20/goroute"
)

// User represents a simple user model
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// in-memory mock database (for demo only)
var userDB = make(map[string]User)

// welcomeHandler handles root endpoint
func welcomeHandler(c *goroute.Context) {
	c.String(http.StatusOK, "Welcome to GoRoute 🚀")
}

// infoHandler returns framework metadata
func infoHandler(c *goroute.Context) {
	c.JSON(http.StatusOK, map[string]interface{}{
		"framework": "goroute",
		"version":   "1.0",
		"author":    "Sahadat Hossain",
	})
}

// createNewUser creates a new user (demo only, uses in-memory storage)
func createNewUser(c *goroute.Context) {
	var newUser User

	// bind JSON body into struct
	if err := c.BindJSON(&newUser); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// store user (static key for demo)
	userDB["1"] = newUser

	c.JSON(http.StatusCreated, map[string]string{
		"message": "User created successfully",
	})
}

// getUser returns the stored user
func getUser(c *goroute.Context) {
	user, exists := userDB["1"]

	if !exists {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

// updateUser updates existing user data
func updateUser(c *goroute.Context) {
	if _, exists := userDB["1"]; !exists {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	var updatedUser User

	// bind request body
	if err := c.BindJSON(&updatedUser); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data")
		return
	}

	userDB["1"] = updatedUser

	c.JSON(http.StatusOK, map[string]string{
		"message": "User updated successfully",
	})
}

// deleteUser removes user from memory
func deleteUser(c *goroute.Context) {
	delete(userDB, "1")

	c.String(http.StatusOK, "User deleted successfully")
}

func main() {
	// create new GoRoute instance
	app := goroute.New()

	// basic routes
	app.GET("/", welcomeHandler)
	app.GET("/api/info", infoHandler)

	// CRUD routes (demo API)
	app.POST("/user", createNewUser)
	app.GET("/user", getUser)
	app.PUT("/user", updateUser)
	app.DELETE("/user", deleteUser)

	// start server
	addr := ":8081"

	fmt.Println("===================================")
	fmt.Println("🚀 GoRoute Server Starting...")
	fmt.Println("📍 URL: http://localhost" + addr)
	fmt.Println("===================================")

	if err := app.Run(addr); err != nil {
		fmt.Println("❌ Failed to start server:", err)
	}
}
