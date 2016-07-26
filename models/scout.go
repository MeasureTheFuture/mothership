/*
 * Copyright (C) 2016 Clinton Freeman
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

package models

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	_ "github.com/lib/pq"
)

type ScoutState string

const (
	IDLE        ScoutState = "idle"
	CALIBRATING ScoutState = "calibrating"
	CALIBRATED  ScoutState = "calibrated"
	MEASURING   ScoutState = "measuring"
)

func (s *ScoutState) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Unable to deserialise ScoutState")
	}

	*s = ScoutState(string(asBytes))
	return nil
}

func (s ScoutState) Value() (driver.Value, error) {
	return string(s), nil
}

type Scout struct {
	Id               int64
	UUID             string
	IpAddress        string
	Authorised       bool
	Name             string
	CalibrationFrame []byte
	State            ScoutState
}

func GetScoutById(db *sql.DB, id int64) (*Scout, error) {
	const query = `SELECT uuid, ip_address, authorised, name, calibration_frame, state
				   FROM scouts WHERE id = $1`
	var result Scout
	err := db.QueryRow(query, id).Scan(&result.UUID, &result.IpAddress, &result.Authorised,
		&result.Name, &result.CalibrationFrame, &result.State)
	result.Id = id

	return &result, err
}

func GetScoutByUUID(db *sql.DB, uuid string) (*Scout, error) {
	const query = `SELECT id, ip_address, authorised, name, calibration_frame, state
				   FROM scouts WHERE uuid = $1`
	var result Scout
	err := db.QueryRow(query, uuid).Scan(&result.Id, &result.IpAddress, &result.Authorised,
		&result.Name, &result.CalibrationFrame, &result.State)
	result.UUID = uuid

	return &result, err
}

func NumScouts(db *sql.DB) (int64, error) {
	const query = `SELECT COUNT(*) FROM scouts`
	var result int64
	err := db.QueryRow(query).Scan(&result)

	return result, err
}

func (s *Scout) Insert(db *sql.DB) error {
	const query = `INSERT INTO scouts (uuid, ip_address, authorised, name, calibration_frame, state)
				   VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	return db.QueryRow(query, s.UUID, s.IpAddress, s.Authorised, s.Name, s.CalibrationFrame,
		s.State).Scan(&s.Id)
}

func (s *Scout) Update(db *sql.DB) error {
	const query = `UPDATE scouts SET uuid = $1, ip_address = $2, authorised = $3, name = $4,
				   calibration_frame = $5, state = $6 WHERE id = $7`
	_, err := db.Exec(query, s.UUID, s.IpAddress, s.Authorised, s.Name, s.CalibrationFrame, s.State, s.Id)
	return err
}
