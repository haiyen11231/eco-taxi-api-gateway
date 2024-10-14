package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/haiyen11231/eco-taxi-api-gateway/internal/handler"
	"github.com/haiyen11231/eco-taxi-api-gateway/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	loadEnv() // Loads environment variables from the .env file

	r := gin.Default() // Creates a new Gin router with default middleware

    r.Use(cors.New(cors.Config{ // Adds CORS middleware
		AllowOrigins:     []string{"http://localhost"}, // Specifies allowed origins
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}, // Specifies allowed HTTP methods
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Specifies allowed headers
		ExposeHeaders:    []string{"Link"}, // Specifies headers exposed to the client
		AllowCredentials: true, // Allows credentials (e.g., cookies, authorization headers)
		MaxAge:           300, // Caches the CORS preflight response for 300 seconds
	}))

    v1 := r.Group("/v1")

    v1.GET("/ping", func(ctx *gin.Context) {
        ctx.JSON(http.StatusAccepted, gin.H{"ok": true})
   })

    user := v1.Group("/user")
    user.POST("/signup", handler.SignUp()) // Endpoint for user signup
    user.POST("/login", handler.LogIn()) // Endpoint for user login
    user.PATCH("/update", handler.UpdateUser()) // Endpoint for user update 
    user.GET("/", handler.GetUser()) // Endpoint for getting user info
    user.POST("/authenticate", handler.Verify()) // Endpoint for verifying user authentication

    trip := v1.Group("/trip")
    trip.POST("/", handler.SearchTripPreview())
    trip.Use(middleware.VerifyToken) // Middleware to verify JWT tokens
    trip.POST("/confirm", handler.ConfirmBooking())
    trip.PATCH("/:id", handler.UpdateBookingStatus())
    trip.GET("/history", handler.GetBookingHistory())
    
    payment := v1.Group("/payment")
    payment.Use(middleware.VerifyToken) // Middleware to verify JWT tokens
    payment.GET("/", handler.GetCards()) // Fetch all cards
    payment.POST("/create", handler.CreateCard()) // Create a new card
    payment.PATCH("/:id", handler.UpdateCard()) // Update card by ID
    payment.DELETE("/:id", handler.DeleteCard()) // Delete card by ID

    r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))

}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}