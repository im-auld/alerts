package client

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"

	pb "github.com/im-auld/alerts/proto"
    "fmt"
)

// AlertClient ...
type AlertClient struct {
	client pb.AlertsClient
}

// GetAlertsForUser gets the alerts for the requested user
func (ac AlertClient) GetAlertsForUser(userID int64, ctx context.Context) (*pb.GetAlertsForUserResponse, error) {
	req := &pb.GetAlertsForUserRequest{UserId: userID}
	resp, err := ac.client.GetAlertsForUser(ctx, req)
	return resp, err
}

//MarkAlertSeen ...
func (ac AlertClient) MarkAlertSeen(userId, uniq int64, ctx context.Context) (*pb.MarkAlertSeenResponse, error) {
	req := &pb.MarkAlertSeenRequest{UserId: userId, Uniq: uniq}
	resp, err := ac.client.MarkAlertSeen(ctx, req)
	return resp, err
}

// SendAlert ...
func (ac AlertClient) SendAlert(alert *pb.Alert, ctx context.Context) (*pb.SendAlertResponse, error) {
	req := &pb.SendAlertRequest{Alert: alert}
	resp, err := ac.client.SendAlert(ctx, req)
	return resp, err
}

func NewAlertClient(svcHost, svcPort string) AlertClient {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithInsecure())
	//conn, err := grpc.Dial("192.168.99.100:32312", opts...)
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", svcHost, svcPort), opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
	}
	client := pb.NewAlertsClient(conn)
	return AlertClient{client: client}
}
