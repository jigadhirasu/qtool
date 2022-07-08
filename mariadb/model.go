package mariadb

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/jigadhirasu/qtool/types"
	"gorm.io/gorm"
)

type Model interface {
	TableName() string
	UU(uuid ...string) string
	Owner() string
}

type Instance interface {
	Model
	Type() Model
}

// 有紀錄
func Create(tags types.Tags, m Model) func(db *gorm.DB) types.Bytes {
	r := &Record{
		OpID:     tags.String("User"),
		OwnerID:  m.Owner(),
		Method:   "create",
		Target:   m.TableName(),
		TargetID: m.UU(),
		OpBefore: []byte(`{}`),
		OpAfter:  types.JSON(m),
	}
	pk := &Pack{Doc: types.JSON(m)}
	return func(db *gorm.DB) types.Bytes {
		tx := db.Table(m.TableName()).Create(pk)
		if tx.Error != nil {
			return types.JSON(Result{Error: types.DBErr(tx.Error)})
		}
		db.Table(r.TableName()).Create(&Pack{Doc: types.JSON(r)})
		return types.JSON(Result{
			UUID:  m.UU(),
			Error: types.DBErr(tx.Error),
		})
	}
}

// 沒紀錄
func Createx(tags types.Tags, m Model) func(db *gorm.DB) types.Bytes {
	pk := &Pack{Doc: types.JSON(m)}
	return func(db *gorm.DB) types.Bytes {
		tx := db.Table(m.TableName()).Create(pk)
		if tx.Error != nil {
			return types.JSON(Result{Error: types.DBErr(tx.Error)})
		}
		return types.JSON(Result{
			UUID:  m.UU(),
			Error: types.DBErr(tx.Error),
		})
	}
}

func Delete(tags types.Tags, m Model) func(db *gorm.DB) types.Bytes {
	r := &Record{
		OpID:     tags.String("User"),
		OwnerID:  m.Owner(),
		Method:   "delete",
		Target:   m.TableName(),
		TargetID: m.UU(),
		OpAfter:  []byte(`{}`),
	}
	return func(db *gorm.DB) types.Bytes {
		r.OpBefore = C{Table: m.TableName(), UUID: m.UU()}.Value(db)
		tx := db.Table(m.TableName()).Where("UUID = ?", m.UU()).Delete(m)
		db.Table(r.TableName()).Create(&Pack{Doc: types.JSON(r)})
		return types.JSON(Result{
			UUID:  m.UU(),
			Error: types.DBErr(tx.Error),
		})
	}
}

// return Result
func Update(tags types.Tags, m Model, fields ...string) func(db *gorm.DB) types.Bytes {
	r := &Record{
		OpID:     tags.String("User"),
		OwnerID:  m.Owner(),
		Method:   "update",
		Field:    strings.Join(fields, ","),
		Target:   m.TableName(),
		TargetID: m.UU(),
	}

	oldPtr := reflect.New(reflect.Indirect(reflect.ValueOf(m)).Type())
	o := oldPtr.Interface()
	newPtr := reflect.New(reflect.Indirect(reflect.ValueOf(m)).Type())
	n := newPtr.Interface()
	if len(fields) > 0 {
		vm := reflect.Indirect(reflect.ValueOf(m))
		vn := reflect.Indirect(newPtr)
		for i := 0; i < vm.NumField(); i++ {
			for _, field := range fields {
				if strings.EqualFold(strings.ToLower(field), strings.ToLower(vm.Type().Field(i).Name)) {
					vn.Field(i).Set(vm.Field(i))
				}
			}
		}
	} else {
		types.JSON(m).Decode(n)
	}

	return func(db *gorm.DB) types.Bytes {
		// *** 檢查異動欄位 ***
		v := C{Table: m.TableName(), UUID: m.UU()}.Value(db)
		if len(v) < 1 {
			return types.JSON(Result{Error: types.NewErr(520, "data not found")})
		}
		v.Decode(o)
		SQL := fmt.Sprintf("UPDATE %s SET Doc = JSON_MERGE_PATCH(Doc, ?) WHERE UUID = ?", m.TableName())
		tx := db.Exec(SQL, types.JSON(n), m.UU())
		if err := diff(o, n); err != nil {
			return types.JSON(Result{Error: types.DBErr(err)})
		}
		r.OpBefore = types.JSON(o)
		r.OpAfter = types.JSON(n)
		// *** 僅記錄有異動的欄位 ***
		if tx.RowsAffected < 1 {
			return types.JSON(Result{
				Error: types.NewErr(520, "no changed"),
			})
		}
		db.Table(r.TableName()).Create(&Pack{Doc: types.JSON(r)})
		return types.JSON(Result{
			RowsAffected: tx.RowsAffected,
			Error:        types.DBErr(tx.Error),
		})
	}
}

