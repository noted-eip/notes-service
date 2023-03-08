package mongo

import (
	"context"
	"notes-service/models"
	"time"

	"github.com/jaevor/go-nanoid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

type groupsRepository struct {
	repository
}

func NewGroupsRepository(db *mongo.Database, logger *zap.Logger) models.GroupsRepository {
	newUUID, err := nanoid.Standard(21)
	if err != nil {
		panic(err)
	}

	return &groupsRepository{
		repository: repository{
			logger:  logger.Named("mongo").Named("groups"),
			coll:    db.Collection("groups"),
			newUUID: newUUID,
		},
	}
}

func (repo *groupsRepository) CreateGroup(ctx context.Context, payload *models.CreateGroupPayload, accountID string) (*models.Group, error) {
	group := &models.Group{
		ID:                 repo.newUUID(),
		Name:               payload.Name,
		Description:        payload.Description,
		AvatarUrl:          payload.AvatarUrl,
		WorkspaceAccountID: nil,
		CreatedAt:          time.Now(),
		ModifiedAt:         time.Now(),
		Conversations: &[]models.GroupConversation{
			{ID: repo.newUUID(), Name: payload.DefaultConversationName, CreatedAt: time.Now()},
		},
		Members: &[]models.GroupMember{
			{AccountID: accountID, IsAdmin: true, JoinedAt: time.Now()},
		},
		Invites:     &[]models.GroupInvite{},
		InviteLinks: &[]models.GroupInviteLink{},
	}

	err := repo.insertOne(ctx, group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (repo *groupsRepository) CreateWorkspace(ctx context.Context, payload *models.CreateWorkspacePayload, accountID string) (*models.Group, error) {
	workspace := &models.Group{
		ID:                 repo.newUUID(),
		Name:               payload.Name,
		Description:        payload.Description,
		AvatarUrl:          payload.AvatarUrl,
		WorkspaceAccountID: &payload.OwnerAccountID,
		CreatedAt:          time.Now(),
		ModifiedAt:         time.Now(),
		Conversations:      nil,
		Members:            nil,
		Invites:            nil,
		InviteLinks:        nil,
	}

	err := repo.insertOne(ctx, workspace)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

func (repo *groupsRepository) GetWorkspaceInternal(ctx context.Context, accountID string) (*models.Group, error) {
	group := &models.Group{}

	query := bson.D{
		{Key: "workspaceAccountId", Value: accountID},
	}

	err := repo.findOne(ctx, query, group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (repo *groupsRepository) GetGroup(ctx context.Context, filter *models.OneGroupFilter, accountID string) (*models.Group, error) {
	group := &models.Group{}

	query := bson.D{
		{Key: "_id", Value: filter.GroupID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "members.accountId", Value: accountID}},
			bson.D{{Key: "workspaceAccountId", Value: accountID}},
		}},
	}

	err := repo.findOne(ctx, query, group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (repo *groupsRepository) GetGroupInternal(ctx context.Context, filter *models.OneGroupFilter) (*models.Group, error) {
	group := &models.Group{}

	query := bson.D{{Key: "_id", Value: filter.GroupID}}

	err := repo.coll.FindOne(ctx, query).Decode(group)
	if err != nil {
		return nil, repo.mongoFindOneErrorToModelsError(query, err)
	}

	return group, nil
}

func (repo *groupsRepository) UpdateGroup(ctx context.Context, filter *models.OneGroupFilter, payload *models.UpdateGroupPayload, accountID string) (*models.Group, error) {
	group := &models.Group{}
	query := bson.D{
		{Key: "_id", Value: filter.GroupID},
		{Key: "members", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "accountId", Value: accountID},
				{Key: "isAdmin", Value: true},
			}},
		}}}
	update := bson.D{{Key: "$set", Value: payload}, {Key: "$set", Value: bson.D{
		{Key: "modifiedAt", Value: time.Now()},
	}}}

	err := repo.findOneAndUpdate(ctx, query, update, group)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (repo *groupsRepository) DeleteGroup(ctx context.Context, filter *models.OneGroupFilter, accountID string) error {
	query := bson.D{
		{Key: "_id", Value: filter.GroupID},
		{Key: "members", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "accountId", Value: accountID},
				{Key: "isAdmin", Value: true},
			}},
		}}}

	return repo.deleteOne(ctx, query)
}

func (repo *groupsRepository) ListGroupsInternal(ctx context.Context, filter *models.ManyGroupsFilter, lo *models.ListOptions) ([]*models.Group, error) {
	groups := make([]*models.Group, 0)

	query := bson.D{}
	if filter != nil && filter.AccountID != "" {
		query = append(query, bson.E{Key: "$or", Value: bson.A{
			bson.D{{Key: "members.accountId", Value: filter.AccountID}},
			bson.D{{Key: "workspaceAccountId", Value: filter.AccountID}},
		}})
	}

	err := repo.find(ctx, query, &groups, lo)
	if err != nil {
		return nil, err
	}

	return groups, nil
}

