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
	adminphysicalgood "github.com/mikhail5545/product-service-go/internal/handlers/admin/physical_good"
	adminseminar "github.com/mikhail5545/product-service-go/internal/handlers/admin/seminar"
	admints "github.com/mikhail5545/product-service-go/internal/handlers/admin/training_session"
	publiccourse "github.com/mikhail5545/product-service-go/internal/handlers/public/course"
	publiccp "github.com/mikhail5545/product-service-go/internal/handlers/public/course_part"
	publicphysicalgood "github.com/mikhail5545/product-service-go/internal/handlers/public/physical_good"
	publicseminar "github.com/mikhail5545/product-service-go/internal/handlers/public/seminar"
	publicts "github.com/mikhail5545/product-service-go/internal/handlers/public/training_session"
	"github.com/mikhail5545/product-service-go/internal/services/course"
	coursepart "github.com/mikhail5545/product-service-go/internal/services/course_part"
	physicalgood "github.com/mikhail5545/product-service-go/internal/services/physical_good"
	"github.com/mikhail5545/product-service-go/internal/services/product"
	"github.com/mikhail5545/product-service-go/internal/services/seminar"
	trainingsession "github.com/mikhail5545/product-service-go/internal/services/training_session"
	"github.com/mikhail5545/product-service-go/internal/util/errors"
)

