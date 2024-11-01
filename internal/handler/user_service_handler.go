package handler

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/haiyen11231/eco-taxi-api-gateway/internal/grpc/pb"
	"github.com/haiyen11231/eco-taxi-api-gateway/internal/model"
	"github.com/haiyen11231/eco-taxi-api-gateway/internal/utils"

	"github.com/gin-gonic/gin"
)

func SignUp() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userData := model.SignUpUserData{}

		// Binding and validating incoming request for signup
		if err := ctx.ShouldBindJSON(&userData); err != nil {
			log.Println("Failed to bind JSON for SignUp:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "Invalid request format")
			return
		}

		// Establishing a gRPC connection
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial gRPC service:", err)
			utils.ResponseError(ctx, http.StatusInternalServerError, "Service unavailable")
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second) 
		defer cancel()

		// Sending a SignUpRequest with user details to the gRPC service
		response, err := client.SignUp(c, &pb.SignUpRequest{
			Name:        userData.Name,
			PhoneNumber: userData.PhoneNumber,
			Email:       userData.Email,
			Password:    userData.Password, 
		})

		if err != nil {
			log.Println("Failed to signup:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "Signup failed")
			return
		}

		utils.ResponseSuccess(ctx, http.StatusCreated, response)
	}
}

func LogIn() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logInUserData := model.LogInUserData{}

		// Binding and validating incoming request for login
		if err := ctx.ShouldBindJSON(&logInUserData); err != nil {
			log.Println("Failed to bind JSON for LogIn:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "Invalid request format")
			return
		}

		// Establishing a gRPC connection
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial gRPC service:", err)
			utils.ResponseError(ctx, http.StatusInternalServerError, "Service unavailable")
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Sending a LogInRequest to the gRPC service for login
		response, err := client.LogIn(c, &pb.LogInRequest{
			PhoneNumber: logInUserData.PhoneNumber,
			Password:    logInUserData.Password, 
		})

		if err != nil {
			log.Println("Failed to login:", err)
			utils.ResponseError(ctx, http.StatusUnauthorized, "Invalid credentials")
			return
		}

		// Setting secure cookie with attributes
		ctx.SetSameSite(http.SameSiteLaxMode)
		ctx.SetCookie("Authorization", response.RefreshToken, 3600*24, "/", "", true, true)

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func LogOut() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userId := ctx.GetUint64("user_id")

		// Establishing a gRPC connection
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial gRPC service:", err)
			utils.ResponseError(ctx, http.StatusInternalServerError, "Service unavailable")
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second) 
		defer cancel()

		// Sending a LogOutRequest to the gRPC service for logout
		response, err := client.LogOut(c, &pb.LogOutRequest{
			Id: userId,
		})

		if err != nil {
			log.Println("Failed to logout:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "Logout failed")
			return
			}

		
		// Invalidate the session by clearing the cookie
		ctx.SetCookie("Authorization", "", -1, "/", "", true, true)
		utils.ResponseSuccess(ctx, http.StatusOK, "Logged out successfully")

		utils.ResponseSuccess(ctx, http.StatusOK, response)
	}
}

func ForgotPassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		forgotPasswordUserData := model.ForgotPasswordUserData{}

		// Binding and validating incoming request for password reset
		if err := ctx.ShouldBindJSON(&forgotPasswordUserData); err != nil {
			log.Println("Failed to bind JSON for ForgotPassword:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "Invalid request format")
			return
		}

		// Establishing a gRPC connection
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial gRPC service:", err)
			utils.ResponseError(ctx, http.StatusInternalServerError, "Service unavailable")
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second) 
		defer cancel()

		// Sending a ForgotPasswordRequest to the gRPC service for resetting password
		response, err := client.ForgotPassword(c, &pb.ForgotPasswordRequest{
			Email:       forgotPasswordUserData.Email,
			NewPassword: forgotPasswordUserData.NewPassword, // Ensure this is hashed before sending
		})

		if err != nil {
			log.Println("Failed to reset password:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "Password reset failed")
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

		// Binding and validating incoming request for user update
		if err := ctx.ShouldBindJSON(&userData); err != nil {
			log.Println("Failed to bind JSON for UpdateUser:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "Invalid request format")
			return
		}

		// Establishing a gRPC connection
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial gRPC service:", err)
			utils.ResponseError(ctx, http.StatusInternalServerError, "Service unavailable")
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second) 
		defer cancel()

		// Sending a UpdateUserRequest with user details to the gRPC service for updating user
		response, err := client.UpdateUser(c, &pb.UpdateUserRequest{
			Id:          userId,
			Name:        userData.Name,
			PhoneNumber: userData.PhoneNumber,
			Email:       userData.Email,
		})

		if err != nil {
			log.Println("Failed to update user:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "User update failed")
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
			log.Println("Failed to dial gRPC service:", err)
			utils.ResponseError(ctx, http.StatusInternalServerError, "Service unavailable")
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second) 
		defer cancel()

		// Sending a GetUserRequest to the gRPC service for getting user
		response, err := client.GetUser(c, &pb.GetUserRequest{
			Id: userId,
		})

		if err != nil {
			log.Println("Failed to get user:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "User fetch failed")
			return
			}

		utils.ResponseSuccess(ctx, http.StatusOK, response)
	}
}

func ChangePassword() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieving the user_id from the context, set previously in middleware
		userId := ctx.GetUint64("user_id")

		changePasswordUserData := model.ChangePasswordUserData{}

		// Binding the incoming request to change password
		if err := ctx.ShouldBindJSON(&changePasswordUserData); err != nil {
			log.Println("Failed to bind JSON:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "Invalid request data")
			return
		}

		// Establishing a gRPC connection
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial gRPC server:", err)
			utils.ResponseError(ctx, http.StatusInternalServerError, "Internal server error")
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Sending a ChangePasswordRequest with user details to the gRPC service for changing password
		response, err := client.ChangePassword(c, &pb.ChangePasswordRequest{
			Id:          userId,
			OldPassword: changePasswordUserData.OldPassword,
			NewPassword: changePasswordUserData.NewPassword,
		})

		// Check for error in changing password
		if err != nil {
			log.Println("Failed to change password:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "Failed to change password")
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
			log.Println("Failed to bind JSON:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "Invalid request data")
			return
		}

		// Establishing a gRPC connection
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial gRPC server:", err)
			utils.ResponseError(ctx, http.StatusInternalServerError, "Internal server error")
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		// Sending an UpdateDistanceTravelledRequest to the gRPC service
		response, err := client.UpdateDistanceTravelled(c, &pb.UpdateDistanceTravelledRequest{
			Id:       userId,
			Distance: updateDistanceUserData.Distance,
		})

		// Check for error in updating distance travelled
		if err != nil {
			log.Println("Failed to update distance travelled:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "Failed to update distance travelled")
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func AuthenticateUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authenticateUserData := model.AuthenticateUserData{}

		// Binding the incoming request to verify user
		if err := ctx.ShouldBindJSON(&authenticateUserData); err != nil {
			log.Println("Failed to bind JSON:", err)
			utils.ResponseError(ctx, http.StatusBadRequest, "Invalid request data")
			return
		}

		// Establishing a gRPC connection
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial gRPC server:", err)
			utils.ResponseError(ctx, http.StatusInternalServerError, "Internal server error")
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

		// Check for error in authentication
		if err != nil {
			log.Println("Failed to authenticate user:", err)
			utils.ResponseError(ctx, http.StatusUnauthorized, "Authentication failed")
			return
		}

		utils.ResponseSuccess(ctx, http.StatusOK, response)
	}
}

func RefreshToken() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		refreshToken := strings.TrimPrefix(ctx.GetHeader("Cookie"), "refreshToken=")
		// refreshToken := ctx.Request.Header.Get("Authorization") // Assuming the refresh token is sent in the Authorization header
		if refreshToken == "" {
			utils.ResponseError(ctx, http.StatusBadRequest, "Missing refresh token")
			return
		}

		// Establishing a gRPC connection
		conn, err := utils.GRPCClient(os.Getenv("GRPC_USER_HOST"))
		if err != nil {
			log.Println("Failed to dial gRPC service:", err)
			utils.ResponseError(ctx, http.StatusInternalServerError, "Service unavailable")
			return
		}
		defer conn.Close()

		client := pb.NewUserServiceClient(conn)
		c, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Sending a RefreshTokenRequest to the gRPC service
		response, err := client.RefreshToken(c, &pb.RefreshTokenRequest{
			RefreshToken: refreshToken,
		})

		if err != nil {
			log.Println("Failed to refresh token:", err)
			utils.ResponseError(ctx, http.StatusUnauthorized, "Invalid refresh token")
			return
		}

		// Respond with new access token
		utils.ResponseSuccess(ctx, http.StatusOK, response)
	}
}