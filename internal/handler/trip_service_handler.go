package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/haiyen11231/eco-taxi-api-gateway/internal/grpc/pb"
	"github.com/haiyen11231/eco-taxi-api-gateway/internal/utils"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// which need authorization -> need userId to identify which one belongs to the user -> through authentication
// which one need id -> id is known through param

// get (single), update, delete -> need id
// create, delete, update -> need userId

type SearchTripPreviewData struct {
	Pickup      string `json:"pickup" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

type ConfirmBookingData struct {
	Pickup               string                 `json:"pickup" binding:"required"`
	Destination          string                 `json:"destination" binding:"required"`
	Distance             float64                `json:"distance" binding:"required"`
	Fare                 float64                `json:"fare" binding:"required"`
	EstimatedArrivalTime *timestamppb.Timestamp `json:"estimated_arrival_time" binding:"required"`
	EstimatedWaitingTime int64                  `json:"estimated_waiting_time" binding:"required"`
	NumOfAvailableTaxis  int64                  `json:"num_of_available_taxis" binding:"required"`
	NearestTaxiLongitude float64                `json:"nearest_taxi_longitude" binding:"required"`
	NearestTaxiLatitude  float64                `json:"nearest_taxi_latitude" binding:"required"`
	DefaultPaymentMethod string                 `json:"default_payment_method" binding:"required"`  
}

type UpdateBookingStatusData struct {
	EstimatedArrivalTime *timestamppb.Timestamp `json:"estimated_arrival_time" binding:"required"`
	EstimatedWaitingTime int64                  `json:"estimated_waiting_time" binding:"required"`
	BookingStatus        pb.BookingStatus          `json:"booking_status"`
}

func SearchTripPreview() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		searchTripPreview := SearchTripPreviewData{}

		// Binding the incoming request to search trip preview.
		if err := ctx.ShouldBindJSON(&searchTripPreview); err != nil {
			log.Println("Failed binding json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		// Establishing a gRPC connection.
        conn, err := utils.GRPCClient(os.Getenv("GRPC_Trip_HOST"))
        if err != nil {
            log.Println("Failed to dial", err)
            utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
            return
        }
        defer conn.Close()
    
        client := pb.NewTripServiceClient(conn)
        c, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()
    
        // Sending a SearchTripPreviewRequest to the gRPC service for searching trip preview.
		response, err := client.SearchTripPreview(c, &pb.SearchTripPreviewRequest{
			Pickup: searchTripPreview.Pickup,
			Destination: searchTripPreview.Destination,
		})


		// If searching trip preview fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to search trip preview", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func ConfirmBooking() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		confirmBooking := ConfirmBookingData{}
		userId := ctx.GetUint64("user_id")

		// Binding the incoming request to confirm booking.
		if err := ctx.ShouldBindJSON(&confirmBooking); err != nil {
			log.Println("Failed binding json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		
		// Establishing a gRPC connection.
        conn, err := utils.GRPCClient(os.Getenv("GRPC_Trip_HOST"))
        if err != nil {
            log.Println("Failed to dial", err)
            utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
            return
        }
        defer conn.Close()
    
        client := pb.NewTripServiceClient(conn)
        c, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()
    
        // Sending a ConfirmBookingRequest to the gRPC service for confirming booking.
		response, err := client.ConfirmBooking(c, &pb.ConfirmBookingRequest{
			Pickup: confirmBooking.Pickup,
			Destination: confirmBooking.Destination,
			Distance: confirmBooking.Distance,
			Fare: confirmBooking.Fare,
			EstimatedArrivalTime: confirmBooking.EstimatedArrivalTime,
			EstimatedWaitingTime: confirmBooking.EstimatedWaitingTime,
			NumOfAvailableTaxis: confirmBooking.NumOfAvailableTaxis,
			NearestTaxiLongitude: confirmBooking.NearestTaxiLongitude,
			NearestTaxiLatitude: confirmBooking.NearestTaxiLongitude,
			DefaultPaymentMethod: confirmBooking.DefaultPaymentMethod,
			UserId: userId,
		})


		// If confirming booking fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to confirm booking", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func UpdateBookingStatus() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Params.ByName("id"))
		if err != nil {
			log.Println("Failed to convert params", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		updateBookingStatus := UpdateBookingStatusData{}
		userId := ctx.GetUint64("user_id")

		// Binding the incoming request to update booking status.
		if err := ctx.ShouldBindJSON(&updateBookingStatus); err != nil {
			log.Println("Failed binding json", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		
		// Establishing a gRPC connection.
        conn, err := utils.GRPCClient(os.Getenv("GRPC_Trip_HOST"))
        if err != nil {
            log.Println("Failed to dial", err)
            utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
            return
        }
        defer conn.Close()
    
        client := pb.NewTripServiceClient(conn)
        c, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()
    
        // Sending a UpdateBookingRequest to the gRPC service for updating booking status.
		response, err := client.UpdateBookingStatus(c, &pb.UpdateBookingRequest{
			Id: uint64(id),
			EstimatedArrivalTime: updateBookingStatus.EstimatedArrivalTime,
			EstimatedWaitingTime: updateBookingStatus.EstimatedWaitingTime,
			BookingStatus: updateBookingStatus.BookingStatus,
			UserId: userId,
		})


		// If updating booking status fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to update booking status", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, response)
	}
}

func GetBookingHistory() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Retrieves query parameters
		p := ctx.Query("page")
		l := ctx.Query("limit")

		if p == "" {
			p = "1"
		}

		if l == "" {
			l = "10"
		}

		page, err := strconv.Atoi(p)
		if err != nil {
			log.Println("Failed to convert query page", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		limit, err := strconv.Atoi(l)
		if err != nil {
			log.Println("Failed to convert query limit", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}
		
		userId := ctx.GetUint64("user_id")

        bookingStatusHeader := ctx.GetHeader("booking-status")
        orderAscHeader := ctx.GetHeader("order-asc")

        // Default values
        bookingStatuses := []pb.BookingStatus{pb.BookingStatus_INCOMPLETED, pb.BookingStatus_COMPLETED, pb.BookingStatus_CANCELED} // Default to all booking statuses

        // Parse headers for booking statuses
        if bookingStatusHeader != "" {
            statusArray := strings.Split(bookingStatusHeader, ",")
            bookingStatuses = bookingStatusesFromStrings(statusArray)
        }

        // Set default value for orderAsc
        orderAsc := true // Default to ascending order
        if orderAscHeader != "" {
            orderAsc = orderAscHeader == "true" // Convert to boolean if provided
        }

		
		// Establishing a gRPC connection.
        conn, err := utils.GRPCClient(os.Getenv("GRPC_Trip_HOST"))
        if err != nil {
            log.Println("Failed to dial", err)
            utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
            return
        }
        defer conn.Close()
    
        client := pb.NewTripServiceClient(conn)
        c, cancel := context.WithTimeout(context.Background(), time.Second)
        defer cancel()
    
        // Sending a GetBookingHistoryRequest to the gRPC service for getting booking history.
		response, err := client.GetBookingHistory(c, &pb.GetBookingHistoryRequest{
			Page: uint64(page),
			Limit: uint64(limit),
			UserId: userId,
			BookingStatuses: bookingStatuses,
			OrderAsc: orderAsc,
		})


		// If getting booking history fails, logs the error and returns a 400 Bad Request error. On success, it sends a success response with http.StatusAccepted.
		if err != nil {
			log.Println("Failed to get booking history", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		b, err := protojson.Marshal(response)
		if err != nil {
			log.Println("Failed to marshal response", err)
			utils.ResponseError(ctx, http.StatusBadRequest, err.Error())
			return
		}

		bookings := map[string]any{}
		json.Unmarshal(b, &bookings)
 
		if bookings["result"] == nil {
			bookings["result"] = []interface{}{}
		}

		utils.ResponseSuccess(ctx, http.StatusAccepted, bookings)
	}
}

// Function to convert string to BookingStatus enum
func bookingStatusFromString(statusStr string) pb.BookingStatus {
    switch statusStr {
    case "INCOMPLETED":
        return pb.BookingStatus_INCOMPLETED
    case "COMPLETED":
        return pb.BookingStatus_COMPLETED
    case "CANCELLED":
        return pb.BookingStatus_CANCELED
    default:
        return pb.BookingStatus_COMPLETED
    }
}

// Function to convert slice of strings to slice of BookingStatus enums
func bookingStatusesFromStrings(statuses []string) []pb.BookingStatus {
    var enumStatuses []pb.BookingStatus
    for _, status := range statuses {
        enumStatus := bookingStatusFromString(status)
        enumStatuses = append(enumStatuses, enumStatus)
    }
    return enumStatuses
}