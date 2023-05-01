package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/arpitchauhan/simple-database/database"
)

const testDatabasePath = "database_test.csv"

func Test_server_Get(t *testing.T) {
	tests := []struct {
		name             string
		databaseContents [][]string
		inputKey         string
		want             string
		wantErr          bool
		wantErrCode      codes.Code
		wantErrMsg       string
	}{
		{
			name:             "Database with just one value for the input key",
			inputKey:         "key",
			databaseContents: [][]string{{"key", "value"}},
			want:             "value",
			wantErr:          false,
		},
		{
			name:             "Database with two values for the input key",
			inputKey:         "key",
			databaseContents: [][]string{{"key", "value1"}, {"key", "value2"}},
			want:             "value2",
			wantErr:          false,
		},
		{
			name:             "Database with two keys - key 1",
			inputKey:         "key",
			databaseContents: [][]string{{"key", "value"}, {"key2", "value2"}},
			want:             "value",
			wantErr:          false,
		},
		{
			name:             "Database with two keys - key 2",
			inputKey:         "key2",
			databaseContents: [][]string{{"key", "value"}, {"key2", "value2"}},
			want:             "value2",
			wantErr:          false,
		},
		{
			name:             "Input key blank",
			inputKey:         "",
			databaseContents: [][]string{},
			want:             "",
			wantErr:          true,
			wantErrCode:      codes.InvalidArgument,
			wantErrMsg:       "Key cannot be empty",
		},
		{
			name:             "Key not present in database",
			inputKey:         "nonexistent_key",
			databaseContents: [][]string{{"key", "value"}},
			want:             "",
			wantErr:          true,
			wantErrCode:      codes.NotFound,
			wantErrMsg:       "Key was not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(deleteDatabase)
			createDatabase(tt.databaseContents)

			s := getServer()
			getRequest := &pb.GetRequest{Key: tt.inputKey}
			got, err := s.Get(context.Background(), getRequest)

			if err == nil {
				if got.Value != tt.want {
					t.Errorf("got = %v, want %v", got, tt.want)
				}
			} else {
				if !tt.wantErr {
					t.Errorf("error = %v, did not want error", err)
					return
				}

				s, _ := status.FromError(err)

				if s.Code() != tt.wantErrCode {
					t.Errorf("code = %v, want = %v", err, tt.wantErrCode)
					return
				}

				if s.Message() != tt.wantErrMsg {
					t.Errorf("error message = %v, want = %v", s.Message(), tt.wantErrMsg)
					return
				}
			}
		})
	}
}

func Test_server_Set(t *testing.T) {
	tests := []struct {
		name        string
		input       [][]string // a bunch of key-value pairs with which Set is called
		want        string     // end state of database file
		wantErr     bool
		wantErrCode codes.Code
		wantErrMsg  string
	}{
		{
			name:    "One key-value pair",
			input:   [][]string{{"key", "value"}},
			want:    "key,value\n",
			wantErr: false,
		},
		{
			name:    "Two key-value pairs (different keys)",
			input:   [][]string{{"key1", "value1"}, {"key2", "value2"}},
			want:    "key1,value1\nkey2,value2\n",
			wantErr: false,
		},
		{
			name:    "Two key-value pairs (same key)",
			input:   [][]string{{"key", "value"}, {"key", "value2"}},
			want:    "key,value\nkey,value2\n",
			wantErr: false,
		},
		{
			name:    "Blank value for a key",
			input:   [][]string{{"key", ""}},
			want:    "key,\n",
			wantErr: false,
		},
		{
			name:        "Empty key",
			input:       [][]string{{"", "value"}},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
			wantErrMsg:  "Key cannot be empty",
		},
		{
			name:        "Blank key",
			input:       [][]string{{"  ", "value"}},
			wantErr:     true,
			wantErrCode: codes.InvalidArgument,
			wantErrMsg:  "Key cannot be empty",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Cleanup(deleteDatabase)
			s := getServer()

			for _, kv := range tt.input {
				setRequest := &pb.SetRequest{Key: kv[0], Value: kv[1]}
				_, err := s.Set(context.Background(), setRequest)
				if err != nil {
					if !tt.wantErr {
						t.Errorf("error = %v, did not want error", err)
						return
					}

					s, _ := status.FromError(err)

					if s.Code() != tt.wantErrCode {
						t.Errorf("code = %v, want = %v", err, tt.wantErrCode)
					}

					if s.Message() != tt.wantErrMsg {
						t.Errorf("error message = %v, want = %v", s.Message(), tt.wantErrMsg)
					}
				}
			}

			if !tt.wantErr {
				dbContents, err := os.ReadFile(testDatabasePath)
				if err != nil {
					t.Fatal(err)
				}

				if string(dbContents) != tt.want {
					t.Errorf(
						"The content of database file is not as expected. got = %v, want = %v",
						string(dbContents),
						tt.want,
					)
				}
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	b.Cleanup(deleteDatabase)

	keyValuePairs := [][]string{}
	var testKeys []string

	rowsCount := 100000

	for i := 0; i < rowsCount; i++ {
		key := randStringBytes(5)

		if i == 0 || i == rowsCount/2 || i == rowsCount-1 {
			testKeys = append(testKeys, key)
		}

		value := randStringBytes(15)
		keyValuePairs = append(keyValuePairs, []string{key, value})
	}

	createDatabase(keyValuePairs)

	b.ResetTimer()

	s := getServer()

	log.SetOutput(ioutil.Discard) // skip logging

	for i, testKey := range testKeys {
		b.Run(fmt.Sprintf("key_%d", i), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				getRequest := &pb.GetRequest{Key: testKey}
				_, err := s.Get(context.Background(), getRequest)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func BenchmarkSet(b *testing.B) {
	b.Cleanup(deleteDatabase)
	s := getServer()

	log.SetOutput(ioutil.Discard) // skip logging

	for n := 0; n < b.N; n++ {
		setRequest := &pb.SetRequest{Key: "k", Value: "v"}
		// fmt.Println(n)
		_, err := s.Set(context.Background(), setRequest)
		if err != nil {
			b.Fatalf("Error: %s", err)
		}
	}
}

func getServer() *server {
	db := &database{filepath: testDatabasePath, initialized: false}
	s := &server{
		UnimplementedDatabaseServer: pb.UnimplementedDatabaseServer{},
		db:                          db,
	}
	s.initialize()
	return s
}

func createDatabase(keyValuePairs [][]string) error {
	db, err := os.Create(testDatabasePath)
	if err != nil {
		return err
	}

	defer db.Close()

	csvWriter := csv.NewWriter(db)

	for _, kv := range keyValuePairs {
		err = csvWriter.Write(kv)
		if err != nil {
			return err
		}
	}

	csvWriter.Flush()
	err = csvWriter.Error()
	if err != nil {
		return err
	}

	return nil
}

func deleteDatabase() {
	os.Remove(testDatabasePath)
}

func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
