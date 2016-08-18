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
	"strconv"
	"strings"
	"time"
)

type RealArray []float32
type Path [][2]int

type Waypoint struct {
	XPixels          int     // x-coordinate of waypoint centroid in pixels
	YPixels          int     // y-coordinate of waypoint centroid in pixels
	HalfWidthPixels  int     // Half the width of the waypoint in pixels
	HalfHeightPixels int     // Half the height of the waypoint in pixels
	T                float32 // The number of seconds elapsed since the beginning of the interaction
}

type Interaction struct {
	UUID     string     // The UUID for the scout that detected the interaction.
	Version  string     // The Version of the protocol used for transmitting data to the mothership
	Entered  time.Time  // The time the interaction started (rounded to nearest half hour)
	Duration float32    // The total duration of the interaction.
	Path     []Waypoint // The pathway of the interaction through the scene.
}

type ScoutInteraction struct {
	ScoutId        int64
	Duration       float32
	Waypoints      Path
	WaypointWidths Path
	WaypointTimes  RealArray
	Processed      bool
	CreatedAt      time.Time
	EnteredAt      time.Time
}

func (a *RealArray) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Unable to deserialise RealArray")
	}

	asString := string(asBytes)
	asString = asString[1 : len(asString)-1]
	elements := strings.Split(asString, ",")

	res := make([]float32, len(elements))
	for i, v := range elements {
		t, err := strconv.ParseFloat(v, 32)
		if err != nil {
			return err
		}

		res[i] = float32(t)
	}

	*a = res
	return nil
}

func (p *Path) Scan(value interface{}) error {
	asBytes, ok := value.([]byte)
	if !ok {
		return errors.New("Unable to deserialise Path")
	}

	asString := string(asBytes)
	asString = asString[2 : len(asString)-2]
	elements := strings.Split(asString, "),(")

	res := make([][2]int, len(elements))
	for i, v := range elements {
		wp := strings.Split(v, ",")
		wpv, err := strconv.ParseInt(wp[0], 10, 32)
		if err != nil {
			return err
		}
		res[i][0] = int(wpv)

		wpv, err = strconv.ParseInt(wp[1], 10, 32)
		if err != nil {
			return err
		}
		res[i][1] = int(wpv)
	}

	*p = res
	return nil
}

func (p Path) Value() (driver.Value, error) {
	res := "["
	for i, v := range p {
		res = res + "(" + strconv.FormatInt(int64(v[0]), 10) + "," + strconv.FormatInt(int64(v[1]), 10) + ")"

		if i < len(p)-1 {
			res = res + ","
		}
	}
	res = res + "]"

	return res, nil
}

func (a RealArray) Value() (driver.Value, error) {
	res := "{"
	for i, v := range a {
		res = res + strconv.FormatFloat(float64(v), 'f', -1, 32)
		if i < len(a)-1 {
			res = res + ","
		}
	}
	res = res + "}"

	return res, nil
}

func CreateScoutInteraction(i *Interaction) ScoutInteraction {
	var result ScoutInteraction

	result.ScoutId = -1
	result.Duration = i.Duration

	wpLength := len(i.Path)
	result.Waypoints = make([][2]int, wpLength)
	result.WaypointWidths = make([][2]int, wpLength)
	result.WaypointTimes = make([]float32, wpLength)

	for i, wp := range i.Path {
		result.Waypoints[i] = [2]int{wp.XPixels, wp.YPixels}
		result.WaypointWidths[i] = [2]int{wp.HalfWidthPixels, wp.HalfHeightPixels}
		result.WaypointTimes[i] = wp.T
	}

	result.Processed = false
	result.EnteredAt = i.Entered
	result.CreatedAt = time.Now().UTC()

	return result
}

func GetScoutInteractionById(db *sql.DB, scoutId int64, t time.Time) (*ScoutInteraction, error) {
	const query = `SELECT duration, waypoints, waypoint_widths, waypoint_times, processed, entered_at FROM scout_interactions WHERE scout_id = $1 AND created_at = $2`

	var result ScoutInteraction
	var et time.Time
	err := db.QueryRow(query, scoutId, t).Scan(&result.Duration, &result.Waypoints, &result.WaypointWidths, &result.WaypointTimes, &result.Processed, &et)
	result.ScoutId = scoutId
	result.EnteredAt = et.UTC()
	result.CreatedAt = t

	return &result, err
}

func GetLastScoutInteraction(db *sql.DB, scoutId int64) (*ScoutInteraction, error) {
	const query = `SELECT duration, waypoints, waypoint_widths, waypoint_times, processed, entered_at, created_at FROM scout_interactions WHERE scout_id = $1 ORDER BY created_at DESC LIMIT 1`

	var result ScoutInteraction
	err := db.QueryRow(query, scoutId).Scan(&result.Duration, &result.Waypoints, &result.WaypointWidths, &result.WaypointTimes, &result.Processed, &result.EnteredAt, &result.CreatedAt)
	result.ScoutId = scoutId

	return &result, err
}

func NumScoutInteractions(db *sql.DB) (int64, error) {
	const query = `SELECT COUNT(*) FROM scout_interactions`
	var result int64
	err := db.QueryRow(query).Scan(&result)

	return result, err
}

func (si *ScoutInteraction) Insert(db *sql.DB) error {
	const query = `INSERT INTO scout_interactions (scout_id, duration, waypoints, waypoint_widths,
		waypoint_times, processed, entered_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(query, si.ScoutId, si.Duration, si.Waypoints, si.WaypointWidths, si.WaypointTimes, si.Processed, si.EnteredAt, si.CreatedAt)

	return err
}
