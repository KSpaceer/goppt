package ioadapters_test

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/KSpaceer/goppt/internal/ioadapters"
)

func TestAdaptedReaderAt(t *testing.T) {
	type operation struct {
		offset          int64
		count           int
		expected        []byte
		expectedWritten int
		expectedErr     error
	}

	type tcase struct {
		name       string
		input      io.Reader
		operations []operation
	}

	tcases := []tcase{
		{
			name:  "convertible to ReaderAt",
			input: bytes.NewReader([]byte("abcdef")),
			operations: []operation{
				{
					offset:          0,
					count:           3,
					expected:        []byte("abc"),
					expectedWritten: 3,
				},
				{
					offset:          3,
					count:           3,
					expected:        []byte("def"),
					expectedWritten: 3,
				},
				{
					offset:          5,
					count:           3,
					expected:        []byte("f"),
					expectedWritten: 1,
					expectedErr:     io.EOF,
				},
			},
		},
		{
			name:  "non-convertible to ReaderAt",
			input: bytes.NewBuffer([]byte("abcdef")),
			operations: []operation{
				{
					offset:          0,
					count:           3,
					expected:        []byte("abc"),
					expectedWritten: 3,
				},
				{
					offset:          3,
					count:           3,
					expected:        []byte("def"),
					expectedWritten: 3,
				},
				{
					offset:          5,
					count:           3,
					expected:        []byte("f"),
					expectedWritten: 1,
					expectedErr:     io.EOF,
				},
			},
		},
		{
			name:  "bytes gap",
			input: bytes.NewBuffer([]byte("FIRSTxxxxxxxxxxxxxxxxWORLDxxxxxxxxxxxxxxHELLO")),
			operations: []operation{
				{
					offset:          0,
					count:           5,
					expected:        []byte("FIRST"),
					expectedWritten: 5,
				},
				{
					offset:          40,
					count:           5,
					expected:        []byte("HELLO"),
					expectedWritten: 5,
				},
				{
					offset:          21,
					count:           5,
					expected:        []byte("WORLD"),
					expectedWritten: 5,
				},
				{
					offset:          43,
					count:           4,
					expected:        []byte("LO"),
					expectedWritten: 2,
					expectedErr:     io.EOF,
				},
			},
		},
	}

	for _, tc := range tcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			ra := ioadapters.ToReaderAt(tc.input)
			for i, op := range tc.operations {
				buf := make([]byte, op.count)
				written, err := ra.ReadAt(buf, op.offset)
				if !errors.Is(err, op.expectedErr) {
					t.Errorf("operation %d: expected error %v, got %v", i, op.expectedErr, err)
				}
				if written != op.expectedWritten {
					t.Errorf("operation %d: expected written count %d, got %d", i, op.expectedWritten, written)
				} else if !bytes.Equal(op.expected, buf[:written]) {
					t.Errorf("operation %d: expected bytes %v, got %v", i, op.expected, buf)
				}
			}
		})
	}
}
