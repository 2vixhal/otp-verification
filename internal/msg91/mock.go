package msg91

import (
    "fmt"
    "time"
)

type MockSender struct{}

func NewMockSender() *MockSender { return &MockSender{} }

func (m *MockSender) SendSMS(phone, message string) (string, error) {
    // just simulate a send and return fake id
    return fmt.Sprintf("mock-sms-%d", time.Now().UnixNano()), nil
}

func (m *MockSender) SendOTP(phone, otp string) (string, error) {
    return fmt.Sprintf("mock-otp-%d", time.Now().UnixNano()), nil
}

func (m *MockSender) SendWhatsApp(phone, templateName string, params map[string]string) (string, error) {
    return fmt.Sprintf("mock-wa-%d", time.Now().UnixNano()), nil
}
