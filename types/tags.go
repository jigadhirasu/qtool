package types

import "strconv"

type Tags map[string]string

func (tg Tags) String(key string) string {
	return tg[key]
}
func (tg Tags) Bytes(key string) []byte {
	return []byte(tg.String(key))
}
func (tg Tags) Int(key string) int {
	i, _ := strconv.Atoi(tg.String(key))
	return i
}
func (tg Tags) Int64(key string) int64 {
	i, _ := strconv.ParseInt(tg.String(key), 10, 64)
	return i
}
func (tg Tags) Float64(key string) float64 {
	f, _ := strconv.ParseFloat(tg.String(key), 64)
	return f
}
