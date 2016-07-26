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

package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"github.com/labstack/echo"
	"mothership/models"
	"net/http"
	"time"
)

func isScoutAuthorised(db *sql.DB, c echo.Context) (*models.Scout, error) {
	// No authorization signature.
	if !c.Request().Header().Contains("Mothership-Authorization") {
		return nil, nil
	}

	uuid := c.Request().Header().Get("Mothership-Authorization")
	s, err := models.GetScoutByUUID(db, uuid)
	if err != nil {
		// Scout doesn't exist, create it and mark it as un-authorized.
		ns := models.Scout{-1, uuid, "0.0.0.0", false, "Unknown location", "idle"}
		err = ns.Insert(db)
		return nil, err
	}

	if !s.Authorised {
		return nil, nil
	}

	return s, nil
}

func ScoutCalibrated(db *sql.DB, c echo.Context) error {
	return c.HTML(http.StatusNotFound, "")
}

func ScoutInteraction(db *sql.DB, c echo.Context) error {
	return c.HTML(http.StatusNotFound, "")
}

func ScoutLog(db *sql.DB, c echo.Context) error {
	s, err := isScoutAuthorised(db, c)
	if err != nil {
		return err
	}
	if s == nil {
		return c.HTML(http.StatusNotFound, "")
	}

	data, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := data.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Store scout log.
	var buff bytes.Buffer
	_, err = buff.ReadFrom(src)
	sl := models.ScoutLog{s.Id, buff.Bytes(), time.Now().UTC()}
	err = sl.Insert(db)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "Log received")
}

type Heartbeat struct {
	UUID    string     // The UUID for the scout.
	Version string     // The version of the protocol used used for transmitting data to the mothership.
	Health  HealthData // The current health status of the scout.

}

type HealthData struct {
	IpAddress   string  // The current IP address of the scout.
	CPU         float32 // The amount of CPU load currently being consumed on the scout. 0.0 - no load, 1.0 - full load.
	Memory      float32 // The amount of memory consumed on the scout. 0.0 - no memory used, 1.0 no memory available.
	TotalMemory float32 // The total number of gigabytes of virtual memory currently available.
	Storage     float32 // The amount of storage consumed on the scout. 0.0 - disk unused, 1.0 disk full.
}

func ScoutHeartbeat(db *sql.DB, c echo.Context) error {
	s, err := isScoutAuthorised(db, c)
	if err != nil {
		return err
	}
	if s == nil {
		return c.HTML(http.StatusNotFound, "")
	}

	// Deserialise the heartbeat from JSON.
	data, err := c.FormFile("file")
	if err != nil {
		return err
	}

	src, err := data.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	var hb Heartbeat
	err = json.NewDecoder(src).Decode(&hb)
	if err != nil {
		return err
	}

	// Create new record for the scout health.
	sh := models.ScoutHealth{s.Id, hb.Health.CPU, hb.Health.Memory, hb.Health.TotalMemory, hb.Health.Storage, time.Now().UTC()}
	err = sh.Insert(db)
	if err != nil {
		return err
	}

	// Update IP address of scout.
	s.IpAddress = hb.Health.IpAddress
	err = s.Update(db)
	if err != nil {
		return err
	}

	return c.HTML(http.StatusOK, "Heartbeat received")
}
