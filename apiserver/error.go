/*
 * error.go --- JSON error structure.
 *
 * Copyright (c) 2021 Paul Ward <asmodai@gmail.com>
 *
 * Author:     Paul Ward <asmodai@gmail.com>
 * Maintainer: Paul Ward <asmodai@gmail.com>
 *
 * This program is free software; you can redistribute it and/or
 * modify it under the terms of the GNU Lesser General Public License
 * as published by the Free Software Foundation; either version 3
 * of the License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this program; if not, see <http://www.gnu.org/licenses/>.
 */

package apiserver

type ErrorDocument struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

func NewError(status int, msg string) *ErrorDocument {
	return &ErrorDocument{
		Status:  status,
		Message: msg,
	}
}

func NewErrorDocument(status int, msg string) *Document {
	doc := NewDocument(status, nil)

	doc.SetError(
		&ErrorDocument{
			Status:  status,
			Message: msg,
		},
	)

	return doc
}

/* error.go ends here. */
