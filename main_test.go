package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
	"time"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "A",
			args: []string{"get", "key"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = append([]string{"get"}, tt.args...)

			var buffOut bytes.Buffer
			var buffErr bytes.Buffer

			run(&buffOut, &buffErr)

			time.Sleep(2 * time.Second)
			fmt.Print(buffOut.String())
			fmt.Print(buffErr.String())
		})
	}
}
