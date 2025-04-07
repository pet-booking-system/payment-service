package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"payment-service/internal/models"

	"github.com/google/uuid"
)

type PaymentRepository interface {
	Create(ctx context.Context, p *models.Payment) error
	UpdateStatus(ctx context.Context, paymentID uuid.UUID, status models.PaymentStatus) error
	GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error)
}

type paymentRepo struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) PaymentRepository {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) Create(ctx context.Context, p *models.Payment) error {
	query := `
		INSERT INTO payments (
			id, booking_id, user_id, amount, status, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	p.ID = uuid.New()
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now
	p.Status = "pending"

	_, err := r.db.ExecContext(ctx, query,
		p.ID, p.BookingID, p.UserID, p.Amount, p.Status, p.CreatedAt, p.UpdatedAt,
	)
	return err
}

func (r *paymentRepo) UpdateStatus(ctx context.Context, paymentID uuid.UUID, status models.PaymentStatus) error {
	query := `UPDATE payments SET status = $1, updated_at = now() WHERE id = $2`

	res, err := r.db.ExecContext(ctx, query, status, paymentID)
	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no payment found with given ID")
	}

	return nil
}

func (r *paymentRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.Payment, error) {
	query := `SELECT id, booking_id, user_id, amount, status, created_at, updated_at
	          FROM payments WHERE id = $1`

	var p models.Payment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&p.ID, &p.BookingID, &p.UserID, &p.Amount, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &p, nil
}
