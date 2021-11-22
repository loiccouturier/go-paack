package paack

type Label struct {
	CodCurrency string `json:"cod_currency"`
	CodValue    int    `json:"cod_value"`
	Customer    struct {
		Address        Address `json:"address"`
		CustomerType   string  `json:"customer_type"`
		Email          string  `json:"email"`
		ExternalId     string  `json:"external_id"`
		FirstName      string  `json:"first_name"`
		LastName       string  `json:"last_name"`
		HasGdprConsent bool    `json:"has_gdpr_consent"`
		OrderRef       string  `json:"order_ref"`
		Phone          string  `json:"phone"`
		Language       string  `json:"language"`
	} `json:"customer"`
	DeliveryAddress      Address `json:"delivery_address"`
	DeliveryInstructions string  `json:"delivery_instructions"`
	DeliveryType         string  `json:"delivery_type"`
	ExpectedDeliveryTs   struct {
		End struct {
			Date string `json:"date"`
			Time string `json:"time"`
		} `json:"end"`
		Start struct {
			Date string `json:"date"`
			Time string `json:"time"`
		} `json:"start"`
	} `json:"expected_delivery_ts"`
	ExpectedPickUpTs struct {
		End struct {
			Date string `json:"date"`
			Time string `json:"time"`
		} `json:"end"`
		Start struct {
			Date string `json:"date"`
			Time string `json:"time"`
		} `json:"start"`
	} `json:"expected_pick_up_ts"`
	ExternalId           string   `json:"external_id"`
	InsuredCurrency      string   `json:"insured_currency"`
	InsuredValue         float64  `json:"insured_value"`
	Parcels              []Parcel `json:"parcels"`
	PickUpAddress        Address  `json:"pick_up_address"`
	PickUpInstructions   string   `json:"pick_up_instructions"`
	ServiceType          string   `json:"service_type"`
	UndeliverableAddress Address  `json:"undeliverable_address"`
}
