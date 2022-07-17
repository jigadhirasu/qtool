package mariadb

import (
	"time"

	"github.com/jigadhirasu/qtool/types"
)

type Pack struct {
	Doc       types.Bytes `gorm:"column:Doc;"`
	ID        int64       `gorm:"column:ID;"`
	CreatedAt time.Time   `gorm:"column:CreatedAt; default:CURRENT_TIMESTAMP;"`
	UpdatedAt time.Time   `gorm:"column:UpdatedAt; default:CURRENT_TIMESTAMP;"`
}
