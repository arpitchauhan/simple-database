package main

import (
	"context"
	"log"
	"net"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/arpitchauhan/simple-database/database"
)

type server struct {
	pb.UnimplementedDatabaseServer
	db *database
}

const (
	addr = "localhost:50051"
)

var (
	internalErr  = status.Error(codes.Internal, "Internal error")
	databasePath = "database.csv"
)

func (s *server) initialize() {
	s.db.initialize()
}

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	gs := grpc.NewServer()

	d := &database{filepath: databasePath, initialized: false}
	s := &server{db: d}
	s.initialize()

	pb.RegisterDatabaseServer(gs, s)

	log.Printf("server listening at %v", lis.Addr())

	if err := gs.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	log.Printf("Get: received key: %v", in.Key)

	keyValid, errmsg := isKeyValid(in.Key)

	if !keyValid {
		return nil, status.Error(codes.InvalidArgument, errmsg)
	}

	value, errCode := s.db.getKey(in.Key)

	if errCode == KeyNotFound {
		return nil, status.Error(codes.NotFound, "Key was not found")
	}

	if errCode != OK {
		return nil, internalErr
	}

	return &pb.GetReply{Value: value}, nil
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	log.Printf("Set: received key: %v, value: %v", in.Key, in.Value)

	keyValid, errmsg := isKeyValid(in.Key)

	if !keyValid {
		return nil, status.Error(codes.InvalidArgument, errmsg)
	}

	code := s.db.setKey(in.Key, in.Value)

	if code != OK {
		return nil, internalErr
	}

	return &pb.SetReply{}, nil
}

func isKeyValid(key string) (bool, string) {
	if len(strings.TrimSpace(key)) == 0 {
		return false, "Key cannot be empty"
	}

	return true, ""
}
