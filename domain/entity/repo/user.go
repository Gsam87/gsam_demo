package repo_entity

import "time"

// User struct 對應到資料表欄位
type User struct {
	Email    string    `db:"email"`
	Password string    `db:"password"`
	Created  time.Time `db:"created"`
	Updated  time.Time `db:"updated"`
}
