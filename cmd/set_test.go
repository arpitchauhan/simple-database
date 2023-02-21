package cmd

import (
	"os"
	"testing"

	"github.com/spf13/cobra"
)

const testDatabaseFile = "database_test.csv"

func TestSetValidInput(t *testing.T) {
	setTestDatabaseFile()

	type testCase struct {
		input [][]string
		want  string
	}
	cases := []testCase{
		{
			input: [][]string{{"key", "value"}},
			want:  "key,value\n",
		},
		{
			input: [][]string{{"key1", "value1"}, {"key2", "value2"}},
			want:  "key1,value1\nkey2,value2\n",
		},
		{
			input: [][]string{{"key", "value"}, {"key", "value2"}},
			want:  "key,value\nkey,value2\n",
		},
		{
			input: [][]string{{"key", ""}},
			want:  "key,\n",
		},
	}

	for _, tc := range cases {
		for _, kv := range tc.input {
			err := execute(t, setCmd, kv)
			if err != nil {
				t.Fatal(err)
			}
		}

		dbContents, err := os.ReadFile(testDatabaseFile)

		if err != nil {
			t.Fatal(err)
		}

		if string(dbContents) != tc.want {
			t.Errorf(
				"The content of database is not as expected. Want: %s, got: %s",
				tc.want,
				dbContents,
			)
		}

		os.Remove(testDatabaseFile)
	}
}

func TestSetInvalidInput(t *testing.T) {
	setTestDatabaseFile()

	type testCase struct {
		input     []string
		wantError string
	}

	cases := []testCase{
		{
			input:     []string{"", "value"},
			wantError: "key cannot be empty",
		},
		{
			input:     []string{"  ", "value"},
			wantError: "key cannot be empty",
		},
		{
			input:     []string{"value"},
			wantError: "requires exactly 2 args",
		},
	}
	for _, tc := range cases {
		err := execute(t, setCmd, tc.input)

		if err == nil {
			t.Fatal("got no error, want one")
		}

		if err.Error() != tc.wantError {
			t.Errorf("got error \"%s\", want \"%s\"", err.Error(), tc.wantError)
		}

		os.Remove(testDatabaseFile)
	}
}

func execute(t *testing.T, c *cobra.Command, args []string) error {
	t.Helper()

	os.Args = append([]string{"", "set"}, args...)
	err := c.Execute()

	return err
}

func setTestDatabaseFile() {
	rootCmd.PersistentFlags().Set("database", testDatabaseFile)
}
