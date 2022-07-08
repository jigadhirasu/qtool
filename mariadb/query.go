package mariadb

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

type IQuery interface {
	TableName() string
	Where(tx *gorm.DB) *gorm.DB
	OrderBy(tx *gorm.DB) *gorm.DB
	Limit(tx *gorm.DB) *gorm.DB
}

type Query struct {
	Asc   []string `gorm:"-"`
	Desc  []string `gorm:"-"`
	Size  int      `gorm:"-"`
	Index int      `gorm:"-"`
}

func (q Query) OrderBy(tx *gorm.DB) *gorm.DB {
	sortStr := []string{}
	sortStr = append(sortStr, q.Asc...)
	for _, desc := range q.Desc {
		sortStr = append(sortStr, fmt.Sprintf("%s DESC", desc))
	}
	if len(sortStr) > 0 {
		tx = tx.Order(strings.Join(sortStr, ", "))
	}
	return tx
}

func (q Query) Limit(tx *gorm.DB) *gorm.DB {
	size := q.Size
	if size == 0 {
		size = 10
	}
	index := q.Index
	return tx.Limit(size).Offset(index * size)
}
