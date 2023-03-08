package main

import (
	"context"

	"notes-service/models"
	notesv1 "notes-service/protorepo/noted/notes/v1"
	"notes-service/validators"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (srv *groupsAPI) GetMember(ctx context.Context, req *notesv1.GetMemberRequest) (*notesv1.GetMemberResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateGetMemberRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	group, err := srv.groups.GetGroup(ctx, &models.OneGroupFilter{GroupID: req.GroupId}, token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.GetMemberResponse{Member: modelsMemberToProtobufMember(group.FindMember(req.AccountId))}, nil
}

func (srv *groupsAPI) UpdateMember(ctx context.Context, req *notesv1.UpdateMemberRequest) (*notesv1.UpdateMemberResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateUpdateMemberRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// TODO: We're not making use of the req.UpdateMask because for now there's only
	// one field you can update.

	member, err := srv.groups.UpdateGroupMember(ctx,
		&models.OneMemberFilter{GroupID: req.GroupId, AccountID: req.AccountId},
		&models.UpdateMemberPayload{IsAdmin: &req.Member.IsAdmin},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.UpdateMemberResponse{Member: modelsMemberToProtobufMember(member)}, nil
}

func (srv *groupsAPI) RemoveMember(ctx context.Context, req *notesv1.RemoveMemberRequest) (*notesv1.RemoveMemberResponse, error) {
	token, err := srv.authenticate(ctx)
	if err != nil {
		return nil, err
	}

	err = validators.ValidateRemoveMemberRequest(req)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	// Get every notes in the group
	notes, err := srv.notes.ListAllNotesInternal(ctx, &models.ManyNotesFilter{AuthorAccountID: req.AccountId, GroupID: req.GroupId})
	if err != nil {
		return nil, statusFromModelError(err)
	}
	memberWorkspace, err := srv.groups.GetWorkspaceInternal(ctx, req.AccountId)
	processedUsersCache := make(map[string]struct{})
	if err == nil {
		for _, note := range notes {
			if _, ok := processedUsersCache[note.AuthorAccountID]; ok {
				continue
			}
			processedUsersCache[note.AuthorAccountID] = struct{}{}
			_, err = srv.notes.UpdateNotesInternal(
				ctx,
				&models.ManyNotesFilter{GroupID: req.GroupId, AuthorAccountID: note.AuthorAccountID},
				models.UpdateNoteGroupPayload{GroupID: memberWorkspace.ID},
			)
			if err != nil {
				return nil, statusFromModelError(err) // NOTE: Should we stop everything ?
			}
		}
	} else if err == models.ErrNotFound {
		err = srv.notes.DeleteNotes(ctx, &models.ManyNotesFilter{
			GroupID:         req.GroupId,
			AuthorAccountID: req.AccountId,
		})
		if err != nil && err != models.ErrNotFound {
			return nil, statusFromModelError(err)
		}
	} else {
		return nil, statusFromModelError(err)
	}

	err = srv.groups.RemoveGroupMember(ctx,
		&models.OneMemberFilter{GroupID: req.GroupId, AccountID: req.AccountId},
		token.AccountID)
	if err != nil {
		return nil, statusFromModelError(err)
	}

	return &notesv1.RemoveMemberResponse{}, nil
}

func modelsMemberToProtobufMember(member *models.GroupMember) *notesv1.GroupMember {
	if member == nil {
		return nil
	}
	return &notesv1.GroupMember{
		AccountId: member.AccountID,
		IsAdmin:   member.IsAdmin,
		JoinedAt:  timestamppb.New(member.JoinedAt),
	}
}
