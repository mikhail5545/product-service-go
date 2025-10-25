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
	admincourse "github.com/mikhail5545/product-service-go/internal/handlers/admin/course"
	admincp "github.com/mikhail5545/product-service-go/internal/handlers/admin/course_part"
	adminproduct "github.com/mikhail5545/product-service-go/internal/handlers/admin/product"
	adminseminar "github.com/mikhail5545/product-service-go/internal/handlers/admin/seminar"
	admints "github.com/mikhail5545/product-service-go/internal/handlers/admin/training_session"
	publiccourse "github.com/mikhail5545/product-service-go/internal/handlers/public/course"
	publiccp "github.com/mikhail5545/product-service-go/internal/handlers/public/course_part"
	publicproduct "github.com/mikhail5545/product-service-go/internal/handlers/public/product"
	publicseminar "github.com/mikhail5545/product-service-go/internal/handlers/public/seminar"
	publicts "github.com/mikhail5545/product-service-go/internal/handlers/public/training_session"
	"github.com/mikhail5545/product-service-go/internal/services/course"
	coursepart "github.com/mikhail5545/product-service-go/internal/services/course_part"
	"github.com/mikhail5545/product-service-go/internal/services/product"
	"github.com/mikhail5545/product-service-go/internal/services/seminar"
	trainingsession "github.com/mikhail5545/product-service-go/internal/services/training_session"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
)

func Setup(
	e *echo.Echo,
	productService *product.Service,
	cpService *coursepart.Service,
	tsService *trainingsession.Service,
	courseService *course.Service,
	seminarService *seminar.Service,
) {
	e.HTTPErrorHandler = errors.HTTPErrorHandler

	api := e.Group("/api")
	ver := api.Group("/v0")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --- Public handlers ---
	productHandler := publicproduct.New(productService)
	cpHandler := publiccp.New(cpService)
	tsHandler := publicts.New(tsService)
	courseHandler := publiccourse.New(courseService)
	seminarHandler := publicseminar.New(seminarService)

	// --- Admin handlers ---
	adminProductHandler := adminproduct.New(productService)
	admincpHandler := admincp.New(cpService)
	admintsHandler := admints.New(tsService)
	adminCourseHandler := admincourse.New(courseService)
	adminSeminarHandler := adminseminar.New(seminarService)

	products := ver.Group("/products")
	{
		products.GET("", productHandler.List)
		products.GET("/:id", productHandler.Get)
	}
	trainingSesssions := ver.Group("/training-sessions")
	{
		trainingSesssions.GET("", tsHandler.List)
		trainingSesssions.GET("/:id", tsHandler.Get)
	}
	courses := ver.Group("/courses")
	{
		courses.GET("", courseHandler.List)
		courses.GET("/:id", courseHandler.Get)
	}
	course_parts := ver.Group("/course-parts")
	{
		course_parts.GET("/:cid", cpHandler.List)
		course_parts.GET("/:id", cpHandler.Get)
	}
	seminars := ver.Group("/seminars")
	{
		seminars.GET("", seminarHandler.List)
		seminars.GET("/:id", seminarHandler.Get)
	}

	admin := ver.Group("/admin")
	{
		adminProducts := admin.Group("/products")
		{
			adminProducts.GET("", adminProductHandler.List)
			adminProducts.GET("/:id", adminProductHandler.Get)
			adminProducts.POST("", adminProductHandler.Create)
			adminProducts.PATCH("/:id", adminProductHandler.Update)
			adminProducts.DELETE("/:id", adminProductHandler.Delete)
		}
		adminTrainingSessions := admin.Group("/training-sessions")
		{
			adminTrainingSessions.GET("", admintsHandler.List)
			adminTrainingSessions.GET("/:id", admintsHandler.Get)
			adminTrainingSessions.POST("", admintsHandler.Create)
			adminTrainingSessions.PATCH("/:id", admintsHandler.Update)
			adminTrainingSessions.DELETE("/:id", admintsHandler.Delete)
		}
		adminCourses := admin.Group("/courses")
		{
			adminCourses.GET("", adminCourseHandler.List)
			adminCourses.GET("/:id", adminCourseHandler.Get)
			adminCourses.POST("", adminCourseHandler.Create)
			adminCourses.PATCH("/:id", adminCourseHandler.Update)
		}
		adminCourseParts := admin.Group("/course-parts")
		{
			adminCourseParts.GET("", admincpHandler.List)
			adminCourseParts.GET("/:id", admincpHandler.Get)
		}
		adminSeminars := admin.Group("/seminars")
		{
			adminSeminars.GET("", adminSeminarHandler.List)
			adminSeminars.GET("/:id", adminSeminarHandler.Get)
			adminSeminars.POST("", adminSeminarHandler.Create)
			adminSeminars.PATCH("/:id", adminSeminarHandler.Update)
			adminSeminars.DELETE("/:id", adminSeminarHandler.Delete)
		}
	}
}
