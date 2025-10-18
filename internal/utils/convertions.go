// vitainmove.com/product-service-go
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

package utils

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"vitainmove.com/product-service-go/internal/models"
	productpb "vitainmove.com/product-service-go/proto/product/v0"
	trainingsessionpb "vitainmove.com/product-service-go/proto/training_session/v0"
)

func ConvertToProtobufProduct(product *models.Product) *productpb.Product {
	pbProduct := &productpb.Product{
		Id:               product.ID,
		CreatedAt:        timestamppb.New(product.CreatedAt),
		UpdatedAt:        timestamppb.New(product.UpdatedAt),
		Name:             product.Name,
		Description:      product.Description,
		Price:            product.Price,
		ImageUrl:         product.ImageUrl,
		Amount:           int32(product.Amount),
		ProductType:      product.ProductType,
		ShippingRequired: product.ShippingRequired,
	}
	return pbProduct
}

func ConvertToProtobufTrainingSession(ts *models.TrainingSession) *trainingsessionpb.TrainingSession {
	pbTs := &trainingsessionpb.TrainingSession{
		Id:              ts.ID,
		CreatedAt:       timestamppb.New(ts.CreatedAt),
		UpdatedAt:       timestamppb.New(ts.UpdatedAt),
		DurationMinutes: int32(ts.DurationMinutes),
		Format:          ts.Format,
		Product:         ConvertToProtobufProduct(ts.Product),
	}
	return pbTs
}
