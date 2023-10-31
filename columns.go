package chorm

import "reflect"

type Field struct {
	Table      string
	PrimaryKey bool
	Name       string
	Default    string
	Type       string
	GoField    reflect.StructField
}

/* func GetColumnReaders(columns ...string) []column.ColumnBasic {
	var columnReaders []column.ColumnBasic
	for _, fieldName := range b.GetFieldsByOrder() {
		if !stringInSlice(fieldName, columns) {
			continue
		}

		field := b.fields[fieldName]
		col, err := b.GetColumnByType(b, field, columnReaders, false)
		if err != nil {
			zap.L().Error(
				"error getting column by type",
				zap.Error(err),
				zap.String("field", field.Name),
			)
			continue
		}

		columnReaders = append(columnReaders, col)
	}
	return columnReaders
}

func (b *BaseModel[T]) GetColumnByType(m any, field Field, columns []column.ColumnBasic, fill bool) (column.ColumnBasic, error) {
	element := reflect.ValueOf(m).Elem()
	switch field.GoField.Type.Kind() {
	case reflect.Int:
		col := column.New[int64]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(element.FieldByName(field.GoField.Name).Int())
		}
		return col, nil
	case reflect.Int8:
		col := column.New[int8]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(int8(element.FieldByName(field.GoField.Name).Int()))
		}
		return col, nil
	case reflect.Int16:
		col := column.New[int16]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(int16(element.FieldByName(field.GoField.Name).Int()))
		}
		return col, nil
	case reflect.Int32:
		col := column.New[int32]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(int32(element.FieldByName(field.GoField.Name).Int()))
		}
		return col, nil
	case reflect.Int64:
		col := column.New[int64]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(int64(element.FieldByName(field.GoField.Name).Int()))
		}
		return col, nil
	case reflect.Uint:
		col := column.New[uint64]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(uint64(element.FieldByName(field.GoField.Name).Uint()))
		}
		return col, nil
	case reflect.Uint8:
		col := column.New[uint8]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(uint8(element.FieldByName(field.GoField.Name).Uint()))
		}
		return col, nil
	case reflect.Uint16:
		col := column.New[uint16]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(uint16(element.FieldByName(field.GoField.Name).Uint()))
		}
		return col, nil
	case reflect.Uint32:
		col := column.New[uint32]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(uint32(element.FieldByName(field.GoField.Name).Uint()))
		}
		return col, nil
	case reflect.Uint64:
		col := column.New[uint64]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(uint64(element.FieldByName(field.GoField.Name).Uint()))
		}
		return col, nil
	case reflect.String:
		col := column.NewString()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(element.FieldByName(field.GoField.Name).String())
		}
		return col, nil
	case reflect.Bool:
		col := column.New[bool]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(element.FieldByName(field.GoField.Name).Bool())
		}
		return col, nil
	case reflect.Float32:
		col := column.New[float32]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(float32(element.FieldByName(field.GoField.Name).Float()))
		}
		return col, nil
	case reflect.Float64:
		col := column.New[float64]()
		col.SetName([]byte(field.Name))
		if fill {
			col.Append(element.FieldByName(field.GoField.Name).Float())
		}
		return col, nil
	case reflect.Struct:
		switch field.GoField.Type.String() {
		case "time.Time":
			col := column.NewDate[types.DateTime]()
			col.SetName([]byte(field.Name))
			if fill {
				col.Append(element.FieldByName(field.GoField.Name).Interface().(time.Time))
			}
			return col, nil
		default:
			return nil, fmt.Errorf("unsupported struct data type: %s", field.GoField.Type.Kind())
		}
	case reflect.Array:
		switch field.GoField.Type.String() {
		case "uuid.UUID":
			col := column.New[uuid.UUID]()
			col.SetName([]byte(field.Name))
			if fill {
				col.Append(element.FieldByName(field.GoField.Name).Interface().(uuid.UUID))
			}
			return col, nil
		default:
			return nil, fmt.Errorf("unsupported struct data type: %s", field.GoField.Type.Kind())
		}
	default:
		return nil, fmt.Errorf("unsupported data type: %s", field.GoField.Type.Kind())
	}
}
*/