// return Result
func Push(tags types.Tags, m Model, field string) func(db *gorm.DB) types.Bytes {
	r := &Record{
		OpID:     tags.String("User"),
		OwnerID:  m.Owner(),
		Method:   "push",
		Field:    field,
		Target:   m.TableName(),
		TargetID: m.UU(),
	}

	vm := reflect.ValueOf(m)
	depth := strings.Split(field, ".")
	f := vm
	for _, d := range depth {
		f = reflect.Indirect(f).FieldByName(d)
	}

	newPtr := reflect.New(reflect.Indirect(reflect.ValueOf(m)).Type())

	return func(db *gorm.DB) types.Bytes {
		SQL := fmt.Sprintf(`UPDATE %s SET doc = JSON_SET(doc, '$.%s', COALESCE(JSON_VALUE(doc, '$.%s'), 0) + ?) WHERE UUID = ?`, m.TableName(), field, field)
		tx := db.Exec(SQL, f.Interface(), m.UU())
		v := C{Table: m.TableName(), UUID: m.UU(), Field: []string{"UUID"}, Attr: []string{field}}.Value(db)
		if len(v) < 1 {
			return types.JSON(Result{Error: types.NewErr(520, "data not found")})
		}
		v.Decode(newPtr.Interface())
		r.OpAfter = types.JSON(newPtr.Interface())

		oldPtr := reflect.New(reflect.Indirect(reflect.ValueOf(m)).Type())
		n := newPtr
		o := oldPtr
		for _, d := range depth {
			n = reflect.Indirect(n).FieldByName(d)
			o = reflect.Indirect(o).FieldByName(d)
		}
		if n.Int() < 0 {
			return types.JSON(Result{Error: types.NewErr(400, "after calculate need greater than or equal to 0")})
		}
		o.SetInt(n.Int() - f.Int())
		r.OpBefore = types.JSON(oldPtr.Interface())

		db.Table(r.TableName()).Create(&Pack{Doc: types.JSON(r)})
		return types.JSON(Result{
			RowsAffected: tx.RowsAffected,
			Error:        types.DBErr(tx.Error),
		})
	}
}

// // return Result
// func Find(m IQuery) func(db *gorm.DB) types.Bytes {
// 	return func(db *gorm.DB) types.Bytes {
// 		tx := m.Where(db).Table(m.TableName())

// 		var count int64
// 		tx.Count(&count)

// 		bb := [][]byte{}
// 		tx = m.OrderBy(tx)
// 		tx = m.Limit(tx)
// 		tx = tx.Pluck(`JSON_MERGE_PATCH(Doc, JSON_OBJECT('CreatedAt',CreatedAt,'UpdatedAt',UpdatedAt)) as Doc`, &bb)

// 		b := bytes.Join(bb, []byte(","))
// 		return types.JSON(Result{
// 			Data:   bytes.Join([][]byte{[]byte("["), b, []byte("]")}, []byte("")),
// 			Length: count,
// 			Error:  types.DBErr(tx.Error),
// 		})
// 	}
// }

// // return Result
// func Findx(m IQuery) func(db *gorm.DB) types.Bytes {
// 	return func(db *gorm.DB) types.Bytes {
// 		tx := m.Where(db).Table(m.TableName())

// 		var count int64
// 		tx.Count(&count)

// 		bb := [][]byte{}
// 		tx = m.OrderBy(tx)
// 		tx = m.Limit(tx)
// 		tx = tx.Pluck(`Doc`, &bb)

// 		b := bytes.Join(bb, []byte(","))
// 		return types.JSON(Result{
// 			Data:   bytes.Join([][]byte{[]byte("["), b, []byte("]")}, []byte("")),
// 			Length: count,
// 			Error:  types.DBErr(tx.Error),
// 		})
// 	}
// }
