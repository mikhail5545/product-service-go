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

package types

import (
	"github.com/mikhail5545/product-service-go/internal/models"
	coursepb "github.com/mikhail5545/proto-go/proto/course/v0"
	coursepartpb "github.com/mikhail5545/proto-go/proto/course_part/v0"
	muxpb "github.com/mikhail5545/proto-go/proto/mux_upload/v0"
	productpb "github.com/mikhail5545/proto-go/proto/product/v0"
	seminarpb "github.com/mikhail5545/proto-go/proto/seminar/v0"
	trainingsessionpb "github.com/mikhail5545/proto-go/proto/training_session/v0"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ProductToProtobuf(product *models.Product) *productpb.Product {
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

func ProductToProtobufUpdate(updates map[string]any) *productpb.UpdateResponse {
	resp := &productpb.UpdateResponse{}
	for k, v := range updates {
		switch k {
		case "name":
			if val, ok := v.(string); ok {
				resp.Name = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "name")
			}
		case "description":
			if val, ok := v.(string); ok {
				resp.Description = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "description")
			}
		case "price":
			if val, ok := v.(float32); ok {
				resp.Price = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "price")
			}
		case "amount":
			if val, ok := v.(int32); ok {
				resp.Amount = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "amount")
			}
		case "shipping_required":
			if val, ok := v.(bool); ok {
				resp.ShippingRequired = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "shipping_required")
			}
		}
	}
	return resp
}

func TrainingSessionToProtobuf(ts *models.TrainingSession) *trainingsessionpb.TrainingSession {
	pbTs := &trainingsessionpb.TrainingSession{
		Id:              ts.ID,
		CreatedAt:       timestamppb.New(ts.CreatedAt),
		UpdatedAt:       timestamppb.New(ts.UpdatedAt),
		DurationMinutes: int32(ts.DurationMinutes),
		Format:          ts.Format,
		Product:         ProductToProtobuf(ts.Product),
	}
	return pbTs
}

func TrainingSessionToProtobufUpdate(updates map[string]any, productUpdates map[string]any) *trainingsessionpb.UpdateResponse {
	resp := &trainingsessionpb.UpdateResponse{}
	for k, v := range updates {
		switch k {
		case "duration_minutes":
			if val, ok := v.(int32); ok {
				resp.DurationMinutes = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "duration_minutes")
			}
		case "format":
			if val, ok := v.(string); ok {
				resp.Format = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "format")
			}
		}
	}
	for k, v := range productUpdates {
		switch k {
		case "name":
			if val, ok := v.(string); ok {
				resp.Product.Name = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "product.name")
			}
		case "description":
			if val, ok := v.(string); ok {
				resp.Product.Description = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "product.description")
			}
		case "price":
			if val, ok := v.(float32); ok {
				resp.Product.Price = val
				resp.Updated.Paths = append(resp.Updated.Paths, "product.price")
			}
		}
	}
	return resp
}

func ConvertToProtobufCoursePart(part *models.CoursePart) *coursepartpb.CoursePart {
	pbPart := &coursepartpb.CoursePart{
		Id:          part.ID,
		Name:        part.Name,
		Description: part.Description,
		Number:      int32(part.Number),
		CourseId:    part.CourseID,
	}
	if part.MUXVideoID != nil {
		pbPart.MuxVideoId = part.MUXVideoID
	}
	if part.MUXVideo != nil {
		pbPart.MuxVideo = ConvertToProtobufMUXUpload(part.MUXVideo)
	}
	return pbPart
}

func CourseToProtobuf(course *models.Course) *coursepb.Course {
	pbCourse := &coursepb.Course{
		Id:             course.ID,
		Name:           course.Name,
		Description:    course.Description,
		Topic:          course.Topic,
		ProductId:      course.ProductID,
		AccessDuration: int32(course.AccessDuration),
		Product:        ProductToProtobuf(course.Product),
	}
	for _, part := range course.CourseParts {
		pbCourse.CourseParts = append(pbCourse.CourseParts, ConvertToProtobufCoursePart(part))
	}
	return pbCourse
}

func CourseToProtobufUpdate(updates map[string]any, productUpdates map[string]any) *coursepb.UpdateResponse {
	resp := &coursepb.UpdateResponse{}
	for k, v := range updates {
		switch k {
		case "name":
			if val, ok := v.(string); ok {
				resp.Name = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "name")
			}
		case "description":
			if val, ok := v.(string); ok {
				resp.Description = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "description")
			}
		case "topic":
			if val, ok := v.(string); ok {
				resp.Topic = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "topic")
			}
		case "access_duration":
			if val, ok := v.(int32); ok {
				resp.AccessDuration = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "access_duration")
			}
		}
	}
	for k, v := range productUpdates {
		switch k {
		case "name":
			if val, ok := v.(string); ok {
				resp.Product.Name = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "product.name")
			}
		case "description":
			if val, ok := v.(string); ok {
				resp.Product.Description = &val
				resp.Updated.Paths = append(resp.Updated.Paths, "product.description")
			}
		case "Price":
			if val, ok := v.(float32); ok {
				resp.Product.Price = val
				resp.Updated.Paths = append(resp.Updated.Paths, "product.price")
			}
		}
	}
	return resp
}

