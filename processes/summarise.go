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

package processes

import (
	"database/sql"
	"github.com/MeasureTheFuture/mothership/configuration"
	"github.com/MeasureTheFuture/mothership/models"
	"log"
	"time"
)

func Summarise(db *sql.DB, c configuration.Configuration) {
	poll := time.NewTicker(time.Millisecond * time.Duration(c.SummariseInterval)).C

	for {
		select {
		case <-poll:
			log.Printf("INFO: Generating summary of interaction data.")
			updateUnprocessed(db)
		}
	}
}

func updateUnprocessed(db *sql.DB) {
	up, err := models.GetUnprocessed(db)
	if err != nil {
		log.Printf("ERROR: Summarise unable to get unprocessed scout interactions.")
		return
	}

	for _, si := range up {
		err := models.IncrementVisitorCount(db, si.ScoutId)
		if err != nil {
			log.Printf("ERROR: Summarise unable to increment visitor count")
			log.Print(err)
			return
		}

		err = models.MarkProcessed(db, si.Id)
		if err != nil {
			log.Printf("ERROR: Summarise unable to make scout interaction as processed")
			log.Print(err)
			return
		}
	}
}
