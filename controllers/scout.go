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
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/MeasureTheFuture/mothership/models"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetScouts(db *sql.DB, c echo.Context) error {
	s, err := models.GetAllScouts(db)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, s)
}

func GetScoutFrame(db *sql.DB, c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	s, err := models.GetScoutById(db, id)
	if err != nil {
		return err
	}

	frame, err := s.GetCalibrationFrame(db)
	if err != nil {
		return err
	}

	c.Response().Header().Set(echo.HeaderContentType, "image/jpeg")
	c.Response().WriteHeader(http.StatusOK)
	_, err = c.Response().Write(frame)

	return err
}

func GetScout(db *sql.DB, c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	s, err := models.GetScoutById(db, id)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, s)
}

func UpdateScout(db *sql.DB, c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}

	body, err := ioutil.ReadAll(c.Request().Body())
	if err != nil {
		return err
	}

	var ns models.Scout
	err = json.Unmarshal(body, &ns)
	if err != nil {
		return err
	}

	if id != ns.Id {
		return errors.New("Mismatched Ids")
	}

	err = ns.Update(db)
	if err != nil {
		return err
	}

	if ns.State == models.CALIBRATING {
		go http.Get("http://" + ns.IpAddress + "/calibrate")

	} else if ns.State == models.MEASURING {
		go http.Get("http://" + ns.IpAddress + "/measure/start")

	}

	c.Request()
	return c.HTML(http.StatusOK, "updated succesfully")
}
