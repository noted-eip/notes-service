package auth_test

import (
	"context"
	"notes-service/auth"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestTestService(t *testing.T) {
	service := &auth.TestService{}

	uid := uuid.Must(uuid.NewRandom())
	token, err := service.ContextWithToken(context.TODO(), &auth.Token{UserID: uid, Role: "admin"})
	require.NoError(t, err)

	info, err := service.TokenFromContext(token)
	require.NoError(t, err)
	require.Equal(t, info.UserID, uid)
	require.Equal(t, info.Role, "admin")
}
