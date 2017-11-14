package server

import (
	"net"
	"os"

	pb "github.com/im-auld/alerts/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
    "google.golang.org/grpc/metadata"
)

var logger = grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout)

// AlertService ...
type AlertService struct {
	DB     DB
	logger grpclog.LoggerV2
}

func instrumentRequest(ctx context.Context, logger grpclog.LoggerV2) {
    md, _ := metadata.FromIncomingContext(ctx)
    endpoint := md["endpoint"]
    userAgent := md["user-agent"]
    caller := md["caller"]
    template := "endpoint: %s user-agent: %s response_code: [200] caller: %s"
    logger.Infof(template, endpoint, userAgent, caller)
}

// GetAlertsForUser gets the alerts for the requested user
func (as AlertService) GetAlertsForUser(ctx context.Context, req *pb.GetAlertsForUserRequest) (*pb.GetAlertsForUserResponse, error) {
    defer instrumentRequest(ctx, as.logger)
	var alertsPb []*pb.Alert
	alerts, _ := as.DB.GetAlertsForRecipient(req.UserId)
	for _, a := range alerts {
		alertsPb = append(alertsPb, AlertToProto(a))
	}
	resp := &pb.GetAlertsForUserResponse{Alerts: alertsPb}
	return resp, nil
}

// MarkAlertSeen: Marks an alert as seen.
func (as AlertService) MarkAlertSeen(ctx context.Context, req *pb.MarkAlertSeenRequest) (*pb.MarkAlertSeenResponse, error) {
    defer instrumentRequest(ctx, as.logger)
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
    defer instrumentRequest(ctx, as.logger)
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
