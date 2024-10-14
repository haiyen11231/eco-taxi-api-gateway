package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/haiyen11231/eco-taxi-api-gateway/internal/grpc/pb"
	"github.com/haiyen11231/eco-taxi-api-gateway/internal/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

type CreateCardData struct {
    CardNumber string                 `json:"card_number" binding:"required"`
	CardHolder string                 `json:"card_holder" binding:"required"`
	ExpiryDate *timestamppb.Timestamp `json:"expiry_date" binding:"required"`
	Cvv        uint64                 `json:"cvv" binding:"required"`
	IsDefault  bool                   `json:"is_default" binding:"required"`
}

type UpdateCardData struct {
	CardNumber string                 `json:"card_number"`
	CardHolder string                 `json:"card_holder"`
	ExpiryDate *timestamppb.Timestamp `json:"expiry_date"`
	Cvv        uint64                 `json:"cvv"`
	IsDefault  bool                   `json:"is_default"`
}
  
func GetCards() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        userId := ctx.GetUint64("user_id")
    
        // Establishes a gRPC connection.
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
    
        // Sends a GetCardsRequest to the gRPC service for getting cards.
        response, err := client.GetCards(c, &pb.GetCardsRequest{
            UserId: userId,
        })

        // If getting cards fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
        if err != nil {
            log.Println("Failed to get cards", err)
            utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
            return
        }
    
		b, err := protojson.Marshal(response)
		if err != nil {
			log.Println("Failed to marshal response", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		cards := map[string]any{}
		json.Unmarshal(b, &cards)

		if cards["result"] == nil {
			cards["result"] = []interface{}{}
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, cards)
    }
}

func CreateCard() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		createCard := CreateCardData{}
		userId := ctx.GetUint64("user_id")

        // Binds the incoming request to create card.
		if err := ctx.ShouldBindJSON(&createCard); err != nil {
			log.Println("Failed binding json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishes a gRPC connection.
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

        // Sends a GetCardsRequest to the gRPC service for getting cards.
		response, err := client.CreateCard(c, &pb.CreateCardRequest{
			CardNumber: createCard.CardNumber,
            CardHolder: createCard.CardHolder,
            ExpiryDate: createCard.ExpiryDate,
            Cvv: createCard.Cvv,
            IsDefault: createCard.IsDefault,
            UserId: userId,
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
		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			log.Println("Failed to convert params", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		updateCard := UpdateCardData{}
		userId := ctx.GetUint64("user_id")

        // Binds the incoming request to update card.
		if err := ctx.ShouldBindJSON(&updateCard); err != nil {
			log.Println("Failed binding json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishes a gRPC connection.
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

        // Sends a UpdateCardRequest to the gRPC service for updating card.
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
		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			log.Println("Failed to convert params", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishes a gRPC connection.
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

        // Sends a DeleteCardRequest to the gRPC service for deleting card.
		response, err := client.DeleteCard(c, &pb.DeleteCardRequest{
			Id:     uint64(id),
			UserId: ctx.GetUint64("user_id"),
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