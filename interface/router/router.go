package router

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/qww83728/gsam_demo/controller"
	repo "github.com/qww83728/gsam_demo/domain/repository"
	cryptionSvc "github.com/qww83728/gsam_demo/domain/service/cryption"
	"github.com/qww83728/gsam_demo/handler"
	middlerware "github.com/qww83728/gsam_demo/interface/middleware"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func Router(r *gin.Engine) {
	dsn := "root:root@tcp(127.0.0.1:3306)/testdb?parseTime=true"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln("連線資料庫失敗:", err)
	}
	// defer db.Close()
	fmt.Println("✅ MySQL 連線成功 (sqlx)")

	middlerware := middlerware.NewMiddleware()

	userRepo := repo.NewUserRepo(db)
	weatherRepo := repo.NewWeatherRepo(db)

	cryptionService := cryptionSvc.NewCryptionService() // 初始化 service

	// controller
	userController := controller.NewUserController(cryptionService, userRepo)
	weatherController := controller.NewWeatherController(weatherRepo)

	// handler
	userHandler := handler.NewUserHandler(middlerware, userController)
	weatherHandler := handler.NewWeatherHandler(weatherController)

	r.GET("/hello/", middlerware.JWTMiddleware(), HelloWorld)
	// r.POST("/login", func(c *gin.Context) {
	// 	// 假設 login 成功
	// 	token, _ := middlerware.GenerateToken("12345")
	// 	c.JSON(http.StatusOK, gin.H{"token": token})
	// })

	userGroup := r.Group("/user")
	{
		userGroup.POST("", userHandler.AddUser)
		userGroup.PATCH("/pwd", userHandler.ModifyUserPassword)
		userGroup.POST("/login", userHandler.Login)
	}

	weatherGroup := r.Group("/weather")
	{
		weatherGroup.GET("", weatherHandler.GetTodayInfo)
		weatherGroup.GET("/auth", middlerware.JWTMiddleware(), weatherHandler.GetTodayInfo)
	}

}

var balance = 1000

func HelloWorld(context *gin.Context) {
	var msg = "您的帳戶內有:" + strconv.Itoa(balance) + "元"
	context.JSON(http.StatusOK, gin.H{
		"amount":  balance,
		"status":  "ok",
		"message": msg,
	})

}
