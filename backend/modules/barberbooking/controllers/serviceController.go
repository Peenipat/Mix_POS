// barberbooking/controller/service_controller.go
package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	barberBookingService "myapp/modules/barberbooking/services"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type ServiceController struct {
	ServiceService *barberBookingService.ServiceService
}

func NewServiceController(svc *barberBookingService.ServiceService) *ServiceController {
	return &ServiceController{ServiceService: svc}
}

func (ctrl *ServiceController) GetAllServices(c *fiber.Ctx) error {
	services, err := ctrl.ServiceService.GetAllServices()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch services",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Services retrieved",
		"data":    services,
	})
}

func (ctrl *ServiceController) GetServiceByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service ID",
		})
	}

	service, err := ctrl.ServiceService.GetServiceByID(uint(id))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch service",
			"error":   err.Error(),
		})
	}
	if service == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Service not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Service retrieved",
		"data":    service,
	})
}


func (ctrl *ServiceController) CreateService(c *fiber.Ctx) error {
	var payload barberBookingModels.Service
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if payload.Name == "" || payload.Duration <= 0 || payload.Price < 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service data",
		})
	}

	payload.CreatedAt = time.Now()
	payload.UpdatedAt = time.Now()

	if err := ctrl.ServiceService.CreateService(&payload); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create service",
			"error":   err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Service created successfully",
		"data":    payload,
	})
}

func (ctrl *ServiceController) UpdateService(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service ID",
		})
	}

	var payload barberBookingModels.Service
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	updated, err := ctrl.ServiceService.UpdateService(uint(id), &payload)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update service",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Service updated",
		"data":    updated,
	})
}

func (ctrl *ServiceController) DeleteService(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service ID",
		})
	}

	if err := ctrl.ServiceService.DeleteService(uint(id)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete service",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Service deleted successfully",
	})
}