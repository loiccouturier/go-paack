package paack

type Customer struct {
	Address        Address `json:"address,omitempty"`
	CustomerType   string  `json:"customer_type,omitempty"`
	Email          string  `json:"email,omitempty"`
	ExternalId     string  `json:"external_id,omitempty"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	HasGdprConsent bool    `json:"has_gdpr_consent,omitempty"`
	OrderRef       string  `json:"order_ref,omitempty"`
	Phone          string  `json:"phone,omitempty"`
	Language       string  `json:"language"`
}
