package infra

import (
	"context"
	"time"
)

// EmailService interface defines email operations (if needed for user verification)
type EmailService interface {
	// SendVerificationEmail sends email verification
	SendVerificationEmail(ctx context.Context, email, token string) error

	// SendPasswordResetEmail sends password reset email
	SendPasswordResetEmail(ctx context.Context, email, token string) error

	// SendWelcomeEmail sends welcome email to new users
	SendWelcomeEmail(ctx context.Context, email, name string) error

	// SendNotificationEmail sends general notification
	SendNotificationEmail(ctx context.Context, email, subject, body string) error
}

// TokenService interface defines JWT token operations (if implementing JWT in usecase)
type TokenService interface {
	// GenerateAccessToken generates JWT access token
	GenerateAccessToken(ctx context.Context, userID int, role string) (string, error)

	// GenerateRefreshToken generates JWT refresh token
	GenerateRefreshToken(ctx context.Context, userID int) (string, error)

	// ValidateToken validates and parses JWT token
	ValidateToken(ctx context.Context, token string) (*TokenClaims, error)

	// RefreshToken generates new access token from refresh token
	RefreshToken(ctx context.Context, refreshToken string) (string, error)

	// RevokeToken adds token to blacklist
	RevokeToken(ctx context.Context, token string) error

	// IsTokenRevoked checks if token is blacklisted
	IsTokenRevoked(ctx context.Context, token string) (bool, error)
}

// TokenClaims represents JWT token claims
type TokenClaims struct {
	UserID    int    `json:"user_id"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	TokenType string `json:"token_type"` // "access" or "refresh"
}

// FileStorageService interface for file operations (if needed for user profiles)
type FileStorageService interface {
	// Upload uploads a file and returns URL
	Upload(ctx context.Context, filename string, data []byte) (string, error)

	// Delete deletes a file
	Delete(ctx context.Context, filename string) error

	// GetURL gets signed URL for file access
	GetURL(ctx context.Context, filename string, expiry time.Duration) (string, error)

	// Exists checks if file exists
	Exists(ctx context.Context, filename string) (bool, error)
}

// PasswordService interface for password operations
type PasswordService interface {
	// Hash hashes a password
	Hash(password string) (string, error)

	// Verify verifies password against hash
	Verify(password, hash string) bool

	// GenerateRandomPassword generates secure random password
	GenerateRandomPassword(length int) (string, error)

	// ValidatePasswordStrength checks password strength
	ValidatePasswordStrength(password string) error

	// IsCommonPassword checks if password is in common passwords list
	IsCommonPassword(password string) bool
}

// RateLimiterService interface for rate limiting
type RateLimiterService interface {
	// Allow checks if action is allowed for key
	Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error)

	// Remaining returns remaining attempts
	Remaining(ctx context.Context, key string, limit int, window time.Duration) (int, error)

	// Reset resets rate limit for key
	Reset(ctx context.Context, key string) error

	// Block blocks key for specified duration
	Block(ctx context.Context, key string, duration time.Duration) error

	// IsBlocked checks if key is blocked
	IsBlocked(ctx context.Context, key string) (bool, error)
}

// AuditService interface for audit logging
type AuditService interface {
	// LogUserAction logs user action for audit
	LogUserAction(ctx context.Context, action *AuditAction) error

	// GetUserAuditLog gets audit log for user
	GetUserAuditLog(ctx context.Context, userID int, limit, offset int) ([]*AuditAction, error)

	// GetSystemAuditLog gets system-wide audit log
	GetSystemAuditLog(ctx context.Context, limit, offset int) ([]*AuditAction, error)
}

// AuditAction represents an auditable action
type AuditAction struct {
	ID         string                 `json:"id"`
	UserID     int                    `json:"user_id"`
	Action     string                 `json:"action"`
	Resource   string                 `json:"resource"`
	ResourceID string                 `json:"resource_id,omitempty"`
	Details    map[string]interface{} `json:"details,omitempty"`
	IPAddress  string                 `json:"ip_address,omitempty"`
	UserAgent  string                 `json:"user_agent,omitempty"`
	Success    bool                   `json:"success"`
	Error      string                 `json:"error,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	SessionID  string                 `json:"session_id,omitempty"`
}

// MetricsService interface for application metrics
type MetricsService interface {
	// Counter increments a counter metric
	Counter(name string, value float64, tags map[string]string)

	// Gauge sets a gauge metric
	Gauge(name string, value float64, tags map[string]string)

	// Histogram records a histogram value
	Histogram(name string, value float64, tags map[string]string)

	// Timing records timing information
	Timing(name string, duration time.Duration, tags map[string]string)
}

// NotificationService interface for push notifications
type NotificationService interface {
	// SendPushNotification sends push notification to user
	SendPushNotification(ctx context.Context, userID int, title, body string, data map[string]interface{}) error

	// SendBulkNotification sends notification to multiple users
	SendBulkNotification(ctx context.Context, userIDs []int, title, body string, data map[string]interface{}) error

	// RegisterDevice registers device for push notifications
	RegisterDevice(ctx context.Context, userID int, deviceToken, platform string) error

	// UnregisterDevice removes device from notifications
	UnregisterDevice(ctx context.Context, deviceToken string) error
}
