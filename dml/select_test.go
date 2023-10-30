package dml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type DummyModel struct {
	Name string
}

func (d *DummyModel) TableName() string {
	return "dummy"
}

func (d *DummyModel) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name": d.Name,
	}
}

func TestSelectBuilder(t *testing.T) {
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

			builder := NewSelectBuilder()
			tAssert.NotNil(builder)

			builder.Model((*DummyModel)(nil))
			builder.Select(
				"count()",
				NewSelectBuilder().Model((*DummyModel)(nil)).Select("count()"),
			)

			stmt, err := builder.Build()
			tAssert.Equal(tt.expectedErr, err)
			tAssert.Equal("", stmt)
		})
	}
}
