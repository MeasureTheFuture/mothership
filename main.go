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

package main

import (
	"database/sql"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	_ "github.com/lib/pq"
	"mothership/controllers"
	//"net/http"
)

func main() {
	db, err := sql.Open("postgres", "user=cfreeman dbname=mothership")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	e := echo.New()
	e.Static("/", "public")
	e.Static("/css", "public/css")
	e.Static("/js", "public/js")
	// e.GET("/", func(c echo.Context) error {
	// 	return c.String(http.StatusOK, "Hello, World!")
	// })
	e.POST("/scout_api/calibrated", func(c echo.Context) error {
		return controllers.ScoutCalibrated(db, c)
	})
	e.POST("/scout_api/interaction", func(c echo.Context) error {
		return controllers.ScoutInteraction(db, c)
	})
	e.POST("/scout_api/log", func(c echo.Context) error {
		return controllers.ScoutLog(db, c)
	})
	e.POST("/scout_api/heartbeat", func(c echo.Context) error {
		return controllers.ScoutHeartbeat(db, c)
	})

	e.Run(fasthttp.New(":1323"))
}
