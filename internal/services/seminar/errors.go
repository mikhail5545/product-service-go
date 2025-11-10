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

package seminar

import "errors"

var (
	// ErrInvalidArgument invalid request payload error
	ErrInvalidArgument = errors.New("invalid argument")
	// ErrNotFound seminar not found error
	ErrNotFound = errors.New("seminar not found")
	// ErrIncompleteData seminar missing one or more required product IDs error
	ErrIncompleteData = errors.New("seminar record is missing one or more required product IDs")
	// ErrProductsNotFound unable to find all products for seminar error
	ErrProductsNotFound = errors.New("could not find all products for seminar")
	// ErrImageLimitExceeded can't upload more images error
	ErrImageLimitExceeded = errors.New("maximum number of uploaded images is 5 per item")
	// ErrImageNotFoundOnOwner can't find image on seminar error
	ErrImageNotFoundOnOwner = errors.New("image not found on seminar")
)
