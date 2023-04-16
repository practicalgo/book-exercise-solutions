package main

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestHandleCommand(t *testing.T) {
	testConfigs := []struct {
		args         []string
		goldenOutput string
		err          error
	}{
		{
			args:         []string{},
			err:          errInvalidSubCommand,
			goldenOutput: "expectedGolden.0",
		},
		{
			args:         []string{"-h"},
			err:          nil,
			goldenOutput: "expectedGolden.1",
		},
		{
			args:         []string{"foo"},
			err:          errInvalidSubCommand,
			goldenOutput: "expectedGolden.2",
		},
	}

	byteBuf := new(bytes.Buffer)
	for i, tc := range testConfigs {
		err := handleCommand(byteBuf, tc.args)
		if tc.err == nil && err != nil {
			t.Fatalf("Expected nil error, got %v", err)
		}

		if tc.err != nil && err.Error() != tc.err.Error() {
			t.Fatalf("Expected error %v, got %v", tc.err, err)
		}

		gotOutput := byteBuf.String()
		expectedOutput, err := os.ReadFile("testdata/" + tc.goldenOutput)
		if err != nil {
			t.Fatalf("error reading expected golden output: %s:%s\n", tc.goldenOutput, err)
		}
		if string(expectedOutput) != gotOutput {
			gotOutputFilename := fmt.Sprintf("testdata/gotOutput.%d", i)
			t.Errorf(
				"Expected output to be:\n%s\n\nGot:\n%s\n\n"+
					"Writing expected data to file: %s",
				string(expectedOutput), gotOutput,
				gotOutputFilename,
			)
			if ok := os.WriteFile(gotOutputFilename, []byte(gotOutput), 0666); ok != nil {
				t.Fatal("Error writing expected output to file", err)
			}
		}
		byteBuf.Reset()
	}
}
