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

package vec

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestUtil(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Util Suite")
}

var _ = Describe("Util", func() {
	Context("Max", func() {
		It("Should return the maximum of two values", func() {
			Ω(Max(5, 2)).Should(Equal(5))
			Ω(Max(2, 5)).Should(Equal(5))
		})
	})

	Context("Min", func() {
		It("Should return the minimum of two values", func() {
			Ω(Min(5, 2)).Should(Equal(2))
			Ω(Min(2, 5)).Should(Equal(2))
		})
	})
})
