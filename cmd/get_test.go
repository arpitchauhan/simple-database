package cmd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	setTestDatabaseFile()
}

func createDatabase(keyValuePairs [][]string) error {
	db, err := os.Create(testDatabaseFile)
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
	os.Remove(testDatabaseFile)
}

func TestGetValidInput(t *testing.T) {
	t.Cleanup(deleteDatabase)

	type testCase struct {
		databaseContents [][]string
		inputKey string
		expectedResult string
	}

	cases := []testCase{
		{
			databaseContents: [][]string{{"key", "value"}},
			inputKey: "key",
			expectedResult: "Answer: value",
		},
		{
			databaseContents: [][]string{{"key", "value1"}, {"key", "value2"}},
			inputKey: "key",
			expectedResult: "Answer: value2",
		},
		{
			databaseContents: [][]string{{"k", "v"}, {"k2", "v2"}},
			inputKey: "k2",
			expectedResult: "Answer: v2",
		},
	}

	for _, tc := range cases {
		err := createDatabase(tc.databaseContents)
		if err != nil {
			t.Fatal("Failed to create database")
		}

		output, err := executeGetCmd(t, []string{tc.inputKey})

		if err != nil {
			t.Fatal(err)
		}

		if output != tc.expectedResult {
			t.Errorf("Want: %s, got: %s", tc.expectedResult, output)
		}

		deleteDatabase()
	}
}

func TestGetInvalidInput(t *testing.T) {
	type testCase struct {
		input []string
		expectedResult string
	}

	cases := []testCase{
		{
			input: []string{""},
			expectedResult: "Error: key cannot be empty",
		},
		{
			input: []string{"arg1", "arg2"},
			expectedResult: "Error: requires exactly one argument",
		},
		{
			input: []string{"nonexistent_key"},
			expectedResult: "The key is not present in the database",
		},
	}

	t.Cleanup(deleteDatabase)
	createDatabase([][]string{})
	for _, tc := range cases {
		_, err := executeGetCmd(t, tc.input)

		if err == nil {
			if err.Error() != tc.expectedResult {
				t.Errorf("Want: %s, got: %s", tc.expectedResult, err)
			}
		}
	}
}
func BenchmarkGet(b *testing.B) {
	b.Cleanup(deleteDatabase)

	keyValuePairs := [][]string{}
	var testKeys []string

	rowsCount := 100000

	for i := 0; i < rowsCount; i++ {
		key := randStringBytes(5)

		if i == 0 || i == rowsCount / 2 || i == rowsCount - 1 {
			testKeys = append(testKeys, key)
		}

		value := randStringBytes(15)
		keyValuePairs = append(keyValuePairs, []string{key, value})
	}

	createDatabase(keyValuePairs)

	b.ResetTimer()

	for i, testKey := range testKeys {
		b.Run(fmt.Sprintf("key_%d", i), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				_, err := executeGetCmd(b, []string{testKey})
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

func executeGetCmd(t testing.TB, args []string) (string, error) {
	t.Helper()

	b := bytes.NewBufferString("")
	getCmd.SetOut(b)
	os.Args = append([]string{"", "get"}, args...)
	err := getCmd.Execute()

	if err != nil {
		return "", err
	}

	out, err := ioutil.ReadAll(b)
	return string(out), err
}

func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
