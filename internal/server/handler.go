package server

import (
	"context"
	"payment-service/internal/models"
	"payment-service/internal/service"

	paymentpb "github.com/pet-booking-system/proto-definitions/payment"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PaymentHandler struct {
	paymentpb.UnimplementedPaymentServiceServer
	service service.PaymentService
}

func NewPaymentHandler(s service.PaymentService) *PaymentHandler {
	return &PaymentHandler{service: s}
}

func (h *PaymentHandler) ProcessPayment(ctx context.Context, req *paymentpb.ProcessPaymentRequest) (*paymentpb.ProcessPaymentResponse, error) {
	bookingID, err := uuid.Parse(req.GetBookingId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid booking_id: %v", err)
	}

	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user_id: %v", err)
	}

	input := models.CreatePaymentInput{
		BookingID: bookingID,
		UserID:    userID,
		Amount:    req.GetAmount(),
	}

	payment, err := h.service.ProcessPayment(ctx, input)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to process payment: %v", err)
	}

	return &paymentpb.ProcessPaymentResponse{
		PaymentId: payment.ID.String(),
		Status:    string(payment.Status),
	}, nil
}

func (h *PaymentHandler) GetPaymentStatus(ctx context.Context, req *paymentpb.GetPaymentStatusRequest) (*paymentpb.GetPaymentStatusResponse, error) {
	paymentID, err := uuid.Parse(req.GetPaymentId())
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid payment_id: %v", err)
	}

	statusStr, err := h.service.GetPaymentStatus(ctx, paymentID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get payment status: %v", err)
	}

	return &paymentpb.GetPaymentStatusResponse{
		Status: string(statusStr),
	}, nil
}
