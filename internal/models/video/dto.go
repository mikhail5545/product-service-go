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

package video

type AddRequest struct {
	OwnerID        string `json:"owner_id"`
	MediaServiceID string `json:"media_service_id"`
}

type RemoveRequest struct {
	OwnerID        string `json:"owner_id"`
	MediaServiceID string `json:"media_service_id"`
}

type Owner struct {
	ID      string  `json:"id"`
	VideoID *string `json:"video_id"`
}
