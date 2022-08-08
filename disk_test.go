package main

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSaveToYandexDisk(t *testing.T) {
	type testCase struct {
		Name string

		File       io.Reader
		RemotePath string

		ExpectedError error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			actualError := saveToYandexDisk(tc.File, tc.RemotePath)

			assert.Equal(t, tc.ExpectedError, actualError)
		})
	}

	f, _ := os.Open("go.sum")

	validate(t, &testCase{
		Name:          "go mod",
		File:          f,
		RemotePath:    time.Nanosecond.String() + f.Name(),
		ExpectedError: nil,
	})
}
