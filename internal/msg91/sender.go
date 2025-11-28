package msg91

type Sender interface {
    SendSMS(phone, message string) (string, error)       // returns provider message id
    SendOTP(phone, otp string) (string, error)           // optional: uses SendOTP API
    SendWhatsApp(phone, templateName string, params map[string]string) (string, error)
}
