package main

import (
	"context"
	"notes-service/grpc/notespb"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var _ notespb.NotesServiceServer = &notesService{}

func (srv *notesService) AddBlock(ctx context.Context, in *notespb.AddBlockRequest) (*emptypb.Empty, error) {

	// the Id and Index are hard codded because the last PR protorepo is not merged
	//objID, err := primitive.ObjectIDFromHex(in.Note_id)
	var note_id = "62795c015f122130bb469849"

	objID, err := primitive.ObjectIDFromHex(note_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error : during get note 1")
	}

	var note *notespb.Note

	err = NotesCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&note)

	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	if (len(note.Blocks) - 1) >= int(in.Index) {
		if in.Index == 0 {
			note.Blocks = append([]*notespb.Block{in.Block}, note.Blocks...)
		} else {
			fstPartNote := note.Blocks[0 : in.Index-1]
			secPartNote := note.Blocks[in.Index:len(note.Blocks)]
			fstPartNote = append(fstPartNote, in.Block)
			note.Blocks = append(fstPartNote, secPartNote...)
		}
	} else {
		// the index doesn't exist so we add it at the end
		note.Blocks = append(note.Blocks, in.Block)
	}

	id, _ := primitive.ObjectIDFromHex(note_id)
	_, err = NotesCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{
				{"blocks", note.Blocks},
			}},
		},
	)
	if err != nil {
		return &emptypb.Empty{}, nil
	}

	return nil, status.Errorf(codes.OK, "Success : Block well created")
}

func (srv *notesService) UpdateBlock(ctx context.Context, in *notespb.UpdateBlockRequest) (*emptypb.Empty, error) {

	// the Id and Index are hard codded because the last PR protorepo is not merged
	//objID, err := primitive.ObjectIDFromHex(in.Note_id)
	var index = 3
	var note_id = "62795c015f122130bb469849"

	objID, err := primitive.ObjectIDFromHex(note_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error : during get note 1")
	}
	var note *notespb.Note
	err = NotesCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&note)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var exist = true
	var currentIndex = 0
	for currentIndex = 0; currentIndex < len(note.Blocks); currentIndex++ {
		if note.Blocks[currentIndex].Id == in.Id {
			exist = true
			break
		}
	}
	if !exist {
		return nil, status.Errorf(codes.NotFound, "No such block with requested id")
	}

	//modifier le contenu + la position
	if (len(note.Blocks) - 1) >= int(index) {
		if index == 0 {
			//suprimer l'endroit ou y a l'id dans note
			//tout décaler de 1 depuis currentIndex
			fstPartNote := note.Blocks[0:currentIndex]
			secPartNote := note.Blocks[currentIndex+1 : len(note.Blocks)]
			note.Blocks = append(fstPartNote, secPartNote...)
			//ajouter en 1er
			note.Blocks = append([]*notespb.Block{in.Block}, note.Blocks...)
		} else {
			//suprimmer l'ancien endroit ou il etait
			fstPartNote := note.Blocks[0:currentIndex]
			secPartNote := note.Blocks[currentIndex+1 : len(note.Blocks)]
			note.Blocks = append(fstPartNote, secPartNote...)
			//le rajouter a son nouvel index
			fstPartNote = note.Blocks[0:index]
			secPartNote = note.Blocks[index+1 : len(note.Blocks)]
			fstPartNote = append(fstPartNote, in.Block)
			note.Blocks = append(fstPartNote, secPartNote...)
		}
	} else {
		//suprimer l'ednroit ou y a l'id dans note
		//tout décaler de 1 depuis currentIndex
		fstPartNote := note.Blocks[0:currentIndex]
		secPartNote := note.Blocks[currentIndex+1 : len(note.Blocks)]
		note.Blocks = append(fstPartNote, secPartNote...)
		//ajouter en dernier
		note.Blocks = append(note.Blocks, in.Block)
	}

	id, _ := primitive.ObjectIDFromHex(note_id)
	_, err = NotesCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{
				{"blocks", note.Blocks},
			}},
		},
	)
	if err != nil {
		return &emptypb.Empty{}, nil
	}

	return nil, status.Errorf(codes.OK, "Success : Block well updated")
}

func (srv *notesService) DeleteBlock(ctx context.Context, in *notespb.DeleteBlockRequest) (*emptypb.Empty, error) {

	// the Id is hard codded because the last PR protorepo is not merged
	//objID, err := primitive.ObjectIDFromHex(in.Note_id)
	var note_id = "62795c015f122130bb469849"

	objID, err := primitive.ObjectIDFromHex(note_id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Error : during get note 1")
	}
	var note *notespb.Note
	err = NotesCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&note)
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	var exist = true
	var currentIndex = 0
	for currentIndex = 0; currentIndex < len(note.Blocks); currentIndex++ {
		if note.Blocks[currentIndex].Id == in.Id {
			exist = true
			break
		}
	}
	if !exist {
		return nil, status.Errorf(codes.NotFound, "No such block with requested id")
	}

	//suprimer l'ancien endroit ou il etait
	fstPartNote := note.Blocks[0:currentIndex]
	secPartNote := note.Blocks[currentIndex+1 : len(note.Blocks)]
	note.Blocks = append(fstPartNote, secPartNote...)

	id, _ := primitive.ObjectIDFromHex(note_id)
	_, err = NotesCollection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.D{
			{"$set", bson.D{
				{"blocks", note.Blocks},
			}},
		},
	)
	if err != nil {
		return &emptypb.Empty{}, nil
	}

	return nil, status.Errorf(codes.OK, "Success : Block well deleted")
}
