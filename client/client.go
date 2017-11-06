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
	req := &pb.GetAlertsForUserRequest{UserId: userID, VisibleOnly: true, Cached: true}
	resp, _ := ac.client.GetAlertsForUser(context.Background(), req)
	return resp, nil
}

// ArchiveAlert - Archive an alert.
func (ac AlertClient) ArchiveAlert(userID, uniq int64) (*pb.ArchiveAlertResponse, error) {
	req := &pb.ArchiveAlertRequest{UserId: userID, Uniq: uniq}
	resp, _ := ac.client.ArchiveAlert(context.Background(), req)
	return resp, nil
}

// UnarchiveAlert - Archive an alert.
func (ac AlertClient) UnarchiveAlert(userID, uniq int64) (*pb.UnarchiveAlertResponse, error) {
	req := &pb.UnarchiveAlertRequest{UserId: userID, Uniq: uniq}
	resp, _ := ac.client.UnarchiveAlert(context.Background(), req)
	return resp, nil
}

// SendAlert ...
func (ac AlertClient) SendAlert(alert *pb.Alert) (*pb.SendAlertResponse, error) {
	req := &pb.SendAlertRequest{Alert: alert}
	resp, _ := ac.client.SendAlert(context.Background(), req)
	return resp, nil
}

func NewAlertClient() AlertClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	conn, err := grpc.Dial("192.168.99.100:32312", opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewAlertsClient(conn)
	return AlertClient{client: client}
}
