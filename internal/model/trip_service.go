package model

import (
	"github.com/haiyen11231/eco-taxi-api-gateway/internal/grpc/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type SearchTripPreviewData struct {
	Pickup      string `json:"pickup" binding:"required"`
	Destination string `json:"destination" binding:"required"`
}

type ConfirmBookingData struct {
	Pickup                   string                 `json:"pickup" binding:"required"`
	Destination              string                 `json:"destination" binding:"required"`
	Distance                 float64                `json:"distance" binding:"required"`
	Fare                     float64                `json:"fare" binding:"required"`
	CardNumber               string                 `json:"card_number" binding:"required"`
	EstimatedArrivalDateTime *timestamppb.Timestamp `json:"estimated_arrival_date_time" binding:"required"`
	EstimatedWaitingTime     int64                  `json:"estimated_waiting_time" binding:"required"`
	BookingStatus            pb.BookingStatus       `json:"booking_status" binding:"required"`
}

type UpdateBookingStatusData struct {
	Pickup                   string                 `json:"pickup"`
	Destination              string                 `json:"destination"`
	Distance                 float64                `json:"distance"`
	Fare                     float64                `json:"fare"`
	CardNumber               string                 `json:"card_number"`
	EstimatedArrivalDateTime *timestamppb.Timestamp `json:"estimated_arrival_date_time"`
	EstimatedWaitingTime     int64                  `json:"estimated_waiting_time" binding:"required"`
	BookingStatus            pb.BookingStatus       `json:"booking_status" binding:"required"`
}