package mariadb

import "github.com/jigadhirasu/qtool/types"

type Result struct {
	Version      string       `json:",omitempty"` // 版本號
	Timestamp    int          `json:",omitempty"` // 伺服器時間戳
	UUID         string       `json:",omitempty"`
	RowsAffected int64        `json:",omitempty"`
	Error        *types.Error `json:",omitempty"`
	Length       int64        `json:",omitempty"`
	Data         types.Bytes  `json:",omitempty"`
}

func Sum[T int | float64](a, b T) T {
	return a + b
}
