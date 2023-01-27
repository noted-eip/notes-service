package mongo

import (
	"context"
	"notes-service/models"
	"testing"

	"github.com/jaevor/go-nanoid"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestCreateGroup(t *testing.T) {
	ctx := context.TODO()
	newUUID, err := nanoid.Standard(21)
	require.NoError(t, err)
	logger, err := zap.NewDevelopment()
	require.NoError(t, err)
	db, err := NewDatabase(ctx, "mongodb://localhost:27017", "repository-test", logger)
	require.NoError(t, err)
	db.DB.Collection("groups").Drop(ctx)

	repo := groupsRepository{
		repository: repository{
			logger:  logger,
			coll:    db.DB.Collection("groups"),
			newUUID: newUUID,
		},
	}

	group, err := repo.CreateGroup(ctx, &models.CreateGroupPayload{
		Name:                    "New Group",
		Description:             "My Description",
		AvatarUrl:               "",
		DefaultConversationName: "General",
	}, "456")
	require.NoError(t, err)
	require.NotNil(t, group)

	group, err = repo.CreateGroup(ctx, &models.CreateGroupPayload{
		Name:                    "Mine Group",
		Description:             "My Description",
		AvatarUrl:               "",
		DefaultConversationName: "General",
	}, "123")
	require.NoError(t, err)
	require.NotNil(t, group)

	group, err = repo.GetGroupInternal(ctx, &models.OneGroupFilter{GroupID: group.ID})
	require.NoError(t, err)
	require.NotNil(t, group)

	group, err = repo.UpdateGroup(ctx, &models.OneGroupFilter{GroupID: group.ID}, &models.UpdateGroupPayload{Name: "Modified"}, "123")
	require.NoError(t, err)
	require.NotNil(t, group)

	groups, err := repo.ListGroupsInternal(ctx, nil, nil)
	require.NoError(t, err)
	require.Len(t, groups, 2)

	groups, err = repo.ListGroupsInternal(ctx, &models.ManyGroupsFilter{AccountID: "123"}, nil)
	require.NoError(t, err)
	require.Len(t, groups, 1)
}
