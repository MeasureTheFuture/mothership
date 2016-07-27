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

package configuration

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	DBUserName   string // The name of the user with read/write privileges on DBName
	DBName       string // The name of the database that holds the production data.
	DBTestName   string // The name of the database that holds testing data.
	Address      string // The address and port that the mothership is accessible on.
	StaticAssets string // The path to the static assets rendered by the mothership.
}

func Parse(configFile string) (c Configuration, err error) {
	c = Configuration{"mtf", "mothership", "mothership_test", ":80", "public"}

	// Open the configuration file.
	file, err := os.Open(configFile)
	if err != nil {
		return c, err
	}
	defer file.Close()

	// Parse JSON in the configuration file.
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&c)
	return c, err
}