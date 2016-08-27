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
	"github.com/MeasureTheFuture/mothership/models"
)

type Shaft struct {
	Bounds   AABB
	MinPlane [2]Vec
	MaxPlane [2]Vec
}

func ShaftFromWaypoints(a models.Waypoint, b models.Waypoint, maxW int, maxH int) Shaft {
	bb := AABBFromWaypoints(a, b, maxW, maxH)

	// Default to B being the bottom, left edge of the shaft.
	la := [2]Vec{Vec{(b.XPixels - b.HalfWidthPixels), (b.YPixels + b.HalfHeightPixels)},
		Vec{(a.XPixels - a.HalfWidthPixels), (a.YPixels + a.HalfHeightPixels)}}
	ra := [2]Vec{Vec{(b.XPixels + b.HalfWidthPixels), (b.YPixels - b.HalfHeightPixels)},
		Vec{(a.XPixels + a.HalfWidthPixels), (b.YPixels - b.HalfHeightPixels)}}

	lb := [2]Vec{Vec{(b.XPixels - b.HalfWidthPixels), (b.YPixels - b.HalfHeightPixels)},
		Vec{(a.XPixels - a.HalfWidthPixels), (a.YPixels - a.HalfHeightPixels)}}
	rb := [2]Vec{Vec{(b.XPixels + b.HalfWidthPixels), (b.YPixels + b.HalfHeightPixels)},
		Vec{(a.XPixels + a.HalfWidthPixels), (a.YPixels + a.HalfHeightPixels)}}

	if a.XPixels < b.XPixels {
		if a.YPixels < b.YPixels {
			// a is bottom, left edge of shaft.
			return Shaft{bb, la, ra}

		} else {
			// a is top, left edge of shaft.
			return Shaft{bb, lb, rb}
		}
	} else {
		// b is left edge.
		if a.YPixels < b.YPixels {
			// a is bottom, right edge of shaft.
			return Shaft{bb, lb, rb}

		} else {
			// a is top, right edge of shaft.
			return Shaft{bb, la, ra}
		}
	}
}

func (s *Shaft) Intersects(b *AABB) bool {
	if !b.Intersects(&s.Bounds) {
		return false
	}

	// if b is left of MinPlane - return false.

	// if b is right of MaxPlane - return false.

	return true
}
