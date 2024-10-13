package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/haiyen11231/eco-taxi-api-gateway/internal/handlers"
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
    user.POST("/signup", handlers.Register()) // Endpoint for user signup
    user.POST("/login", handlers.Login()) // Endpoint for user login
    user.PATCH("/update", handlers.Login()) // Endpoint for user update 
    user.GET("/", handlers.Login()) // Endpoint for getting user info
    user.POST("/verify", handlers.Verify()) // Endpoint for verifying user authentication

//     rpc GetTripPreview(GetTripPreviewRequest) returns(GetTripPreviewResponse);
//   rpc ConfirmBooking(ConfirmBookingRequest) returns(ConfirmBookingResponse);
//   rpc UpdateBookingStatus(UpdateBookingRequest) returns (UpdateBookingResponse);
//   rpc GetBookingHistory(GetBookingHistoryRequest) returns (GetBookingHistoryResponse);
    trip := v1.Group("/trip")
    trip.GET("/", handlers.GetAllProduct()) // Fetch all products
    trip.GET("/:id", handlers.GetProduct()) // Fetch a single product by ID

// rpc GetCards(GetCardsRequest) returns(GetCardsResponse);
//   rpc CreateCard(CreateCardRequest) returns(CreateCardResponse);
//   rpc UpdatecCard(UpdateCardRequest) returns(UpdateCardResponse);
//   rpc DeleteCard(DeleteCardRequest) returns(DeleteCardResponse);
    payment := v1.Group("/payment")
    payment.GET("/", handlers.GetAllProduct()) // Fetch all products
    payment.GET("/:id", handlers.GetProduct()) // Fetch a single product by ID

    trip.Use(middleware.VerifyToken) // Middleware to verify JWT tokens
    trip.POST("/", handlers.CreateProduct()) // Create a new product
    trip.PATCH("/:id", handlers.UpdateProduct()) // Update product by ID
    trip.DELETE("/:id", handlers.DeleteProduct()) // Delete product by ID

    r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))

}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}