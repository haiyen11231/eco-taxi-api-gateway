package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/haiyen11231/eco-taxi-api-gateway/internal/grpc/pb"
	"github.com/haiyen11231/eco-taxi-api-gateway/internal/model"
	"github.com/haiyen11231/eco-taxi-api-gateway/internal/utils"

	"github.com/gin-gonic/gin"
)

func SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userData := model.SignUpUserData{}

		// Binding the incoming request to signup
		if err := ctx.ShouldBindJSON(&userData); err != nil {
			log.Println("Failed to bind json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishing a gRPC connection
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

		// Sending a SignUpRequest with user details to the gRPC service for signup
		response, err := client.SignUp(c, &pb.SignUpRequest{
			Name: userData.Name,
			PhoneNumber: userData.PhoneNumber,
			Email: userData.Email,
			Password: userData.Password,
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
		logInUserData := model.LogInUserData{}

		// Binding the incoming request to login
		if err := ctx.ShouldBindJSON(&logInUserData); err != nil {
			log.Println("Failed to bind json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishing a gRPC connection
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

		// Sending a LogInRequest with user details to the gRPC service for login
		response, err := client.LogIn(c, &pb.LogInRequest{
			PhoneNumber: logInUserData.PhoneNumber,
			Password: logInUserData.Password,
		})

		// If signup fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to login", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		ctx.SetSameSite(http.SameSiteLaxMode)
		ctx.SetCookie("Authorization", response.Token, 3600 * 24 * 2, "", "", false, true)

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func ForgotPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		forgotPasswordUserData := model.ForgotPasswordUserData{}

		// Binding the incoming request to reset password
		if err := ctx.ShouldBindJSON(&forgotPasswordUserData); err != nil {
			log.Println("Failed to bind JSON", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishing a gRPC connection
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

		// Sending a ForgotPasswordRequest with user details to the gRPC service for resetting password
		response, err := client.ForgotPassword(c, &pb.ForgotPasswordRequest{
			Email: forgotPasswordUserData.Email,
			NewPassword: forgotPasswordUserData.NewPassword,
		})

		// If resetting password fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to reset password", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func UpdateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieving the user_id from the context, set previously in middleware
        userId := ctx.GetUint64("user_id")

		userData := model.UpdateUserData{}

		// Binding the incoming request to update user
		if err := ctx.ShouldBindJSON(&userData); err != nil {
			log.Println("Failed to bind json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishing a gRPC connection
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

		// Sending a UpdateUserRequest with user details to the gRPC service for updating user
		response, err := client.UpdateUser(c, &pb.UpdateUserRequest{
			Id: userId,
			Name: userData.Name,
			PhoneNumber: userData.PhoneNumber,
			Email: userData.Email,
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

		// Establishing a gRPC connection
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

		// Sending a GetUserRequest to the gRPC service for getting user
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

func ChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieving the user_id from the context, set previously in middleware
        userId := ctx.GetUint64("user_id")

		changePasswordUserData := model.ChangePasswordUserData{}

		// Binding the incoming request to change password
		if err := ctx.ShouldBindJSON(&changePasswordUserData); err != nil {
			log.Println("Failed to bind json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishing a gRPC connection
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

		// Sending a ChangePasswordRequest with user details to the gRPC service for changing password
		response, err := client.ChangePassword(c, &pb.ChangePasswordRequest{
			Id: userId,
			OldPassword: changePasswordUserData.OldPassword,
			NewPassword: changePasswordUserData.NewPassword,
		})

		// If changing password fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to change password", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func UpdateDistanceTravelled() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieving the user_id from the context, set previously in middleware
        userId := ctx.GetUint64("user_id")

		updateDistanceUserData := model.UpdateDistanceUserData{}

		// Binding the incoming request to update distance travelled
		if err := ctx.ShouldBindJSON(&updateDistanceUserData); err != nil {
			log.Println("Failed to bind json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishing a gRPC connection
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

		// Sending a UpdateDistanceTravelledRequest with user details to the gRPC service for updating distance travelled
		response, err := client.UpdateDistanceTravelled(c, &pb.UpdateDistanceTravelledRequest{
			Id: userId,
			Distance: updateDistanceUserData.Distance,
		})

		// If updating distance travelled fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to update distance travelled", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)

	}
}

func AuthenticateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authenticateUserData := model.AuthenticateUserData{}

		// Binding the incoming request to verify user.
		if err := ctx.ShouldBindJSON(&authenticateUserData); err != nil {
			log.Println("Failed to bind json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishing a gRPC connection
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

		// Sending an AuthenticateUserRequest with user token to the gRPC service for authentication
		response, err := client.AuthenticateUser(c, &pb.AuthenticateUserRequest{
			Token: authenticateUserData.Token,
		})

		// If authentication fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to authenticate", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}