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

package models

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type ScoutHealth struct {
	ScoutId     int64
	CPU         float32
	Memory      float32
	TotalMemory float32
	Storage     float32
	CreatedAt   time.Time
}

func (s *ScoutHealth) Insert(db *sql.DB) error {
	const query = `INSERT INTO scout_healths (scout_id, cpu, memory, total_memory, storage,
		created_at) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.Exec(query, s.ScoutId, s.CPU, s.Memory, s.TotalMemory, s.Storage, s.CreatedAt)
	return err
}

// TODO: Prune/Delete. Reduce fidielity of data storage the older it becomes.
// TODO: Get last Health.
// TODO: Get Health history summary.
