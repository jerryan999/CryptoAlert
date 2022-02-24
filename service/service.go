package service

import (
	"database/sql"
	"errors"

	"github.com/go-redis/redis"
	"github.com/jerryan999/CryptoAlert/model"
	"github.com/jerryan999/CryptoAlert/utils"
	"github.com/labstack/gommon/log"
)

type AlertService struct {
	db  *sql.DB
	rdb *redis.Client
}

func NewAlertService(db *sql.DB, rdb *redis.Client) *AlertService {
	return &AlertService{db: db, rdb: rdb}
}

func (s *AlertService) GetAlertByID(id int64) (*model.Alert, error) {
	return model.GetAlertByID(s.db, id)
}

func (s *AlertService) AddAlert(alert *model.Alert) (*model.Alert, error) {
	// Add alert to database
	alert, err := model.SaveAlert(s.db, alert)
	if err != nil {
		return alert, errors.New("saving alert to DB Failed")
	}
	log.Info("Alert saved to DB")

	// save alert to redis sorted set
	key := utils.GetAlertQueueKey(alert.Crypto, alert.Direction)
	s.rdb.ZAdd(key, redis.Z{
		Score:  alert.Price,
		Member: alert.ID,
	})
	log.Info("Alert saved to redis")

	return alert, nil
}

func (s *AlertService) RemoveAlert(alert *model.Alert) (*model.Alert, error) {
	// Remove alert from database
	alert, err := model.RemoveAlert(s.db, alert.ID)
	if err != nil {
		return alert, nil
	}
	log.Info("Alert removed from DB")

	// remove alert from redis sortedset
	key := utils.GetAlertQueueKey(alert.Crypto, alert.Direction)
	s.rdb.ZRem(key, alert.ID)
	log.Info("Alert removed from redis")
	return alert, nil

}

func (s *AlertService) UpdateAlert(alert *model.Alert) (*model.Alert, error) {
	// Update alert in database
	_, err := model.UpdateAlert(s.db, alert)
	if err != nil {
		return alert, errors.New("updating alert in DB Failed")
	}
	log.Info("Alert update into DB")

	// Update alert in redis sortedset
	key := utils.GetAlertQueueKey(alert.Crypto, alert.Direction)
	s.rdb.ZAdd(key, redis.Z{
		Score:  alert.Price,
		Member: alert.ID,
	})
	log.Info("Alert update into redis")

	return alert, nil
}

func (s *AlertService) GetAlerts() ([]*model.Alert, error) {
	// Get alerts from database
	alerts, err := model.GetAlerts(s.db)
	if err != nil {
		return alerts, errors.New("getting alerts from DB Failed")
	}

	return alerts, nil
}