func (repo *groupsRepository) SendInvite(ctx context.Context, filter *models.OneGroupFilter, payload *models.SendInvitePayload, accountID string) (*models.GroupInvite, error) {
	group := &models.Group{}
	query := bson.D{
		{Key: "_id", Value: filter.GroupID},
		// Sender is member.
		{Key: "members.accountId", Value: accountID},
		// Recipient is not a member.
		{Key: "members.accountId", Value: bson.D{
			{Key: "$ne", Value: payload.RecipientAccountID},
		}},
		// No duplicate invites.
		{Key: "invites", Value: bson.D{
			{Key: "$not", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "recipientAccountId", Value: payload.RecipientAccountID},
					{Key: "senderAccountId", Value: accountID},
				}},
			}},
		}},
	}
	inviteID := repo.newUUID()
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "invites", Value: &models.GroupInvite{
				ID:                 inviteID,
				SenderAccountID:    accountID,
				RecipientAccountID: payload.RecipientAccountID,
				CreatedAt:          time.Now(),
				ValidUntil:         payload.ValidUntil,
			}},
		}}}

	err := repo.findOneAndUpdate(ctx, query, update, group)
	if err != nil {
		return nil, err
	}

	return group.FindInvite(inviteID), nil
}

func (repo *groupsRepository) AcceptInvite(ctx context.Context, filter *models.OneInviteFilter, accountID string) (*models.GroupMember, error) {
	group := &models.Group{}
	query := bson.D{
		{Key: "_id", Value: filter.GroupID},
		{Key: "invites", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "id", Value: filter.InviteID},
				{Key: "recipientAccountId", Value: accountID},
			}},
		}}}
	update := bson.D{
		{Key: "$push", Value: bson.D{
			{Key: "members", Value: &models.GroupMember{
				AccountID: accountID,
				IsAdmin:   false,
				JoinedAt:  time.Now(),
			}},
		}},
		{Key: "$pull", Value: bson.D{
			{Key: "invites", Value: bson.D{
				{Key: "recipientAccountId", Value: accountID},
			}},
		}}}

	err := repo.findOneAndUpdate(ctx, query, update, group)
	if err != nil {
		return nil, err
	}

	return group.FindMember(accountID), nil
}

func (repo *groupsRepository) DenyInvite(ctx context.Context, filter *models.OneInviteFilter, accountID string) error {
	group := &models.Group{}
	query := bson.D{
		{Key: "_id", Value: filter.GroupID},
		{Key: "invites", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "id", Value: filter.InviteID},
				{Key: "recipientAccountId", Value: accountID},
			}},
		}}}

	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "invites", Value: bson.D{
				{Key: "id", Value: filter.InviteID},
			}},
		}}}

	return repo.findOneAndUpdate(ctx, query, update, group)
}

func (repo *groupsRepository) GetInvite(ctx context.Context, filter *models.OneInviteFilter, accountID string) (*models.GroupInvite, error) {
	group := &models.Group{}

	query := bson.D{
		{Key: "_id", Value: filter.GroupID},
		{Key: "$or", Value: bson.A{
			bson.D{
				{Key: "invites", Value: bson.D{
					{Key: "$elemMatch", Value: bson.D{
						{Key: "id", Value: filter.InviteID},
						{Key: "$or", Value: bson.A{
							bson.D{{Key: "recipientAccountId", Value: accountID}},
							bson.D{{Key: "senderAccountId", Value: accountID}}, // idk better be sure u know
						}},
					}},
				}},
			},
			bson.D{
				{Key: "members.accountId", Value: accountID},
			},
		}},
	}

	err := repo.findOne(ctx, query, group)
	if err != nil {
		return nil, err
	}
	if len(*group.Invites) == 0 {
		return nil, models.ErrNotFound
	}

	return group.FindInvite(filter.InviteID), nil
}

