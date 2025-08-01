package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq" // PostgreSQL driver
	"golang.org/x/crypto/bcrypt"
)

var db *sql.DB

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Spending struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

func main() {
	// Connect to the database
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ping the database to ensure connectivity
	err = db.Ping()
	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	log.Println("Successfully connected to the database!")

	// Initialize the database schema
	initDB()

	// Set up Gin router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost"}, // Nginx is on port 80
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	api := router.Group("/api")
	{
		api.POST("/register", registerHandler)
		api.POST("/login", loginHandler)
		api.GET("/spending-summary", spendingSummaryHandler)
		api.POST("/generate-summary", generateSummaryHandler)
	}

	log.Fatal(router.Run(":8080"))
}

func initDB() {
	// Create users table
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            username TEXT UNIQUE NOT NULL,
            password_hash TEXT NOT NULL
        );
    `)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	// Create spending table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS spending (
            id SERIAL PRIMARY KEY,
            category TEXT NOT NULL,
            amount NUMERIC(10, 2) NOT NULL
        );
    `)
	if err != nil {
		log.Fatalf("Failed to create spending table: %v", err)
	}

	// Insert mock spending data if the table is empty
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM spending").Scan(&count)
	if err == nil && count == 0 {
		mockData := []Spending{
			{Category: "Food", Amount: 150.50},
			{Category: "Transport", Amount: 75.00},
			{Category: "Entertainment", Amount: 200.00},
			{Category: "Utilities", Amount: 120.00},
		}
		for _, item := range mockData {
			db.Exec("INSERT INTO spending (category, amount) VALUES ($1, $2)", item.Category, item.Amount)
		}
		log.Println("Inserted mock spending data.")
	}
}

func registerHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	// Insert user into the database
	_, err = db.Exec("INSERT INTO users (username, password_hash) VALUES ($1, $2)", user.Username, string(hashedPassword))
	if err != nil {
		// Check for unique constraint violation
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func loginHandler(c *gin.Context) {
	var user User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	var storedHash string
	err := db.QueryRow("SELECT password_hash FROM users WHERE username = $1", user.Username).Scan(&storedHash)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func spendingSummaryHandler(c *gin.Context) {
	rows, err := db.Query("SELECT category, amount FROM spending")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch spending data"})
		return
	}
	defer rows.Close()

	var spendingData []Spending
	for rows.Next() {
		var item Spending
		if err := rows.Scan(&item.Category, &item.Amount); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan spending data"})
			return
		}
		spendingData = append(spendingData, item)
	}

	c.JSON(http.StatusOK, spendingData)
}

func generateSummaryHandler(c *gin.Context) {
	// In a real application, you would make a request to a third-party AI service here.
	// You would use the AI_API_KEY from the .env file.
	// For this example, we'll just return a mock summary.
	c.JSON(http.StatusOK, gin.H{"summary": "Based on your spending, you are doing great! Keep it up."})
}