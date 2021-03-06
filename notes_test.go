package main

import (
	"context"
	"notes-service/grpc/notespb"
	"testing"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestNotesServiceCreateNote(t *testing.T) {
	srv := notesService{}

	res, err := srv.CreateNote(context.TODO(), &notespb.Note{})
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
