package utils

import (
	"google.golang.org/protobuf/types/known/timestamppb"
	"vitainmove.com/product-service-go/internal/models"
	productpb "vitainmove.com/product-service-go/proto/product/v0"
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
