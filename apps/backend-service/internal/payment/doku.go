package payment

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
)

type DokuClient interface {
	CreateCheckoutSession(ctx context.Context, req DokuCheckoutRequest) (*DokuCheckoutResponse, error)
	VerifyWebhookSignature(headers http.Header, rawBody []byte, targetPath string) bool
	IsSimulateMode() bool
}

type dokuClient struct {
	clientID     string
	secretKey    string
	baseURL      string
	simulateMode bool
	httpClient   *http.Client
}

func NewDokuClient(clientID, secretKey string, isProduction, simulateMode bool) DokuClient {
	baseURL := "https://api-sandbox.doku.com"
	if isProduction {
		baseURL = "https://api.doku.com"
	}

	return &dokuClient{
		clientID:     clientID,
		secretKey:    secretKey,
		baseURL:      baseURL,
		simulateMode: simulateMode,
		httpClient:   &http.Client{Timeout: 15 * time.Second},
	}
}

func (c *dokuClient) IsSimulateMode() bool {
	return c.simulateMode
}

// ComputeDigest generates Base64(SHA256(requestBody))
func ComputeDigest(body []byte) string {
	hash := sha256.Sum256(body)
	return base64.StdEncoding.EncodeToString(hash[:])
}

// GenerateSignature generates DOKU HMAC-SHA256 signature
// Formula: Client-Id:[id]\nRequest-Id:[reqId]\nRequest-Timestamp:[time]\nRequest-Target:[target]\nDigest:[digest]
func GenerateSignature(clientID, requestID, timestamp, target, digest, secretKey string) string {
	rawSignature := fmt.Sprintf(
		"Client-Id:%s\nRequest-Id:%s\nRequest-Timestamp:%s\nRequest-Target:%s\nDigest:%s",
		clientID, requestID, timestamp, target, digest,
	)

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(rawSignature))
	b64 := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	return "HMACSHA256=" + b64
}

// VerifySignature validates an incoming webhook or request signature
func VerifySignature(clientID, requestID, timestamp, target, digest, incomingSig, secretKey string) bool {
	expectedSig := GenerateSignature(clientID, requestID, timestamp, target, digest, secretKey)
	return hmac.Equal([]byte(strings.TrimSpace(incomingSig)), []byte(strings.TrimSpace(expectedSig)))
}

func (c *dokuClient) VerifyWebhookSignature(headers http.Header, rawBody []byte, targetPath string) bool {
	if c.simulateMode {
		return true // Always allow in simulation mode
	}

	clientID := headers.Get("Client-Id")
	requestID := headers.Get("Request-Id")
	timestamp := headers.Get("Request-Timestamp")
	signature := headers.Get("Signature")

	if clientID == "" || requestID == "" || timestamp == "" || signature == "" {
		return false
	}

	// Verify client ID matches
	if clientID != c.clientID {
		return false
	}

	digest := ComputeDigest(rawBody)
	return VerifySignature(clientID, requestID, timestamp, targetPath, digest, signature, c.secretKey)
}

func (c *dokuClient) CreateCheckoutSession(ctx context.Context, req DokuCheckoutRequest) (*DokuCheckoutResponse, error) {
	// If in simulate mode, return a mock checkout response pointing to local simulation
	if c.simulateMode {
		resp := &DokuCheckoutResponse{}
		resp.Response.Order.InvoiceNumber = req.Order.InvoiceNumber
		// In demo mode, redirect to order detail page with payment demo query param
		resp.Response.Payment.URL = fmt.Sprintf(
			"%s?doku_checkout=1&invoice=%s&amount=%.0f",
			req.Order.CallbackURL,
			req.Order.InvoiceNumber,
			req.Order.Amount,
		)
		return resp, nil
	}

	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal doku request: %w", err)
	}

	requestID := uuid.New().String()
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05Z")
	target := "/checkout/v1/payment"
	digest := ComputeDigest(bodyBytes)
	signature := GenerateSignature(c.clientID, requestID, timestamp, target, digest, c.secretKey)

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+target, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Client-Id", c.clientID)
	httpReq.Header.Set("Request-Id", requestID)
	httpReq.Header.Set("Request-Timestamp", timestamp)
	httpReq.Header.Set("Signature", signature)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("doku checkout request failed: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read doku response body: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("doku api error (%d): %s", resp.StatusCode, string(respBytes))
	}

	var dokuResp DokuCheckoutResponse
	if err := json.Unmarshal(respBytes, &dokuResp); err != nil {
		return nil, fmt.Errorf("failed to decode doku response: %w", err)
	}

	return &dokuResp, nil
}

