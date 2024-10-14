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

func VerifyToken (ctx *gin.Context) {
	// Extracting and validating the Bearer token from incoming requests
	token := strings.TrimPrefix(ctx.GetHeader("Authorization"), "Bearer ")

	if token == "" {
		log.Println("Token required")
		utils.ResponseError(ctx, http.StatusUnauthorized, "Unauthorized!")
		return
	}

	// Using a gRPC call to verify the token with user-service
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

	response, err := client.Verify(c, &pb.VerifyRequest{
		Token: token,
	})
	
	if err != nil {
		log.Println("Error Verify", err)
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