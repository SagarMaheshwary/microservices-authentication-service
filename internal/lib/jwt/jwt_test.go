package jwt_test

import (
	"context"
	"errors"
	"testing"
	"time"

	jwtlib "github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/sagarmaheshwary/microservices-authentication-service/internal/config"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/constant"
	"github.com/sagarmaheshwary/microservices-authentication-service/internal/lib/jwt"
)

func TestJWTManager_NewAndParseToken(t *testing.T) {
	mockRedis := new(MockRedisClient)
	cfg := &config.JWT{Secret: "test-secret", Expiry: time.Hour}
	manager := jwt.NewJWTManager(cfg, mockRedis)

	token, err := manager.NewToken(123, "alice")
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := manager.ParseToken(token)
	assert.NoError(t, err)
	assert.Equal(t, float64(123), claims["id"])
	assert.Equal(t, "alice", claims["username"])
	assert.Contains(t, claims, "exp")
	assert.Contains(t, claims, "jti")
}

func TestJWTManager_ParseTokenErrors(t *testing.T) {
	mockRedis := new(MockRedisClient)
	cfg := &config.JWT{Secret: "test-secret", Expiry: time.Hour}
	manager := jwt.NewJWTManager(cfg, mockRedis)

	tests := []struct {
		name      string
		token     string
		expectErr bool
	}{
		{
			name:      "invalid token string",
			token:     "not-a-jwt",
			expectErr: true,
		},
		{
			name: "wrong signing method",
			token: func() string {
				tok := jwtlib.New(jwtlib.SigningMethodRS256)
				ss, _ := tok.SignedString([]byte("different-secret"))
				return ss
			}(),
			expectErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := manager.ParseToken(tt.token)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, claims)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, claims)
			}
		})
	}
}

func TestJWTManager_AddToBlacklist(t *testing.T) {
	mockRedis := new(MockRedisClient)
	cfg := &config.JWT{Secret: "test-secret", Expiry: time.Hour}
	manager := jwt.NewJWTManager(cfg, mockRedis)

	jti := "test-jti"
	expiry := time.Now().Add(10 * time.Second).Unix()
	key := constant.RedisTokenBlacklist + ":" + jti

	mockRedis.On("Set", mock.Anything, key, "", mock.Anything).Return(nil)

	err := manager.AddToBlacklist(context.Background(), jti, expiry)
	assert.NoError(t, err)
	mockRedis.AssertExpectations(t)
}

func TestJWTManager_AddToBlacklist_ExpiredToken(t *testing.T) {
	mockRedis := new(MockRedisClient)
	cfg := &config.JWT{Secret: "test-secret", Expiry: time.Hour}
	manager := jwt.NewJWTManager(cfg, mockRedis)

	// Expired token, should not call redis.Set
	jti := "expired-jti"
	expiry := time.Now().Add(-1 * time.Second).Unix()

	err := manager.AddToBlacklist(context.Background(), jti, expiry)
	assert.NoError(t, err)
	mockRedis.AssertNotCalled(t, "Set")
}

func TestJWTManager_IsBlacklisted(t *testing.T) {
	tests := []struct {
		name      string
		redisResp error
		expected  bool
	}{
		{
			name:      "token is blacklisted",
			redisResp: nil,
			expected:  true,
		},
		{
			name:      "token not blacklisted",
			redisResp: errors.New("not found"),
			expected:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRedis := new(MockRedisClient)
			cfg := &config.JWT{Secret: "test-secret", Expiry: time.Hour}
			manager := jwt.NewJWTManager(cfg, mockRedis)

			jti := "some-jti"
			key := constant.RedisTokenBlacklist + ":" + jti

			mockRedis.On("Get", mock.Anything, key).Return("", tt.redisResp)

			result := manager.IsBlacklisted(context.Background(), jti)
			assert.Equal(t, tt.expected, result)
			mockRedis.AssertExpectations(t)
		})
	}
}
