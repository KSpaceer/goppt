package goppt_test

import (
	"bytes"
	"embed"
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/KSpaceer/goppt"
)

//go:embed testdata
var testdataFS embed.FS

const (
	simplePresPath = "testdata/simplepres.ppt"
)

func TestExtractText(t *testing.T) {
	type tcase struct {
		name         string
		r            io.Reader
		expectedText string
		expectErr    error
	}

	tcases := []tcase{
		{
			name: "simple reader at",
			r: func() io.Reader {
				f, err := testdataFS.Open(simplePresPath)
				if err != nil {
					panic(err)
				}
				t.Cleanup(func() {
					f.Close()
				})
				return f
			}(),
			expectedText: "12345",
		},
		{
			name: "simple no reader at",
			r: func() io.Reader {
				f, err := testdataFS.Open(simplePresPath)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				data, err := io.ReadAll(f)
				if err != nil {
					panic(err)
				}
				return bytes.NewBuffer(data)
			}(),
			expectedText: "12345",
		},
	}

	for _, tc := range tcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			text, err := goppt.ExtractText(tc.r)
			if !errors.Is(err, tc.expectErr) {
				t.Fatalf("expected error %v, but got %v", tc.expectErr, err)
			}
			text = strings.TrimSpace(text)
			expected := strings.TrimSpace(tc.expectedText)
			if text != expected {
				t.Fatalf("expected text %q, but got %q", expected, text)
			}
		})
	}
}
