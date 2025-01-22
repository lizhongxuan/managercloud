package model

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMiddleware(t *testing.T) {
	mw := Middleware{
		Name:        "test-redis",
		Type:        "redis",
		Version:     "6.2",
		Status:      "running",
		Host:        "localhost",
		Port:        "6379",
		Credentials: "password",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	assert.Equal(t, "test-redis", mw.Name)
	assert.Equal(t, "redis", mw.Type)
	assert.Equal(t, "6.2", mw.Version)
	assert.Equal(t, "running", mw.Status)
	assert.Equal(t, "localhost", mw.Host)
	assert.Equal(t, "6379", mw.Port)
	assert.Equal(t, "password", mw.Credentials)
} 