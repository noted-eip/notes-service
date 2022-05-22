package main

import (
	"context"
	"fmt"
	"notes-service/grpc/notespb"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.mongodb.org/mongo-driver/bson"
)

type notesService struct {
	notespb.UnimplementedNotesServiceServer
}

var _ notespb.NotesServiceServer = &notesService{}

func (srv *notesService) CreateNote(ctx context.Context, in *notespb.Note) (*emptypb.Empty, error) {

	_, err := NotesCollection.InsertOne(
		ctx,
		bson.D{
			{"authorId", in.AuthorId},
			{"title", in.Title},
			{"blocks", in.Blocks},
		})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error : during saving note")
	}

	//fmt.Printf("Inserted %v documents into episode collection!\n", len(noteResult.InsertedID))
	return nil, status.Errorf(codes.OK, "Success : Note well saved")
}

func (srv *notesService) GetNote(ctx context.Context, in *notespb.GetNoteRequest) (*notespb.Note, error) {

	objID, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error : during get note 1")
	}

	var note *notespb.Note

	err = NotesCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&note)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return note, status.Errorf(codes.OK, "Success : Note found")
}

func (srv *notesService) ListNotes(ctx context.Context, in *notespb.ListNotesRequest) (*notespb.Notes, error) {

	/*filterCursor, err := NotesCollection.Find(ctx, bson.M{"author_id": in.AuthorId})
	if err != nil {
		return nil, err
	}
	var episodesFiltered []bson.M
	err = filterCursor.All(ctx, &episodesFiltered)
	fmt.Println(episodesFiltered)
	*/

	filter := bson.M{"authorId": in.AuthorId}
	cur, err := NotesCollection.Find(context.TODO(), filter)

	if err != nil {
		return nil, err
	}
	fmt.Println(cur)

	var results *notespb.Notes = nil
	for cur.Next(context.TODO()) {
		//error here for decoding
		var elem *notespb.Note
		err := cur.Decode(&elem)
		if err != nil {
			return nil, err
		}
		results.Notes = append(results.Notes, elem)
	}
	err = cur.Err()

	if err != nil {
		return nil, err
	}

	fmt.Println(results)
	return nil, status.Errorf(codes.OK, "Note found")
}

func (srv *notesService) UpdateNote(ctx context.Context, in *notespb.UpdateNoteRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (srv *notesService) DeleteNote(ctx context.Context, in *notespb.DeleteNoteRequest) (*emptypb.Empty, error) {

	objID, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error : during get note 1")
	}

	_, err = NotesCollection.DeleteOne(ctx, bson.M{"_id": objID})

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return nil, status.Errorf(codes.OK, "Success : %s note(s) deleted", 1)
}
