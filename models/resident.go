package models

type ResidentRegisterForm struct {
	ResidentName            string                  `json:"residentName" validate:"omitempty,min=1,max=64"`
	ContactNumber           string                  `json:"contactNumber" validate:"omitempty,min=1,max=12"`
	ContactEmail            string                  `json:"contactEmail" validate:"omitempty,email,min=5,max=32"`
	ResidentAddressLine1    string                  `json:"residentAddressLine1" validate:"omitempty,min=5,max=48"`
	ResidentAddressLine2    string                  `json:"residentAddressLine2" validate:"omitempty,max=48"`
	TinNumber               string                  `json:"tinNumber" validate:"omitempty,min=1,max=16"`
	IsTenant                bool                    `json:"isTenant"`
	ResidentPlates          []ResidentPlates        `json:"residentPlates"`
	ResidentSupportingFiles ResidentSupportingFiles `json:"residentSupportingFiles"`
	IDType                  string                  `json:"idType" validate:"omitempty"`
	IDNumber                string                  `json:"idNumber" validate:"omitempty,min=5,max=16"`
	VehiclePassType         string                  `json:"vehiclePassType" validate:"omitempty"`
}

type ResidentPlates struct {
	PlateNumber string `json:"plateNumber" validate:"omitempty,min=1,max=20"`
	VehicleType string `json:"vehicleType"`
	VehiclePath string `json:"vehiclePath"`
}

type ResidentSupportingFiles struct {
	SPAPath             string `json:"spaPath"`
	ElectricBillPath    string `json:"electricBillPath"`
	TenantAgreementPath string `json:"tenantAgreementPath"`
}
