package domain


import (
	"time"
	"github.com/google/uuid"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusProcessing PaymentStatus = "processing"
	PaymentStatusCompleted PaymentStatus = "completed"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
	PaymentStatusRefunded  PaymentStatus = "refunded"
)

// TransactionType represents the type of transaction
type TransactionType string

const (
	TransactionTypePayment    TransactionType = "payment"
	TransactionTypeRefund     TransactionType = "refund"
	TransactionTypeDeposit    TransactionType = "deposit"
	TransactionTypeWithdrawal TransactionType = "withdrawal"
	TransactionTypeFee        TransactionType = "fee"
	TransactionTypeBonus      TransactionType = "bonus"
)

// PaymentMethod represents supported payment methods
type PaymentMethod string

const (
	PaymentMethodZaloPay    PaymentMethod = "zalopay"
	PaymentMethodWallet     PaymentMethod = "wallet"
	PaymentMethodBankCard   PaymentMethod = "bank_card"
	PaymentMethodCash       PaymentMethod = "cash"
)

// Currency represents supported currencies
type Currency string

const (
	CurrencyVND Currency = "VND"
	CurrencyUSD Currency = "USD"
)

// Payment represents a payment transaction
type Payment struct {
	ID                uuid.UUID       `json:"id"`
	UserID            uuid.UUID       `json:"user_id"`
	TripID            *uuid.UUID      `json:"trip_id,omitempty"`
	Amount            int64           `json:"amount"` // Amount in smallest currency unit (VND)
	Currency          Currency        `json:"currency"`
	PaymentMethod     PaymentMethod   `json:"payment_method"`
	Status            PaymentStatus   `json:"status"`
	Description       string          `json:"description"`
	ExternalRef       string          `json:"external_ref,omitempty"` // ZaloPay transaction ID
	InternalRef       string          `json:"internal_ref"` // Our internal reference
	Metadata          map[string]interface{} `json:"metadata,omitempty"`
	FailureReason     string          `json:"failure_reason,omitempty"`
	ProcessedAt       *time.Time      `json:"processed_at,omitempty"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}

// Wallet represents a user's digital wallet
type Wallet struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	Balance      int64     `json:"balance"` // Balance in smallest currency unit (VND)
	Currency     Currency  `json:"currency"`
	IsActive     bool      `json:"is_active"`
	LastActivity time.Time `json:"last_activity"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Transaction represents a wallet transaction
type Transaction struct {
	ID              uuid.UUID       `json:"id"`
	WalletID        uuid.UUID       `json:"wallet_id"`
	PaymentID       *uuid.UUID      `json:"payment_id,omitempty"`
	Type            TransactionType `json:"type"`
	Amount          int64           `json:"amount"` // Positive for credit, negative for debit
	BalanceBefore   int64           `json:"balance_before"`
	BalanceAfter    int64           `json:"balance_after"`
	Description     string          `json:"description"`
	Reference       string          `json:"reference,omitempty"`
	Metadata        map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
}

// PaymentRequest represents a payment request from external service
type PaymentRequest struct {
	UserID        uuid.UUID     `json:"user_id" validate:"required"`
	TripID        *uuid.UUID    `json:"trip_id,omitempty"`
	Amount        int64         `json:"amount" validate:"required,min=1000"` // Minimum 1000 VND
	Currency      Currency      `json:"currency" validate:"required"`
	PaymentMethod PaymentMethod `json:"payment_method" validate:"required"`
	Description   string        `json:"description" validate:"required"`
	ReturnURL     string        `json:"return_url,omitempty"`
	NotifyURL     string        `json:"notify_url,omitempty"`
}

// RefundRequest represents a refund request
type RefundRequest struct {
	PaymentID     uuid.UUID `json:"payment_id" validate:"required"`
	Amount        int64     `json:"amount,omitempty"` // If empty, refund full amount
	Reason        string    `json:"reason" validate:"required"`
	RequestedBy   uuid.UUID `json:"requested_by" validate:"required"`
}

// WalletTopupRequest represents a wallet top-up request
type WalletTopupRequest struct {
	UserID        uuid.UUID     `json:"user_id" validate:"required"`
	Amount        int64         `json:"amount" validate:"required,min=10000"` // Minimum 10,000 VND
	PaymentMethod PaymentMethod `json:"payment_method" validate:"required"`
	ReturnURL     string        `json:"return_url,omitempty"`
}

// WalletTransferRequest represents a wallet-to-wallet transfer request
type WalletTransferRequest struct {
	FromUserID  uuid.UUID `json:"from_user_id" validate:"required"`
	ToUserID    uuid.UUID `json:"to_user_id" validate:"required"`
	Amount      int64     `json:"amount" validate:"required,min=1000"`
	Description string    `json:"description" validate:"required"`
	PIN         string    `json:"pin" validate:"required,len=6"` // 6-digit PIN
}

// PaymentCallback represents callback data from payment providers
type PaymentCallback struct {
	ExternalRef   string                 `json:"external_ref"`
	Status        string                 `json:"status"`
	Amount        int64                  `json:"amount"`
	Currency      string                 `json:"currency"`
	Signature     string                 `json:"signature"`
	Timestamp     int64                  `json:"timestamp"`
	Data          map[string]interface{} `json:"data,omitempty"`
}

// ZaloPayOrder represents a ZaloPay order request
type ZaloPayOrder struct {
	AppID       string `json:"appid"`
	AppTransID  string `json:"apptransid"`
	AppUser     string `json:"appuser"`
	Amount      int64  `json:"amount"`
	AppTime     int64  `json:"apptime"`
	EmbedData   string `json:"embeddata"`
	Item        string `json:"item"`
	Description string `json:"description"`
	BankCode    string `json:"bankcode"`
	MAC         string `json:"mac"`
}

// ZaloPayResponse represents ZaloPay API response
type ZaloPayResponse struct {
	ReturnCode    int    `json:"returncode"`
	ReturnMessage string `json:"returnmessage"`
	SubReturnCode int    `json:"subreturncode"`
	SubReturnMessage string `json:"subreturnmessage"`
	OrderURL      string `json:"orderurl,omitempty"`
	ZpTransToken  string `json:"zptranstoken,omitempty"`
}

// BusinessRules contains payment business rules and limits
type BusinessRules struct {
	MinPaymentAmount     int64 // Minimum payment amount (1,000 VND)
	MaxPaymentAmount     int64 // Maximum payment amount (50,000,000 VND)
	MinWalletTopup       int64 // Minimum wallet top-up (10,000 VND)
	MaxWalletTopup       int64 // Maximum wallet top-up (10,000,000 VND)
	MaxWalletBalance     int64 // Maximum wallet balance (100,000,000 VND)
	MinTransferAmount    int64 // Minimum transfer amount (1,000 VND)
	MaxTransferAmount    int64 // Maximum transfer amount (5,000,000 VND)
	DailyTransferLimit   int64 // Daily transfer limit per user
	MonthlyTransferLimit int64 // Monthly transfer limit per user
}

// DefaultBusinessRules returns default business rules for Vietnamese market
func DefaultBusinessRules() BusinessRules {
	return BusinessRules{
		MinPaymentAmount:     1000,      // 1,000 VND
		MaxPaymentAmount:     50000000,  // 50,000,000 VND
		MinWalletTopup:       10000,     // 10,000 VND
		MaxWalletTopup:       10000000,  // 10,000,000 VND
		MaxWalletBalance:     100000000, // 100,000,000 VND
		MinTransferAmount:    1000,      // 1,000 VND
		MaxTransferAmount:    5000000,   // 5,000,000 VND
		DailyTransferLimit:   10000000,  // 10,000,000 VND per day
		MonthlyTransferLimit: 100000000, // 100,000,000 VND per month
	}
}

// IsValidAmount checks if an amount is valid for payment
func (br BusinessRules) IsValidAmount(amount int64) bool {
	return amount >= br.MinPaymentAmount && amount <= br.MaxPaymentAmount
}

// IsValidTopupAmount checks if an amount is valid for wallet top-up
func (br BusinessRules) IsValidTopupAmount(amount int64) bool {
	return amount >= br.MinWalletTopup && amount <= br.MaxWalletTopup
}

// IsValidTransferAmount checks if an amount is valid for wallet transfer
func (br BusinessRules) IsValidTransferAmount(amount int64) bool {
	return amount >= br.MinTransferAmount && amount <= br.MaxTransferAmount
}

// CanAddToWallet checks if amount can be added to wallet without exceeding limits
func (br BusinessRules) CanAddToWallet(currentBalance, amount int64) bool {
	return (currentBalance + amount) <= br.MaxWalletBalance
}