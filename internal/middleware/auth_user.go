package middleware

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/haiyen11231/eco-taxi-api-gateway/internal/grpc/pb"
	"github.com/haiyen11231/eco-taxi-api-gateway/internal/utils"

	"github.com/gin-gonic/gin"
)

func AuthenticateUser (ctx *gin.Context) {
	// Extracting and validating the Bearer token from incoming requests
	token := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")

	if token == "" {
		log.Println("Token required")
		utils.ResponseError(ctx, http.StatusUnauthorized, "Unauthorized!")
		return
	}

	// Establishing a gRPC connection
	conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
	if err != nil {
		log.Println("Error Connection to GRPC", err)
		utils.ResponseError(ctx, http.StatusUnauthorized, "Unauthorized!")
		return
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)
	c, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// Sending an AuthenticateUserRequest with user token to the gRPC service for authentication
	response, err := client.AuthenticateUser(c, &pb.AuthenticateUserRequest{
		Token: token,
	})
	
	// If authentication fails, logs the error and returns a 401 Unauthorized error. On success, it sends a success response.
	if err != nil {
		log.Println("Failed to authenticate", err)
		utils.ResponseError(ctx, http.StatusUnauthorized, "Unauthorized!")
		return
	}
	
	if !response.Valid {
		log.Println("Error ", err)
		utils.ResponseError(ctx, http.StatusUnauthorized, "Unauthorized!")
		return
	}

	ctx.Set("user_id", response.UserId)
	ctx.Next()
}