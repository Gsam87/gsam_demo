package controller

import (
	"errors"

	"github.com/qww83728/gsam_demo/domain/entity"
	repo_entity "github.com/qww83728/gsam_demo/domain/entity/repo"
	repo "github.com/qww83728/gsam_demo/domain/repository"
	cryptionSvc "github.com/qww83728/gsam_demo/domain/service/cryption"
)

type UserController interface {
	AddUser(user entity.AddUser) error
	ModifyUserPassword(user entity.ModifyUserPassword) error
	GetUserWithPassword(user entity.GetUser) (repo_entity.User, error)
}

type UserControllerImpl struct {
	cryptionService cryptionSvc.CryptionService
	userRepo        repo.UserRepo
}

func NewUserController(
	cryptionService cryptionSvc.CryptionService,
	userRepo repo.UserRepo,
) UserController {
	return &UserControllerImpl{
		cryptionService: cryptionService,
		userRepo:        userRepo,
	}
}

func (c *UserControllerImpl) AddUser(
	user entity.AddUser,
) error {
	// 確認使用者存在
	_, err := c.userRepo.GetUserByEmail(
		user.Email,
	)
	if err != entity.ErrNotFound {
		if err == nil {
			return errors.New("user email has been signup")
		}
		return err
	}

	// password encode
	hashpwd, err := c.cryptionService.BcryptEncode(user.Password)
	if err != nil {
		return err
	}

	if err := c.userRepo.AddUser(
		repo_entity.User{
			Email:    user.Email,
			Password: hashpwd,
		},
	); err != nil {
		return err
	}

	return nil
}

func (c *UserControllerImpl) ModifyUserPassword(
	user entity.ModifyUserPassword,
) error {
	// 確認使用者存在
	userGet, err := c.userRepo.GetUserByEmail(
		user.Email,
	)
	if err != nil {
		return err
	}
	// password encode
	if !c.cryptionService.BcryptCheck(userGet.Password, user.OldPassword) {
		return errors.New("password invalid")
	}

	// password encode
	hashpwd, err := c.cryptionService.BcryptEncode(user.NewPassword)
	if err != nil {
		return err
	}
	if err := c.userRepo.UpdateUserPassword(
		user.Email,
		hashpwd,
	); err != nil {
		return err
	}

	return nil
}

func (c *UserControllerImpl) GetUserWithPassword(
	user entity.GetUser,
) (repo_entity.User, error) {
	result, err := c.userRepo.GetUserByEmail(
		user.Email,
	)
	if err != nil {
		return repo_entity.User{}, err
	}

	// password encode
	if !c.cryptionService.BcryptCheck(result.Password, user.Password) {
		return repo_entity.User{}, errors.New("password invalid")
	}

	return result, nil
}
