package auth_test

import (
	"context"
	"notes-service/auth"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
)

func TestIncomingToOutgoingContext(t *testing.T) {
	// Given
	inCtx := metadata.NewIncomingContext(context.TODO(), metadata.Pairs(auth.AuthorizationHeaderKey, "token"))

	// When
	outCtx := auth.IncomingToOutgoingContext(inCtx)

	// Then
	md, _ := metadata.FromOutgoingContext(outCtx)
	require.Equal(t, len(md), 1)
	require.Equal(t, md.Get(auth.AuthorizationHeaderKey)[0], "token")
}
