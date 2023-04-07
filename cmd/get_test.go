package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Get(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		receivedCode codes.Code
		want         string
	}{
		{
			name:         "Value returned for key",
			key:          "key",
			receivedCode: codes.OK,
			want:         "Answer: value",
		},
		{
			name:         "Server not running",
			receivedCode: codes.Unavailable,
			want:         "Error: the server is not running",
		},
		{
			name:         "Key not present on server",
			receivedCode: codes.NotFound,
			want:         "Error: the key was not found",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedKey string

			// override the fn used to get key from server
			getValueForKey = func(k string) (string, error) {
				receivedKey = k

				return "value", status.Error(tt.receivedCode, "")
			}

			out := executeGetCmd(t, []string{tt.key})

			if receivedKey != tt.key {
				t.Errorf(
					"Server called with wrong key, got = %v, want = %v",
					receivedKey,
					tt.key,
				)
				return
			}

			if out != tt.want {
				t.Errorf("got = %v, want = %v", out, tt.want)
				return
			}
		})
	}
}

func executeGetCmd(t *testing.T, args []string) string {
	t.Helper()

	b := bytes.NewBufferString("")
	getCmd.SetOut(b)
	os.Args = append([]string{"", "get"}, args...)
	err := getCmd.Execute()
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("Error reading output of command: %v", err)
	}

	return string(out)
}
