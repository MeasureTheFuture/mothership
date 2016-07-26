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
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestScout(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scout Suite")
}

var db *sql.DB
var err error

var _ = Describe("Scout Model", func() {

	BeforeSuite(func() {
		db, err = sql.Open("postgres", "user=cfreeman dbname=mothership_test")
		Ω(err).Should(BeNil())
	})

	AfterEach(func() {
		const query = `DELETE FROM scouts`
		_, err := db.Exec(query)
		Ω(err).Should(BeNil())
	})

	Context("Insert", func() {
		It("should insert a valid scout into the DB.", func() {
			s := Scout{-1, "800fd548-2d2b-4185-885d-6323ccbe88a0", "192.168.0.1", true, "foo", "calibrated"}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			s2, err := GetScoutById(db, s.Id)
			Ω(err).Should(BeNil())
			Ω(&s).Should(Equal(s2))
		})

		It("should return an error when an invalid scout is inserted into the DB.", func() {
			s := Scout{-1, "aa", "192.168.0.1", true, "foo", "calibrating"}
			err := s.Insert(db)
			Ω(err).ShouldNot(BeNil())
			Ω(s.Id).Should(Equal(int64(-1)))
		})
	})

	Context("GetByUUID", func() {
		It("should be able to get a scout by UUID", func() {
			s, err := GetScoutByUUID(db, "800fd548-2d2b-4185-885d-6323ccbe88a0")
			Ω(err).ShouldNot(BeNil())

			s2 := Scout{-1, "800fd548-2d2b-4185-885d-6323ccbe88a0", "192.168.0.1", true, "foo", "calibrated"}
			err = s2.Insert(db)
			Ω(err).Should(BeNil())

			s, err = GetScoutByUUID(db, "800fd548-2d2b-4185-885d-6323ccbe88a0")
			Ω(err).Should(BeNil())
			Ω(s).Should(Equal(&s2))

			c, err := NumScouts(db)
			Ω(err).Should(BeNil())
			Ω(c).Should(Equal(int64(1)))
		})
	})

	Context("Update", func() {
		It("should be able to update a scout in the DB", func() {
			s := Scout{-1, "800fd548-2d2b-4185-885d-6323ccbe88a0", "192.168.0.1", true, "foo", "measuring"}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			s.IpAddress = "192.168.0.2"
			err = s.Update(db)
			Ω(err).Should(BeNil())
			s2, err := GetScoutById(db, s.Id)
			Ω(err).Should(BeNil())

			Ω(&s).Should(Equal(s2))
		})
	})
})
