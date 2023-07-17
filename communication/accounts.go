package communication

import (
	accountsv1 "notes-service/protorepo/noted/accounts/v1"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AccountsServiceClient struct {
	conn     *grpc.ClientConn
	Accounts accountsv1.AccountsAPIClient
}

func NewAccountsServiceClient(address string) (*AccountsServiceClient, error) {
	res := AccountsServiceClient{}

	err := res.Init(address)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *AccountsServiceClient) Init(address string) error {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	c.conn = conn
	c.Accounts = accountsv1.NewAccountsAPIClient(c.conn)

	return nil
}

func (c *AccountsServiceClient) Close() error {
	return c.conn.Close()
}
