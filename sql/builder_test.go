package sql

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuilder(t *testing.T) {
	tests := []struct {
		name        string
		expectedErr error
	}{
		{
			name:        "Valid Config",
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tAssert := assert.New(t)
			_ = tAssert
		})
	}
}
