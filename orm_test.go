package gchm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewORM(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		arg     *Config
		wantErr bool
	}{
		{
			name:    "Basic Test Without Provided Config",
			ctx:     context.TODO(),
			arg:     nil,
			wantErr: true,
		},
		{
			name: "Basic Test With Provided Config",
			ctx:  context.TODO(),
			arg: &Config{
				Host:     "localhost",
				Port:     9000,
				Username: "default",
				Password: "local12345",
				Database: "unpack",
				Insecure: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tAssert := assert.New(t)
			orm, err := NewORM(tt.ctx, tt.arg)
			if tt.wantErr {
				tAssert.Error(err)
			} else {
				tAssert.NoError(err)
				tAssert.NotNil(orm)
			}
		})
	}
}
