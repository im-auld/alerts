package server

import (
	"fmt"
	"net"
	"os"

	pb "github.com/im-auld/alerts/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var logger = grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stderr)

// AlertService ...
type AlertService struct {
	DB     DB
	logger grpclog.LoggerV2
}

// GetAlertsForUser gets the alerts for the requested user
func (as AlertService) GetAlertsForUser(ctx context.Context, req *pb.GetAlertsForUserRequest) (*pb.GetAlertsForUserResponse, error) {
	logger.Infof("Getting alerts for: %d", req.UserId)
	archived, unarchived := []*pb.Alert{}, []*pb.Alert{}
	// for _, alert := range as.DB[req.UserId] {
	// 	if alert.Archived {
	// 		archived = append(archived, alert)
	// 	} else {
	// 		unarchived = append(unarchived, alert)
	// 	}
	// }
	as.logger.Infof("Got %d unarchived alerts", len(unarchived))
	resp := &pb.GetAlertsForUserResponse{ArchivedAlerts: archived, UnarchivedAlerts: unarchived}
	return resp, nil
}

// ArchiveAlert - Archive an alert.
func (as AlertService) ArchiveAlert(ctx context.Context, req *pb.ArchiveAlertRequest) (*pb.ArchiveAlertResponse, error) {
	resp := &pb.ArchiveAlertResponse{Error: nil}
	return resp, nil
}

// UnarchiveAlert - Archive an alert.
func (as AlertService) UnarchiveAlert(ctx context.Context, req *pb.UnarchiveAlertRequest) (*pb.UnarchiveAlertResponse, error) {
	resp := &pb.UnarchiveAlertResponse{Error: nil}
	return resp, nil
}

// SendAlert ...
func (as AlertService) SendAlert(ctx context.Context, req *pb.SendAlertRequest) (*pb.SendAlertResponse, error) {
	var alertError *pb.AlertError
	logger.Info(req.Alert)
	alert := AlertFromProto(req.Alert)
	err := as.DB.SaveAlert(alert)
	if err != nil {
		alertError = &pb.AlertError{ErrorCode: pb.AlertErrorCode_SERVER_ERROR, Message: fmt.Sprintf(err.Error())}
	}
	logger.Infof("Set alert with ID: %d", alert.Uniq)
	resp := &pb.SendAlertResponse{Error: alertError}
	return resp, nil
}

func newServer() *AlertService {
	db := NewDB()
	grpclog.SetLoggerV2(logger)
	return &AlertService{DB: db, logger: logger}
}

// StartAlertsServer ...
func StartAlertsServer() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		grpclog.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAlertsServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
