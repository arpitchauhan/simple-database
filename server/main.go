package main

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"net"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/arpitchauhan/simple-database/database"
)

type server struct {
	pb.UnimplementedDatabaseServer
	databasePath string
}

const (
	addr = "localhost:50051"
)

var (
	internalErr  = status.Error(codes.Internal, "Internal error")
	databasePath = "database.csv"
)

func main() {
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterDatabaseServer(s, &server{databasePath: databasePath})

	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *server) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetReply, error) {
	log.Printf("Received key: %v", in.Key)

	keyValid, errmsg := isKeyValid(in.Key)

	if !keyValid {
		return nil, status.Error(codes.InvalidArgument, errmsg)
	}

	db, err := os.Open(s.databasePath)
	if err != nil {
		log.Printf("Failed to open the database file: %v", err)
		return nil, internalErr
	}

	defer db.Close()

	lookupKey := in.Key
	csvReader := csv.NewReader(db)
	keyFound := false
	var value string

	for {
		record, err := csvReader.Read()

		// reached the end of the file
		if err == io.EOF {
			break
		} else if err != nil {
			log.Printf("Error while reading the database file: %v", err)
			return nil, internalErr
		}

		key := record[0]

		if key == lookupKey {
			keyFound = true
			value = record[1]
		}
	}

	if keyFound {
		return &pb.GetReply{Value: value}, nil
	} else {
		return nil, status.Error(codes.NotFound, "Key was not found")
	}
}

func (s *server) Set(ctx context.Context, in *pb.SetRequest) (*pb.SetReply, error) {
	db, err := os.OpenFile(
		s.databasePath,
		os.O_CREATE|os.O_APPEND|os.O_WRONLY,
		0o644,
	)
	if err != nil {
		log.Printf("Error while opening the database file: %v", err)
		return nil, internalErr
	}

	defer db.Close()

	csvWriter := csv.NewWriter(db)
	err = csvWriter.Write([]string{in.Key, in.Value})

	if err != nil {
		log.Printf("Error while writing to the database file: %v", err)
		return nil, internalErr
	}

	csvWriter.Flush()

	err = csvWriter.Error()

	if err != nil {
		log.Printf("Error while writing to the database file: %v", err)
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
