package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/haiyen11231/eco-taxi-api-gateway/internal/grpc/pb"
	"github.com/haiyen11231/eco-taxi-api-gateway/internal/model"
	"github.com/haiyen11231/eco-taxi-api-gateway/internal/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)
  
func GetCards() gin.HandlerFunc {
    return func(ctx *gin.Context) {
		// Retrieving the user_id from the context, set previously in middleware
        userId := ctx.GetUint64("user_id")
    
        // Establishing a gRPC connection
        conn, err := utils.GRPCClient(os.Getenv("GRPC_PAYMENT_HOST"))
        if err != nil {
            log.Println("Failed to dial", err)
            utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
            return
        }
        defer conn.Close()
    
		// Creating a new PaymentService gRPC client from the connection
        client := pb.NewPaymentServiceClient(conn)

		// Setting a timeout for the gRPC request to avoid long-running calls
        c, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()
    
        // Sending a GetCardsRequest to the gRPC service for getting cards
        response, err := client.GetCards(c, &pb.GetCardsRequest{
            UserId: userId,
        })

        // If getting cards fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
        if err != nil {
            log.Println("Failed to get cards", err)
            utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
            return
        }
    
		// Marshalling the gRPC response into JSON format
		b, err := protojson.Marshal(response)
		if err != nil {
			log.Println("Failed to marshal response", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Converting the marshalled JSON data to a generic map for further processing
		cards := map[string]any{}
		json.Unmarshal(b, &cards)

		// If the "result" key in the response is empty, initializes it as an empty array.
		if cards["result"] == nil {
			cards["result"] = []interface{}{}
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, cards)
    }
}

func CreateCard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		createCard := model.CreateCardData{}
		userId := ctx.GetUint64("user_id")

        // Binding the incoming request to create card
		if err := ctx.ShouldBindJSON(&createCard); err != nil {
			log.Println("Failed to bind json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishing a gRPC connection
        conn, err := utils.GRPCClient(os.Getenv("GRPC_PAYMENT_HOST"))
        if err != nil {
            log.Println("Failed to dial", err)
            utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
            return
        }
        defer conn.Close()
    
        client := pb.NewPaymentServiceClient(conn)
        c, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

// 		expiryDate := time.Unix(createCard.ExpiryDate, 0) // Convert Unix timestamp to time.Time
// timestampProto := timestamppb.New(expiryDate) // Convert to timestamppb.Timestamp


        // Sending a GetCardsRequest to the gRPC service for getting cards
		log.Println("Request: ", userId, createCard)
		response, err := client.CreateCard(c, &pb.CreateCardRequest{
			UserId: userId,
			CardNumber: createCard.CardNumber,
            CardHolder: createCard.CardHolder,
            ExpiryDate: createCard.ExpiryDate,
            Cvv: createCard.Cvv,
            IsDefault: createCard.IsDefault,
		})

        // If creating card fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to create card", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func UpdateCard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieving the card ID from the URL parameters and converting it to an integer.
		// id, err := strconv.Atoi(ctx.Params.ByName("id"))
		// if err != nil {
		// 	log.Println("Failed to convert params", err)
		// 	utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
		// 	return
		// }

		idStr := ctx.Query("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            log.Println("Failed to convert query param id", err)
            utils.ResponseError(ctx, http.StatusBadRequest, "Invalid ID")
            return
        }

		updateCard := model.UpdateCardData{}
		userId := ctx.GetUint64("user_id")

        // Binding the incoming request to update card
		if err := ctx.ShouldBindJSON(&updateCard); err != nil {
			log.Println("Failed to bind json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishing a gRPC connection
        conn, err := utils.GRPCClient(os.Getenv("GRPC_PAYMENT_HOST"))
        if err != nil {
            log.Println("Failed to dial", err)
            utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
            return
        }
        defer conn.Close()
    
        client := pb.NewPaymentServiceClient(conn)
        c, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        // Sending a UpdateCardRequest to the gRPC service for updating card
		log.Println("Request: ", userId, id, updateCard)
		response, err := client.UpdatecCard(c, &pb.UpdateCardRequest{
			Id:          uint64(id),
			CardNumber: updateCard.CardNumber,
            CardHolder: updateCard.CardHolder,
            ExpiryDate: updateCard.ExpiryDate,
            Cvv: updateCard.Cvv,
            IsDefault: updateCard.IsDefault,
			UserId:      userId,
		})

        // If updating card fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to update card", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func DeleteCard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieving the card ID from the URL parameters and converting it to an integer.
		idStr := ctx.Query("id")
        id, err := strconv.Atoi(idStr)
        if err != nil {
            log.Println("Failed to convert query param id", err)
            utils.ResponseError(ctx, http.StatusBadRequest, "Invalid ID")
            return
        }

		userId := ctx.GetUint64("user_id")

		// Establishing a gRPC connection.
        conn, err := utils.GRPCClient(os.Getenv("GRPC_PAYMENT_HOST"))
        if err != nil {
            log.Println("Failed to dial", err)
            utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
            return
        }
        defer conn.Close()
    
        client := pb.NewPaymentServiceClient(conn)
        c, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()

        // Sending a DeleteCardRequest to the gRPC service for deleting card
		response, err := client.DeleteCard(c, &pb.DeleteCardRequest{
			Id:     uint64(id),
			UserId: userId,
		})

		// If deleting card fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to delete card", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}