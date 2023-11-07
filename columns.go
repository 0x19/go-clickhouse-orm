package chorm

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/0x19/go-clickhouse-orm/models"
	"github.com/vahid-sohrabloo/chconn/v3/column"
	"github.com/vahid-sohrabloo/chconn/v3/types"
)

type Field struct {
	Table      string
	PrimaryKey bool
	Name       string
	Default    string
	Type       string
	GoField    reflect.StructField
}

func (s *SelectBuilder[T]) GetColumnsByChType(b models.Model) ([]column.ColumnBasic, error) {
	columns := make([]column.ColumnBasic, len(b.GetDeclaration().GetColumns()))
	for i, col := range b.GetDeclaration().GetColumns() {
		columnByType, err := s.ColumnByType(col.Type(), 0, false, false)
		if err != nil {
			return nil, err
		}
		columnByType.SetName(col.Name())
		columnByType.SetType(col.Type())
		err = columnByType.Validate()
		if err != nil {
			return nil, err
		}
		columns[i] = columnByType
	}
	return columns, nil
}

//nolint:funlen,gocyclo
func (s *SelectBuilder[T]) ColumnByType(chType []byte, arrayLevel int, nullable, lc bool) (column.ColumnBasic, error) {
	switch {
	case string(chType) == "Int8" || IsEnum8(chType):
		return column.New[int8]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "Int16" || IsEnum16(chType):
		return column.New[int16]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "Int32":
		return column.New[int32]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "Int64":
		return column.New[int64]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "Int128":
		return column.New[types.Int128]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "Int256":
		return column.New[types.Int256]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "UInt8":
		return column.New[uint8]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "UInt16":
		return column.New[uint16]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "UInt32":
		return column.New[uint32]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "UInt64":
		return column.New[uint64]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "UInt128":
		return column.New[types.Uint128]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "UInt256":
		return column.New[types.Uint256]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "Float32":
		return column.New[float32]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "Float64":
		return column.New[float64]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "String":
		return column.NewString().Elem(arrayLevel, nullable, lc), nil
	case IsFixedString(chType):
		strLen, err := strconv.Atoi(string(chType[FixedStringStrLen : len(chType)-1]))
		if err != nil {
			return nil, fmt.Errorf("invalid fixed string length: %s: %w", string(chType), err)
		}
		return getFixedType(strLen, arrayLevel, nullable, lc)
	case string(chType) == "Date":
		return column.NewDate[types.Date]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "Date32":
		return column.NewDate[types.Date32]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "DateTime" || IsDateTimeWithParam(chType):
		var params [][]byte
		if bytes.HasPrefix(chType, []byte("DateTime(")) {
			params = bytes.Split(chType[len("DateTime("):len(chType)-1], []byte(", "))
		}
		col := column.NewDate[types.DateTime]()
		if len(params) > 0 && len(params[0]) >= 3 {
			if loc, err := time.LoadLocation(string(params[0][1 : len(params[0])-1])); err == nil {
				col.SetLocation(loc)
			} else if loc, err := time.LoadLocation(s.serverInfo.Timezone); err == nil {
				col.SetLocation(loc)
			}
		}
		return col.Elem(arrayLevel, nullable, lc), nil
	case IsDateTime64(chType):
		params := bytes.Split(chType[DateTime64StrLen:len(chType)-1], []byte(", "))
		if len(params) == 0 {
			panic("DateTime64 invalid params")
		}
		precision, err := strconv.Atoi(string(params[0]))
		if err != nil {
			panic("DateTime64 invalid precision: " + err.Error())
		}
		col := column.NewDate[types.DateTime64]()
		col.SetPrecision(precision)
		if len(params) > 1 && len(params[1]) >= 3 {
			if loc, err := time.LoadLocation(string(params[1][1 : len(params[1])-1])); err == nil {
				col.SetLocation(loc)
			} else if loc, err := time.LoadLocation(s.serverInfo.Timezone); err == nil {
				col.SetLocation(loc)
			}
		}
		return col.Elem(arrayLevel, nullable, lc), nil

	case IsDecimal(chType):
		params := bytes.Split(chType[DecimalStrLen:len(chType)-1], []byte(", "))
		precision, _ := strconv.Atoi(string(params[0]))

		if precision <= 9 {
			return column.New[types.Decimal32]().Elem(arrayLevel, nullable, lc), nil
		}
		if precision <= 18 {
			return column.New[types.Decimal64]().Elem(arrayLevel, nullable, lc), nil
		}
		if precision <= 38 {
			return column.New[types.Decimal128]().Elem(arrayLevel, nullable, lc), nil
		}
		if precision <= 76 {
			return column.New[types.Decimal256]().Elem(arrayLevel, nullable, lc), nil
		}
		panic("Decimal invalid precision: " + string(chType))

	case string(chType) == "UUID":
		return column.New[types.UUID]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "IPv4":
		return column.New[types.IPv4]().Elem(arrayLevel, nullable, lc), nil
	case string(chType) == "IPv6":
		return column.New[types.IPv6]().Elem(arrayLevel, nullable, lc), nil

	case IsNullable(chType):
		return s.ColumnByType(chType[LenNullableStr:len(chType)-1], arrayLevel, true, lc)

	case bytes.HasPrefix(chType, []byte("SimpleAggregateFunction(")):
		return s.ColumnByType(FilterSimpleAggregate(chType), arrayLevel, nullable, lc)
	case IsArray(chType):
		if arrayLevel == 3 {
			return nil, fmt.Errorf("max array level is 3")
		}
		if nullable {
			return nil, fmt.Errorf("array is not allowed in nullable")
		}
		if lc {
			return nil, fmt.Errorf("LowCardinality is not allowed in nullable")
		}
		return s.ColumnByType(chType[LenArrayStr:len(chType)-1], arrayLevel+1, nullable, lc)
	case IsLowCardinality(chType):
		return s.ColumnByType(chType[LenLowCardinalityStr:len(chType)-1], arrayLevel, nullable, true)
	case IsTuple(chType):
		columnsTuple, err := TypesInParentheses(chType[LenTupleStr : len(chType)-1])
		if err != nil {
			return nil, fmt.Errorf("tuple invalid types: %w", err)
		}
		columns := make([]column.ColumnBasic, len(columnsTuple))
		for i, c := range columnsTuple {
			col, err := s.ColumnByType(c.ChType, 0, false, false)
			if err != nil {
				return nil, err
			}
			col.SetName(c.Name)
			columns[i] = col
		}
		return column.NewTuple(columns...).Elem(arrayLevel), nil
	case IsMap(chType):
		columnsMap, err := TypesInParentheses(chType[LenMapStr : len(chType)-1])
		if err != nil {
			return nil, fmt.Errorf("map invalid types: %w", err)
		}
		if len(columnsMap) != 2 {
			return nil, fmt.Errorf("map must have 2 columns")
		}
		columns := make([]column.ColumnBasic, len(columnsMap))
		for i, col := range columnsMap {
			col, err := s.ColumnByType(col.ChType, arrayLevel, nullable, lc)
			if err != nil {
				return nil, err
			}
			columns[i] = col
		}
		return column.NewMapBase(columns[0], columns[1]), nil
	case IsNested(chType):
		return s.ColumnByType(NestedToArrayType(chType), arrayLevel, nullable, lc)
	}
	return nil, fmt.Errorf("unknown type: %s", chType)
}

