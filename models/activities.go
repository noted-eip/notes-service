package models

import (
	"context"
	"time"
)

type ActivityType string

const (
	NoteAdded     ActivityType = "ADD-NOTE"
	MemberJoined  ActivityType = "ADD-MEMBER"
	MemberRemoved ActivityType = "REMOVE-MEMBER"
)

type ActivityPayload struct {
	GroupID string
	Type    ActivityType
	Event   string
}

type OneActivityFilter struct {
	GroupID    string
	ActivityId string
}

type ManyActivitiesFilter struct {
	GroupID   string
	AccountID string
}

type Activity struct {
	ID        string    `json:"id" bson:"_id"`
	GroupID   string    `json:"groupId" bson:"groupId"`
	Type      string    `json:"type" bson:"type"`
	Event     string    `json:"event" bson:"event"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
}

type ActivitiesRepository interface {
	ListActivitiesInternal(ctx context.Context, filter *ManyActivitiesFilter, lo *ListOptions) ([]*Activity, error)
	GetActivityInternal(ctx context.Context, filter *OneActivityFilter) (*Activity, error)
	CreateActivityInternal(ctx context.Context, payload *ActivityPayload) (*Activity, error)
}
