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
    user.POST("/signup", handler.SignUp()) 
    user.PATCH("/reset-password", handler.ForgotPassword())
    user.Use(middleware.AuthenticateUser)
    user.PATCH("/update", handler.UpdateUser()) 
    user.GET("/", handler.GetUser()) 
    user.PATCH("/change-password", handler.ChangePassword())
    user.PATCH("/update-distance", handler.UpdateDistanceTravelled())
    user.POST("/authenticate", handler.AuthenticateUser()) 

    trip := v1.Group("/trip")
    trip.POST("/", handler.SearchTripPreview())
    trip.Use(middleware.AuthenticateUser) 
    trip.POST("/confirm", handler.ConfirmBooking())
    trip.GET("/incompleted-booking", handler.GetIncompletedBooking())
    trip.PATCH("/:id", handler.UpdateBookingStatus())
    trip.GET("/history", handler.GetBookingHistory())
    
    payment := v1.Group("/payment")
    payment.Use(middleware.AuthenticateUser)
    payment.GET("/", handler.GetCards())
    payment.POST("/create", handler.CreateCard()) 
    payment.PATCH("/:id", handler.UpdateCard()) 
    payment.DELETE("/:id", handler.DeleteCard())

    r.Run(fmt.Sprintf(":%s", os.Getenv("PORT")))

}

func loadEnv() {
	err := godotenv.Load("app.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}