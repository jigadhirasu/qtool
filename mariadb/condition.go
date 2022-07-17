package mariadb

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jigadhirasu/qtool/types"
	"gorm.io/gorm"
)

func column(field ...string) string {
	cols := []string{}
	for _, f := range field {
		if !strings.Contains(f, "$") {
			cols = append(cols, fmt.Sprintf(`'$.%s'`, f))
			continue
		}
		cols = append(cols, fmt.Sprintf(`'%s'`, f))
	}

	if len(cols) < 1 {
		return ""
	}
	return fmt.Sprintf(`JSON_EXTRACT(Doc, %s)`, strings.Join(cols, ", "))
}

type C struct {
	Table string   // 資料表明
	UUID  string   // 唯一編號
	Field []string // 有索引的屬性
	Attr  []string // 沒索引的屬性
	Index KV       // 有索引的條件
	// Condition KV       // 沒索引條件
}

func (c C) Select() string {
	col := "Doc"
	switch len(c.Field) + len(c.Attr) {
	case 0:
	case 1:
		col = strings.Join(c.Field, ",") + column(c.Attr...)
	default:
		ss := []string{}
		for _, f := range c.Field {
			ss = append(ss, fmt.Sprintf("'%s'", f), column(f))
		}
		for _, f := range c.Attr {
			ss = append(ss, fmt.Sprintf("'%s'", f), column(f))
		}
		col = fmt.Sprintf(`JSON_OBJECT(%s)`, strings.Join(ss, ", "))
	}
	return col
}

func (c C) Where(tx *gorm.DB) *gorm.DB {
	tx = tx.Table(c.Table)
	if c.UUID != "" {
		tx = tx.Where("UUID = ?", c.UUID)
	}
	// 僅索引可以搜尋
	// for k, v := range c.Index {
	// 	tx = tx.Where(k+" = ?", v)
	// }
	if c.Index != nil {
		tx = c.Index.Where(tx)
	}
	return tx
}

func (c C) Count(tx *gorm.DB) int64 {
	tx = c.Where(tx)
	var count int64
	tx.Count(&count)
	return count
}

func (c C) Value(tx *gorm.DB) types.Bytes {
	tx = c.Where(tx)
	bb := [][]byte{}
	tx.Limit(1).Pluck(c.Select(), &bb)
	if len(bb) > 0 {
		return bb[0]
	}
	return types.Bytes{}
}

func (c C) Values(tx *gorm.DB) types.Bytes {
	tx = c.Where(tx)
	// if c.Query != nil {
	// 	tx = c.Query.Limit(tx)
	// 	tx = c.Query.OrderBy(tx)
	// }
	bb := [][]byte{}
	tx.Limit(1).Pluck(fmt.Sprintf(`JSON_ARRAYAGG(%s)`, c.Select())+" as Doc", &bb)
	if len(bb) > 0 {
		return bb[0]
	}
	return types.Bytes{}
}

type KV map[string]interface{}

func (kv KV) Where(tx *gorm.DB) *gorm.DB {

	isDep, _ := regexp.Compile("^\\$\\.")

	for k, v := range kv {
		if ok := isDep.MatchString(k); ok {
			k = fmt.Sprintf("JSON_VALUE(Doc,'%s')", k)
		}

		switch t := v.(type) {
		case int:
			tx = tx.Where(k+" = ?", t)
		case []string:
			tx = tx.Where(k+" IN ?", t)
		case string:
			if p := kv.EQ(t); p != nil {
				tx = tx.Where(k+" = ?", p)
			}
			if p := kv.LIKE(t); p != nil {
				tx = tx.Where(k+" LIKE ?", p)
			}
			if p := kv.LT(t); p != nil {
				tx = tx.Where(k+" < ?", p)
			}
			if p := kv.LTE(t); p != nil {
				tx = tx.Where(k+" <= ?", p)
			}
			if p := kv.GT(t); p != nil {
				tx = tx.Where(k+" > ?", p)
			}
			if p := kv.GTE(t); p != nil {
				tx = tx.Where(k+" >= ?", p)
			}
			if p := kv.IN(t); p != nil {
				tx = tx.Where(k+" IN ?", p)
			}
			if p := kv.NOT(t); p != nil {
				tx = tx.Where(k+" NOT ?", p)
			}
		}
	}

	return tx
}

func (KV) EQ(str string) interface{} {
	if ok, _ := regexp.MatchString(":", str); !ok {
		return str
	}
	if ok, _ := regexp.MatchString("^EQ:", strings.ToUpper(str)); ok {
		return str[3:]
	}
	return nil
}
func (KV) LIKE(str string) interface{} {
	if ok, _ := regexp.MatchString("^LIKE:", strings.ToUpper(str)); ok {
		str = strings.ReplaceAll(str, "$", "%")
		return str[5:]
	}
	return nil
}

func (KV) LT(str string) interface{} {
	if ok, _ := regexp.MatchString("^LT:", strings.ToUpper(str)); ok {
		return str[3:]
	}
	return nil
}
func (KV) LTE(str string) interface{} {
	if ok, _ := regexp.MatchString("^LTE:", strings.ToUpper(str)); ok {
		return str[4:]
	}
	return nil
}
func (KV) GT(str string) interface{} {
	if ok, _ := regexp.MatchString("^GT:", strings.ToUpper(str)); ok {
		return str[3:]
	}
	return nil
}
func (KV) GTE(str string) interface{} {
	if ok, _ := regexp.MatchString("^GTE:", strings.ToUpper(str)); ok {
		return str[4:]
	}
	return nil
}
func (KV) IN(str string) interface{} {
	if ok, _ := regexp.MatchString("^IN:", strings.ToUpper(str)); ok {
		str = str[3:]
		return strings.Split(str, ",")
	}
	return nil
}
func (KV) NOT(str string) interface{} {
	if ok, _ := regexp.MatchString("^NOT:", strings.ToUpper(str)); ok {
		return str[4:]
	}
	return nil
}
