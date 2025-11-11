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

// Package types provides utility functions to convert from internal types to gRPC protobuf messages.
package types

import (
	"time"

	coursemodel "github.com/mikhail5545/product-service-go/internal/models/course"
	coursepartmodel "github.com/mikhail5545/product-service-go/internal/models/course_part"
	imagemodel "github.com/mikhail5545/product-service-go/internal/models/image"
	physicalgoodmodel "github.com/mikhail5545/product-service-go/internal/models/physical_good"
	productmodel "github.com/mikhail5545/product-service-go/internal/models/product"
	seminarmodel "github.com/mikhail5545/product-service-go/internal/models/seminar"
	trainingsessionmodel "github.com/mikhail5545/product-service-go/internal/models/training_session"
	coursepb "github.com/mikhail5545/proto-go/proto/product_service/course/v1"
	coursepartpb "github.com/mikhail5545/proto-go/proto/product_service/course_part/v0"
	imagepb "github.com/mikhail5545/proto-go/proto/product_service/image/v0"
	physicalgoodpb "github.com/mikhail5545/proto-go/proto/product_service/physical_good/v1"
	productpb "github.com/mikhail5545/proto-go/proto/product_service/product/v0"
	seminarpb "github.com/mikhail5545/proto-go/proto/product_service/seminar/v1"
	trainingsessionpb "github.com/mikhail5545/proto-go/proto/product_service/training_session/v1"
	"google.golang.org/protobuf/types/known/fieldmaskpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ImageToProtobuf is a helper function to convert from internal model
// to gRPC protobuf message.
func ImageToProtobuf(image *imagemodel.Image) *imagepb.Image {
	return &imagepb.Image{
		PublicId:       image.PublicID,
		Url:            image.URL,
		SecureUrl:      image.SecureURL,
		MediaServiceId: image.MediaServiceID,
	}
}

// CourseDetaisToProtobuf is a helper function to convert from internal model
// to gRPC protobuf message.
func CourseDetaisToProtobuf(details *coursemodel.CourseDetails) *coursepb.CourseDetails {
	pbdetails := &coursepb.CourseDetails{
		Course: &coursepb.Course{
			Id:               details.ID,
			CreatedAt:        timestamppb.New(details.CreatedAt),
			UpdatedAt:        timestamppb.New(details.UpdatedAt),
			Name:             details.Name,
			ShortDescription: details.ShortDescription,
			LongDescription:  details.LongDescription,
			Topic:            details.Topic,
			Tags:             details.Tags,
			AccessDuration:   int32(details.AccessDuration),
			InStock:          details.InStock,
		},
		Price:     details.Price,
		ProductId: details.ProductID,
	}
	if details.DeletedAt.Valid {
		pbdetails.Course.DeletedAt = timestamppb.New(details.DeletedAt.Time)
	}
	if len(details.CourseParts) > 0 {
		for _, p := range details.CourseParts {
			pbdetails.Course.CourseParts = append(pbdetails.Course.CourseParts, CoursePartToProtobuf(p))
		}
	}
	if len(details.Course.Images) > 0 {
		for _, img := range details.Course.Images {
			pbdetails.Course.Images = append(pbdetails.Course.Images, ImageToProtobuf(&img))
		}
	}
	return pbdetails
}

// CourseToProtobufUpdate is a helper function to convert from internal service response
// in the map[string]any representation to gRPC protobuf message and populate response.
func CourseToProtobufUpdate(updates map[string]any) *coursepb.UpdateResponse {
	resp := &coursepb.UpdateResponse{}
	resp.Updated = &fieldmaskpb.FieldMask{}
	if courseUpdates, ok := updates["course"].(map[string]any); ok {
		for k, v := range courseUpdates {
			switch k {
			case "name":
				if val, ok := v.(string); ok {
					resp.Name = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.name")
				}
			case "short_description":
				if val, ok := v.(string); ok {
					resp.ShortDescription = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.short_description")
				}
			case "long_description":
				if val, ok := v.(string); ok {
					resp.LongDescription = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.long_description")
				}
			case "access_duration":
				if val, ok := v.(int32); ok {
					resp.AccessDuration = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.access_duration")
				}
			case "topic":
				if val, ok := v.(string); ok {
					resp.Topic = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.topic")
				}
			case "tags":
				if val, ok := v.([]string); ok {
					resp.Tags = val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.tags")
				}
			}
		}
	}
	if productUpdates, ok := updates["product"].(map[string]any); ok {
		if price, ok := productUpdates["price"].(float32); ok {
			resp.Price = &price
			resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.price")
		}
	}
	return resp
}

// CoursePartToProtobuf is a helper function to convert from internal model
// to gRPC protobuf message.
func CoursePartToProtobuf(part *coursepartmodel.CoursePart) *coursepartpb.CoursePart {
	pbpart := &coursepartpb.CoursePart{
		Id:               part.ID,
		CourseId:         part.CourseID,
		CreatedAt:        timestamppb.New(part.CreatedAt),
		UpdatedAt:        timestamppb.New(part.UpdatedAt),
		Name:             part.Name,
		ShortDescription: part.ShortDescription,
		LongDescription:  part.LongDescription,
		Tags:             part.Tags,
		Number:           int32(part.Number),
		Published:        part.Published,
	}
	return pbpart
}

// CoursePartToProtobufUpdate is a helper function to convert from internal service response
// in the map[string]any representation to gRPC protobuf message and populate response.
func CoursePartToProtobufUpdate(resp *coursepartpb.UpdateResponse, updates map[string]any) *coursepartpb.UpdateResponse {
	resp.Updated = &fieldmaskpb.FieldMask{}
	for k, v := range updates {
		switch k {
		case "name":
			if val, ok := v.(string); ok {
				resp.Name = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.name")
			}
		case "short_description":
			if val, ok := v.(string); ok {
				resp.ShortDescription = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.short_description")
			}
		case "long_description":
			if val, ok := v.(string); ok {
				resp.LongDescription = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.long_description")
			}
		case "number":
			if val, ok := v.(int32); ok {
				resp.Number = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.number")
			}
		case "tags":
			if val, ok := v.([]string); ok {
				resp.Tags = val
				resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.tags")
			}
		}
	}
	return resp
}

// SeminarDetaisToProtobuf is a helper function to convert from internal model
// to gRPC protobuf message.
func SeminarDetailsToProtobuf(details *seminarmodel.SeminarDetails) *seminarpb.SeminarDetails {
	pbdetails := &seminarpb.SeminarDetails{
		Seminar: &seminarpb.Seminar{
			Id:                      details.ID,
			CreatedAt:               timestamppb.New(details.CreatedAt),
			UpdatedAt:               timestamppb.New(details.UpdatedAt),
			Name:                    details.Name,
			ShortDescription:        details.ShortDescription,
			LongDescription:         details.LongDescription,
			Tags:                    details.Tags,
			Place:                   details.Place,
			Date:                    timestamppb.New(details.Date),
			EndingDate:              timestamppb.New(details.EndingDate),
			LatePaymentDate:         timestamppb.New(details.LatePaymentDate),
			ReservationProductId:    *details.ReservationProductID,
			EarlyProductId:          *details.EarlyProductID,
			LateProductId:           *details.LateProductID,
			EarlySurchargeProductId: *details.EarlySurchargeProductID,
			LateSurchargeProductId:  *details.LateSurchargeProductID,
			InStock:                 details.InStock,
		},
		ReservationPrice:               details.ReservationPrice,
		EarlyPrice:                     details.EarlyPrice,
		LatePrice:                      details.LatePrice,
		EarlySurchargePrice:            details.EarlySurchargePrice,
		LateSurchargePrice:             details.LateSurchargePrice,
		CurrentPrice:                   details.CurrentPrice,
		CurrentPriceProductId:          details.CurrentPriceProductID,
		CurrentSurchargePrice:          details.CurrentSurchargePrice,
		CurrentSurchargePriceProductId: details.CurrentSurchargePriceProductID,
	}
	if details.DeletedAt.Valid {
		pbdetails.Seminar.DeletedAt = timestamppb.New(details.DeletedAt.Time)
	}
	if len(details.Seminar.Images) > 0 {
		for _, img := range details.Seminar.Images {
			pbdetails.Seminar.Images = append(pbdetails.Seminar.Images, ImageToProtobuf(&img))
		}
	}
	return pbdetails
}

// SeminarToProtobufUpdate is a helper function to convert from internal service response
// in the map[string]any representation to gRPC protobuf message and populate response.
func SeminarToProtobufUpdate(resp *seminarpb.UpdateResponse, updates map[string]any) *seminarpb.UpdateResponse {
	resp.Updated = &fieldmaskpb.FieldMask{}
	if seminarUpdates, ok := updates["seminar"].(map[string]any); ok {
		for k, v := range seminarUpdates {
			switch k {
			case "name":
				if val, ok := v.(string); ok {
					resp.Name = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.name")
				}
			case "place":
				if val, ok := v.(string); ok {
					resp.Place = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.place")
				}
			case "tags":
				if val, ok := v.([]string); ok {
					resp.Tags = val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.tags")
				}
			case "short_description":
				if val, ok := v.(string); ok {
					resp.ShortDescription = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.short_description")
				}
			case "long_description":
				if val, ok := v.(string); ok {
					resp.LongDescription = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.long_description")
				}
			case "date":
				if val, ok := v.(time.Time); ok {
					resp.Date = timestamppb.New(val)
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.date")
				}
			case "ending_date":
				if val, ok := v.(time.Time); ok {
					resp.EndingDate = timestamppb.New(val)
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.ending_date")
				}
			case "late_payment_date":
				if val, ok := v.(time.Time); ok {
					resp.LatePaymentDate = timestamppb.New(val)
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.late_payment_date")
				}
			}
		}
	}
	if reservationProductUpdates, ok := updates["reservation_product"].(map[string]any); ok {
		if price, ok := reservationProductUpdates["price"].(float32); ok {
			resp.ReservationPrice = &price
			resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.reservation_price")
		}
	}
	if earlyProductUpdates, ok := updates["early_product"].(map[string]any); ok {
		if price, ok := earlyProductUpdates["price"].(float32); ok {
			resp.EarlyPrice = &price
			resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.early_price")
		}
	}
	if lateProductUpdates, ok := updates["late_product"].(map[string]any); ok {
		if price, ok := lateProductUpdates["price"].(float32); ok {
			resp.LatePrice = &price
			resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.late_price")
		}
	}
	if earlySurchargeProductUpdates, ok := updates["early_surcharge_product"].(map[string]any); ok {
		if price, ok := earlySurchargeProductUpdates["price"].(float32); ok {
			resp.EarlySurchargePrice = &price
			resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.early_surcharge_price")
		}
	}
	if lateSurchargeProductUpdates, ok := updates["late_surcharge_product"].(map[string]any); ok {
		if price, ok := lateSurchargeProductUpdates["price"].(float32); ok {
			resp.LateSurchargePrice = &price
			resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.late_surcharge_price")
		}
	}
	return resp
}

// TrainingSessionDetailsToProtobuf is a helper function to convert from internal model
// to gRPC protobuf message.
func TrainingSessionDetailsToProtobuf(details *trainingsessionmodel.TrainingSessionDetails) *trainingsessionpb.TrainingSessionDetails {
	pbdetails := &trainingsessionpb.TrainingSessionDetails{
		TrainingSession: &trainingsessionpb.TrainingSession{
			Id:               details.ID,
			CreatedAt:        timestamppb.New(details.CreatedAt),
			UpdatedAt:        timestamppb.New(details.UpdatedAt),
			Name:             details.Name,
			ShortDescription: details.ShortDescription,
			LongDescription:  details.LongDescription,
			Format:           details.Format,
			DurationMinutes:  int32(details.DurationMinutes),
			Tags:             details.Tags,
			InStock:          details.InStock,
		},
		Price:     details.Price,
		ProductId: details.ProductID,
	}
	if details.DeletedAt.Valid {
		pbdetails.TrainingSession.DeletedAt = timestamppb.New(details.DeletedAt.Time)
	}
	if len(details.TrainingSession.Images) > 0 {
		for _, img := range details.TrainingSession.Images {
			pbdetails.TrainingSession.Images = append(pbdetails.TrainingSession.Images, ImageToProtobuf(&img))
		}
	}
	return pbdetails
}

// TrainingSessionToProtobufUpdate is a helper function to convert from internal service response
// in the map[string]any representation to gRPC protobuf message and populate response.
func TrainingSessionToProtobufUpdate(resp *trainingsessionpb.UpdateResponse, updates map[string]any) *trainingsessionpb.UpdateResponse {
	resp.Updated = &fieldmaskpb.FieldMask{}
	if tsUpdates, ok := updates["training_session"].(map[string]any); ok {
		for k, v := range tsUpdates {
			switch k {
			case "name":
				if val, ok := v.(string); ok {
					resp.Name = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.name")
				}
			case "short_description":
				if val, ok := v.(string); ok {
					resp.ShortDescription = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.short_description")
				}
			case "long_description":
				if val, ok := v.(string); ok {
					resp.LongDescription = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.long_description")
				}
			case "format":
				if val, ok := v.(string); ok {
					resp.Format = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.format")
				}
			case "duration_minutes":
				if val, ok := v.(int32); ok {
					resp.DurationMinutes = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.duration_minutes")
				}
			case "tags":
				if val, ok := v.([]string); ok {
					resp.Tags = val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.tags")
				}
			}
		}
	}
	if productUpdates, ok := updates["product"].(map[string]any); ok {
		if price, ok := productUpdates["price"].(float32); ok {
			resp.Price = &price
			resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.price")
		}
	}
	return resp
}

// PhysicalGoodDetailsToProtobuf is a helper function to convert from internal model
// to gRPC transport message.
func PhysicalGoodDetailsToProtobuf(details *physicalgoodmodel.PhysicalGoodDetails) *physicalgoodpb.PhysicalGoodDetails {
	pbdetails := &physicalgoodpb.PhysicalGoodDetails{
		PhysicalGood: &physicalgoodpb.PhysicalGood{
			Id:               details.ID,
			CreatedAt:        timestamppb.New(details.CreatedAt),
			UpdatedAt:        timestamppb.New(details.UpdatedAt),
			Name:             details.Name,
			ShortDescription: details.ShortDescription,
			LongDescription:  details.LongDescription,
			Tags:             details.Tags,
			Amount:           int32(details.Amount),
			ShippingRequired: details.ShippingRequired,
			InStock:          details.InStock,
		},
		Price:     details.Price,
		ProductId: details.ProductID,
	}
	if details.DeletedAt.Valid {
		pbdetails.PhysicalGood.DeletedAt = timestamppb.New(details.DeletedAt.Time)
	}
	if len(details.PhysicalGood.Images) > 0 {
		for _, img := range details.PhysicalGood.Images {
			pbdetails.PhysicalGood.Images = append(pbdetails.PhysicalGood.Images, ImageToProtobuf(&img))
		}
	}
	return pbdetails
}

// PhysicalGoodToProtobufUpdate is a helper function to convert from internal service response
// in the map[string]any representation to gRPC protobuf message and populate response.
func PhysicalGoodToProtobufUpdate(resp *physicalgoodpb.UpdateResponse, updates map[string]any) *physicalgoodpb.UpdateResponse {
	resp.Updated = &fieldmaskpb.FieldMask{}
	if pgUpdates, ok := updates["physical_good"].(map[string]any); ok {
		for k, v := range pgUpdates {
			switch k {
			case "name":
				if val, ok := v.(string); ok {
					resp.Name = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.name")
				}
			case "short_description":
				if val, ok := v.(string); ok {
					resp.ShortDescription = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.short_description")
				}
			case "long_description":
				if val, ok := v.(string); ok {
					resp.LongDescription = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.long_description")
				}
			case "shipping_required":
				if val, ok := v.(bool); ok {
					resp.ShippingRequired = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.shipping_required")
				}
			case "amount":
				if val, ok := v.(int32); ok {
					resp.Amount = &val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.amount")
				}
			case "tags":
				if val, ok := v.([]string); ok {
					resp.Tags = val
					resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.tags")
				}
			}
		}
	}
	if productUpdates, ok := updates["product"].(map[string]any); ok {
		if price, ok := productUpdates["price"].(float32); ok {
			resp.Price = &price
			resp.Updated.Paths = append(resp.Updated.Paths, "updateresponse.price")
		}
	}
	return resp
}

// ProductToProtobuf is a helper function to convert from internal model
// to gRPC protobuf message.
func ProductToProtobuf(product *productmodel.Product) *productpb.Product {
	pbproduct := productpb.Product{
		Id:          product.ID,
		CreatedAt:   timestamppb.New(product.CreatedAt),
		UpdatedAt:   timestamppb.New(product.UpdatedAt),
		Price:       product.Price,
		InStock:     product.InStock,
		DetailsId:   product.DetailsID,
		DetailsType: product.DetailsType,
	}
	if product.DeletedAt.Valid {
		pbproduct.DeletedAt = timestamppb.New(product.DeletedAt.Time)
	}
	return &pbproduct
}
