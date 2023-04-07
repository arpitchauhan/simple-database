package client

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/arpitchauhan/simple-database/database"
)

const addr = "localhost:50051"

// Create a new connection, a client that uses that connection, executes
// the passed-in request, and then returns the result of the execution
func executeRequest(
	requestFn func(pb.DatabaseClient, context.Context) (string, error),
) (string, error) {
	// Use insecure credentials as this project is not meant for real-world use
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()

	client := pb.NewDatabaseClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10000*time.Second)
	defer cancel()

	return requestFn(client, ctx)
}

func GetValueForKey(key string) (string, error) {
	requestFn := func(client pb.DatabaseClient, ctx context.Context) (string, error) {
		reply, err := client.Get(ctx, &pb.GetRequest{Key: key})
		if err != nil {
			return "", err
		}

		return reply.Value, nil
	}

	value, err := executeRequest(requestFn)
	if err != nil {
		return "", err
	}

	return value, nil
}

func SetValueForKey(key string, value string) error {
	requestFn := func(client pb.DatabaseClient, ctx context.Context) (string, error) {
		_, err := client.Set(ctx, &pb.SetRequest{Key: key, Value: value})
		if err != nil {
			return "", err
		}

		return "", nil
	}

	_, err := executeRequest(requestFn)

	return err
}
