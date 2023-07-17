package main

import (
	"fmt"

	mailing "github.com/noted-eip/noted/mailing-service"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func SendGroupInviteMailContent(recipientId string, groupName string, validUntil *timestamppb.Timestamp) *mailing.SendEmailsRequest {
	body := fmt.Sprintf(`<span>Bonjour,<br/>Vous avez été invité à rejoindre le groupe %s.
	<br/>Veuillez vous connecter à votre profil pour accepter ou refuser l'invitation.
	<a href="https://noted-eip.vercel.app/profile" style="color: blue">Visitez mon profil (https://noted-eip.vercel.app/profile) </a>
	<br/>Attention, cette invitation est valable seulement 2 semaines</span>`, groupName)

	return &mailing.SendEmailsRequest{
		To:      []string{recipientId},
		Sender:  "noted.organisation@gmail.com",
		Title:   "Invitation à rejoindre un groupe",
		Subject: "Vous avez été invité à rejoindre le groupe " + groupName,
		Body:    body,
	}
}
