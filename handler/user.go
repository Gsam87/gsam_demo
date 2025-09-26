package handler

import (
	"net/http"

	"github.com/qww83728/gsam_demo/controller"
	"github.com/qww83728/gsam_demo/domain/entity"
	middlerware "github.com/qww83728/gsam_demo/interface/middleware"
	"github.com/qww83728/gsam_demo/util"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	Middlerware    middlerware.Middleware
	UserController controller.UserController
}

func NewUserHandler(
	middlerware middlerware.Middleware,
	userController controller.UserController,
) *UserHandler {
	return &UserHandler{
		Middlerware:    middlerware,
		UserController: userController,
	}
}

type AddUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (ctrl *UserHandler) AddUser(c *gin.Context) {
	var input AddUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.UserController.AddUser(
		entity.AddUser{
			Email:    input.Email,
			Password: input.Password,
		},
	); err != nil {
		c.JSON(http.StatusInternalServerError, util.MakeFailResponse(
			http.StatusInternalServerError,
			"Internal Server Error",
			err,
		))
		return
	}

	c.JSON(http.StatusOK, util.MakeSuceessResponse(http.StatusOK, nil))
}

type ModifyUserPassword struct {
	Email       string `json:"email"`
	NewPassword string `json:"newPassword"`
	OldPassword string `json:"oldPassword"`
}

func (ctrl *UserHandler) ModifyUserPassword(c *gin.Context) {
	var input ModifyUserPassword
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.UserController.ModifyUserPassword(
		entity.ModifyUserPassword{
			Email:       input.Email,
			OldPassword: input.OldPassword,
			NewPassword: input.NewPassword,
		},
	); err != nil {
		if err == entity.ErrNotFound {
			// 是否透漏訊息給使用者可再討論
			c.JSON(http.StatusBadRequest, util.MakeFailResponse(
				http.StatusBadRequest,
				err.Error(),
				err,
			))
		} else {
			c.JSON(http.StatusInternalServerError, util.MakeFailResponse(
				http.StatusInternalServerError,
				"Internal Server Error",
				err,
			))
		}
		return
	}

	c.JSON(http.StatusOK, util.MakeSuceessResponse(http.StatusOK, nil))
}

type LoginInfo struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginData struct {
	Jwt string `json:"jwt"`
}

func (ctrl *UserHandler) Login(c *gin.Context) {
	var input LoginInfo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := ctrl.UserController.GetUserWithPassword(
		entity.GetUser{
			Email:    input.Email,
			Password: input.Password,
		},
	)
	if err != nil {
		if err == entity.ErrNotFound {
			// 是否透漏訊息給使用者可再討論
			c.JSON(http.StatusNotFound, util.MakeFailResponse(
				http.StatusNotFound,
				err.Error(),
				err,
			))
		} else {
			c.JSON(http.StatusInternalServerError, util.MakeFailResponse(
				http.StatusInternalServerError,
				"Internal Server Error",
				err,
			))
		}
		return
	}

	jwt, err := ctrl.Middlerware.GenerateToken(
		result.Email,
		result.Updated,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.MakeFailResponse(
			http.StatusInternalServerError,
			"Generate Jwt Token Error",
			err,
		))
	}

	c.JSON(http.StatusOK, util.MakeSuceessResponse(http.StatusOK, LoginData{Jwt: jwt}))
}
