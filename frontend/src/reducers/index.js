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
const initialState = {
  locations:[],
  active:0
}

function selectLocation(value, l) {
  return value.id == l.id;
}

function mothership(state, action) {
  if (state === undefined) {
    return initialState;
  }

  switch (action.type) {
    case 'UPDATE_LOCATIONS':
      return {
        locations: action.locations,
        active: 0
      }

    case 'SET_ACTIVE':
      return {
        locations: state.locations,
        active: Math.min(state.locations.length - 1, Math.max(0, action.active))
      }

    default:
      return state;
  }
}

module.exports = mothership