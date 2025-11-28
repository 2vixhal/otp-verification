package main

import (
	"context"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yourusername/otp-service/internal/store"
	"github.com/yourusername/otp-service/internal/msg91"

)

var redisStore *store.RedisStore
var messenger msg91.Sender
type SendOTPReq struct {
	Phone string `json:"phone" binding:"required"`
}

type VerifyOTPReq struct {
	Phone string `json:"phone" binding:"required"`
	OTP   string `json:"otp" binding:"required,len=6"`
}

func genOTP() string {
	n := 100000 + rand.Intn(900000)
	return strconv.Itoa(n)
}

func sendOTPHandler(c *gin.Context) {
	var req SendOTPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: missing or malformed phone"})
		return
	}

	otp := genOTP()

	// Save OTP in Redis with 5 minute TTL
	if err := redisStore.SetOTP(context.Background(), req.Phone, otp, 5*time.Minute); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to store otp"})
		return
	}

	// For now we return the OTP so you can test easily. REMOVE in production.
	c.JSON(http.StatusOK, gin.H{
		"status":          "otp_sent",
		"expires_in":      300,
		"otp_for_testing": otp,
	})
}

func verifyOTPHandler(c *gin.Context) {
	var req VerifyOTPReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	stored, err := redisStore.GetOTP(context.Background(), req.Phone)
	if err != nil {
		// Redis returns an error when key not found or other failures
		c.JSON(http.StatusUnauthorized, gin.H{"error": "otp not found or expired"})
		return
	}

	if stored != req.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "incorrect otp"})
		return
	}

	// one-time use: delete from Redis
	_ = redisStore.DeleteOTP(context.Background(), req.Phone)

	c.JSON(http.StatusOK, gin.H{"status": "verified"})
}

func main() {
	// seed random
	rand.Seed(time.Now().UnixNano())
	messenger = msg91.NewMockSender()
	// create Redis store (expects Redis on localhost:6379)
	redisStore = store.NewRedisStore("localhost:6379")

	// Basic health / readiness check (optional)
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := r.Group("/v1")
	{
		v1.POST("/otp/send", sendOTPHandler)
		v1.POST("/otp/verify", verifyOTPHandler)
	}

	// run server on :8080
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}
