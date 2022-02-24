package model

import (
	"database/sql"
	"errors"

	"github.com/labstack/gommon/log"
)

type Alert struct {
	ID        int64   `json:"id"`
	Crypto    string  `json:"crypto"`
	Direction bool    `json:"direction"`
	Price     float64 `json:"price"`
}

func SaveAlert(db *sql.DB, alert *Alert) (*Alert, error) {
	// Insert alert into database
	sql := `
		INSERT INTO alerts(crypto, direction, price)
		VALUES(?, ?, ?)
		`
	res, err := db.Exec(sql, alert.Crypto, alert.Direction, alert.Price)
	if err != nil {
		return nil, err
	}
	alert.ID, _ = res.LastInsertId()
	return alert, nil
}

func RemoveAlert(db *sql.DB, id int64) (*Alert, error) {
	// Remove alert from database
	record, err := GetAlertByID(db, id)
	if err != nil {
		return nil, errors.New("Alert not found")
	}

	sql := `
		DELETE FROM alerts
		WHERE id = ?
		`
	_, err = db.Exec(sql, id)
	if err != nil {
		return nil, errors.New("Alert Removal Failed")
	}
	return record, nil
}

func GetAlertByID(db *sql.DB, id int64) (*Alert, error) {
	// Get alert from database
	sql := `
		SELECT id, crypto, price, direction
		FROM alerts
		WHERE id = ?
		`
	row := db.QueryRow(sql, id)
	alert := new(Alert)
	err := row.Scan(&alert.ID, &alert.Crypto, &alert.Price, &alert.Direction)
	if err != nil {
		return nil, err
	}
	return alert, nil
}

func GetAlerts(db *sql.DB) ([]*Alert, error) {
	// Get alerts from database
	sql := `
		SELECT id, crypto, price, direction
		FROM alerts
		`
	rows, err := db.Query(sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	alerts := make([]*Alert, 0)
	for rows.Next() {
		alert := new(Alert)
		err := rows.Scan(&alert.ID, &alert.Crypto, &alert.Price, &alert.Direction)
		if err != nil {
			return nil, err
		}
		alerts = append(alerts, alert)
	}
	return alerts, nil
}

func UpdateAlert(db *sql.DB, alert *Alert) (*Alert, error) {
	// Update alert in database
	sql := `
		UPDATE alerts
		SET price = ?, direction = ?
		WHERE id = ?
		`
	_, err := db.Exec(sql, alert.Price, alert.Direction, alert.ID)
	if err != nil {
		log.Debug(err)
		return nil, errors.New("Alert Update Failed")
	}
	return alert, nil
}
