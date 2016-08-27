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
	"github.com/MeasureTheFuture/mothership/vec"
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

const (
	FrameW = 1024
	FrameH = 720
)

func maxTravelTime(a models.Waypoint, b models.Waypoint, ss *models.ScoutSummary) float32 {
	return float32(0.0)
}

func updateTimeBuckets(db *sql.DB, ss *models.ScoutSummary, si *models.ScoutInteraction) {

	// For each segment in an interaction.
	// Generate shaft AABB from the two waypoints.
	// Work out direction of travel (vec) and total travel duration for segment.
	// Work out maximum travel time that can be spent in a bucket.

	// For each bucket in summary
	// Intersect test bucket against the shaft AABB.
	// if intersect bucket
	// work out how much the bucket intersects shaft.
	// bt is the ratio of this as multiple of max bucket travel time.
	// increment bucket time by bt above
}
