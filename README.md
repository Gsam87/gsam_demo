# Gsam



## Mysql Created 語法
理論上已經包在 docker compose up 裡面，下方資續供參考。
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