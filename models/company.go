package models

type CompanyRegisterForm struct {
	EmployerID             string                 `json:"employerID" validate:"omitempty"`
	TinNumber              string                 `json:"tinNumber" validate:"omitempty,min=5,max=16"`
	CompanyName            string                 `json:"companyName" validate:"omitempty,min=5,max=64"`
	ContactPerson          string                 `json:"contactPerson" validate:"omitempty,min=5,max=64"`
	ContactNumber          string                 `json:"contactNumber" validate:"omitempty,min=5,max=12"`
	ContactEmail           string                 `json:"contactEmail" validate:"omitempty,email,min=5,max=32"`
	CompanyAddressLine1    string                 `json:"companyAddressLine1" validate:"omitempty,min=5,max=48"`
	CompanyAddressLine2    string                 `json:"companyAddressLine2" validate:"omitempty,max=48"`
	IsTenant               bool                   `json:"isTenant"`
	CompanyPlates          []CompanyPlates        `json:"companyPlates"`
	CompanySupportingFiles CompanySupportingFiles `json:"companySupportingFiles"`
	IDType                 string                 `json:"idType" validate:"omitempty"`
	IDNumber               string                 `json:"idNumber" validate:"omitempty,min=5,max=16"`
	VehiclePassType        string                 `json:"vehiclePassType" validate:"omitempty"`
	CardType               string                 `json:"cardType" validate:"omitempty"`
}

type CompanyPlates struct {
	PlateNumber string `json:"plateNumber" validate:"omitempty,min=1,max=20"`
	VehicleType string `json:"vehicleType"`
	VehiclePath string `json:"vehiclePath"`
}

type CompanySupportingFiles struct {
	SSMPath             string `json:"ssmPath"`
	ElectricBillPath    string `json:"electricBillPath"`
	TenantAgreementPath string `json:"tenantAgreementPath"`
}
