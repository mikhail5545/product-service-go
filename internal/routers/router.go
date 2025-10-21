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

package routers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/mikhail5545/product-service-go/internal/handlers"
	adminhandlers "github.com/mikhail5545/product-service-go/internal/handlers/admin"
	publichandlers "github.com/mikhail5545/product-service-go/internal/handlers/public"
	"github.com/mikhail5545/product-service-go/internal/services"
)

func SetupRouter(
	e *echo.Echo, productService *services.ProductService,
	tsService *services.TrainingSessionService,
	courseService *services.CourseService,
	seminarService *services.SeminarService,
) {
	e.HTTPErrorHandler = handlers.HTTPErrorHandler

	api := e.Group("/api")
	ver := api.Group("/v0")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --- Public handlers ---
	productHandler := publichandlers.NewProductHandler(productService)
	tsHandler := publichandlers.NewTrainingSessionHandler(tsService)
	courseHandler := publichandlers.NewCourseHandler(courseService)
	seminarHandler := publichandlers.NewSeminarHandler(seminarService)

	// --- Admin handlers ---
	adminProductHandler := adminhandlers.NewProductHandler(productService)
	adminTrainingSessionHandler := adminhandlers.NewTrainingSessionHandler(tsService)
	adminCourseHandler := adminhandlers.NewCourseHandler(courseService)
	adminSeminarHandler := adminhandlers.NewSeminarService(seminarService)

	products := ver.Group("/products")
	{
		products.GET("", productHandler.GetProducts)
		products.GET("/:id", productHandler.GetProduct)
	}
	trainingSesssions := ver.Group("/training-sessions")
	{
		trainingSesssions.GET("", tsHandler.GetTrainingSessions)
		trainingSesssions.GET("/:id", tsHandler.GetTrainingSession)
	}
	courses := ver.Group("/courses")
	{
		courses.GET("/:id", courseHandler.GetCourse)
	}
	seminars := ver.Group("/seminars")
	{
		seminars.GET("", seminarHandler.GetSeminars)
		seminars.GET("/:id", seminarHandler.GetSeminar)
	}

	admin := ver.Group("/admin")
	{
		adminProducts := admin.Group("/products")
		{
			adminProducts.GET("", adminProductHandler.GetProducts)
			adminProducts.GET("/:id", adminProductHandler.GetProduct)
			adminProducts.POST("", adminProductHandler.CreateProduct)
			adminProducts.PUT("/:id", adminProductHandler.UpdateProduct)
			adminProducts.DELETE("/:id", adminProductHandler.DeleteProduct)
		}
		adminTrainingSessions := admin.Group("/training-sessions")
		{
			adminTrainingSessions.GET("", adminTrainingSessionHandler.GetTrainingSessions)
			adminTrainingSessions.GET("/:id", adminTrainingSessionHandler.GetTrainingSession)
			adminTrainingSessions.POST("", adminTrainingSessionHandler.CreateTrainingSession)
			adminTrainingSessions.PUT("/:id", adminTrainingSessionHandler.UpdateTrainingSession)
			adminTrainingSessions.DELETE("/:id", adminTrainingSessionHandler.DeleteTrainingSession)
		}
		adminCourses := admin.Group("/courses")
		{
			adminCourses.GET("", adminCourseHandler.GetCourses)
			adminCourses.GET("/:id", adminCourseHandler.GetCourse)
			adminCourses.POST("", adminCourseHandler.CreateCourse)
			adminCourses.PUT("/:id", adminCourseHandler.UpdateCourse)
		}
		adminSeminars := admin.Group("/seminars")
		{
			adminSeminars.GET("", adminSeminarHandler.GetSeminars)
			adminSeminars.GET("/:id", adminSeminarHandler.GetSeminar)
			adminSeminars.POST("", adminSeminarHandler.CreateSeminar)
			adminSeminars.PUT("/:id", adminSeminarHandler.UpdateSeminar)
			adminSeminars.DELETE("/:id", adminSeminarHandler.DeleteSeminar)
		}
	}
}