func CourseToProtobufListItem(course *models.Course) *coursepb.CourseListItem {
	pbCourseListItem := &coursepb.CourseListItem{
		Id:          course.ID,
		Name:        course.Name,
		Description: course.Description,
		Topic:       course.Topic,
		ProductId:   course.ProductID,
		Product:     ProductToProtobuf(course.Product),
	}
	return pbCourseListItem
}

func ConvertToProtobufMUXUpload(muxUpload *models.MUXUpload) *muxpb.MuxUpload {
	if muxUpload == nil {
		return nil
	}
	pbMuxUpload := &muxpb.MuxUpload{
		Id:                    muxUpload.ID,
		CreatedAt:             timestamppb.New(muxUpload.CreatedAt),
		UpdatedAt:             timestamppb.New(muxUpload.UpdatedAt),
		VideoProcessingStatus: muxUpload.VideoProcessingStatus,
	}
	// Handle optional string pointers
	pbMuxUpload.MuxUploadId = muxUpload.MUXUploadID
	pbMuxUpload.MuxAssetId = muxUpload.MUXAssetID
	pbMuxUpload.MuxPlaybackId = muxUpload.MUXPlaybackID
	pbMuxUpload.AspectRatio = muxUpload.AspectRatio

	// Handle optional scalar pointers with type conversion
	if muxUpload.Duration != nil {
		val := float32(*muxUpload.Duration) // Convert float64 to float32
		pbMuxUpload.Duration = &val
	}
	if muxUpload.MaxWidth != nil {
		val := int32(*muxUpload.MaxWidth) // Convert int to int32
		pbMuxUpload.Width = &val
	}
	if muxUpload.MaxHeight != nil {
		val := int32(*muxUpload.MaxHeight) // Convert int to int32
		pbMuxUpload.Height = &val
	}

	// Handle optional time.Time pointer to timestamppb.Timestamp pointer
	if muxUpload.AssetCreatedAt != nil {
		pbMuxUpload.AssetCreatedAt = timestamppb.New(*muxUpload.AssetCreatedAt)
	}

	return pbMuxUpload
}

func SeminarToProtobuf(seminar *models.Seminar) *seminarpb.Seminar {
	if seminar == nil {
		return nil
	}
	pbSeminar := &seminarpb.Seminar{
		Id:                 seminar.ID,
		CreatedAt:          timestamppb.New(seminar.CreatedAt),
		UpdatedAt:          timestamppb.New(seminar.UpdatedAt),
		Name:               seminar.Name,
		Description:        seminar.Description,
		Place:              seminar.Place,
		Date:               timestamppb.New(seminar.Date),
		EndingDate:         timestamppb.New(seminar.EndingDate),
		Details:            seminar.Details,
		ReservationProduct: ProductToProtobuf(seminar.ReservationProduct),
		EarlyProduct:       ProductToProtobuf(seminar.EarlyProduct),
		LateProduct:        ProductToProtobuf(seminar.LateProduct),
	}
	if seminar.EarlySurchargeProduct != nil {
		pbSeminar.EarlySurchargeProduct = ProductToProtobuf(seminar.EarlySurchargeProduct)
	}
	if seminar.LateSurchargeProduct != nil {
		pbSeminar.LateSurchargeProduct = ProductToProtobuf(seminar.LateSurchargeProduct)
	}
	return pbSeminar
}

