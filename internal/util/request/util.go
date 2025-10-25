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

// Package request provides shared utility functions for HTTP handlers.
//
// It works with requests, parsing and validating various data from echo.Context:
//
//   - Prams
//
//   - Query Search Params
//
//   - Request JSON payload
package request

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// BindAndValidateJSON binds the request body to the given struct and handles validation errors.
func BindAndValidateJSON(c echo.Context, req any) error {
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request JSON payload."})
	}
	// Here you could add validation logic using a library like 'validator'
	return nil
}

// GetIDParam extracts a required ID from the path parameters.
func GetIDParam(c echo.Context, paramName, errorMsg string) (string, error) {
	id := c.Param(paramName)
	if id == "" {
		return "", c.JSON(http.StatusBadRequest, map[string]string{"error": errorMsg})
	}
	return id, nil
}

// GetPaginationParams extracts 'limit' and 'offset' from query parameters with default values.
func GetPaginationParams(c echo.Context, defaultLimit, defaultOffset int) (int, int) {
	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limit < -1 {
		limit = defaultLimit
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offset < 0 {
		offset = defaultOffset
	}

	return limit, offset
}
