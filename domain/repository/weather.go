package repo

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/qww83728/gsam_demo/domain/entity"
	repo_entity "github.com/qww83728/gsam_demo/domain/entity/repo"
)

type WeatherRepo interface {
	GetTodayInfo() ([]repo_entity.Weather, error)
}

type WeatherRepoImpl struct {
	db *sqlx.DB
}

func NewWeatherRepo(db *sqlx.DB) WeatherRepo {
	return &WeatherRepoImpl{
		db: db,
	}
}

func (r *WeatherRepoImpl) GetTodayInfo() ([]repo_entity.Weather, error) {
	var weathers []repo_entity.Weather

	// 今天日期字串
	today := time.Now().Format("2006-01-02")

	query := `
		SELECT id, city, min_t, max_t, period, date, created, updated
		FROM Weather
		WHERE date = ?
	`

	if err := r.db.Select(&weathers, query, today); err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	return weathers, nil
}
