package zalopay


import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/southern-martin/zride/backend/services/payment-service/internal/domain"
	"github.com/southern-martin/zride/backend/shared/errors"
	"github.com/southern-martin/zride/backend/shared/logger"
)

const (
	// ZaloPay Sandbox URLs (use production URLs in production)
	ZaloPayCreateOrderURL = "https://sb-openapi.zalopay.vn/v2/create"
	ZaloPayQueryURL       = "https://sb-openapi.zalopay.vn/v2/query"
	ZaloPayRefundURL      = "https://sb-openapi.zalopay.vn/v2/refund"
)

// ZaloPayConfig holds ZaloPay configuration
type ZaloPayConfig struct {
	AppID     string `json:"app_id"`
	Key1      string `json:"key1"`      // For MAC calculation
	Key2      string `json:"key2"`      // For callback verification
	Endpoint  string `json:"endpoint"`  // ZaloPay API endpoint
	IsTestMode bool  `json:"is_test_mode"`
}

// ZaloPayGateway implements the PaymentGateway interface for ZaloPay
type ZaloPayGateway struct {
	config     ZaloPayConfig
	httpClient *http.Client
	logger     logger.Logger
}

// NewZaloPayGateway creates a new ZaloPay gateway instance
func NewZaloPayGateway(config ZaloPayConfig, logger logger.Logger) domain.PaymentGateway {
	return &ZaloPayGateway{
		config: config,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// CreatePayment creates a payment order in ZaloPay
func (z *ZaloPayGateway) CreatePayment(ctx context.Context, request *domain.PaymentRequest) (*domain.PaymentGatewayResponse, error) {
	// Generate unique transaction ID
	appTransID := z.generateAppTransID()
	
	// Create ZaloPay order
	order := &domain.ZaloPayOrder{
		AppID:       z.config.AppID,
		AppTransID:  appTransID,
		AppUser:     request.UserID.String(),
		Amount:      request.Amount,
		AppTime:     time.Now().UnixMilli(),
		EmbedData:   z.createEmbedData(request),
		Item:        z.createItemDescription(request),
		Description: request.Description,
		BankCode:    "", // Empty for ZaloPay wallet
	}
	
	// Calculate MAC
	order.MAC = z.calculateMAC(order)
	
	// Send request to ZaloPay
	response, err := z.sendCreateOrderRequest(ctx, order)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeExternalError, "failed to create ZaloPay order", err)
	}
	
	if response.ReturnCode != 1 {
		return nil, errors.NewAppError(errors.CodeExternalError, 
			fmt.Sprintf("ZaloPay error: %s", response.ReturnMessage), nil)
	}
	
	z.logger.Info(ctx, "ZaloPay order created successfully", map[string]interface{}{
		"app_trans_id": appTransID,
		"user_id":      request.UserID,
		"amount":       request.Amount,
	})
	
	return &domain.PaymentGatewayResponse{
		ExternalRef: appTransID,
		PaymentURL:  response.OrderURL,
		Status:      "pending",
		ExpiresAt:   time.Now().Add(15 * time.Minute).Unix(), // 15 minutes expiry
	}, nil
}

// GetPaymentStatus queries payment status from ZaloPay
func (z *ZaloPayGateway) GetPaymentStatus(ctx context.Context, externalRef string) (*domain.PaymentStatusResponse, error) {
	// Create query request
	queryData := map[string]interface{}{
		"appid":       z.config.AppID,
		"apptransid":  externalRef,
		"mac":         z.calculateQueryMAC(externalRef),
	}
	
	// Send request
	response, err := z.sendQueryRequest(ctx, queryData)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeExternalError, "failed to query ZaloPay status", err)
	}
	
	status := "pending"
	if response.ReturnCode == 1 {
		status = "completed"
	} else if response.ReturnCode == 2 {
		status = "failed"
	}
	
	return &domain.PaymentStatusResponse{
		ExternalRef: externalRef,
		Status:      status,
		Amount:      0, // ZaloPay doesn't return amount in query response
		Currency:    "VND",
	}, nil
}

// ProcessRefund processes a refund through ZaloPay
func (z *ZaloPayGateway) ProcessRefund(ctx context.Context, request *domain.RefundRequest, originalPayment *domain.Payment) (*domain.RefundResponse, error) {
	refundID := z.generateRefundID()
	
	refundData := map[string]interface{}{
		"appid":       z.config.AppID,
		"mrefundid":   refundID,
		"timestamp":   time.Now().UnixMilli(),
		"amount":      request.Amount,
		"zptransid":   originalPayment.ExternalRef,
		"description": request.Reason,
	}
	
	// Calculate MAC for refund
	refundData["mac"] = z.calculateRefundMAC(refundData)
	
	// Send refund request
	response, err := z.sendRefundRequest(ctx, refundData)
	if err != nil {
		return nil, errors.NewAppError(errors.CodeExternalError, "failed to process ZaloPay refund", err)
	}
	
	status := "pending"
	if response.ReturnCode == 1 {
		status = "completed"
	} else {
		status = "failed"
	}
	
	z.logger.Info(ctx, "ZaloPay refund processed", map[string]interface{}{
		"refund_id":   refundID,
		"payment_id":  request.PaymentID,
		"amount":      request.Amount,
		"status":      status,
	})
	
	return &domain.RefundResponse{
		RefundID:    refundID,
		Status:      status,
		Amount:      request.Amount,
		ProcessedAt: time.Now().Unix(),
	}, nil
}

