/*
 * state.go --- Internal process state.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package process

// Internal state for processes.
type State struct {
	parent *Process
}

// Send data from a process to an external entity.
func (ps *State) Send(data interface{}) bool {
	select {
	case ps.parent.chanFromState <- data:
		return true

	default:
	}

	return false
}

// Send data from a process to an external entity with blocking.
func (ps *State) SendBlocking(data interface{}) {
	ps.parent.chanFromState <- data
}

// Read data from an external entity.
func (ps *State) Receive() (interface{}, bool) {
	select {
	case data := <-ps.parent.chanToState:
		return data, true

	default:
	}

	return nil, false
}

// Read data from an external entity with blocking.
func (ps *State) ReceiveBlocking() interface{} {
	return <-ps.parent.chanToState
}

/* state.go ends here. */