func (repo *groupsRepository) ListInvites(ctx context.Context, filter *models.ManyInvitesFilter, lo *models.ListOptions) ([]*models.ListInvitesResult, error) {
	invites := make([]*models.ListInvitesResult, 0)

	mongoDocumentMatch := bson.D{}

	if filter != nil {
		if filter.SenderAccountID != "" {
			mongoDocumentMatch = append(mongoDocumentMatch, bson.E{Key: "invites.senderAccountId", Value: filter.SenderAccountID})
		}
		if filter.RecipientAccountID != "" {
			mongoDocumentMatch = append(mongoDocumentMatch, bson.E{Key: "invites.recipientAccountId", Value: filter.RecipientAccountID})
		}
		if filter.GroupID != "" {
			mongoDocumentMatch = append(mongoDocumentMatch, bson.E{Key: "_id", Value: filter.GroupID})
		}
	}

	idToMongoCondition := func(id *string, varIdentifier string) interface{} {
		if *id != "" {
			return bson.D{{Key: "$eq", Value: []interface{}{varIdentifier, id}}}
		}
		return "true"
	}

	invitesArrayFilterCondition := bson.D{
		{Key: "$and",
			Value: bson.A{
				idToMongoCondition(&filter.SenderAccountID, "$$invite.senderAccountId"),
				idToMongoCondition(&filter.RecipientAccountID, "$$invite.recipientAccountId"),
			},
		},
	}

	// Match only groups that matches the filter (which match the invites that were asked for)
	matchQuery := bson.D{{Key: "$match", Value: mongoDocumentMatch}}

	// In those documents, filter the invites array to have only the invites that matches the filter
	filterQuery := bson.D{{
		Key: "$project", Value: bson.D{
			{Key: "invites", Value: bson.D{
				{Key: "$filter", Value: bson.D{
					{Key: "input", Value: "$invites"},
					{Key: "as", Value: "invite"},
					{Key: "cond", Value: invitesArrayFilterCondition},
				}},
			},
			}}}}

	// Separate every invite in it's specific array element
	unwindQuery := bson.D{{
		Key:   "$unwind",
		Value: "$invites",
	}}

	paginationSkip := bson.D{{ // NOTE: Not putting it in repo.aggregate because it can reduce the work done after, like in this example, we skip and offset before projecting
		Key:   "$skip",
		Value: lo.Offset,
	}}
	paginationLimit := bson.D{{
		Key:   "$limit",
		Value: lo.Limit,
	}}

	// Reorder every variables to have a concise a pertinent element
	projectionQuery := bson.D{{
		Key: "$project",
		Value: bson.D{
			{Key: "recipientAccountId", Value: "$invites.recipientAccountId"},
			{Key: "senderAccountId", Value: "$invites.senderAccountId"},
			{Key: "id", Value: "$invites.id"},
			{Key: "groupId", Value: "$_id"},
			{Key: "_id", Value: 0},
		},
	}}

	err := repo.aggregate(ctx, mongo.Pipeline{matchQuery, filterQuery, unwindQuery, paginationSkip, paginationLimit, projectionQuery}, &invites)
	if err != nil {
		return nil, repo.mongoFindErrorToModelsError(mongoDocumentMatch, lo, err)
	}

	return invites, nil
}

func (repo *groupsRepository) RevokeGroupInvite(ctx context.Context, filter *models.OneInviteFilter, accountID string) error {
	group := &models.Group{}
	query := bson.D{
		{Key: "_id", Value: filter.GroupID},
		{Key: "$or", Value: bson.A{
			bson.D{{Key: "invites", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "id", Value: filter.InviteID},
					{Key: "senderAccountId", Value: accountID},
				}},
			}}},
			bson.D{{Key: "members", Value: bson.D{
				{Key: "$elemMatch", Value: bson.D{
					{Key: "accountId", Value: accountID},
					{Key: "isAdmin", Value: true},
				}},
			}}},
		}},
	}

	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "invites", Value: bson.D{
				{Key: "id", Value: filter.InviteID},
			}},
		}}}

	return repo.findOneAndUpdate(ctx, query, update, group)
}

func (repo *groupsRepository) GetConversation(ctx context.Context, filter *models.OneConversationFilter, accountID string) (*models.GroupConversation, error) {
	return nil, nil
}

func (repo *groupsRepository) UpdateConversation(ctx context.Context, filter *models.OneConversationFilter, payload *models.UpdateGroupConversationPayload, accountID string) (*models.GroupConversation, error) {
	return nil, nil
}

func (repo *groupsRepository) SendConversationMessage(ctx context.Context, filter *models.OneConversationFilter, accountID string) (*models.ConversationMessage, error) {
	return nil, nil
}

func (repo *groupsRepository) GetConversationMessage(ctx context.Context, filter *models.OneConversationMessageFilter, accountID string) (*models.ConversationMessage, error) {
	return nil, nil
}

func (repo *groupsRepository) UpdateConversationMessage(ctx context.Context, filter *models.OneConversationMessageFilter, payload *models.UpdateGroupConversationMessagePayload, accountID string) (*models.ConversationMessage, error) {
	return nil, nil
}

func (repo *groupsRepository) DeleteConversationMessage(ctx context.Context, filter *models.OneConversationMessageFilter, accountID string) error {
	return nil
}

func (repo *groupsRepository) ListConversationMessages(ctx context.Context, filter *models.OneConversationFilter, accountID string) ([]*models.ConversationMessage, error) {
	return nil, nil
}

