package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"maxpark_opp_registration/config"
	"maxpark_opp_registration/models"
	"maxpark_opp_registration/services"
	"maxpark_opp_registration/utils"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type ResidentHandler struct {
	natsService *services.NATSService
	config      *config.Config
}

func NewResidentHandler(natsService *services.NATSService, cfg *config.Config) *ResidentHandler {
	return &ResidentHandler{natsService: natsService, config: cfg}
}

func (h *ResidentHandler) HandleCreateResidentRegister(c echo.Context) error {
	form := new(models.ResidentRegisterForm)

	if err := c.Bind(form); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid input format", err.Error())
	}

	if err := c.Validate(form); err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", he.Message)
		}
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	natsPayload := map[string]interface{}{
		"vehicleNum": strings.ToUpper(strings.ReplaceAll(form.ResidentPlates.PlateNumber, " ", "")),
	}

	data, err := json.Marshal(natsPayload)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to marshal NATS payload", err.Error())
	}

	subject := fmt.Sprintf("%s.%s", h.config.RegisterIndividualSubject, h.config.SiteCode)
	natsResponse, err := h.natsService.GetConnection().Request(subject, data, 10*time.Second)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send NATS message", err.Error())
	}

	var natsResponseData map[string]interface{}
	if err := json.Unmarshal(natsResponse.Data, &natsResponseData); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to parse NATS response", err.Error())
	}

	// Check if validation failed from NATS
	if success, ok := natsResponseData["success"].(bool); ok && !success {
		if errObj, ok := natsResponseData["error"].(map[string]interface{}); ok {
			if message, ok := errObj["message"].(string); ok {
				// Transform the error to match validation error format
				validationError := map[string]interface{}{
					"message": "validate failed",
					"errors": []map[string]string{
						{
							"field":   "ResidentPlates.PlateNumber",
							"message": message,
						},
					},
				}
				return utils.ErrorResponse(c, 400, "Validation failed", validationError)
			}
		}
	}

	registerID := uuid.New().String()
	responseData := map[string]interface{}{
		"registerID": registerID,
	}
	return utils.SuccessResponse(c, "Vehicle validation completed", responseData)
}

func (h *ResidentHandler) HandleCreateResidentRegisterFinalize(c echo.Context) error {
	form := new(models.ResidentRegisterForm)

	if err := c.Bind(form); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid input format", err.Error())
	}

	if err := c.Validate(form); err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", he.Message)
		}
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	natsPayload := map[string]interface{}{
		"nric":             form.NricNumber,
		"tinNumber":        form.TinNumber,
		"fullName":         strings.ToUpper(form.ResidentName),
		"email":            form.ContactEmail,
		"contactNumber":    form.ContactNumber,
		"address1":         form.ResidentAddressLine1,
		"address2":         form.ResidentAddressLine2,
		"vehicleNum":       strings.ToUpper(strings.ReplaceAll(form.ResidentPlates.PlateNumber, " ", "")),
		"vehicleClass":     form.ResidentPlates.VehicleType,
		"vehiclePath":      form.ResidentPlates.VehiclePath,
		"spaPath":          form.ResidentSupportingFiles.SPAPath,
		"electricBillPath": form.ResidentSupportingFiles.ElectricBillPath,
	}

	data, err := json.Marshal(natsPayload)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to marshal NATS payload", err.Error())
	}

	subject := fmt.Sprintf("%s.%s", h.config.RegisterIndividualFinalizeSubject, h.config.SiteCode)
	natsResponse, err := h.natsService.GetConnection().Request(subject, data, 10*time.Second)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send NATS message", err.Error())
	}

	var natsResponseData interface{}
	if err := json.Unmarshal(natsResponse.Data, &natsResponseData); err != nil {
		natsResponseData = string(natsResponse.Data)
	}

	responseData := map[string]interface{}{
		"residentName": form.ResidentName,
		"natsResponse": natsResponseData,
	}
	return utils.SuccessResponse(c, "Resident registration finalized successfully", responseData)
}
