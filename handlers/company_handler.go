package handlers

import (
	"encoding/json"
	"fmt"
	"log"
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

type CompanyHandler struct {
	natsService    *services.NATSService
	config         *config.Config
	receiptService *services.ReceiptService
}

func NewCompanyHandler(natsService *services.NATSService, cfg *config.Config) *CompanyHandler {
	return &CompanyHandler{natsService: natsService, config: cfg}
}

func (h *CompanyHandler) HandleCreateCompanyRegister(c echo.Context) error {
	form := new(models.CompanyRegisterForm)

	if err := c.Bind(form); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid input format", err.Error())
	}

	if err := c.Validate(form); err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", he.Message)
		}
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	// Normalize data before creating payload
	form.CompanyName = strings.ToUpper(strings.TrimSpace(form.CompanyName))
	form.ContactPerson = strings.ToUpper(strings.TrimSpace(form.ContactPerson))

	registrationPayload := map[string]interface{}{
		"employerID":             uuid.New().String(),
		"employerName":           form.CompanyName,
		"contactPerson":          form.ContactPerson,
		"contactNumber":          form.ContactNumber,
		"address1":               form.CompanyAddressLine1,
		"address2":               form.CompanyAddressLine2,
		"tinNumber":              form.TinNumber,
		"email":                  form.ContactEmail,
		"companySupportingFiles": form.CompanySupportingFiles,
		"individuals":            []map[string]interface{}{},
		"idType":                 form.IDType,
		"idNumber":               form.IDNumber,
		"vehiclePassType":        form.VehiclePassType,
	}

	for _, plate := range form.CompanyPlates {
		individual := map[string]interface{}{
			"fullName":         form.ContactPerson,
			"email":            form.ContactEmail,
			"contactNumber":    form.ContactNumber,
			"address1":         form.CompanyAddressLine1,
			"address2":         form.CompanyAddressLine2,
			"nric":             plate.NricNumber,
			"vehicleNum":       strings.ToUpper(strings.ReplaceAll(plate.PlateNumber, " ", "")),
			"tinNumber":        form.TinNumber,
			"vehicleClass":     plate.VehicleType,
			"spaPath":          plate.SPAPath,
			"electricBillPath": plate.ElectricBillPath,
			"vehiclePath":      plate.VehiclePath,
		}
		registrationPayload["individuals"] = append(registrationPayload["individuals"].([]map[string]interface{}), individual)
	}

	registrationData, err := json.Marshal(registrationPayload)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to marshal NATS payload", err.Error())
	}

	registrationSubject := fmt.Sprintf("%s.%s", h.config.RegisterEmployerSubject, h.config.SiteCode)
	registrationResponse, err := h.natsService.GetConnection().Request(registrationSubject, registrationData, 10*time.Second)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send NATS message", err.Error())
	}

	var registrationResponseData map[string]interface{}
	if err := json.Unmarshal(registrationResponse.Data, &registrationResponseData); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to parse NATS response", err.Error())
	}

	success, ok := registrationResponseData["success"].(bool)
	if !ok || !success {
		// Handle failure response
		errorData := registrationResponseData["error"].(map[string]interface{})
		errorCode := ""
		errorMsg := ""

		if code, ok := errorData["code"].(string); ok {
			errorCode = code
		}
		if msg, ok := errorData["message"].(string); ok {
			errorMsg = msg
		}

		// Extract vehicle plates from data.errors
		var failedPlates []string
		if dataField, ok := registrationResponseData["data"].(map[string]interface{}); ok {
			if errors, ok := dataField["errors"].([]interface{}); ok {
				for _, err := range errors {
					if errMap, ok := err.(map[string]interface{}); ok {
						if plates, ok := errMap["plates"].([]interface{}); ok {
							for _, plate := range plates {
								if plateStr, ok := plate.(string); ok {
									failedPlates = append(failedPlates, plateStr)
								}
							}
						}
					}
				}
			}
		}

		log.Printf("NATS registration failed - Code: %s, Message: %s, Failed Plates: %v", errorCode, errorMsg, failedPlates)

		// Include failed plates in error response
		errorDetail := map[string]interface{}{
			"message": errorMsg,
			"plates":  failedPlates,
		}
		return utils.ErrorResponse(c, 400, "Company registration failed", errorDetail)
	}

	employerIDsubject := fmt.Sprintf("%s.%s", h.config.RegisterEmployerIDSubject, h.config.SiteCode)
	employerIDResponse, err := h.natsService.GetConnection().Request(employerIDsubject, nil, 10*time.Second)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send NATS message", err.Error())
	}

	var response map[string]interface{}
	if err := json.Unmarshal(employerIDResponse.Data, &response); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to parse NATS response", err.Error())
	}

	employerID, ok := response["data"].(string)
	if !ok {
		return utils.ErrorResponse(c, 500, "Invalid response format", "Missing or invalid ID")
	}

	log.Printf("Created company name %s path %s", form.CompanyName, employerID)
	for _, plate := range form.CompanyPlates {
		log.Printf(" - with plate number: %s", plate.PlateNumber)
	}

	registerID := uuid.New().String()
	responseData := map[string]string{
		"registerID": registerID,
		"employerID": employerID,
	}
	return utils.SuccessResponse(c, "Company registration initiated successfully", responseData)
}

