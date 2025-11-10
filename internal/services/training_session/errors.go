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

// Package trainingsession provides service-layer business logic for training sessions.
package trainingsession

import "errors"

var (
	// ErrInvalidArgument invalid request payload error
	ErrInvalidArgument = errors.New("invalid argument")
	// ErrNotFound training session or it's product not found error
	ErrNotFound = errors.New("training session not found")
	// ErrImageLimitExceeded can't upload more images error
	ErrImageLimitExceeded = errors.New("maximum number of uploaded images is 5 per item")
	// ErrImageNotFoundOnOwner can't find image on training session error
	ErrImageNotFoundOnOwner = errors.New("image not found on training session")
)