// TODO: Improve the implementation of this method because it is not going to
// work well if in the future we need to update fields other than `isAdmin`.
func (repo *groupsRepository) UpdateGroupMember(ctx context.Context, filter *models.OneMemberFilter, payload *models.UpdateMemberPayload, accountID string) (*models.GroupMember, error) {
	group := &models.Group{}

	// NOTE: There's something very weird about the ordering of these fields.
	// Something about matching with '$' and array operators.
	query := bson.D{
		{Key: "_id", Value: filter.GroupID},
		// Caller is admin.
		{Key: "members", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "accountId", Value: accountID},
				{Key: "isAdmin", Value: true},
			}},
		}},
		// Target is in group.
		{Key: "members", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "accountId", Value: filter.AccountID},
				{Key: "isAdmin", Value: false},
			}},
		}},
	}

	// Forbidden operations, either results in no-op or demoting.
	if payload == nil || payload.IsAdmin == nil || !*payload.IsAdmin {
		return nil, models.ErrForbidden
	}

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "members.$.isAdmin", Value: true}}}}

	err := repo.findOneAndUpdate(ctx, query, update, group)
	if err != nil {
		return nil, err
	}

	return group.FindMember(filter.AccountID), nil
}

func (repo *groupsRepository) RemoveGroupMember(ctx context.Context, filter *models.OneMemberFilter, accountID string) error {
	group := &models.Group{}
	condition := bson.E{Key: "$and", Value: bson.A{
		// Caller is admin.
		bson.D{{Key: "members", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "accountId", Value: accountID},
				{Key: "isAdmin", Value: true},
			}},
		}}},
		// Target is a regular member.
		bson.D{{Key: "members", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "accountId", Value: filter.AccountID},
				{Key: "isAdmin", Value: false},
			}},
		}}},
	}}

	// Caller is trying to remove themselves from the group.
	if filter.AccountID == accountID {
		condition = bson.E{Key: "members.accountId", Value: accountID}
	}

	query := bson.D{
		{Key: "_id", Value: filter.GroupID},
		condition,
	}
	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "members", Value: bson.D{
				{Key: "accountId", Value: filter.AccountID},
			}},
		}}}

	err := repo.findOneAndUpdate(ctx, query, update, group)
	if err != nil {
		return err
	}

	return nil
}

func (repo *groupsRepository) GenerateGroupInviteLink(ctx context.Context, filter *models.OneGroupFilter, payload *models.GenerateGroupInviteLinkPayload, accountID string) (*models.GroupInviteLink, error) {
	return nil, nil
}

func (repo *groupsRepository) GetInviteLink(ctx context.Context, filter *models.OneInviteLinkFilter, accountID string) (*models.GroupInviteLink, error) {
	return nil, nil
}

func (repo *groupsRepository) RevokeInviteLink(ctx context.Context, filter *models.OneInviteLinkFilter, accountID string) error {
	return nil
}

func (repo *groupsRepository) UseInviteLink(ctx context.Context, filter *models.OneInviteLinkFilter, accountID string) (*models.GroupMember, error) {
	return nil, nil
}

func (repo *groupsRepository) deleteEveryInviteOfAccount(ctx context.Context, accountID string) error {
	query := bson.D{
		{Key: "invites", Value: bson.D{
			{Key: "$elemMatch", Value: bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "recipientAccountId", Value: accountID}},
					bson.D{{Key: "senderAccountId", Value: accountID}},
				}},
			}},
		}},
	}

	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "invites", Value: bson.D{
				{Key: "$or", Value: bson.A{
					bson.D{{Key: "recipientAccountId", Value: accountID}},
					bson.D{{Key: "senderAccountId", Value: accountID}},
				}},
			}},
		}},
	}

	_, err := repo.updateMany(ctx, query, update)
	if err != nil {
		return err
	}

	return nil
}

func (repo *groupsRepository) deleteEveryMemberReferenceOfAccount(ctx context.Context, accountID string) error {
	query := bson.D{
		{Key: "members.accountId", Value: accountID},
	}

	update := bson.D{
		{Key: "$pull", Value: bson.D{
			{Key: "members", Value: bson.D{
				{Key: "accountId", Value: accountID},
			}},
		}},
	}

	_, err := repo.updateMany(ctx, query, update)
	if err != nil {
		return err
	}

	return nil
}

func (repo *groupsRepository) deleteWorkspaces(ctx context.Context, accountID string) error {
	query := bson.D{
		{Key: "workspaceAccountId", Value: accountID},
	}

	return repo.deleteOne(ctx, query)
}

func (repo *groupsRepository) OnAccountDelete(ctx context.Context, accountID string) error {
	err := repo.deleteEveryInviteOfAccount(ctx, accountID)
	if err != nil {
		return err
	}

	err = repo.deleteWorkspaces(ctx, accountID)
	if err != nil {
		return err
	}

	err = repo.deleteEveryMemberReferenceOfAccount(ctx, accountID)
	if err != nil {
		return err
	}

	return nil
}
