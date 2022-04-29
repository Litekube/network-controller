package utils

import (
	"testing"
)

func TestGetUniqueToken(t *testing.T) {
	GetLogger()
	token := GetUniqueToken()
	logger.Info(token)
}
