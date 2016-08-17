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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestScoutInteraction(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scout Interaction Suite")
}

var _ = Describe("Scout Interaction Model", func() {
	AfterEach(func() {
		const query = `DELETE FROM scout_interactions`
		_, err := db.Exec(query)
		Ω(err).Should(BeNil())
	})

	Context("CreateScoutInteraction", func() {
		It("Should be able to create a scout interaction", func() {
			t := time.Now().UTC()

			wp := []Waypoint{Waypoint{1, 2, 3, 4, 0.1}}
			i := Interaction{"abc", "0.1", t, 0.1, wp}

			si := CreateScoutInteraction(&i)
			Ω(si.ScoutId).Should(Equal(int64(-1)))
			Ω(si.Duration).Should(Equal(i.Duration))
			Ω(si.EnteredAt).Should(Equal(t))
			Ω(si.Waypoints).Should(Equal(Path{[2]int{1, 2}}))
			Ω(si.WaypointWidths).Should(Equal(Path{[2]int{3, 4}}))
			Ω(si.WaypointTimes).Should(Equal(RealArray{0.1}))
		})
	})

	Context("Insert", func() {
		It("Should be able to insert a scout interaction", func() {
			s := Scout{-1, "800fd548-2d2b-4185-885d-6323ccbe88a0", "192.168.0.1", 8080, true, "foo", "idle"}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			t := time.Now().UTC()
			et := t.Round(15 * time.Minute)
			si := ScoutInteraction{s.Id, 0.2, [][2]int{[2]int{1, 2}, [2]int{5, 6}}, [][2]int{[2]int{3, 4}}, []float32{0.1}, false, t, et}
			err = si.Insert(db)
			Ω(err).Should(BeNil())

			si2, err := GetScoutInteractionById(db, s.Id, t)
			Ω(err).Should(BeNil())
			Ω(si2).Should(Equal(&si))
		})
	})
})
