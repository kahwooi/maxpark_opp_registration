package models

type CompanyRegisterForm struct {
	EmployerID             string                 `json:"employerID"`
	TinNumber              string                 `json:"tinNumber" validate:"min=5,max=16"`
	CompanyName            string                 `json:"companyName" validate:"required,min=1,max=64"`
	ContactPerson          string                 `json:"contactPerson" validate:"required,min=5,max=64"`
	ContactNumber          string                 `json:"contactNumber" validate:"required,min=12,max=12"`
	ContactEmail           string                 `json:"contactEmail" validate:"required,email,min=5,max=32"`
	CompanyAddressLine1    string                 `json:"companyAddressLine1" validate:"required,min=5,max=48"`
	CompanyAddressLine2    string                 `json:"companyAddressLine2" validate:"max=48"`
	CompanyPlates          []CompanyPlates        `json:"companyPlates"`
	CompanySupportingFiles CompanySupportingFiles `json:"companySupportingFiles"`
	IDType                 string                 `json:"idType" validate:"required"`
	IDNumber               string                 `json:"idNumber" validate:"required,min=5,max=16"`
	VehiclePassType        string                 `json:"vehiclePassType"`
}

type CompanyPlates struct {
	NricNumber       string `json:"nricNumber"`
	PlateNumber      string `json:"plateNumber" validate:"min=1,max=20"`
	VehicleType      string `json:"vehicleType"`
	SPAPath          string `json:"spaPath"`
	ElectricBillPath string `json:"electricBillPath"`
	VehiclePath      string `json:"vehiclePath"`
}

type CompanySupportingFiles struct {
	SSMPath          string `json:"ssmPath"`
	ElectricBillPath string `json:"electricBillPath"`
	VehiclePath      string `json:"vehiclePath"`
}
