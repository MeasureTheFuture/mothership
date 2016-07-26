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
	_ "github.com/lib/pq"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"time"
)

func TestScoutHealth(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scout Health Suite")
}

var _ = Describe("Scout Health Model", func() {
	AfterEach(func() {
		const query = `DELETE FROM scout_healths`
		_, err := db.Exec(query)
		Ω(err).Should(BeNil())
	})

	Context("Insert", func() {
		It("should insert a valid scouthealth into the DB.", func() {
			s := Scout{-1, "800fd548-2d2b-4185-885d-6323ccbe88a0", "192.168.0.1", true, "foo", "idle"}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			t := time.Now()
			sh := ScoutHealth{s.Id, 0.1, 0.2, 0.3, 0.4, t}
			err = sh.Insert(db)
			Ω(err).Should(BeNil())

			sh2, err := GetScoutHealthById(db, s.Id, t)
			Ω(err).Should(BeNil())
			Ω(&sh).Should(Equal(sh2))

		})

		It("should return an error when an invalid scout health is inserted into the DB.", func() {
			sh := ScoutHealth{-1, 0.1, 0.2, 0.3, 0.4, time.Now()}
			err := sh.Insert(db)
			Ω(err).ShouldNot(BeNil())
		})
	})
})