func SeminarToProtobufUpdate(updates map[string]map[string]any) *seminarpb.UpdateResponse {
	resp := &seminarpb.UpdateResponse{}
	for k, v := range updates {
		switch k {
		case "seminar":
			for sk, sv := range v {
				switch sk {
				case "name":
					if val, ok := sv.(string); ok {
						resp.Name = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "name")
					}
				case "description":
					if val, ok := sv.(string); ok {
						resp.Description = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "description")
					}
				case "date":
					if val, ok := sv.(*timestamppb.Timestamp); ok {
						resp.Date = val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "date")
					}
				case "ending_date":
					if val, ok := sv.(*timestamppb.Timestamp); ok {
						resp.EndingDate = val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "ending_date")
					}
				case "late_payment_date":
					if val, ok := sv.(*timestamppb.Timestamp); ok {
						resp.LatePaymentDate = val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "late_payment_date")
					}
				case "place":
					if val, ok := sv.(string); ok {
						resp.Place = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "place")
					}
				case "details":
					if val, ok := sv.(string); ok {
						resp.Details = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "details")
					}
				}
			}
		case "reservation_product":
			for rk, rv := range v {
				switch rk {
				case "name":
					if val, ok := rv.(string); ok {
						resp.ReservationProduct.Name = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "reservation_product.name")
					}
				case "description":
					if val, ok := rv.(string); ok {
						resp.ReservationProduct.Description = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "reservation_product.description")
					}
				case "price":
					if val, ok := rv.(float32); ok {
						resp.ReservationProduct.Price = val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "reservation_product.price")
					}
				}
			}
		case "early_product":
			for ek, ev := range v {
				switch ek {
				case "name":
					if val, ok := ev.(string); ok {
						resp.EarlyProduct.Name = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "early_product.name")
					}
				case "description":
					if val, ok := ev.(string); ok {
						resp.EarlyProduct.Description = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "early_product.description")
					}
				case "price":
					if val, ok := ev.(float32); ok {
						resp.EarlyProduct.Price = val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "early_product.price")
					}
				}
			}
		case "late_product":
			for lk, lv := range v {
				switch lk {
				case "name":
					if val, ok := lv.(string); ok {
						resp.LateProduct.Name = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "late_product.name")
					}
				case "description":
					if val, ok := lv.(string); ok {
						resp.LateProduct.Description = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "late_product.description")
					}
				case "price":
					if val, ok := lv.(float32); ok {
						resp.LateProduct.Price = val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "late_product.price")
					}
				}
			}
		case "early_surcharge_product":
			for esk, esv := range v {
				switch esk {
				case "name":
					if val, ok := esv.(string); ok {
						resp.EarlySurchargeProduct.Name = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "early_surcharge_product.name")
					}
				case "description":
					if val, ok := esv.(string); ok {
						resp.EarlySurchargeProduct.Description = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "early_surcharge_product.description")
					}
				case "price":
					if val, ok := esv.(float32); ok {
						resp.EarlySurchargeProduct.Price = val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "early_surcharge_product.price")
					}
				}
			}
		case "late_surcharge_product":
			for lsk, lsv := range v {
				switch lsk {
				case "name":
					if val, ok := lsv.(string); ok {
						resp.LateSurchargeProduct.Name = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "late_surcharge_product.name")
					}
				case "description":
					if val, ok := lsv.(string); ok {
						resp.LateSurchargeProduct.Description = &val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "late_surcharge_product.description")
					}
				case "price":
					if val, ok := lsv.(float32); ok {
						resp.LateSurchargeProduct.Price = val
						resp.UpdateMask.Paths = append(resp.UpdateMask.Paths, "late_surcharge_product.price")
					}
				}
			}
		}
	}
	return resp
}

func ConvertFromProtobufMUXUpload(muxUpload *muxpb.MuxUpload) *models.MUXUpload {
	if muxUpload == nil {
		return nil
	}
	upload := &models.MUXUpload{
		ID:                    muxUpload.GetId(),                    // Non-optional string
		CreatedAt:             muxUpload.GetCreatedAt().AsTime(),    // Non-optional Timestamp
		UpdatedAt:             muxUpload.GetUpdatedAt().AsTime(),    // Non-optional Timestamp
		VideoProcessingStatus: muxUpload.GetVideoProcessingStatus(), // Non-optional string
	}
	// Optional string fields (pointers) can be directly assigned
	upload.MUXUploadID = muxUpload.MuxUploadId
	upload.MUXAssetID = muxUpload.MuxAssetId
	upload.MUXPlaybackID = muxUpload.MuxPlaybackId
	upload.AspectRatio = muxUpload.AspectRatio

	// Optional scalar fields require HasX() check and type conversion
	if muxUpload.Duration != nil {
		duration := float64(muxUpload.GetDuration()) // Convert float32 to float64
		upload.Duration = &duration
	}
	if muxUpload.Width != nil {
		width := int(muxUpload.GetWidth()) // Convert int32 to int
		upload.MaxWidth = &width
	}
	if muxUpload.Height != nil {
		height := int(muxUpload.GetHeight()) // Convert int32 to int
		upload.MaxHeight = &height
	}

	// Optional Timestamp field requires nil check and conversion
	if muxUpload.AssetCreatedAt != nil {
		assetCreatedAt := muxUpload.GetAssetCreatedAt().AsTime()
		upload.AssetCreatedAt = &assetCreatedAt
	}
	return upload
}
