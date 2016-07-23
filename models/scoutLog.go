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
	_ "github.com/lib/pq"
	"time"
)

type ScoutLog struct {
	ScoutId   int64
	Log       []byte
	CreatedAt time.Time
}

func GetScoutLogById(db *sql.DB, scoutId int64, time time.Time) (*ScoutLog, error) {
	const query = `SELECT log FROM scout_logs WHERE scout_id = $1 AND created_at = $2`

	var result ScoutLog
	err := db.QueryRow(query, scoutId, time).Scan(&result.Log)
	result.ScoutId = scoutId
	result.CreatedAt = time

	return &result, err
}

func (s *ScoutLog) Insert(db *sql.DB) error {
	const query = `INSERT INTO scout_logs (scout_id, log, created_at) VALUES ($1, $2, $3)`
	_, err := db.Exec(query, s.ScoutId, s.Log, s.CreatedAt)

	return err
}