func Setup(
	e *echo.Echo,
	productService product.Service,
	cpService coursepart.Service,
	tsService trainingsession.Service,
	courseService course.Service,
	seminarService seminar.Service,
	phgService physicalgood.Service,
) {
	e.HTTPErrorHandler = errors.HTTPErrorHandler

	api := e.Group("/api")
	ver := api.Group("/v0")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// --- Public handlers ---
	phgHandler := publicphysicalgood.New(phgService)
	cpHandler := publiccp.New(cpService)
	tsHandler := publicts.New(tsService)
	courseHandler := publiccourse.New(courseService)
	seminarHandler := publicseminar.New(seminarService)

	// --- Admin handlers ---
	adminphgHandler := adminphysicalgood.New(phgService)
	admincpHandler := admincp.New(cpService)
	admintsHandler := admints.New(tsService)
	adminCourseHandler := admincourse.New(courseService)
	adminSeminarHandler := adminseminar.New(seminarService)

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
	physicalGoods := ver.Group("/physical-good")
	{
		physicalGoods.GET("", phgHandler.List)
		physicalGoods.GET("/:id", phgHandler.Get)
	}
	admin := ver.Group("/admin")
	{
		adminPhysicalGoods := admin.Group("/physical-good")
		{
			adminPhysicalGoods.GET("", adminphgHandler.List)
			adminPhysicalGoods.GET("/deleted", adminphgHandler.ListDeleted)
			adminPhysicalGoods.GET("/unpublished", adminphgHandler.ListUnpublished)
			adminPhysicalGoods.GET("/:id", adminphgHandler.Get)
			adminPhysicalGoods.GET("/deleted/:id", adminphgHandler.GetWithDeleted)
			adminPhysicalGoods.GET("/unpublished/:id", adminphgHandler.GetWithUnpublished)
			adminPhysicalGoods.POST("", adminphgHandler.Create)
			adminPhysicalGoods.PATCH("/:id", adminphgHandler.Update)
			adminPhysicalGoods.POST("/publish/:id", adminphgHandler.Publish)
			adminPhysicalGoods.POST("/unpublish/:id", adminphgHandler.Unpublish)
			adminPhysicalGoods.POST("/restore/:id", adminphgHandler.Restore)
			adminPhysicalGoods.DELETE("/:id", adminphgHandler.Delete)
			adminPhysicalGoods.DELETE("/permanent/:id", adminphgHandler.DeletePermanent)
		}
		adminTrainingSessions := admin.Group("/training-sessions")
		{
			adminTrainingSessions.GET("", admintsHandler.List)
			adminTrainingSessions.GET("/deleted", admintsHandler.ListDeleted)
			adminTrainingSessions.GET("/unpublished", admintsHandler.ListUnpublished)
			adminTrainingSessions.GET("/:id", admintsHandler.Get)
			adminTrainingSessions.GET("/deleted/:id", admintsHandler.GetWithDeleted)
			adminTrainingSessions.GET("/unpublished/:id", admintsHandler.GetWithUnpublished)
			adminTrainingSessions.POST("", admintsHandler.Create)
			adminTrainingSessions.PATCH("/:id", admintsHandler.Update)
			adminTrainingSessions.POST("/publish/:id", admintsHandler.Publish)
			adminTrainingSessions.POST("/unpublish/:id", admintsHandler.Unpublish)
			adminTrainingSessions.POST("/restore/:id", admintsHandler.Restore)
			adminTrainingSessions.DELETE("/:id", admintsHandler.Delete)
			adminTrainingSessions.DELETE("/permanent/:id", admintsHandler.DeletePermanent)
		}
		adminCourses := admin.Group("/courses")
		{
			adminCourses.GET("", adminCourseHandler.List)
			adminCourses.GET("/deleted", adminCourseHandler.ListDeleted)
			adminCourses.GET("/unpublished", adminCourseHandler.ListUnpublished)
			adminCourses.GET("/:id", adminCourseHandler.Get)
			adminCourses.GET("/deleted/:id", adminCourseHandler.GetWithDeleted)
			adminCourses.GET("/unpublished/:id", adminCourseHandler.GetWithUnpublished)
			adminCourses.POST("", adminCourseHandler.Create)
			adminCourses.PATCH("/:id", adminCourseHandler.Update)
			adminCourses.POST("/publish/:id", adminCourseHandler.Publish)
			adminCourses.POST("/unpublish/:id", adminCourseHandler.Unpublish)
			adminCourses.DELETE("/:id", adminCourseHandler.Delete)
			adminCourses.DELETE("/permanent/:id", adminCourseHandler.DeletePermanent)
			adminCourses.POST("restore/:id", adminCourseHandler.Restore)
			// --- Course parts assigned to the course ---
			adminCourses.GET("/:cid/parts/", admincpHandler.List)
			adminCourses.GET("/:cid/parts/deleted", admincpHandler.ListDeleted)
			adminCourses.GET("/:cid/parts/unpublished", admincpHandler.ListUnpublished)
			adminCourses.POST("/:cid/parts", admincpHandler.Create)
		}
		adminCourseParts := admin.Group("/course-parts")
		{
			adminCourseParts.GET("/:id", admincpHandler.Get)
			adminCourseParts.GET("/deleted/:id", admincpHandler.GetWithDeleted)
			adminCourseParts.GET("/unpublished/:id", admincpHandler.GetWithUnpublished)
			adminCourseParts.POST("/publish/:id", admincpHandler.Publish)
			adminCourseParts.POST("/unpublish/:id", admincpHandler.Unpublish)
			adminCourseParts.POST("/restore/:id", admincpHandler.Restore)
			adminCourseParts.PATCH("/:id", admincpHandler.Update)
			adminCourseParts.DELETE("/:id", admincpHandler.Delete)
			adminCourseParts.DELETE("/permanent/:id", admincpHandler.DeletePermanent)
		}
		adminSeminars := admin.Group("/seminars")
		{
			adminSeminars.GET("", adminSeminarHandler.List)
			adminSeminars.GET("/deleted", adminSeminarHandler.ListDeleted)
			adminSeminars.GET("/unpublished", adminSeminarHandler.ListUnpublished)
			adminSeminars.GET("/:id", adminSeminarHandler.Get)
			adminSeminars.GET("/deleted/:id", adminSeminarHandler.GetWithDeleted)
			adminSeminars.GET("/unpublished/:id", adminSeminarHandler.GetWithUnpublished)
			adminSeminars.POST("", adminSeminarHandler.Create)
			adminSeminars.PATCH("/:id", adminSeminarHandler.Update)
			adminSeminars.POST("/publish/:id", adminSeminarHandler.Publish)
			adminSeminars.POST("/unpublish/:id", adminSeminarHandler.Unpublish)
			adminSeminars.POST("/restore/:id", adminSeminarHandler.Restore)
			adminSeminars.DELETE("/:id", adminSeminarHandler.Delete)
			adminSeminars.DELETE("/permanent/:id", adminSeminarHandler.DeletePermanent)
		}
	}
}
