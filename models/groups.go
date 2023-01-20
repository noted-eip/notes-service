package models

import (
	"context"
	"time"
)

type ConversationMessage struct {
	ID              string    `json:"id" bson:"_id"`
	GroupID         string    `json:"groupId" bson:"groupId"`
	ConversationID  string    `json:"conversationId" bson:"conversationId"`
	SenderAccountID string    `json:"senderAccountId" bson:"senderaccountId"`
	Content         string    `json:"content" bson:"content"`
	CreatedAt       time.Time `json:"createdAt" bson:"createdAt"`
	ModifiedAt      time.Time `json:"modifiedAt" bson:"modifiedAt"`
}

type GroupConversation struct {
	ID        string    `json:"id" bson:"_id"`
	Name      string    `json:"name" bson:"name"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type GroupMember struct {
	AccountID string    `json:"accountId" bson:"accountId"`
	IsAdmin   bool      `json:"isAdmin" bson:"isAdmin"`
	JoinedAt  time.Time `json:"joinedAt" bson:"joinedAt"`
}

type GroupInvite struct {
	ID                 string    `json:"id" bson:"_id"`
	SenderAccountID    string    `json:"senderAccountId" bson:"senderAccountId"`
	RecipientAccountID string    `json:"recipientAccountId" bson:"recipientAccountId"`
	CreatedAt          time.Time `json:"createdAt" bson:"createdAt"`
	ValidUntil         time.Time `json:"validUntil" bson:"validUntil"`
}

type GroupInviteLink struct {
	Code                 string    `json:"code" bson:"code"`
	GeneratedByAccountID string    `json:"generatedByAccountId" bson:"generatedByAccountId"`
	CreatedAt            time.Time `json:"createdAt" bson:"createdAt"`
	ValidUntil           time.Time `json:"validUntil" bson:"validUntil"`
}

type Group struct {
	ID                 string               `json:"id" bson:"_id"`
	Name               string               `json:"name" bson:"name"`
	Description        string               `json:"description" bson:"description"`
	AvatarUrl          string               `json:"avatarUrl" bson:"avatarUrl"`
	WorkspaceAccountID *string              `json:"workspaceAccountId" bson:"workspaceAccountId"`
	CreatedAt          time.Time            `json:"createdAt" bson:"createdAt"`
	ModifiedAt         time.Time            `json:"modifiedAt" bson:"modifiedAt"`
	Conversations      *[]GroupConversation `json:"conversations,omitempty" bson:"conversations,omitempty"`
	Members            *[]GroupMember       `json:"members,omitempty" bson:"members,omitempty"`
	Invites            *[]GroupInvite       `json:"invites,omitempty" bson:"invites,omitempty"`
	InviteLinks        *[]GroupInviteLink   `json:"inviteLinks,omitempty" bson:"inviteLinks,omitempty"`
}

func (group *Group) FinConversation(id string) *GroupConversation {
	if group.Conversations == nil {
		return nil
	}
	for i := 0; i < len(*group.Conversations); i++ {
		if (*group.Conversations)[i].ID == id {
			return &(*group.Conversations)[i]
		}
	}
	return nil
}

func (group *Group) FindMember(accountId string) *GroupMember {
	if group.Members == nil {
		return nil
	}
	for i := 0; i < len(*group.Members); i++ {
		if (*group.Members)[i].AccountID == accountId {
			return &(*group.Members)[i]
		}
	}
	return nil
}

func (group *Group) FindInviteByRecipient(recipientAccountId string) *GroupInvite {
	if group.Invites == nil {
		return nil
	}
	for i := 0; i < len(*group.Invites); i++ {
		if (*group.Invites)[i].SenderAccountID == recipientAccountId {
			return &(*group.Invites)[i]
		}
	}
	return nil
}

func (group *Group) FindInviteLinkByCode(code string) *GroupInviteLink {
	if group.InviteLinks == nil {
		return nil
	}

	for i := 0; i < len(*group.InviteLinks); i++ {
		if (*group.InviteLinks)[i].Code == code {
			return &(*group.InviteLinks)[i]
		}
	}
	return nil
}

type CreateGroupPayload struct {
	Name        string
	Description string
	AvatarUrl   string
	// AccountID of the user who has created the group.
	// Will be added to the members list as admin.
	OwnerAccountID string
	// Upon creation a group has a single conversation.
	// This defines its name.
	DefaultConversationName string
}

type CreateWorkspacePayload struct {
	Name        string
	Description string
	AvatarUrl   string
	// AccountID of the user who has created the group.
	// Will be set as `workspaceAccountId`.
	OwnerAccountID string
}

type UpdateGroupPayload struct {
	Name        *string
	Description *string
	AvatarUrl   *string
}

// Optionally filter the search by looking for a
// group that has at least one of these filters.
//
// Example:
// 	filter := &OneGroupFilter{
// 		MemberAccountID: accountID,
// 		InviteeAccountID: accountID,
// 	}
// 	// Will return the group only if it has either accountID
//	// has a member or as an invitee.
// 	repository.GetGroup(ctx, id, filter)
type OneGroupFilter struct {
	MemberAccountID  *string
	InviteeAccountID *string
	InviteLinkCode   *string
}

type ManyGroupsFilter struct {
	// List all groups to which this user belongs.
	AccountID *string
}

type AddMemberPayload struct {
	AccountID string
	IsAdmin   bool
}

type UpdateMemberPayload struct {
	IsAdmin *bool
}

type SendInvitePayload struct {
	SenderAccountID    string
	RecipientAccountID string
	ValidUntil         time.Time
}

type ManyInvitesFilter struct {
	// (Optional) List all invites sent by this user.
	SenderAccountID *string
	// (Optional) List all invites destined to this user.
	RecipientAccountID *string
	// (Optional) List all invites in this group.
	GroupID *string
}

type UpdateGroupConversationPayload struct {
	Name *string
}

type GenerateGroupInviteLinkPayload struct {
	GeneratedByAccountID string
	ValidUntil           time.Time
}

type GroupConversationUUID struct {
	GroupID        string
	ConversationID string
}

type GroupMemberUUID struct {
	GroupID   string
	AccountID string
}

type ConversationMessageUUID struct {
	GroupID        string
	ConversationID string
	MessageID      string
}

type GroupInviteUUID struct {
	GroupID  string
	InviteID string
}

type GroupInviteLinkUUID struct {
	GroupID        string
	InviteLinkCode string
}

// GroupsRepository encapsulates the persistence layer that stores groups.
// Every endpoint is translated into a single (hopefully) optimized query on
// groups.
type GroupsRepository interface {
	// Groups
	CreateGroup(ctx context.Context, group *CreateGroupPayload) (*Group, error)
	CreateWorkspace(ctx context.Context, workspace *CreateWorkspacePayload) (*Group, error)
	GetGroup(ctx context.Context, groupID string, filter *OneGroupFilter) (*Group, error)
	UpdateGroup(ctx context.Context, groupID string, group *UpdateGroupPayload) (*Group, error)
	DeleteGroup(ctx context.Context, groupID string) error
	ListGroups(ctx context.Context, filter *ManyGroupsFilter, opts *ListOptions) ([]*Group, error)

	// Conversations
	GetConversation(ctx context.Context, conversationUUID *GroupConversationUUID) (*GroupConversation, error)
	UpdateConversation(ctx context.Context, conversationUUID *GroupConversationUUID, conversation *UpdateGroupConversationPayload) (*GroupConversation, error)

	// Messages
	SendConversationMessage(ctx context.Context, conversationUUID *GroupConversationUUID) (*ConversationMessage, error)
	GetConversationMessage(ctx context.Context, messageUUID *ConversationMessageUUID) (*ConversationMessage, error)
	UpdateConversationMessage(ctx context.Context, messageUUID *ConversationMessageUUID) (*ConversationMessage, error)
	DeleteConversationMessage(ctx context.Context, messageUUID *ConversationMessageUUID) error
	ListConversationMessages(ctx context.Context, conversationUUID *GroupConversationUUID) (*[]ConversationMessage, error)

	// Members
	AddGroupMember(ctx context.Context, groupID string, member *AddMemberPayload) (*GroupMember, error)
	UpdateGroupMember(ctx context.Context, memberUUID *GroupMemberUUID, member *UpdateMemberPayload) (*GroupMember, error)
	RemoveGroupMember(ctx context.Context, memberUUID *GroupMemberUUID) error

	// Invites
	SendInvite(ctx context.Context, invite *SendInvitePayload) (*GroupInvite, error)
	AcceptInvite(ctx context.Context, inviteUUID *GroupInviteUUID) error
	DenyInvite(ctx context.Context, inviteUUID *GroupInviteUUID) error
	ListInvites(ctx context.Context, groupID string, filter *ManyInvitesFilter) (*GroupInvite, error)
	RevokeGroupInvite(ctx context.Context, inviteUUID *GroupInviteUUID) error

	// Invite Links
	GenerateGroupInviteLink(ctx context.Context, groupID string, inviteLink *GenerateGroupInviteLinkPayload) (*GroupInviteLink, error)
	GetInviteLink(ctx context.Context, inviteLinkUUID *GroupInviteLinkUUID) (*GroupInviteLinkUUID, error)
	RevokeInviteLink(ctx context.Context, inviteLinkUUID *GroupInviteLinkUUID) error
	UseInviteLink(ctx context.Context, inviteLinkUUID *GroupInviteLinkUUID) (*GroupMember, error)
}
