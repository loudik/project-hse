package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sync"
	"time"
)

// EmailAttachment is one file to attach to an outgoing email.
type EmailAttachment struct {
	Filename    string
	ContentType string
	Data        []byte
}

var (
	graphTokenMu     sync.Mutex
	cachedGraphToken string
	graphTokenExpiry time.Time
)

// getGraphAccessToken fetches (and caches) an app-only access token via the
// OAuth2 client credentials flow. Reuses the same Azure AD app registration
// already configured for Microsoft sign-in (MS_CLIENT_ID/MS_TENANT_ID/
// MS_CLIENT_SECRET) - it just needs the Mail.Send Application permission
// (with admin consent) added in Azure Portal.
func getGraphAccessToken() (string, error) {
	graphTokenMu.Lock()
	defer graphTokenMu.Unlock()

	if cachedGraphToken != "" && time.Now().Before(graphTokenExpiry) {
		return cachedGraphToken, nil
	}

	tenantID := os.Getenv("MS_TENANT_ID")
	clientID := os.Getenv("MS_CLIENT_ID")
	clientSecret := os.Getenv("MS_CLIENT_SECRET")
	if tenantID == "" || clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("Microsoft Graph is not configured (MS_TENANT_ID/MS_CLIENT_ID/MS_CLIENT_SECRET)")
	}

	tokenURL := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", clientID)
	form.Set("client_secret", clientSecret)
	form.Set("scope", "https://graph.microsoft.com/.default")

	resp, err := http.PostForm(tokenURL, form)
	if err != nil {
		return "", fmt.Errorf("failed to request Graph token: %w", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Graph token request failed (%d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse Graph token response: %w", err)
	}

	cachedGraphToken = result.AccessToken
	graphTokenExpiry = time.Now().Add(time.Duration(result.ExpiresIn-60) * time.Second)

	return cachedGraphToken, nil
}

func SendEmail(toEmail, subject, htmlBody string) error {
	return SendEmailWithAttachments(toEmail, subject, htmlBody, nil)
}

// SendEmailWithAttachments sends mail via Microsoft Graph's sendMail API,
// from the mailbox configured in MS_SENDER_EMAIL.
func SendEmailWithAttachments(toEmail, subject, htmlBody string, attachments []EmailAttachment) error {
	senderEmail := os.Getenv("MS_SENDER_EMAIL")
	if senderEmail == "" {
		return nil // not configured - caller falls back to another way of surfacing the content
	}

	token, err := getGraphAccessToken()
	if err != nil {
		return fmt.Errorf("failed to get Graph access token: %w", err)
	}

	type graphAttachment struct {
		ODataType    string `json:"@odata.type"`
		Name         string `json:"name"`
		ContentType  string `json:"contentType"`
		ContentBytes string `json:"contentBytes"`
	}

	graphAttachments := []graphAttachment{}
	for _, att := range attachments {
		graphAttachments = append(graphAttachments, graphAttachment{
			ODataType:    "#microsoft.graph.fileAttachment",
			Name:         att.Filename,
			ContentType:  att.ContentType,
			ContentBytes: base64.StdEncoding.EncodeToString(att.Data),
		})
	}

	payload := map[string]interface{}{
		"message": map[string]interface{}{
			"subject": subject,
			"body": map[string]interface{}{
				"contentType": "HTML",
				"content":     htmlBody,
			},
			"toRecipients": []map[string]interface{}{
				{"emailAddress": map[string]string{"address": toEmail}},
			},
			"attachments": graphAttachments,
		},
		"saveToSentItems": true,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to build email payload: %w", err)
	}

	sendURL := fmt.Sprintf("https://graph.microsoft.com/v1.0/users/%s/sendMail", url.PathEscape(senderEmail))
	req, err := http.NewRequest(http.MethodPost, sendURL, bytes.NewReader(payloadBytes))
	if err != nil {
		return fmt.Errorf("failed to build email request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to call Graph sendMail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Graph sendMail failed (%d): %s", resp.StatusCode, string(body))
	}

	return nil
}
