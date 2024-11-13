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

    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:5173"}, // Allow your frontend origin
        AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}, // Allowed methods
        AllowHeaders:     []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Allowed headers
        ExposeHeaders:    []string{"Link"}, // Headers exposed to the frontend
        AllowCredentials: true, // Allows cookies or Authorization headers
        MaxAge:           300, // Cache duration for preflight responses
    }))

    v1 := r.Group("/v1")

    v1.GET("/ping", func(ctx *gin.Context) {
        ctx.JSON(http.StatusAccepted, gin.H{"ok": true})
   })

    user := v1.Group("/user")
    user.POST("/signup", handler.SignUp()) 
    user.POST("/login", handler.LogIn()) 
    user.PATCH("/reset-password", handler.ForgotPassword())
    user.POST("/refresh-token", handler.RefreshToken()) 
    user.Use(middleware.AuthenticateUser)
    user.PATCH("/update", handler.UpdateUser()) 
    user.GET("/", handler.GetUser()) 
    user.PATCH("/change-password", handler.ChangePassword())
    user.PATCH("/update-distance", handler.UpdateDistanceTravelled())
    user.POST("/authenticate", handler.AuthenticateUser()) 
    user.DELETE("/logout", handler.LogOut()) 

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