package seminar

import (
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

// Validate validates fields of [seminar.CreateRequest].
// All request fields are required for creation.
// Validation rules:
//
//   - Name: required, 3-255 characters, Alpha only.
//   - ShortDescription: required, 3-255 characters.
//   - ReservationPrice: required, >= 1.
//   - EarlyPrice: required, >= 1.
//   - LatePrice: required, >= 1.
//   - EarlySurchargePrice: required, >= 1.
//   - LateSurchargePrice: required, >= 1.
//   - Date: required, at least 48 hours from now.
//   - EndingDate: required, at least 1 hour after Date.
//   - LatePaymentDate: required, at least 24 hours from now, max 24 hours before Date.
//   - Place: required, 3-255 characters.
func (req *CreateRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(
			&req.Name,
			validation.Required,
			validation.Length(3, 255),
			is.Alphanumeric,
		),
		validation.Field(
			&req.ShortDescription,
			validation.Required,
			validation.Length(3, 255),
		),
		validation.Field(
			&req.ReservationPrice,
			validation.Required,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.EarlyPrice,
			validation.Required,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.LatePrice,
			validation.Required,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.EarlySurchargePrice,
			validation.Required,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.LateSurchargePrice,
			validation.Required,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.Date,
			validation.Required,
			validation.Min(time.Now().Add(time.Duration(48)*time.Hour)),
		),
		validation.Field(
			&req.EndingDate,
			validation.Required,
			validation.Min(req.Date.Add(time.Duration(1)*time.Hour)),
		),
		validation.Field(
			&req.LatePaymentDate,
			validation.Required,
			validation.Min(time.Now().Add(time.Duration(24)*time.Hour)),
			validation.By(func(value interface{}) error {
				if date, ok := value.(time.Time); ok {
					if req.Date.Sub(date) < (time.Duration(24) * time.Hour) {
						return errors.New("must be at least 24 hours before the seminar date")
					}
				}
				return nil
			}),
		),
		validation.Field(
			&req.Place,
			validation.Required,
			validation.Length(3, 255),
		),
	)
}

// Validate validates fields of [seminar.UpdateRequest].
// All request fields except ID are optional.
// Validation rules:
//
//   - ID: required, UUID
//   - Name: optional, 3-255 characters, Alpha only.
//   - ShortDescription: optional, 3-255 characters.
//   - LongDescription: optional, 3-3000 characters.
//   - ReservationPrice: optional, >= 1.
//   - EarlyPrice: optional, >= 1.
//   - LatePrice: optional, >= 1.
//   - EarlySurchargePrice: optional, >= 1.
//   - LateSurchargePrice: optional, >= 1.
//   - Date: optional, at least 48 hours from now.
//   - EndingDate: optional, at least 1 hour after Date.
//   - LatePaymentDate: optional, at least 24 hours from now, max 24 hours before Date.
//   - Place: optional, 3-255 characters.
//   - Tags: optional, 1-10 items, 3-20 characters each.
func (req *UpdateRequest) Validate() error {
	return validation.ValidateStruct(&req,
		validation.Field(
			&req.ID,
			validation.Required,
			is.UUID,
		),
		validation.Field(
			&req.Name,
			validation.Length(3, 255),
			is.Alpha,
		),
		validation.Field(
			&req.ShortDescription,
			validation.Length(3, 255),
		),
		validation.Field(
			&req.LongDescription,
			validation.Length(3, 3000),
		),
		validation.Field(
			&req.ReservationPrice,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.EarlyPrice,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.LatePrice,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.EarlySurchargePrice,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.LateSurchargePrice,
			validation.Min(float32(1)),
		),
		validation.Field(
			&req.Date,
			validation.Min(time.Now().Add(time.Duration(48)*time.Hour)),
		),
		validation.Field(
			&req.EndingDate,
			validation.When(req.Date != nil, validation.Min(req.Date.Add(time.Duration(1)*time.Hour))),
		),
		validation.Field(
			&req.LatePaymentDate,
			validation.When(req.LatePaymentDate != nil,
				validation.Min(time.Now().Add(time.Duration(24)*time.Hour))),
			validation.When(req.LatePaymentDate != nil && req.Date != nil,
				validation.By(func(value interface{}) error {
					if latePaymentDate, ok := value.(*time.Time); ok && latePaymentDate != nil {
						if req.Date.Sub(*latePaymentDate) < (time.Duration(24) * time.Hour) {
							return errors.New("must be at least 24 hours before the seminar date")
						}
					}
					return nil
				}),
			),
		),
		validation.Field(
			&req.Place,
			validation.Length(3, 255),
		),
		validation.Field(
			&req.Tags,
			validation.Length(1, 10),
			validation.Each(validation.Length(3, 20), is.Alphanumeric),
		),
	)
}
