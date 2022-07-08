package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func JSON(v interface{}) Bytes {
	jb, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return Bytes(jb)
}

func STRUCT(jb []byte, scan interface{}) error {
	return json.Unmarshal(jb, scan)
}

type Bytes []byte

// MarshalJSON returns m as the JSON encoding of m.
func (v Bytes) MarshalJSON() ([]byte, error) {
	if v == nil {
		return []byte("null"), nil
	}
	return v, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (v *Bytes) UnmarshalJSON(data []byte) error {
	if v == nil {
		return errors.New("json.RawMessage: UnmarshalJSON on nil pointer")
	}
	*v = append((*v)[0:0], data...)
	return nil
}

func (v Bytes) Decode(scan interface{}) error {
	return STRUCT(v, scan)
}
func (v Bytes) Bytes() []byte {
	return []byte(v)
}
func (v Bytes) String() string {
	return string(v)
}
func (v Bytes) Int() int {
	i, _ := strconv.Atoi(v.String())
	return i
}

// Bytesalue return json value, implement driver.Bytesaluer interface
func (v Bytes) Bytesalue() (driver.Value, error) {
	return json.Marshal(v)
}

// Scan scan value into Jsonb, implements sql.Scanner interface
func (v *Bytes) Scan(value interface{}) error {
	if value == nil {
		*v = nil
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	scan := []byte{}
	err := json.Unmarshal(bytes, &scan)
	*v = Bytes(scan)
	return err
}

// GormDataType gorm common data type
func (Bytes) GormDataType() string {
	return "json"
}

// GormDBDataType gorm db data type
func (Bytes) GormDBDataType(db *gorm.DB, field *schema.Field) string {
	switch db.Dialector.Name() {
	case "sqlite":
		return "JSON"
	case "mysql":
		return "JSON"
	case "postgres":
		return "JSONB"
	}
	return ""
}
