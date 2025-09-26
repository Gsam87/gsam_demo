package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type WeatherResponse struct {
	Records Records `json:"records"`
}

type Records struct {
	DatasetDescription string     `json:"datasetDescription"`
	Location           []Location `json:"location"`
}

type Location struct {
	LocationName   string           `json:"locationName"`
	WeatherElement []WeatherElement `json:"weatherElement"`
}

type WeatherElement struct {
	ElementName string `json:"elementName"`
	Time        []Time `json:"time"`
}

type Time struct {
	StartTime string    `json:"startTime"`
	EndTime   string    `json:"endTime"`
	Parameter Parameter `json:"parameter"`
}

type Parameter struct {
	ParameterName string `json:"parameterName"`
	ParameterUnit string `json:"parameterUnit"`
}

// 定義 enum 值
type Period string

const (
	AM Period = "AM"
	PM Period = "PM"
)

type Element string

const (
	MinT Element = "MinT"
	MaxT Element = "MaxT"
)

type LocationInfo struct {
	LocationName string
	MinT         string
	MaxT         string
	Date         string
}

func main() {
	today := time.Now().Format("2006-01-02")
	amTimeStr := "&startTime=" + today + "T06%3A00%3A00"
	pmTimeStr := "&startTime=" + today + "T12%3A00%3A00"
	getPMWeather := true // 顆粒為 12 小時, AM 6:00 有資料則會包含道 18:00，則不需要再去抓 PM 資料

	var amWeather LocationInfo
	amWeatherMinT := getWeather(MinT, amTimeStr)
	if len(amWeatherMinT.Records.Location) != 0 &&
		len(amWeatherMinT.Records.Location[0].WeatherElement) != 0 &&
		len(amWeatherMinT.Records.Location[0].WeatherElement[0].Time) != 0 {
		// 有資料
		amWeather.LocationName = amWeatherMinT.Records.Location[0].WeatherElement[0].ElementName
		amWeather.MinT = amWeatherMinT.Records.Location[0].WeatherElement[0].Time[0].Parameter.ParameterName

		// MaxT
		amWeatherMaxT := getWeather(MaxT, amTimeStr)
		amWeather.MaxT = amWeatherMaxT.Records.Location[0].WeatherElement[0].Time[0].Parameter.ParameterName

		// save
		amWeather.Date = today
		saveWeather(AM, amWeather)

		// 顆粒為 12 小時, AM 6:00 有資料則會包含道 18:00，下午視為同一溫度
		var pmWeather LocationInfo
		pmWeather.LocationName = amWeatherMinT.Records.Location[0].WeatherElement[0].ElementName
		pmWeather.MinT = amWeather.MinT
		pmWeather.MaxT = amWeather.MaxT
		// save
		pmWeather.Date = today
		saveWeather(PM, pmWeather)
		getPMWeather = false
	}

	if getPMWeather {
		var pmWeather LocationInfo
		pmWeatherMinT := getWeather(MinT, pmTimeStr)
		if len(pmWeatherMinT.Records.Location) != 0 &&
			len(pmWeatherMinT.Records.Location[0].WeatherElement) != 0 &&
			len(pmWeatherMinT.Records.Location[0].WeatherElement[0].Time) != 0 {
			// 有資料
			pmWeather.LocationName = pmWeatherMinT.Records.Location[0].WeatherElement[0].ElementName
			pmWeather.MinT = pmWeatherMinT.Records.Location[0].WeatherElement[0].Time[0].Parameter.ParameterName

			// MaxT
			pmWeatherMaxT := getWeather(MaxT, pmTimeStr)
			pmWeather.MaxT = pmWeatherMaxT.Records.Location[0].WeatherElement[0].Time[0].Parameter.ParameterName

			// save
			pmWeather.Date = today
			saveWeather(PM, pmWeather)
		}
	}
}

func getWeather(
	element Element,
	timeStr string,
) WeatherResponse {
	url := "https://opendata.cwa.gov.tw/api/v1/rest/datastore/F-C0032-001?Authorization=CWA-B8759F1C-FFC6-47C4-8502-95A53CDBFEB1&limit=10&format=JSON&locationName=%E6%96%B0%E5%8C%97%E5%B8%82&elementName=" +
		string(element) + "&sort=time" + timeStr

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// 外層 JSON 包含 success, result, records
	var outer WeatherResponse

	if err := json.Unmarshal(body, &outer); err != nil {
		panic(err)
	}

	// 印出解析結果
	fmt.Printf("描述: %s\n", outer.Records.DatasetDescription)
	for _, loc := range outer.Records.Location {
		fmt.Printf("地點: %s\n", loc.LocationName)
		for _, we := range loc.WeatherElement {
			fmt.Printf("  元素: %s\n", we.ElementName)
			for _, t := range we.Time {
				fmt.Printf("    %s ~ %s : %s %s\n",
					t.StartTime, t.EndTime,
					t.Parameter.ParameterName, t.Parameter.ParameterUnit)
			}
		}
	}

	return outer
}

func saveWeather(
	period Period,
	locationInfo LocationInfo,
) error {
	dsn := "root:root@tcp(127.0.0.1:3306)/testdb?parseTime=true"
	db, err := sqlx.Connect("mysql", dsn)
	if err != nil {
		log.Fatalln("連線資料庫失敗:", err)
	}
	defer db.Close()
	fmt.Println("✅ MySQL 連線成功 (sqlx)")

	var args []interface{}
	args = append(args, "新北市", locationInfo.MinT, locationInfo.MaxT, string(period), locationInfo.Date)

	query := `INSERT INTO Weather (city, min_t, max_t, period, date, updated) VALUES (?, ?, ?, ?, ?, NOW())`
	_, err = db.Exec(query, args...)
	if err != nil {
		return err
	}
	fmt.Println("✅ 新增成功")

	return nil
}
