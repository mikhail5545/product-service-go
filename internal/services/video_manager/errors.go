// github.com/mikhail5545/product-service-go
// microservice for vitianmove project family
// Copyright (C) 2025  Mikhail Kulik

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published
// by the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package videomanager

import "errors"

var (
	// ErrInvalidArgument invalid argument error
	ErrInvalidArgument = errors.New("invalid argument")
	// ErrVideoNotFound video not found error
	ErrVideoNotFound = errors.New("video not found")
	// ErrOwnerNotFound owner not found error
	ErrOwnerNotFound = errors.New("owner not found")
	// ErrAlreadyAssociated owner already has associated video error
	ErrAlreadyAssociated = errors.New("owner already has associated video")
	// ErrVideoInUse video is already associated with this owner
	ErrVideoInUse = errors.New("video is already associated with this owner")
)
