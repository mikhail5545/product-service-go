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

package common

import (
	"errors"
	"unicode"
)

// ValidateName is a validation rule that checks if a string starts with a letter
// and contains at least one letter. It can handle both `string` and `*string` types.
func ValidateName(value interface{}) error {
	var name string
	switch v := value.(type) {
	case string:
		name = v
	case *string:
		if v != nil {
			name = *v
		}
	}
	if name != "" && !unicode.IsLetter([]rune(name)[0]) {
		return errors.New("must start with a letter")
	}
	return nil
}
