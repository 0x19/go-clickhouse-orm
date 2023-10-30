package gchm

import "github.com/0x19/go-clickhouse-model/models"

func GetModelKeys(model models.Model) []string {
	var keys []string

	for key, _ := range model.ToMap() {
		keys = append(keys, key)
	}

	return keys
}
