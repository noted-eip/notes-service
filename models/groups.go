package models

import (
	"context"
	"time"
)

type ConversationMessage struct {
	ID              string    `json:"id" bson:"_id"`
	GroupID         string    `json:"groupId" bson:"groupId"`
	ConversationID  string    `json:"conversationId" bson:"conversationId"`
	SenderAccountID string    `json:"senderAccountId" bson:"senderAccountId"`
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
	ID                 string    `json:"id" bson:"id"`
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

func (group *Group) FindConversation(id string) *GroupConversation {
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

func (group *Group) FindMember(accountID string) *GroupMember {
	if group.Members == nil {
		return nil
	}
	for i := 0; i < len(*group.Members); i++ {
		if (*group.Members)[i].AccountID == accountID {
			return &(*group.Members)[i]
		}
	}
	return nil
}

func (group *Group) FindInviteByRecipient(recipientAccountID string) *GroupInvite {
	if group.Invites == nil {
		return nil
	}
	for i := 0; i < len(*group.Invites); i++ {
		if (*group.Invites)[i].SenderAccountID == recipientAccountID {
			return &(*group.Invites)[i]
		}
	}
	return nil
}

func (group *Group) FindInvite(inviteID string) *GroupInvite {
	if group.Invites == nil {
		return nil
	}
	for i := 0; i < len(*group.Invites); i++ {
		if (*group.Invites)[i].ID == inviteID {
			return &(*group.Invites)[i]
		}
	}
	return nil
}

func (group *Group) FindInviteByAccountTuple(recipientAccountID string, senderAccountID string) *GroupInvite {
	if group.Invites == nil {
		return nil
	}
	for i := 0; i < len(*group.Invites); i++ {
		if (*group.Invites)[i].RecipientAccountID == recipientAccountID && (*group.Invites)[i].SenderAccountID == senderAccountID {
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

type OneGroupFilter struct {
	GroupID string
}

type OneConversationFilter struct {
	GroupID        string
	ConversationID string
}

type OneMemberFilter struct {
	GroupID   string
	AccountID string
}

type OneConversationMessageFilter struct {
	GroupID        string
	ConversationID string
	MessageID      string
}

type OneInviteFilter struct {
	GroupID  string
	InviteID string
}

type OneInviteLinkFilter struct {
	GroupID        string
	InviteLinkCode string
}

type ManyGroupsFilter struct {
	// (Optional) List all groups to which this user belongs.
	AccountID string
}

type ManyInvitesFilter struct {
	// (Optional) List all invites sent by this user.
	SenderAccountID *string
	// (Optional) List all invites destined to this user.
	RecipientAccountID *string
	// (Optional) List all invites in this group.
	GroupID *string
}

type CreateGroupPayload struct {
	Name        string
	Description string
	AvatarUrl   string
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

type SendInvitePayload struct {
	RecipientAccountID string
	ValidUntil         time.Time
}

type GenerateGroupInviteLinkPayload struct {
	GeneratedByAccountID string
	ValidUntil           time.Time
}

type AddMemberPayload struct {
	AccountID string
	IsAdmin   bool
}

type UpdateGroupPayload struct {
	Name        string `bson:"name,omitempty"`
	Description string `bson:"description,omitempty"`
	AvatarUrl   string `bson:"avatarUrl,omitempty"`
}

type UpdateMemberPayload struct {
	IsAdmin *bool
}

type UpdateGroupConversationPayload struct {
	Name string
}

type UpdateGroupConversationMessagePayload struct {
	Content string
}

type ListInvitesResult struct {
	GroupInvite
	GroupID string
}

// GroupsRepository encapsulates the persistence layer that stores groups.
// Every endpoint is translated into a single (hopefully) optimized query on
// groups.
type GroupsRepository interface {
	// Groups
	CreateGroup(ctx context.Context, payload *CreateGroupPayload, accountID string) (*Group, error)
	CreateWorkspace(ctx context.Context, payload *CreateWorkspacePayload, accountID string) (*Group, error)
	GetGroup(ctx context.Context, filter *OneGroupFilter, accountID string) (*Group, error)
	GetGroupInternal(ctx context.Context, filter *OneGroupFilter) (*Group, error)
	UpdateGroup(ctx context.Context, filter *OneGroupFilter, payload *UpdateGroupPayload, accountID string) (*Group, error)
	DeleteGroup(ctx context.Context, filter *OneGroupFilter, accountID string) error
	ListGroupsInternal(ctx context.Context, filter *ManyGroupsFilter, opts *ListOptions) ([]*Group, error)

	// Invites
	SendInvite(ctx context.Context, filter *OneGroupFilter, payload *SendInvitePayload, accountID string) (*GroupInvite, error)
	AcceptInvite(ctx context.Context, filter *OneInviteFilter, accountID string) (*GroupMember, error)
	DenyInvite(ctx context.Context, filter *OneInviteFilter, accountID string) error
	GetInvite(ctx context.Context, filter *OneInviteFilter, accountID string) (*GroupInvite, error)
	ListInvites(ctx context.Context, filter *ManyInvitesFilter, accountID string) ([]*ListInvitesResult, error)
	RevokeGroupInvite(ctx context.Context, filter *OneInviteFilter, accountID string) error

	// Conversations
	GetConversation(ctx context.Context, filter *OneConversationFilter, accountID string) (*GroupConversation, error)
	UpdateConversation(ctx context.Context, filter *OneConversationFilter, payload *UpdateGroupConversationPayload, accountID string) (*GroupConversation, error)

	// Messages
	SendConversationMessage(ctx context.Context, filter *OneConversationFilter, accountID string) (*ConversationMessage, error)
	GetConversationMessage(ctx context.Context, filter *OneConversationMessageFilter, accountID string) (*ConversationMessage, error)
	UpdateConversationMessage(ctx context.Context, filter *OneConversationMessageFilter, payload *UpdateGroupConversationMessagePayload, accountID string) (*ConversationMessage, error)
	DeleteConversationMessage(ctx context.Context, filter *OneConversationMessageFilter, accountID string) error
	ListConversationMessages(ctx context.Context, filter *OneConversationFilter, accountID string) ([]*ConversationMessage, error)

	// Members
	UpdateGroupMember(ctx context.Context, filter *OneMemberFilter, payload *UpdateMemberPayload, accountID string) (*GroupMember, error)
	RemoveGroupMember(ctx context.Context, filter *OneMemberFilter, accountID string) error

	// Invite Links
	GenerateGroupInviteLink(ctx context.Context, filter *OneGroupFilter, payload *GenerateGroupInviteLinkPayload, accountID string) (*GroupInviteLink, error)
	GetInviteLink(ctx context.Context, filter *OneInviteLinkFilter, accountID string) (*GroupInviteLink, error)
	RevokeInviteLink(ctx context.Context, filter *OneInviteLinkFilter, accountID string) error
	UseInviteLink(ctx context.Context, filter *OneInviteLinkFilter, accountID string) (*GroupMember, error)
}
