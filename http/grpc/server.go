package grpc

import (
	"context"

	"github.com/janmbaco/Saprocate/config"
	"github.com/janmbaco/Saprocate/http/tlsdial"
	"github.com/janmbaco/go-reverseproxy-ssl/servers"
	pb "github.com/janmbaco/saprocate/core/types/protobuf"
	ps "github.com/janmbaco/saprocate/http/grpc/protoservice"
	"google.golang.org/grpc"
)

type server struct {
	ps.UnimplementedGRPCServiceServer
}

func(s *server) ReservePrevHash(ctx context.Context, in *pb.ReservePrevHashRequest) (*pb.ReservePrevHasResponse, error){
	return &pb.ReservePrevHasResponse{}, nil
}

func (s *server) CreateCurrency(ctx context.Context, in *pb.CreateCurrencyRequest) (*pb.CreateCurrencyResponse, error) {
	return &pb.CreateCurrencyResponse{}, nil
}

func (s *server) CreateTransaction(ctx context.Context, in *pb.CreateTransactionRequest) (*pb.CreateTransactionResponse, error) {
	return &pb.CreateTransactionResponse{}, nil
}

func(s *server) GetBalance(ctx context.Context, in *pb.GetBalanceRequest) (*pb.GetBalanceResponse, error){
	return &pb.GetBalanceResponse{}, nil
}

func StartgRPCServer() {
	servers.NewListener(func(serverSetter *servers.ServerSetter) {
		serverSetter.Addr = config.Config.Port
		serverSetter.TLSConfig = tlsdial.GetConfig()

	}).SetProtobuf(func(grpcServer *grpc.Server) {
		ps.RegisterGRPCServiceServer(grpcServer, &server{})
	}).Start()
}
