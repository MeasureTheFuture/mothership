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
 	_ "github.com/lib/pq"
 	"mothership/models"
 	"log"
)

func main() {
	db, err := sql.Open("postgres", "user=cfreeman dbname=mothership")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	s := models.Scout{-1, "800fd548-2d2b-4185-885d-6323ccbe88a0", "192.168.0.1", true, "foo"}
	err = s.Insert(db)
	log.Printf("ERROR: Unable to insert into DB")
	log.Print(err)
	log.Print(s.Id)
}