package service

import (
	"context"
	"math/rand"
	"time"

	"payment-service/internal/models"
	"payment-service/internal/repository"

	"github.com/google/uuid"
)

type PaymentService interface {
	ProcessPayment(ctx context.Context, input models.CreatePaymentInput) (*models.Payment, error)
	GetPaymentStatus(ctx context.Context, bookingID uuid.UUID) (models.PaymentStatus, error)
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo: repo}
}

func (s *paymentService) ProcessPayment(ctx context.Context, input models.CreatePaymentInput) (*models.Payment, error) {
	payment := &models.Payment{
		BookingID: input.BookingID,
		UserID:    input.UserID,
		Amount:    input.Amount,
	}

	if err := s.repo.Create(ctx, payment); err != nil {
		return nil, err
	}

	result := mockPaymentGateway()

	var status models.PaymentStatus
	if result {
		status = models.StatusSuccess
	} else {
		status = models.StatusFailed
	}

	if err := s.repo.UpdateStatus(ctx, payment.ID, status); err != nil {
		return nil, err
	}

	payment.Status = status
	payment.UpdatedAt = time.Now().UTC()

	return payment, nil
}

func (s *paymentService) GetPaymentStatus(ctx context.Context, paymentID uuid.UUID) (models.PaymentStatus, error) {
	payment, err := s.repo.GetByID(ctx, paymentID)
	if err != nil {
		return "", err
	}
	return payment.Status, nil
}

func mockPaymentGateway() bool {
	return rand.Intn(100) < 80
}
