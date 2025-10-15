package domain

import (
	"context"
	"time"
	"github.com/google/uuid"
)

// PaymentRepository defines the repository interface for payments
type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByID(ctx context.Context, id uuid.UUID) (*Payment, error)
	GetByInternalRef(ctx context.Context, ref string) (*Payment, error)
	GetByExternalRef(ctx context.Context, ref string) (*Payment, error)
	GetByUserID(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Payment, error)
	GetByTripID(ctx context.Context, tripID uuid.UUID) ([]*Payment, error)
	Update(ctx context.Context, payment *Payment) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status PaymentStatus) error
	GetPendingPayments(ctx context.Context, olderThan time.Duration) ([]*Payment, error)
}

// WalletRepository defines the repository interface for wallets
type WalletRepository interface {
	Create(ctx context.Context, wallet *Wallet) error
	GetByID(ctx context.Context, id uuid.UUID) (*Wallet, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) (*Wallet, error)
	Update(ctx context.Context, wallet *Wallet) error
	UpdateBalance(ctx context.Context, walletID uuid.UUID, newBalance int64) error
	GetActiveWallets(ctx context.Context) ([]*Wallet, error)
}

// TransactionRepository defines the repository interface for transactions
type TransactionRepository interface {
	Create(ctx context.Context, transaction *Transaction) error
	GetByID(ctx context.Context, id uuid.UUID) (*Transaction, error)
	GetByWalletID(ctx context.Context, walletID uuid.UUID, limit, offset int) ([]*Transaction, error)
	GetByPaymentID(ctx context.Context, paymentID uuid.UUID) ([]*Transaction, error)
	GetUserTransactions(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*Transaction, error)
	GetTransactionsByDateRange(ctx context.Context, walletID uuid.UUID, start, end time.Time) ([]*Transaction, error)
	GetDailyTransactionSum(ctx context.Context, walletID uuid.UUID, date time.Time) (int64, error)
	GetMonthlyTransactionSum(ctx context.Context, walletID uuid.UUID, year int, month time.Month) (int64, error)
}

// CacheRepository defines the repository interface for caching
type CacheRepository interface {
	SetPaymentStatus(ctx context.Context, paymentID uuid.UUID, status PaymentStatus, ttl time.Duration) error
	GetPaymentStatus(ctx context.Context, paymentID uuid.UUID) (PaymentStatus, error)
	SetWalletBalance(ctx context.Context, userID uuid.UUID, balance int64, ttl time.Duration) error
	GetWalletBalance(ctx context.Context, userID uuid.UUID) (int64, error)
	LockWallet(ctx context.Context, userID uuid.UUID, ttl time.Duration) (bool, error)
	UnlockWallet(ctx context.Context, userID uuid.UUID) error
	SetTransferLimit(ctx context.Context, userID uuid.UUID, amount int64, period string, ttl time.Duration) error
	GetTransferLimit(ctx context.Context, userID uuid.UUID, period string) (int64, error)
}