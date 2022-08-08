package auth_test

import (
	"context"
	"crypto/ed25519"
	"notes-service/auth"
	"testing"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func Test_service_TokenFromContext(t *testing.T) {
	// Given
	pub, priv := genKeyOrFail(t)
	srv := auth.NewService(pub)
	uid := uuid.New()
	ctx := contextWithTokenOrFail(t, context.TODO(), &auth.Token{
		UserID: uid,
		Role:   auth.RoleAdmin,
	}, priv)

	// When
	token, err := srv.TokenFromContext(ctx)

	// Then
	require.NoError(t, err)
	require.Equal(t, token.UserID, uid, "the token should contain the user data")
}

func genKeyOrFail(t *testing.T) (ed25519.PublicKey, ed25519.PrivateKey) {
	pub, priv, err := ed25519.GenerateKey(nil)
	require.NoError(t, err)
	return pub, priv
}

func contextWithTokenOrFail(t *testing.T, parent context.Context, info *auth.Token, key ed25519.PrivateKey) context.Context {
	ss := signTokenOrFail(t, info, key)
	return metadata.AppendToOutgoingContext(parent, auth.AuthorizationHeaderKey, ss)
}

func signTokenOrFail(t *testing.T, info *auth.Token, key ed25519.PrivateKey) string {
	jwtTok := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, info)
	ss, err := jwtTok.SignedString(key)
	require.NoError(t, err)
	return ss
}
