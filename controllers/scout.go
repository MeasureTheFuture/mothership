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
	"archive/zip"
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/MeasureTheFuture/mothership/models"
	"github.com/labstack/echo"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strconv"
)

func DownloadData(db *sql.DB, c echo.Context) error {
	var files []string

	sh, err := models.ScoutHealthsAsJSON(db)
	if err != nil {
		return err
	}
	files = append(files, sh)

	si, err := models.ScoutInteractionsAsJSON(db)
	if err != nil {
		return err
	}
	files = append(files, si)

	ss, err := models.ScoutSummariesAsJSON(db)
	if err != nil {
		return err
	}
	files = append(files, ss)

	sa, err := models.ScoutsAsJSON(db)
	if err != nil {
		return err
	}
	files = append(files, sa[:]...)

	// Zip the data export.
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)
	for _, file := range files {
		dst, err := w.Create(path.Base(file))
		if err != nil {
			return err
		}

		src, err := os.Open(file)
		if err != nil {
			return err
		}

		_, err = io.Copy(dst, src)
		if err != nil {
			return err
		}
	}

	err = w.Close()
	if err != nil {
		return err
	}

	// Write the zip to disk.
	zipFile := path.Dir(sh) + "/download.zip"
	err = ioutil.WriteFile(zipFile, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	// Send zip to zee client.
	return c.File(zipFile)
}

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

	body, err := ioutil.ReadAll(c.Request().Body)
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

	// If the scout is de-authorised/deactivated - clear it all out.
	if !ns.Authorised {
		ns.State = models.IDLE

		err = models.DeleteScoutHealths(db, ns.Id)
		if err != nil {
			return err
		}

		err = models.DeleteScoutInteractions(db, ns.Id)
		if err != nil {
			return err
		}

		err = models.DeleteScoutLogs(db, ns.Id)
		if err != nil {
			return err
		}

		ss, err := models.GetScoutSummaryById(db, ns.Id)
		if err != nil {
			return err
		}

		err = ss.Clear(db)
		if err != nil {
			return err
		}

		err = ns.ClearCalibrationFrame(db)
		if err != nil {
			return err
		}

		// Tell the scout to stop!
		go http.Get("http://" + ns.IpAddress + ":" + strconv.Itoa(int(ns.Port)) + "/measure/stop")
	}

	err = ns.Update(db)
	if err != nil {
		return err
	}

	if ns.State == models.CALIBRATING {
		go http.Get("http://" + ns.IpAddress + ":" + strconv.Itoa(int(ns.Port)) + "/calibrate")

	} else if ns.State == models.MEASURING {
		go http.Get("http://" + ns.IpAddress + ":" + strconv.Itoa(int(ns.Port)) + "/measure/start")

	}

	c.Request()
	return c.HTML(http.StatusOK, "updated succesfully")
}
