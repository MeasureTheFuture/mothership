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
	Id         int64      `json:"id"`
	UUID       string     `json:"uuid"`
	IpAddress  string     `json:"ip_address"`
	Port       int64      `json:"port"`
	Authorised bool       `json:"authorised"`
	Name       string     `json:"name"`
	State      ScoutState `json:"state"`
}

func GetScoutById(db *sql.DB, id int64) (*Scout, error) {
	const query = `SELECT uuid, ip_address, port, authorised, name, state
				   FROM scouts WHERE id = $1`
	var result Scout
	err := db.QueryRow(query, id).Scan(&result.UUID, &result.IpAddress, &result.Port, &result.Authorised,
		&result.Name, &result.State)
	result.Id = id

	return &result, err
}

func GetScoutByUUID(db *sql.DB, uuid string) (*Scout, error) {
	const query = `SELECT id, ip_address, port, authorised, name, state
				   FROM scouts WHERE uuid = $1`
	var result Scout
	err := db.QueryRow(query, uuid).Scan(&result.Id, &result.IpAddress, &result.Port, &result.Authorised,
		&result.Name, &result.State)
	result.UUID = uuid

	return &result, err
}

func GetAllScouts(db *sql.DB) ([]*Scout, error) {
	const query = `SELECT id, uuid, ip_address, port, authorised, name, state FROM scouts`

	var result []*Scout
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s Scout
		err = rows.Scan(&s.Id, &s.UUID, &s.IpAddress, &s.Port, &s.Authorised, &s.Name, &s.State)
		if err != nil {
			return nil, err
		}

		result = append(result, &s)
	}

	return result, nil
}

func NumScouts(db *sql.DB) (int64, error) {
	const query = `SELECT COUNT(*) FROM scouts`
	var result int64
	err := db.QueryRow(query).Scan(&result)

	return result, err
}

func (s *Scout) UpdateCalibrationFrame(db *sql.DB, frame []byte) error {
	const query = `UPDATE scouts SET calibration_frame = $1 WHERE id = $2`
	_, err := db.Exec(query, frame, s.Id)
	return err
}

func (s *Scout) GetCalibrationFrame(db *sql.DB) ([]byte, error) {
	const query = `SELECT calibration_frame FROM scouts WHERE id = $1`
	var result []byte
	err := db.QueryRow(query, s.Id).Scan(&result)

	return result, err
}

func (s *Scout) Insert(db *sql.DB) error {
	const query = `INSERT INTO scouts (uuid, ip_address, port, authorised, name, state)
				   VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := db.QueryRow(query, s.UUID, s.IpAddress, s.Port, s.Authorised, s.Name, s.State).Scan(&s.Id)
	if err != nil {
		return err
	}

	// Create matching empty summary.
	ss := ScoutSummary{s.Id, 0, Buckets{}}
	return ss.Insert(db)
}

func (s *Scout) Update(db *sql.DB) error {
	const query = `UPDATE scouts SET uuid = $1, ip_address = $2, port = $3, authorised = $4, name = $5,
				   state = $6 WHERE id = $7`
	_, err := db.Exec(query, s.UUID, s.IpAddress, s.Port, s.Authorised, s.Name, s.State, s.Id)
	return err
}
