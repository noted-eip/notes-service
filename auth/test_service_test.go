package auth_test

import (
	"context"
	"notes-service/auth"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTestService(t *testing.T) {
	service := &auth.TestService{}

	token, err := service.ContextWithToken(context.TODO(), &auth.Token{AccountID: "123"})
	require.NoError(t, err)

	info, err := service.TokenFromContext(token)
	require.NoError(t, err)
	require.Equal(t, info.AccountID, "123")
}
