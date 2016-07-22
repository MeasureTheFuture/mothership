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
	"github.com/labstack/echo"
)

func ScoutCalibrated(c echo.Context) error {
	return nil
}

func ScoutInteraction(c echo.Context) error {
	return nil
}

func ScoutLog(c echo.Context) error {
	return nil
}

func ScoutHeartbeat(c echo.Context) error {
	// No authorization signature. Drop the request.
	if !c.Request().Header().Contains("Mothership-Authorization") {
		return nil
	}

	//uuid := c.Request().Header().Get("Mothership-Authorization")

	return nil
}
