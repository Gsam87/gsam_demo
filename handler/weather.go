package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/qww83728/gsam_demo/controller"
	"github.com/qww83728/gsam_demo/util"
)

type WeatherHandler struct {
	WeatherController controller.WeatherController
}

func NewWeatherHandler(
	weatherController controller.WeatherController,
) *WeatherHandler {
	return &WeatherHandler{
		WeatherController: weatherController,
	}
}

func (ctrl *WeatherHandler) GetTodayInfo(c *gin.Context) {

	results, err := ctrl.WeatherController.GetTodayInfo()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.MakeFailResponse(
			http.StatusInternalServerError,
			"Internal Server Error",
			err,
		))
	}

	switch len(results) {
	case 0:
		c.JSON(http.StatusNotFound, util.MakeFailResponse(
			http.StatusNotFound,
			"無天氣資訊，可能原因 : F-C0032-001 不提供過去時間資訊(ex.6-12, 12-18)",
			err,
		))
	case 1:
		c.JSON(http.StatusOK, util.MakeSuceessResponseWithMsg(
			http.StatusOK,
			"缺失資料，可能原因 : F-C0032-001 不提供過去時間資訊(ex.6-12, 12-18)",
			results,
		))
	default:
		c.JSON(http.StatusOK, util.MakeSuceessResponse(http.StatusOK, results))
	}
}
