package services

import (
	"context"
	"fmt"
	"github.com/shopspring/decimal"
	api "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/api"
	. "gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/logger"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/model"
	"gitlab.ozon.dev/alex.bogushev/telegram-bot/internal/time_util"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type server struct {
	api.UnimplementedReportServer
	resultCh chan<- *model.Report
}

func (s *server) Send(ctx context.Context, result *api.ReportResult) (*emptypb.Empty, error) {
	Log.Info("get report result")
	start, err := time_util.DateToTime(result.Start)
	if err != nil {
		Log.Error("failed to parse start time", zap.Error(err))
		return nil, err
	}
	end, err := time_util.DateToTime(result.End)
	if err != nil {
		Log.Error("failed to parse end time", zap.Error(err))
		return nil, err
	}
	data := make(map[string]decimal.Decimal, len(result.Data))
	for key, val := range result.Data {
		data[key] = decimal.NewFromFloat(val)
	}

	s.resultCh <- model.NewReport(result.UserId, start, end, data)
	return &emptypb.Empty{}, nil
}

func RunGRPCServer(ctx context.Context, resultCh chan<- *model.Report) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", 50051))
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	api.RegisterReportServer(s, &server{resultCh: resultCh})

	Log.Info(fmt.Sprintf("server listening at %v", lis.Addr()))
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	go func() {
		<-ctx.Done()
		s.Stop()
		Log.Info("stop grpc server")
	}()
	return nil
}
