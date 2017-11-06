package main

import (
	"fmt"

	"github.com/im-auld/alerts/client"
	pb "github.com/im-auld/alerts/proto"
)

func main() {
	c := client.NewAlertClient()
	alert := &pb.Alert{
		RecipientId: 123,
		SenderId:    2,
		Message:     "A new message",
		NoticeType:  pb.AlertTypes_name[pb.AlertTypes_value["NEW_OFFER_MESSAGE"]],
		ActionPath:  "/thread/123",
	}
	fmt.Println(c.SendAlert(alert))
	fmt.Println(c.GetAlertsForUser(123))
	alert2 := &pb.Alert{
		RecipientId: 123,
		SenderId:    4,
		Message:     "Another new message",
		NoticeType:  pb.AlertTypes_name[pb.AlertTypes_value["NEW_OFFER_MESSAGE"]],
		ActionPath:  "/thread/456",
	}
	fmt.Println(c.SendAlert(alert2))
	fmt.Println(c.GetAlertsForUser(123))
}
