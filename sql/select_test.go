package sql

import (
	"testing"

	"github.com/0x19/go-clickhouse-orm/models"
	"github.com/stretchr/testify/assert"
	"github.com/vahid-sohrabloo/chconn/v3"
)

type DummyModel struct {
	Name string
}

func (d *DummyModel) TableName() string {
	return "dummy"
}

func (d *DummyModel) GetDeclaration() *models.Declaration {
	return &models.Declaration{}
}

func (d *DummyModel) ScanRow(row chconn.SelectStmt) error {
	return nil
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
