package cmd

import (
	"bytes"
	"io/ioutil"
	"os"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Set(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		value        string
		receivedCode codes.Code
		want         string
	}{
		{
			name:         "Successful operation",
			receivedCode: codes.OK,
			want:         "Successful!",
		},
		{
			name:         "Server not running",
			receivedCode: codes.Unavailable,
			want:         "Error: the server is not running",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var receivedKey, receivedValue string

			// override the fn used to get key from server
			setValueForKey = func(k, v string) error {
				receivedKey = k
				receivedValue = v

				return status.Error(tt.receivedCode, "")
			}

			out := executeSetCmd(t, []string{tt.key, tt.value})

			if receivedKey != tt.key || receivedValue != tt.value {
				t.Errorf(
					"Server called with wrong key-value pair, got: %v and %v, want: %v and %v",
					receivedKey,
					receivedValue,
					tt.key,
					tt.value,
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

func executeSetCmd(t *testing.T, args []string) string {
	t.Helper()

	b := bytes.NewBufferString("")
	setCmd.SetOut(b)
	os.Args = append([]string{"", "set"}, args...)
	err := setCmd.Execute()
	if err != nil {
		t.Fatalf("Error executing command: %v", err)
	}

	out, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("Error reading output of command: %v", err)
	}

	return string(out)
}
