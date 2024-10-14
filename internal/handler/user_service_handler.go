package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/haiyen11231/eco-taxi-api-gateway/internal/grpc/pb"
	"github.com/haiyen11231/eco-taxi-api-gateway/internal/utils"

	"github.com/gin-gonic/gin"
)

// for signup and update
type User struct {
	Name        string `json:"name" binding:"required"`
	PhoneNumber string `json:"phone_number" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type LogInUser struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type VerifyUser struct {
	Token string `json:"token" binding:"required"`
}

func SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		signUpUser := User{}

		// Binds the incoming request to signUp.
		if err := ctx.ShouldBindJSON(&signUpUser); err != nil {
			log.Println("Failed to binding json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishes a gRPC connection.
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Sends a SignUpRequest with user details to the gRPC service for signup.
		response, err := client.SignUp(c, &pb.SignUpRequest{
			Name: signUpUser.Name,
			PhoneNumber: signUpUser.PhoneNumber,
			Email: signUpUser.Email,
			Password: signUpUser.Password,
		})

		// If signup fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to signup", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func LogIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logInUser := LogInUser{}

		// Binds the incoming request to logIn.
		if err := ctx.ShouldBindJSON(&logInUser); err != nil {
			log.Println("Failed to binding json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishes a gRPC connection.
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Sends a LogInRequest with user details to the gRPC service for login.
		response, err := client.LogIn(c, &pb.LogInRequest{
			PhoneNumber: logInUser.PhoneNumber,
			Password: logInUser.Password,
		})

		// If signup fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to login", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		updateUser := User{}

		// Binds the incoming request to update user.
		if err := ctx.ShouldBindJSON(&updateUser); err != nil {
			log.Println("Failed to binding json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishes a gRPC connection.
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Sends a UpdateUserRequest with user details to the gRPC service for updating user.
		response, err := client.UpdateUser(c, &pb.UpdateUserRequest{
			Name: updateUser.Name,
			PhoneNumber: updateUser.PhoneNumber,
			Email: updateUser.Email,
			Password: updateUser.Password,
		})

		// If updating user fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to update user", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetUint64("user_id")

		// Establishes a gRPC connection.
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Sends a GetUserRequest to the gRPC service for getting user.
		response, err := client.GetUser(c, &pb.GetUserRequest{
			Id: userId,
		})

		// If getting user fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to get user", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func Verify() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		verifyUser := VerifyUser{}

		// Binds the incoming request to verify user.
		if err := ctx.ShouldBindJSON(&verifyUser); err != nil {
			log.Println("Failed to binding", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishes a gRPC connection.
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Sends a VerifyRequest with user token to the gRPC service for validation.
		response, err := client.Verify(c, &pb.VerifyRequest{
			Token: verifyUser.Token,
		})

		// If verifying user fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to verify", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}