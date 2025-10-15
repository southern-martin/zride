package domain

import (
	"context"
	"github.com/google/uuid"
)

// PaymentGateway defines the interface for payment gateway integration
type PaymentGateway interface {
	CreatePayment(ctx context.Context, request *PaymentRequest) (*PaymentGatewayResponse, error)
	GetPaymentStatus(ctx context.Context, externalRef string) (*PaymentStatusResponse, error)
	ProcessRefund(ctx context.Context, request *RefundRequest, originalPayment *Payment) (*RefundResponse, error)
	ValidateCallback(ctx context.Context, callback *PaymentCallback) error
}

// PaymentGatewayResponse represents response from payment gateway
type PaymentGatewayResponse struct {
	ExternalRef string `json:"external_ref"`
	PaymentURL  string `json:"payment_url,omitempty"`
	QRCode      string `json:"qr_code,omitempty"`
	Status      string `json:"status"`
	ExpiresAt   int64  `json:"expires_at,omitempty"`
}

// PaymentStatusResponse represents payment status from gateway
type PaymentStatusResponse struct {
	ExternalRef string `json:"external_ref"`
	Status      string `json:"status"`
	Amount      int64  `json:"amount"`
	Currency    string `json:"currency"`
	ProcessedAt int64  `json:"processed_at,omitempty"`
}

// RefundResponse represents refund response from gateway
type RefundResponse struct {
	RefundID    string `json:"refund_id"`
	Status      string `json:"status"`
	Amount      int64  `json:"amount"`
	ProcessedAt int64  `json:"processed_at,omitempty"`
}

// NotificationService defines the interface for sending notifications
type NotificationService interface {
	SendPaymentConfirmation(ctx context.Context, userID uuid.UUID, payment *Payment) error
	SendPaymentFailure(ctx context.Context, userID uuid.UUID, payment *Payment) error
	SendRefundNotification(ctx context.Context, userID uuid.UUID, payment *Payment, refundAmount int64) error
	SendWalletTopupConfirmation(ctx context.Context, userID uuid.UUID, amount int64) error
	SendLowBalanceWarning(ctx context.Context, userID uuid.UUID, balance int64) error
	SendSuspiciousActivity(ctx context.Context, userID uuid.UUID, activity string) error
}

// ExternalService defines the interface for external service calls
type ExternalService interface {
	GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error)
	GetTripByID(ctx context.Context, tripID uuid.UUID) (*Trip, error)
	UpdateTripPaymentStatus(ctx context.Context, tripID uuid.UUID, status string) error
	ValidateUserPermissions(ctx context.Context, userID uuid.UUID, action string) (bool, error)
}

// FraudDetectionService defines the interface for fraud detection
type FraudDetectionService interface {
	CheckTransaction(ctx context.Context, transaction *TransactionRiskData) (*RiskAssessment, error)
	ReportSuspiciousActivity(ctx context.Context, userID uuid.UUID, activity *SuspiciousActivity) error
	GetUserRiskScore(ctx context.Context, userID uuid.UUID) (float64, error)
}

// User represents a user from external service
type User struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Email        string    `json:"email"`
	IsVerified   bool      `json:"is_verified"`
	IsSuspended  bool      `json:"is_suspended"`
	CreatedAt    int64     `json:"created_at"`
}

// Trip represents a trip from external service
type Trip struct {
	ID              uuid.UUID `json:"id"`
	PassengerID     uuid.UUID `json:"passenger_id"`
	DriverID        uuid.UUID `json:"driver_id"`
	Status          string    `json:"status"`
	EstimatedPrice  int64     `json:"estimated_price"`
	FinalPrice      int64     `json:"final_price"`
	PaymentStatus   string    `json:"payment_status"`
}

// TransactionRiskData represents data for fraud detection
type TransactionRiskData struct {
	UserID        uuid.UUID `json:"user_id"`
	Amount        int64     `json:"amount"`
	PaymentMethod string    `json:"payment_method"`
	IPAddress     string    `json:"ip_address"`
	UserAgent     string    `json:"user_agent"`
	Location      string    `json:"location"`
	TimeOfDay     int       `json:"time_of_day"`
	IsFirstTime   bool      `json:"is_first_time"`
	RecentFailures int      `json:"recent_failures"`
}

// RiskAssessment represents risk assessment result
type RiskAssessment struct {
	RiskScore   float64           `json:"risk_score"`   // 0.0 - 1.0
	RiskLevel   string            `json:"risk_level"`   // low, medium, high
	IsBlocked   bool              `json:"is_blocked"`
	Reasons     []string          `json:"reasons"`
	Recommendations map[string]interface{} `json:"recommendations"`
}

// SuspiciousActivity represents suspicious activity data
type SuspiciousActivity struct {
	UserID      uuid.UUID `json:"user_id"`
	ActivityType string   `json:"activity_type"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"`
	Data        map[string]interface{} `json:"data"`
	Timestamp   int64    `json:"timestamp"`
}

// AuditService defines the interface for audit logging
type AuditService interface {
	LogPaymentEvent(ctx context.Context, event *PaymentAuditEvent) error
	LogWalletEvent(ctx context.Context, event *WalletAuditEvent) error
	LogSecurityEvent(ctx context.Context, event *SecurityAuditEvent) error
}

// PaymentAuditEvent represents payment audit event
type PaymentAuditEvent struct {
	EventType   string    `json:"event_type"`
	PaymentID   uuid.UUID `json:"payment_id"`
	UserID      uuid.UUID `json:"user_id"`
	Amount      int64     `json:"amount"`
	Status      string    `json:"status"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Timestamp   int64     `json:"timestamp"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// WalletAuditEvent represents wallet audit event
type WalletAuditEvent struct {
	EventType     string    `json:"event_type"`
	WalletID      uuid.UUID `json:"wallet_id"`
	UserID        uuid.UUID `json:"user_id"`
	Amount        int64     `json:"amount"`
	BalanceBefore int64     `json:"balance_before"`
	BalanceAfter  int64     `json:"balance_after"`
	Reference     string    `json:"reference"`
	IPAddress     string    `json:"ip_address"`
	Timestamp     int64     `json:"timestamp"`
}

// SecurityAuditEvent represents security audit event
type SecurityAuditEvent struct {
	EventType   string    `json:"event_type"`
	UserID      uuid.UUID `json:"user_id"`
	Action      string    `json:"action"`
	Resource    string    `json:"resource"`
	Result      string    `json:"result"`
	RiskScore   float64   `json:"risk_score"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Timestamp   int64     `json:"timestamp"`
}