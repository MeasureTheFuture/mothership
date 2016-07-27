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
	"flag"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/fasthttp"
	_ "github.com/lib/pq"
	"log"
	"mothership/configuration"
	"mothership/controllers"
	"os"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, "configFile", "mothership.json", "The path to the configuration file")
	flag.Parse()

	f, err := os.OpenFile("mothership.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Unable to open log file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)
	log.Printf("INFO: Starting mothership")

	config, err := configuration.Parse(configFile)
	if err != nil {
		log.Fatalf("ERROR: Can't parse configuration - %s", err)
	}

	db, err := sql.Open("postgres", "user="+config.DBUserName+" dbname="+config.DBName)
	if err != nil {
		log.Fatalf("ERROR: Can't open database - %s", err)
	}
	defer db.Close()

	e := echo.New()
	e.Static("/", config.StaticAssets)
	e.Static("/css", config.StaticAssets+"/css")
	e.Static("/js", config.StaticAssets+"/js")

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

	e.Run(fasthttp.New(config.Address))
}
