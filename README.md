# Gsam



## Mysql Created 語法
理論上已經包在 docker compose up 裡面，下方資訊供參考。
./mysql-init/schema.go 也有語法

```sql
CREATE TABLE User (
    email VARCHAR(50) PRIMARY KEY,
    password VARCHAR(255) NOT NULL,
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE Weather (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
    city VARCHAR(50) NOT NULL,
    min_t DECIMAL(4,1) NOT NULL,
    max_t DECIMAL(4,1) NOT NULL,
    period ENUM('AM','PM') NOT NULL,
    date VARCHAR(50) NOT NULL,
    created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

## Init command

```
// IF NEED, reset db
docker-compose down -v

// docker-compose mysql
docker-compose up -d

go mod tidy

// run api server
go run ./main.go

// run job to get weather info
go run .\cronjob\get_weather.go
```

## API
Postman 可匯入檔案
```
./Gsam_demo.postman_collection.json
```

## 說明
在打 api 測試前請先起好 mysql db，

此為簡易 demo，db 暫且寫死，請跑 `docker-compose up -d` ，理論上裡面有包好初始 table。

然後再起 Server 打 API。

### Weather
需要先取資料做儲存才能打 API 取得資料，
通常排程應交由環境管理，
這邊提供相應程式請跑 `go run .\cronjob\get_weather.go`。

ps. 這邊也會需要 DB Table ，理論上也包在 docker compose 初始好了， 可以檢查是否有 Weather table，這邊多紀錄欄位 `date` 來判斷該筆資料日期，保留後續可擴充性。

然後可以打 `/weather` ，取得今日資料。分為兩隻 `/weather/auth` 需要驗證 auth token，這邊 token 為 API Login 給的，若沒有，請先 `AddUser` 然後 `Login`。

ps.官方 response可參考下列，官方 `/F-C0032-001 ` 為提供 `36小時氣象預測`，所以會有下列兩點問題
* 已過的時段沒有包含在response 內，所以下午才觸發取得資料會無法取得上午的資料
* 上午取資料的話，資料顆粒度為 12 小時，若依要求 db 需儲存 `6:00 - 12:00` , `12:00 - 18:00` 為顆粒單位的資料的話，則兩筆皆視為 `6:00 - 18:00` 的資料。
```json
  // 12:00 前取資料
  "records": {
    "datasetDescription": "三十六小時天氣預報",
    "location": [
      {
        "locationName": "新北市",
        "weatherElement": [
          {
            "elementName": "MinT",
            "time": [
              {
                "startTime": "2025-09-26 06:00:00",
                "endTime": "2025-09-26 18:00:00",
                "parameter": {
                  "parameterName": "31",
                  "parameterUnit": "C"
                }
              }
            ]
          }
        ]
      }
    ]
  }

  
  // 12:00後 18:00 前取資料
  "records": {
    "datasetDescription": "三十六小時天氣預報",
    "location": [
      {
        "locationName": "新北市",
        "weatherElement": [
          {
            "elementName": "MinT",
            "time": [
              {
                "startTime": "2025-09-26 06:00:00",
                "endTime": "2025-09-26 18:00:00",
                "parameter": {
                  "parameterName": "31",
                  "parameterUnit": "C"
                }
              }
            ]
          }
        ]
      }
    ]
  }
```