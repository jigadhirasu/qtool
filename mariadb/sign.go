package mariadb

import (
	"github.com/jigadhirasu/qtool/types"
	"gorm.io/gorm"
)

type SignShell struct {
	UserID   string `gorm:"column:UserID; index; type:varchar(40) AS (COALESCE(JSON_VALUE(Doc, '$.UserID'), ''))"`
	UserType string `gorm:"column:UserType; type:varchar(40) AS (COALESCE(JSON_VALUE(Doc, '$.UserType'), ''))"`
	Method   string `gorm:"column:Method; type:varchar(40) AS (COALESCE(JSON_VALUE(Doc, '$.Method'), ''))"`
	Address  string `gorm:"column:Address; type:varchar(40) AS (COALESCE(JSON_VALUE(Doc, '$.Address'), ''))"`
	Region   string `gorm:"column:Region; type:varchar(40) AS (COALESCE(JSON_VALUE(Doc, '$.Region'), ''))"`
	City     string `gorm:"column:City; type:varchar(40) AS (COALESCE(JSON_VALUE(Doc, '$.City'), ''))"`
	Pack
}

func (SignShell) TableName() string {
	return Sign{}.TableName()
}

type Sign struct {
	UserID    string
	UserType  string // iadmin iagent islot
	UserAgent string
	Address   string
	Region    string       `json:",omitempty"` // 區域
	City      string       `json:",omitempty"` // 城市
	Method    string       `json:",omitempty"` // in password otp
	Reason    *types.Error `json:",omitempty"` // 錯誤訊息
	SignAt    *types.Time  `json:",omitempty"`
	SignOut   *types.Time  `json:",omitempty"`
}

func (a *Sign) UU(uuid ...string) string {
	return ""
}

func (Sign) TableName() string {
	return "signs"
}

type SignQuery struct {
	UserID   string
	UserType string
	Query
}

func (SignQuery) TableName() string {
	return Sign{}.TableName()
}

func (m SignQuery) Where(tx *gorm.DB) *gorm.DB {
	if m.UserID != "" {
		tx = tx.Where("UserID = ?", m.UserID)
	}
	if m.UserType != "" {
		tx = tx.Where("UserType = ?", m.UserType)
	}
	return tx
}