// ValidateCallback validates ZaloPay callback signature
func (z *ZaloPayGateway) ValidateCallback(ctx context.Context, callback *domain.PaymentCallback) error {
	// Reconstruct callback data for MAC verification
	data := fmt.Sprintf("%s|%s|%d|%s|%d",
		callback.ExternalRef,
		callback.Status,
		callback.Amount,
		callback.Currency,
		callback.Timestamp,
	)
	
	expectedMAC := z.calculateCallbackMAC(data)
	
	if callback.Signature != expectedMAC {
		return errors.NewAppError(errors.CodeValidationError, "invalid callback signature", nil)
	}
	
	return nil
}

// Helper methods

func (z *ZaloPayGateway) generateAppTransID() string {
	// Format: YYMMDD_AppID_Sequence
	now := time.Now()
	date := now.Format("060102")
	sequence := now.UnixNano() % 1000000
	return fmt.Sprintf("%s_%s_%06d", date, z.config.AppID, sequence)
}

func (z *ZaloPayGateway) generateRefundID() string {
	return fmt.Sprintf("%s_%d", z.config.AppID, time.Now().UnixNano())
}

func (z *ZaloPayGateway) createEmbedData(request *domain.PaymentRequest) string {
	embedData := map[string]interface{}{
		"redirecturl": request.ReturnURL,
		"user_id":     request.UserID.String(),
	}
	if request.TripID != nil {
		embedData["trip_id"] = request.TripID.String()
	}
	
	data, _ := json.Marshal(embedData)
	return string(data)
}

func (z *ZaloPayGateway) createItemDescription(request *domain.PaymentRequest) string {
	if request.TripID != nil {
		return fmt.Sprintf(`[{"itemid":"trip_%s","itemname":"Zride Trip Payment","itemprice":%d,"itemquantity":1}]`,
			request.TripID.String(), request.Amount)
	}
	return fmt.Sprintf(`[{"itemid":"topup_%s","itemname":"Zride Wallet Topup","itemprice":%d,"itemquantity":1}]`,
		uuid.New().String(), request.Amount)
}

func (z *ZaloPayGateway) calculateMAC(order *domain.ZaloPayOrder) string {
	data := fmt.Sprintf("%s|%s|%s|%d|%d|%s|%s",
		order.AppID,
		order.AppTransID,
		order.AppUser,
		order.Amount,
		order.AppTime,
		order.EmbedData,
		order.Item,
	)
	return z.hmacSHA256(data, z.config.Key1)
}

func (z *ZaloPayGateway) calculateQueryMAC(appTransID string) string {
	data := fmt.Sprintf("%s|%s|%s", z.config.AppID, appTransID, z.config.Key1)
	return z.hmacSHA256(data, z.config.Key1)
}

func (z *ZaloPayGateway) calculateRefundMAC(refundData map[string]interface{}) string {
	data := fmt.Sprintf("%s|%s|%d|%d|%s|%s",
		refundData["appid"],
		refundData["mrefundid"],
		refundData["timestamp"],
		refundData["amount"],
		refundData["zptransid"],
		refundData["description"],
	)
	return z.hmacSHA256(data, z.config.Key1)
}

func (z *ZaloPayGateway) calculateCallbackMAC(data string) string {
	return z.hmacSHA256(data, z.config.Key2)
}

func (z *ZaloPayGateway) hmacSHA256(data, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (z *ZaloPayGateway) sendCreateOrderRequest(ctx context.Context, order *domain.ZaloPayOrder) (*domain.ZaloPayResponse, error) {
	jsonData, err := json.Marshal(order)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", ZaloPayCreateOrderURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := z.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var response domain.ZaloPayResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	
	return &response, nil
}

func (z *ZaloPayGateway) sendQueryRequest(ctx context.Context, queryData map[string]interface{}) (*domain.ZaloPayResponse, error) {
	jsonData, err := json.Marshal(queryData)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", ZaloPayQueryURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := z.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var response domain.ZaloPayResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	
	return &response, nil
}

func (z *ZaloPayGateway) sendRefundRequest(ctx context.Context, refundData map[string]interface{}) (*domain.ZaloPayResponse, error) {
	jsonData, err := json.Marshal(refundData)
	if err != nil {
		return nil, err
	}
	
	req, err := http.NewRequestWithContext(ctx, "POST", ZaloPayRefundURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	
	resp, err := z.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	var response domain.ZaloPayResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	
	return &response, nil
}