func (h *CompanyHandler) HandleCreateCompanyRegisterFinalize(c echo.Context) error {
	form := new(models.CompanyRegisterForm)

	if err := c.Bind(form); err != nil {
		return utils.ErrorResponse(c, 400, "Invalid input format", err.Error())
	}

	if err := c.Validate(form); err != nil {
		if he, ok := err.(*echo.HTTPError); ok {
			return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", he.Message)
		}
		return utils.ErrorResponse(c, http.StatusBadRequest, "Validation failed", err.Error())
	}

	// Normalize data before creating payload
	form.CompanyName = strings.ToUpper(strings.TrimSpace(form.CompanyName))
	form.ContactPerson = strings.ToUpper(strings.TrimSpace(form.ContactPerson))

	natsPayload := map[string]interface{}{
		"employerID":             form.EmployerID,
		"tinNumber":              form.TinNumber,
		"employerName":           form.CompanyName,
		"contactPerson":          form.ContactPerson,
		"contactNumber":          form.ContactNumber,
		"email":                  form.ContactEmail,
		"address1":               form.CompanyAddressLine1,
		"address2":               form.CompanyAddressLine2,
		"idType":                 form.IDType,
		"idNumber":               form.IDNumber,
		"vehiclePassType":        form.VehiclePassType,
		"individuals":            []map[string]interface{}{},
		"companySupportingFiles": form.CompanySupportingFiles,
	}

	for _, plate := range form.CompanyPlates {
		individual := map[string]interface{}{
			"fullName":         form.ContactPerson,
			"email":            form.ContactEmail,
			"contactNumber":    form.ContactNumber,
			"address1":         form.CompanyAddressLine1,
			"address2":         form.CompanyAddressLine2,
			"nric":             plate.NricNumber,
			"vehicleNum":       strings.ToUpper(strings.ReplaceAll(plate.PlateNumber, " ", "")),
			"tinNumber":        form.TinNumber,
			"vehicleClass":     plate.VehicleType,
			"spaPath":          plate.SPAPath,
			"electricBillPath": plate.ElectricBillPath,
			"vehiclePath":      plate.VehiclePath,
		}
		natsPayload["individuals"] = append(natsPayload["individuals"].([]map[string]interface{}), individual)
	}

	data, err := json.Marshal(natsPayload)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to marshal NATS payload", err.Error())
	}

	subject := fmt.Sprintf("%s.%s", h.config.RegisterEmployerFinalizeSubject, h.config.SiteCode)
	natsResponse, err := h.natsService.GetConnection().Request(subject, data, 10*time.Second)
	if err != nil {
		return utils.ErrorResponse(c, 500, "Failed to send NATS message", err.Error())
	}

	var natsResponseData map[string]interface{}
	if err := json.Unmarshal(natsResponse.Data, &natsResponseData); err != nil {
		return utils.ErrorResponse(c, 500, "Failed to parse NATS response", err.Error())
	}

	success, ok := natsResponseData["success"].(bool)
	if !ok || !success {
		errorMsg := ""
		if errData, ok := natsResponseData["error"].(map[string]interface{}); ok {
			if msg, ok := errData["message"].(string); ok {
				errorMsg = msg
			}
		}
		return utils.ErrorResponse(c, 400, "Company registration finalization failed", errorMsg)
	}

	h.receiptService.SendReceipt(h.config.MailerFrom, form.ContactEmail, "Registration Received", form.CompanyName)

	responseData := map[string]interface{}{
		"id":           form.IDNumber,
		"natsResponse": natsResponseData,
	}
	return utils.SuccessResponse(c, "Company registration finalized successfully", responseData)
}
