package model

import "google.golang.org/protobuf/types/known/timestamppb"

type CreateCardData struct {
	CardNumber string                 `json:"card_number" binding:"required"`
	CardHolder string                 `json:"card_holder" binding:"required"`
	ExpiryDate *timestamppb.Timestamp `json:"expiry_date" binding:"required"`
	Cvv        string                 `json:"cvv" binding:"required"`
	IsDefault  bool                   `json:"is_default" binding:"omitempty"`
}

type UpdateCardData struct {
	CardNumber string                 `json:"card_number"`
	CardHolder string                 `json:"card_holder"`
	ExpiryDate *timestamppb.Timestamp `json:"expiry_date"`
	Cvv        string                 `json:"cvv"`
	IsDefault  bool                   `json:"is_default" binding:"omitempty"`
}