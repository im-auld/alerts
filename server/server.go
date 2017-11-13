package server

import (
	"net"
	"os"

	pb "github.com/im-auld/alerts/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var logger = grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)

// AlertService ...
type AlertService struct {
	DB     DB
	logger grpclog.LoggerV2
}

// GetAlertsForUser gets the alerts for the requested user
func (as AlertService) GetAlertsForUser(ctx context.Context, req *pb.GetAlertsForUserRequest) (*pb.GetAlertsForUserResponse, error) {
	var alertsPb []*pb.Alert
	logger.Infof("Getting alerts for: %d", req.UserId)
	alerts, _ := as.DB.GetAlertsForRecipient(req.UserId)
	as.logger.Infof("Got %d alerts", len(alerts))
	for _, a := range alerts {
		alertsPb = append(alertsPb, AlertToProto(a))
	}
	resp := &pb.GetAlertsForUserResponse{Alerts: alertsPb}
	return resp, nil
}

// MarkAlertSeen: Marks an alert as seen.
func (as AlertService) MarkAlertSeen(ctx context.Context, req *pb.MarkAlertSeenRequest) (*pb.MarkAlertSeenResponse, error) {
	var alertError *pb.AlertError
	err := as.DB.MarkAlertSeen(req.UserId, req.Uniq)
	if err != nil {
		alertError = &pb.AlertError{ErrorCode: pb.AlertErrorCode_SERVER_ERROR, Message: err.Error()}
	}
	resp := &pb.MarkAlertSeenResponse{Error: alertError}
	return resp, err
}

// SendAlert ...
func (as AlertService) SendAlert(ctx context.Context, req *pb.SendAlertRequest) (*pb.SendAlertResponse, error) {
	var alertError *pb.AlertError
	alert := AlertFromProto(req.Alert)
	err := as.DB.SaveAlert(alert)
	if err != nil {
		as.logger.Errorf("error sending alert: %s", err)
		alertError = &pb.AlertError{ErrorCode: pb.AlertErrorCode_SERVER_ERROR, Message: err.Error()}
		resp := &pb.SendAlertResponse{Error: alertError}
		return resp, err
	}
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
	lis, err := net.Listen("tcp", "0.0.0.0:8081")
	logger.Info(lis.Addr())
	if err != nil {
		logger.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterAlertsServer(grpcServer, newServer())
	logger.Info("Starting alerts server...")
	grpcServer.Serve(lis)
}
