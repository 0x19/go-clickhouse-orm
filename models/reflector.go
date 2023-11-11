package models

import (
	"errors"
	"reflect"
	"sort"
	"strings"
)

type ModelDetails struct {
	TableName   string
	Engine      string
	PartitionBy string
	OrderBy     []string
}

type Reflector struct {
	databaseName string
}

func (r *Reflector) ReflectModel(modelType reflect.Type) (*Declaration, error) {
	if modelType.Kind() != reflect.Struct {
		return nil, errors.New("modelType must be a struct")
	}

	details, err := r.extractModelDetails(modelType)
	if err != nil {
		return nil, err
	}

	declaration := &Declaration{
		DatabaseName: r.databaseName,
		TableName:    details.TableName,
		Engine:       details.Engine,
		PartitionBy:  details.PartitionBy,
		OrderBy:      details.OrderBy,
		Fields:       make(map[string]Field),
		Settings:     r.getSettings(modelType),
	}

	tags, err := r.extractTags(modelType)
	if err != nil {
		return nil, err
	}

	var sortedFieldNames []string
	var fieldIndex int16 = 0

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag, ok := tags[field.Name]
		if !ok || tag == "-" {
			continue
		}

		goType := field.Type.String()
		if goType == "models.Model" {
			continue
		}

		tokens := strings.Split(tag, ",")
		if len(tokens) == 0 {
			continue
		}

		fieldDetails := Field{
			Index:       fieldIndex,
			Name:        field.Name,
			ReflectType: field.Type,
			GoType:      goType,
		}

		for _, token := range tokens {
			pair := strings.SplitN(strings.TrimSpace(token), ":", 2)
			if len(pair) == 2 {
				key, value := strings.TrimSpace(pair[0]), strings.TrimSpace(pair[1])
				switch key {
				case "name":
					fieldDetails.Name = value
				case "primary":
					fieldDetails.PrimaryKey = true
				case "type":
					fieldDetails.Type = value
				case "nullable":
					fieldDetails.Nullable = value == "true"
				case "default":
					fieldDetails.Default = value
				case "comment":
					fieldDetails.Comment = value
				}
			}
		}

		declaration.Fields[field.Name] = fieldDetails
		sortedFieldNames = append(sortedFieldNames, field.Name)
		fieldIndex++
	}

	sort.Slice(sortedFieldNames, func(i, j int) bool {
		return declaration.Fields[sortedFieldNames[i]].Index < declaration.Fields[sortedFieldNames[j]].Index
	})

	return declaration, nil
}

func (r *Reflector) extractTags(modelType reflect.Type) (map[string]string, error) {
	if modelType.Kind() != reflect.Struct {
		return nil, errors.New("modelType must be a struct")
	}

	tags := make(map[string]string)
	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag := field.Tag.Get("clickhouse")
		if tag != "" {
			tags[field.Name] = tag
		}
	}

	return tags, nil
}

func (r *Reflector) getTableName(modelType reflect.Type) string {
	if method, ok := modelType.MethodByName("TableName"); ok {
		if method.Type.NumOut() == 1 && method.Type.Out(0).Kind() == reflect.String {
			result := method.Func.Call([]reflect.Value{reflect.New(modelType).Elem()})
			if len(result) > 0 {
				return result[0].String()
			}
		}
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		tag := field.Tag.Get("clickhouse")
		if tag != "" {
			tagParts := strings.Split(tag, ",")
			for _, part := range tagParts {
				if strings.HasPrefix(part, "table:") {
					tableName := strings.TrimPrefix(part, "table:")
					return strings.TrimSpace(tableName)
				}
			}
		}
	}

	return strings.ToLower(modelType.Name())
}

func (r *Reflector) getSettings(modelType reflect.Type) []string {
	if modelType.Kind() != reflect.Ptr {
		// If modelType is not a pointer, get the pointer to it
		modelType = reflect.PtrTo(modelType)
	}

	if method, ok := modelType.MethodByName("Settings"); ok {
		if method.Type.NumOut() == 1 && method.Type.Out(0).Kind() == reflect.Slice &&
			method.Type.Out(0).Elem().Kind() == reflect.String {

			result := method.Func.Call([]reflect.Value{reflect.New(modelType.Elem())})
			if len(result) > 0 {
				settings, ok := result[0].Interface().([]string)
				if ok {
					return settings
				}
			}
		}
	}

	return []string{}
}

func (r *Reflector) extractModelDetails(modelType reflect.Type) (ModelDetails, error) {
	var details ModelDetails

	details.TableName = r.getTableName(modelType)

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Interface {
			tag := field.Tag.Get("clickhouse")
			if tag == "" {
				continue
			}

			tagParts := strings.Split(tag, ",")
			for _, part := range tagParts {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "engine:") {
					details.Engine = strings.TrimPrefix(part, "engine:")
				} else if strings.HasPrefix(part, "partition:") {
					details.PartitionBy = strings.TrimPrefix(part, "partition:")
				} else if strings.HasPrefix(part, "order:") {
					details.OrderBy = strings.Split(strings.TrimPrefix(part, "order:"), ",")
				}
			}
			break
		}
	}

	return details, nil
}
