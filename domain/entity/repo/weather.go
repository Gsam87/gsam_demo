package repo_entity

import "time"

// User struct 對應到資料表欄位
type Weather struct {
	ID      int       `db:"id"`
	City    string    `db:"city"`
	MinT    float64   `db:"min_t"`
	MaxT    float64   `db:"max_t"`
	Period  string    `db:"period"`
	Date    string    `db:"date"`
	Created time.Time `db:"created"`
	Updated time.Time `db:"updated"`
}
