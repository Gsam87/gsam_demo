package repo

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/qww83728/gsam_demo/domain/entity"
	repo_entity "github.com/qww83728/gsam_demo/domain/entity/repo"
)

type UserRepo interface {
	AddUser(user repo_entity.User) error
	UpdateUserPassword(email string, password string) error
	GetUserByEmail(email string) (repo_entity.User, error)
}

type UserRepoImpl struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return &UserRepoImpl{
		db: db,
	}
}

func (r *UserRepoImpl) AddUser(user repo_entity.User) error {

	var args []interface{}
	args = append(args, user.Email, user.Password)

	query := `INSERT INTO User (email, password, created, updated) VALUES (?, ?, NOW(), NOW())`
	_, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}
	fmt.Println("✅ 新增成功:", user.Email)

	return nil
}

func (r *UserRepoImpl) UpdateUserPassword(
	email string,
	password string,
) error {

	var args []interface{}
	args = append(args, password, email)

	query := `UPDATE User SET password=?, updated=NOW() WHERE email = ?`
	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	// 取得受影響的列數
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		// 沒有更新到任何資料
		return fmt.Errorf("no rows updated")
	}
	fmt.Println("✅ 修改成功:", email)

	return nil
}

func (r *UserRepoImpl) GetUserByEmail(
	email string,
) (repo_entity.User, error) {
	var user repo_entity.User
	var args []interface{}
	args = append(args, email)

	query := `SELECT email, password, created, updated FROM User WHERE email = ?`
	if err := r.db.Get(&user, query, args...); err != nil {
		if err == sql.ErrNoRows {
			return repo_entity.User{}, entity.ErrNotFound
		}
		return repo_entity.User{}, err
	}

	return user, nil
}
