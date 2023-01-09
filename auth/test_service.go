package auth

import (
	"context"
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/metadata"
)

type TestService struct {
}

var _ Service = &TestService{}

func (srv *TestService) ContextWithToken(parent context.Context, token *Token) (context.Context, error) {
	ss, err := json.Marshal(token)
	if err != nil {
		return nil, err
	}
	return metadata.AppendToOutgoingContext(parent, AuthorizationHeaderKey, fmt.Sprint(AuthorizationHeaderPrefix, " ", string(ss))), nil
}

func (srv *TestService) TokenFromContext(ctx context.Context) (*Token, error) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		return nil, ErrNoMetadataInCtx
	}

	values := md.Get(AuthorizationHeaderKey)
	tokenString := ""

	for i := range values {
		tokenString, ok = TokenFromAuthorizationHeader(values[i])
		if ok {
			break
		}
	}

	if tokenString == "" {
		return nil, ErrNoTokenInCtx
	}

	token := &Token{}
	json.Unmarshal([]byte(tokenString), token)

	return token, nil
}
