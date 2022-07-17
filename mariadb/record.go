package mariadb

import (
	"bytes"
	"encoding/json"
	"reflect"

	"github.com/jigadhirasu/qtool/types"
)

type Record struct {
	OpID     string `json:",omitempty"`
	OwnerID  string `json:",omitempty"`
	TargetID string `json:",omitempty"`
	Target   string `json:",omitempty"`
	Method   string `json:",omitempty"`
	Field    string `json:",omitempty"`
	OpBefore types.Bytes
	OpAfter  types.Bytes
}

func (Record) TableName() string {
	return "records"
}

func (a *Record) UU(uuid ...string) string {
	return ""
}

type RecordShell struct {
	OpID     string   `gorm:"column:OpID; index; type:varchar(40) AS (JSON_VALUE(Doc, '$.OpID'))"`
	OwnerID  string   `gorm:"column:OwnerID; index; type:varchar(40) AS (JSON_VALUE(Doc, '$.OwnerID'))"`
	TargetID string   `gorm:"column:TargetID; index(target); type:varchar(40) AS (JSON_VALUE(Doc, '$.TargetID'))"`
	Target   string   `gorm:"column:Target; index(target); type:varchar(40) AS (JSON_VALUE(Doc, '$.Target'))"`
	Targets  []string `gorm:"-"`
	Method   string   `gorm:"column:Method; type:varchar(40) AS (JSON_VALUE(Doc, '$.Method'))"`
	Field    string   `gorm:"column:Field; type:varchar(40) AS (JSON_VALUE(Doc, '$.Field'))"`
	Pack
	Query
}

func (RecordShell) TableName() string {
	return Record{}.TableName()
}

// func (m RecordShell) Where(tx *gorm.DB) *gorm.DB {
// 	if m.OpID != "" {
// 		tx = tx.Where("OpID = ?", m.OpID)
// 	}
// 	if m.OwnerID != "" {
// 		tx = tx.Where("OwnerID = ?", m.OwnerID)
// 	}
// 	if m.TargetID != "" {
// 		tx = tx.Where("TargetID = ?", m.TargetID)
// 	}
// 	if m.Target != "" {
// 		tx = tx.Where("Target = ?", m.Target)
// 	}
// 	if len(m.Targets) > 0 {
// 		tx = tx.Where("Target IN ?", m.Targets)
// 	}
// 	if m.Method != "" {
// 		tx = tx.Where("JSON_VALUE(DOC, '$.Method') = ?", m.Method)
// 	}
// 	if m.Field != "" {
// 		tx = tx.Where("Field = ?", m.Field)
// 	}
// 	return tx
// }

func diff(old, new interface{}) error {

	o := reflect.ValueOf(old)
	n := reflect.ValueOf(new)
	if o.Kind() != reflect.Ptr {
		return types.NewErr(484, "diff target is not ptr")
	}
	if o.Kind() != n.Kind() {
		return types.NewErr(484, "diff target is not ptr")
	}

	// log4.Debug("old : ", types.JSON(o.Interface()))
	// log4.Debug("new : ", types.JSON(n.Interface()))

	diffStruct(o, n)
	return nil
}

func diffStruct(old, new reflect.Value) {
	o := reflect.Indirect(old)
	n := reflect.Indirect(new)

	for i := 0; i < n.NumField(); i++ {
		name := n.Type().Field(i).Name
		nf := n.Field(i)
		of := o.FieldByName(name)

		if !of.IsValid() {
			of.Set(reflect.Zero(nf.Type()))
		}

		if nf.Kind() == reflect.Struct {
			diffStruct(of, nf)
			continue
		}
		if nf.Kind() == reflect.Map {
			diffMap(of, nf)
			continue
		}

		if of.CanSet() {
			IsCompare := func() bool {
				switch nf.Interface().(type) {
				case []byte, json.RawMessage, types.Bytes:
					return bytes.Equal(nf.Interface().(types.Bytes), of.Interface().(types.Bytes))
				case int, int8, int64, string, float64:
					return nf.Interface() == of.Interface()
				default:
					return bytes.Equal(types.JSON(nf.Interface()), types.JSON(of.Interface()))
				}
			}()
			if IsCompare {
				of.Set(reflect.Zero(of.Type()))
				nf.Set(reflect.Zero(nf.Type()))
			}
		}
	}
}

func diffMap(old, new reflect.Value) {
	o := reflect.Indirect(old)
	n := reflect.Indirect(new)
	// set := reflect.MakeMap(n.Type())
	iter := n.MapRange()
	for iter.Next() {
		key, nf := iter.Key(), iter.Value()
		of := o.MapIndex(key)

		if !of.IsValid() {
			of.Set(reflect.Zero(nf.Type()))
		}

		if nf.Kind() == reflect.Struct {
			diffStruct(of, nf)
			continue
		}
		if nf.Kind() == reflect.Map {
			diffMap(of, nf)
			continue
		}

		if of.CanSet() {
			IsCompare := func() bool {
				switch nf.Interface().(type) {
				case []byte, json.RawMessage, types.Bytes:
					return bytes.Equal(nf.Interface().(types.Bytes), of.Interface().(types.Bytes))
				case []string:
					return bytes.Equal(types.JSON(nf.Interface()), types.JSON(of.Interface()))
				default:
					return nf.Interface() == of.Interface()
				}
			}()
			if IsCompare {
				of.Set(reflect.Zero(of.Type()))
				nf.Set(reflect.Zero(nf.Type()))
			}
		}
	}
}
