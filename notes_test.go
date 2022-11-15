package main

import (
	"context"
	notespb "notes-service/protorepo/noted/notes/v1"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNotesServiceCreateNote(t *testing.T) {
	srv := notesService{}

	res, err := srv.CreateNote(context.TODO(), &notespb.CreateNoteRequest{})
	require.Error(t, err)
	require.Equal(t, status.Code(err), codes.Unimplemented)
	require.Nil(t, res)
}

func TestNotesServiceGetNote(t *testing.T) {
	srv := notesService{}

	res, err := srv.GetNote(context.TODO(), &notespb.GetNoteRequest{})
	require.Error(t, err)
	require.Equal(t, status.Code(err), codes.Unimplemented)
	require.Nil(t, res)
}

func TestNotesServiceUpdateNote(t *testing.T) {
	srv := notesService{}

	res, err := srv.UpdateNote(context.TODO(), &notespb.UpdateNoteRequest{})
	require.Error(t, err)
	require.Equal(t, status.Code(err), codes.Unimplemented)
	require.Nil(t, res)
}

func TestNotesServiceDeleteNote(t *testing.T) {
	srv := notesService{}

	res, err := srv.DeleteNote(context.TODO(), &notespb.DeleteNoteRequest{})
	require.Error(t, err)
	require.Equal(t, status.Code(err), codes.Unimplemented)
	require.Nil(t, res)
}

func TestNotesServiceListNotes(t *testing.T) {
	srv := notesService{}

	res, err := srv.ListNotes(context.TODO(), &notespb.ListNotesRequest{})
	require.Error(t, err)
	require.Equal(t, status.Code(err), codes.Unimplemented)
	require.Nil(t, res)
}
