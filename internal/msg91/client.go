package msg91

import (
    "fmt"
    "time"
    "os"
    "github.com/go-resty/resty/v2"
)

type Client struct {
    http *resty.Client
    authKey string
    baseURL string
}

func NewClient() *Client {
    auth := os.Getenv("MSG91_API_KEY")
    base := os.Getenv("MSG91_BASE_URL")
    if base == "" { base = "https://api.msg91.com/api/v5" } // common base
    c := resty.New()
    c.SetHeader("authkey", auth)
    c.SetTimeout(10 * time.Second)
    // 2 retries on network/5xx with exponential backoff
    c.SetRetryCount(2)
    c.AddRetryCondition(func(r *resty.Response, err error) bool {
        if err != nil {
            return true
        }
        status := r.StatusCode()
        return status >= 500 || status == 429
    })
    return &Client{http: c, authKey: auth, baseURL: base}
}

func (c *Client) SendSMS(phone, message string) (string, error) {
    // Use the SMS endpoint — simple example (check docs for exact JSON/params)
    url := fmt.Sprintf("%s/flow", c.baseURL) // NOTE: MSG91 has multiple endpoints; check docs
    payload := map[string]interface{}{
        "mobiles": phone,
        "message": message,
        "sender":  "SENDER", // replace with approved sender ID if needed
    }
    resp, err := c.http.R().
        SetHeader("Content-Type", "application/json").
        SetBody(payload).
        Post(url)
    if err != nil { return "", err }
    if resp.IsError() {
        return "", fmt.Errorf("msg91 sms error: %s", resp.String())
    }
    return resp.String(), nil
}

// Example using SendOTP API (if you prefer provider OTP)
func (c *Client) SendOTP(phone, otp string) (string, error) {
    // SendOTP API usually accepts mobile, authkey, otp, otp_expiry etc.
    url := "https://control.msg91.com/api/v2/otp" // PLACEHOLDER — check docs; many endpoints differ
    payload := map[string]interface{}{
        "mobile": phone,
        "otp": otp,
        "authkey": c.authKey,
        "expiry": 5, // minutes
    }
    resp, err := c.http.R().SetBody(payload).Post(url)
    if err != nil { return "", err }
    if resp.IsError() { return "", fmt.Errorf("msg91 sendotp error: %s", resp.String()) }
    return resp.String(), nil
}

// WhatsApp - template or session messages
func (c *Client) SendWhatsApp(phone, template string, params map[string]string) (string, error) {
    url := fmt.Sprintf("%s/whatsapp/send", c.baseURL) // placeholder path
    body := map[string]interface{}{
        "phone": phone,
        "template": template,
        "params": params,
    }
    resp, err := c.http.R().SetBody(body).Post(url)
    if err != nil { return "", err }
    if resp.IsError() { return "", fmt.Errorf("msg91 wa error: %s", resp.String()) }
    return resp.String(), nil
}
