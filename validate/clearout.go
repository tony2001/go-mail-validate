package validate

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/tony2001/go-mail-validate/config"
	"net/http"
	"time"
)

var httpClient *http.Client

type ClearoutRequest struct {
	Email string `json:"email"`
}

type ClearoutError struct {
	Code    int
	Message string
}

type ClearoutData struct {
	EmailAddress        string `json:"email_address",omitempty`
	SafeToSend          string `json:"safe_to_send",omitempty`
	Status              string `json:"status",omitempty`
	Disposable          string `json:"disposable",omitempty`
	DeliverabilityScore int    `json:"deliverability_score",omitempty`
}

type ClearoutResponse struct {
	Status string        `json:"status"`
	Error  ClearoutError `json:"error",omitempty`
	Data   ClearoutData  `json:"data",omitempty`
}

func ClearoutEnabled() bool {
	if config.GetClearoutToken() != "" {
		return true
	}
	return false
}

func ClearoutInstantCheck(ctx context.Context, emailStr string) (requestSuccess bool, valid bool, err error) {

	if httpClient == nil {
		clearoutTimeoutMsec := config.GetClearoutTimeout()

		httpClient = &http.Client{
			Timeout: time.Millisecond * time.Duration(clearoutTimeoutMsec),
		}
	}

	apiUrl := "https://api.clearout.io/v2/email_verify/instant"

	cReq := ClearoutRequest{
		Email: emailStr,
	}

	jsonBytes, err := json.Marshal(cReq)
	if err != nil {
		return false, false, fmt.Errorf("failed to marshal request data: %s", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, apiUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return false, false, fmt.Errorf("failed to create new HTTP request: %s", err)
	}

	bearerStr := fmt.Sprintf("Bearer: %s", config.GetClearoutToken())
	httpReq.Header.Add("Authorization", bearerStr)
	httpReq.Header.Add("Content-type", "application/json")

	httpResponse, err := httpClient.Do(httpReq)
	if httpResponse.StatusCode != http.StatusOK {
		return false, false, fmt.Errorf("HTTP request to Clearout API failed: %s", httpResponse.Status)
	}

	cResponse := ClearoutResponse{}

	err = json.NewDecoder(httpResponse.Body).Decode(&cResponse)
	if err != nil {
		return false, false, fmt.Errorf("failed to decode JSON data from Clearout API: %s", err)
	}

	if cResponse.Status != "success" {
		return false, false, fmt.Errorf("Clearout API error: code=%d, message=%s", cResponse.Error.Code, cResponse.Error.Message)
	}

	if cResponse.Data.Status == "valid" {
		return true, true, nil
	}

	return true, false, fmt.Errorf("Clearout API response: status = %s, safe_to_send = %s, disposable = %s, deliverability_score = %d", cResponse.Data.Status, cResponse.Data.SafeToSend, cResponse.Data.Disposable, cResponse.Data.DeliverabilityScore)
}
