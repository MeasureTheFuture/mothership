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
)

func TestScoutSummary(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Scout Summary Suite")
}

var _ = Describe("Scout Summary Model", func() {
	AfterEach(cleaner)

	Context("Insert", func() {
		It("Should be able to insert a scout summary", func() {
			s := Scout{-1, "800fd548-2d2b-4185-885d-6323ccbe88a0", "192.168.0.1", 8080, true, "foo", "idle"}
			err := s.Insert(db)
			Ω(err).Should(BeNil())

			ss := ScoutSummary{s.Id, 4, Buckets{}}
			ss.VisitTimeBuckets[1][5] = 0.1
			err = ss.Insert(db)
			Ω(err).Should(BeNil())

			ss2, err := GetScoutSummaryById(db, s.Id)
			Ω(err).Should(BeNil())
			Ω(ss2).Should(Equal(&ss))
		})
	})
})
