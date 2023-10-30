package chorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	tests := []struct {
		name        string
		config      *Config
		expectedErr error
		expectedDSN string
	}{
		{
			name:        "Valid Config",
			config:      &Config{Host: "localhost", Port: 9000, Username: "user", Password: "pass", Database: "db"},
			expectedErr: nil,
			expectedDSN: "host=localhost port=9000 user=user password=pass dbname=db sslmode=disable",
		},
		{
			name:        "Nil Config",
			config:      nil,
			expectedErr: ErrNoConfigProvided,
			expectedDSN: "",
		},
		{
			name:        "No Host",
			config:      &Config{Port: 9000, Username: "user", Password: "pass", Database: "db"},
			expectedErr: ErrNoHostProvided,
			expectedDSN: "",
		},
		// ... more test cases ...
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tAssert := assert.New(t)
			err := tt.config.Validate()
			tAssert.Equal(tt.expectedErr, err)
			if err == nil {
				dsn := tt.config.GetDSN()
				tAssert.Equal(tt.expectedDSN, dsn)
			}
		})
	}
}