//nolint:funlen,gocyclo
func getFixedType(fixedLen, arrayLevel int, nullable, lc bool) (column.ColumnBasic, error) {
	switch fixedLen {
	case 1:
		return column.New[[1]byte]().Elem(arrayLevel, nullable, lc), nil
	case 2:
		return column.New[[2]byte]().Elem(arrayLevel, nullable, lc), nil
	case 3:
		return column.New[[3]byte]().Elem(arrayLevel, nullable, lc), nil
	case 4:
		return column.New[[4]byte]().Elem(arrayLevel, nullable, lc), nil
	case 5:
		return column.New[[5]byte]().Elem(arrayLevel, nullable, lc), nil
	case 6:
		return column.New[[6]byte]().Elem(arrayLevel, nullable, lc), nil
	case 7:
		return column.New[[7]byte]().Elem(arrayLevel, nullable, lc), nil
	case 8:
		return column.New[[8]byte]().Elem(arrayLevel, nullable, lc), nil
	case 9:
		return column.New[[9]byte]().Elem(arrayLevel, nullable, lc), nil
	case 10:
		return column.New[[10]byte]().Elem(arrayLevel, nullable, lc), nil
	case 11:
		return column.New[[11]byte]().Elem(arrayLevel, nullable, lc), nil
	case 12:
		return column.New[[12]byte]().Elem(arrayLevel, nullable, lc), nil
	case 13:
		return column.New[[13]byte]().Elem(arrayLevel, nullable, lc), nil
	case 14:
		return column.New[[14]byte]().Elem(arrayLevel, nullable, lc), nil
	case 15:
		return column.New[[15]byte]().Elem(arrayLevel, nullable, lc), nil
	case 16:
		return column.New[[16]byte]().Elem(arrayLevel, nullable, lc), nil
	case 17:
		return column.New[[17]byte]().Elem(arrayLevel, nullable, lc), nil
	case 18:
		return column.New[[18]byte]().Elem(arrayLevel, nullable, lc), nil
	case 19:
		return column.New[[19]byte]().Elem(arrayLevel, nullable, lc), nil
	case 20:
		return column.New[[20]byte]().Elem(arrayLevel, nullable, lc), nil
	case 21:
		return column.New[[21]byte]().Elem(arrayLevel, nullable, lc), nil
	case 22:
		return column.New[[22]byte]().Elem(arrayLevel, nullable, lc), nil
	case 23:
		return column.New[[23]byte]().Elem(arrayLevel, nullable, lc), nil
	case 24:
		return column.New[[24]byte]().Elem(arrayLevel, nullable, lc), nil
	case 25:
		return column.New[[25]byte]().Elem(arrayLevel, nullable, lc), nil
	case 26:
		return column.New[[26]byte]().Elem(arrayLevel, nullable, lc), nil
	case 27:
		return column.New[[27]byte]().Elem(arrayLevel, nullable, lc), nil
	case 28:
		return column.New[[28]byte]().Elem(arrayLevel, nullable, lc), nil
	case 29:
		return column.New[[29]byte]().Elem(arrayLevel, nullable, lc), nil
	case 30:
		return column.New[[30]byte]().Elem(arrayLevel, nullable, lc), nil
	case 31:
		return column.New[[31]byte]().Elem(arrayLevel, nullable, lc), nil
	case 32:
		return column.New[[32]byte]().Elem(arrayLevel, nullable, lc), nil
	case 33:
		return column.New[[33]byte]().Elem(arrayLevel, nullable, lc), nil
	case 34:
		return column.New[[34]byte]().Elem(arrayLevel, nullable, lc), nil
	case 35:
		return column.New[[35]byte]().Elem(arrayLevel, nullable, lc), nil
	case 36:
		return column.New[[36]byte]().Elem(arrayLevel, nullable, lc), nil
	case 37:
		return column.New[[37]byte]().Elem(arrayLevel, nullable, lc), nil
	case 38:
		return column.New[[38]byte]().Elem(arrayLevel, nullable, lc), nil
	case 39:
		return column.New[[39]byte]().Elem(arrayLevel, nullable, lc), nil
	case 40:
		return column.New[[40]byte]().Elem(arrayLevel, nullable, lc), nil
	case 41:
		return column.New[[41]byte]().Elem(arrayLevel, nullable, lc), nil
	case 42:
		return column.New[[42]byte]().Elem(arrayLevel, nullable, lc), nil
	case 43:
		return column.New[[43]byte]().Elem(arrayLevel, nullable, lc), nil
	case 44:
		return column.New[[44]byte]().Elem(arrayLevel, nullable, lc), nil
	case 45:
		return column.New[[45]byte]().Elem(arrayLevel, nullable, lc), nil
	case 46:
		return column.New[[46]byte]().Elem(arrayLevel, nullable, lc), nil
	case 47:
		return column.New[[47]byte]().Elem(arrayLevel, nullable, lc), nil
	case 48:
		return column.New[[48]byte]().Elem(arrayLevel, nullable, lc), nil
	case 49:
		return column.New[[49]byte]().Elem(arrayLevel, nullable, lc), nil
	case 50:
		return column.New[[50]byte]().Elem(arrayLevel, nullable, lc), nil
	case 51:
		return column.New[[51]byte]().Elem(arrayLevel, nullable, lc), nil
	case 52:
		return column.New[[52]byte]().Elem(arrayLevel, nullable, lc), nil
	case 53:
		return column.New[[53]byte]().Elem(arrayLevel, nullable, lc), nil
	case 54:
		return column.New[[54]byte]().Elem(arrayLevel, nullable, lc), nil
	case 55:
		return column.New[[55]byte]().Elem(arrayLevel, nullable, lc), nil
	case 56:
		return column.New[[56]byte]().Elem(arrayLevel, nullable, lc), nil
	case 57:
		return column.New[[57]byte]().Elem(arrayLevel, nullable, lc), nil
	case 58:
		return column.New[[58]byte]().Elem(arrayLevel, nullable, lc), nil
	case 59:
		return column.New[[59]byte]().Elem(arrayLevel, nullable, lc), nil
	case 60:
		return column.New[[60]byte]().Elem(arrayLevel, nullable, lc), nil
	case 61:
		return column.New[[61]byte]().Elem(arrayLevel, nullable, lc), nil
	case 62:
		return column.New[[62]byte]().Elem(arrayLevel, nullable, lc), nil
	case 63:
		return column.New[[63]byte]().Elem(arrayLevel, nullable, lc), nil
	case 64:
		return column.New[[64]byte]().Elem(arrayLevel, nullable, lc), nil
	case 65:
		return column.New[[65]byte]().Elem(arrayLevel, nullable, lc), nil
	case 66:
		return column.New[[66]byte]().Elem(arrayLevel, nullable, lc), nil
	case 67:
		return column.New[[67]byte]().Elem(arrayLevel, nullable, lc), nil
	case 68:
		return column.New[[68]byte]().Elem(arrayLevel, nullable, lc), nil
	case 69:
		return column.New[[69]byte]().Elem(arrayLevel, nullable, lc), nil
	case 70:
		return column.New[[70]byte]().Elem(arrayLevel, nullable, lc), nil
	}

	return nil, fmt.Errorf("fixed length %d is not supported", fixedLen)
}
