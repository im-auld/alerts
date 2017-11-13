package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/im-auld/alerts/proto"
)

// AlertClient ...
type AlertClient struct {
	client pb.AlertsClient
}

// GetAlertsForUser gets the alerts for the requested user
func (ac AlertClient) GetAlertsForUser(userID int64) (*pb.GetAlertsForUserResponse, error) {
	req := &pb.GetAlertsForUserRequest{UserId: userID}
	resp, err := ac.client.GetAlertsForUser(context.Background(), req)
	return resp, err
}

// SendAlert ...
func (ac AlertClient) SendAlert(alert *pb.Alert) (*pb.SendAlertResponse, error) {
	req := &pb.SendAlertRequest{Alert: alert}
	resp, err := ac.client.SendAlert(context.Background(), req)
	return resp, err
}

func NewAlertClient() AlertClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	// conn, err := grpc.Dial("192.168.99.100:32312", opts...)
	conn, err := grpc.Dial("127.0.0.1:8081", opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewAlertsClient(conn)
	return AlertClient{client: client}
}
