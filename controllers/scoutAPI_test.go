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
	"mothership/models"
	"testing"
)

var (
	heartbeatJSON = `{
						"UUID":"59ef7180-f6b2-4129-99bf-970eb4312b4b",
						"Version":"0.1",
						"Health":{
							"IpAddress":"10.1.1.1",
							"CPU":0.4,
							"Memory":0.1,
							"TotalMemory":1233312.0,
							"Storage":0.1
						}
					}`
)

func TestScoutHeartbeat(t *testing.T) {

}
