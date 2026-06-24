package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	goroute "github.com/Sahadat20/goroute"
)

// User represents a simple user model
type User struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// in-memory mock database (for demo only)
var userDB = make(map[string]User)
var idCounter = 1

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

	// Generate a new ID as a string
	id := strconv.Itoa(idCounter)
	idCounter++

	userDB[id] = newUser
	c.JSON(http.StatusCreated, map[string]string{
		"id":      id,
		"message": "User created successfully",
	})
}

// getUsers returns the stored all user
func getUsers(c *goroute.Context) {
	// The framework's JSON helper automatically encodes the entire map
	if len(userDB) == 0 {
		// Return an empty object or message if no users exist yet
		c.JSON(http.StatusOK, map[string]string{"message": "No users found"})
		return
	}
	c.JSON(http.StatusOK, userDB)
}

// getUser returns the stored user
func getUser(c *goroute.Context) {
	id := c.Param("id") // Extract ID from the URL

	if user, exists := userDB[id]; exists {
		c.JSON(http.StatusOK, user)
	} else {
		c.String(http.StatusNotFound, "User not found")
	}
}

// updateUser updates existing user data
func updateUser(c *goroute.Context) {
	id := c.Param("id") // Extract ID from the URL

	// Ensure the user exists before updating
	if _, exists := userDB[id]; !exists {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	var updatedUser User
	if err := c.BindJSON(&updatedUser); err != nil {
		c.String(http.StatusBadRequest, "Invalid JSON data")
		return
	}

	// Overwrite the existing data at this specific ID
	userDB[id] = updatedUser
	c.JSON(http.StatusOK, map[string]string{"message": "User " + id + " updated"})
}

// deleteUser removes user from memory
func deleteUser(c *goroute.Context) {
	id := c.Param("id") // Extract ID from the URL

	if _, exists := userDB[id]; exists {
		delete(userDB, id)
		c.String(http.StatusOK, "User "+id+" deleted")
	} else {
		c.String(http.StatusNotFound, "User not found")
	}
}

// Logger is a custome global middleware

func Logger() goroute.RouteHandler {
	return func(c *goroute.Context) {
		// 1. Pre-Procesing
		t := time.Now()
		fmt.Printf("[START] request incoming: %s %s\n", c.Method, c.Path)
		// 2. pass control to the next middleware or handler
		c.Next()

		// 3. post processing (execute after the entire request is handled)
		latency := time.Since(t)
		fmt.Printf("[END] request completed: %s %s in %v\n", c.Method, c.Path, latency)

	}
}

func main() {
	// create new GoRoute instance
	app := goroute.New()

	// register global middleware using .Use()
	app.Use(Logger())
	// basic routes
	app.GET("/", welcomeHandler)
	app.GET("/api/info", infoHandler)

	// CRUD routes (demo API)
	app.POST("/user", createNewUser)
	app.GET("/users", getUsers)
	app.GET("/user/:id", getUser)
	app.PUT("/user/:id", updateUser)
	app.DELETE("/user/:id", deleteUser)
	// (Optional) Keep the wildcard route from earlier to show multiple features coexisting
	app.GET("/static/*filepath", func(c *goroute.Context) {
		file := c.Param("filepath")
		c.String(http.StatusOK, "Simulating serving file: "+file)
	})
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
