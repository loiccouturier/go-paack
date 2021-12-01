package paack

type Label struct {
	CodCurrency          string       `json:"cod_currency,omitempty"`
	CodValue             int          `json:"cod_value,omitempty"`
	Customer             Customer     `json:"customer,omitempty"`
	DeliveryAddress      Address      `json:"delivery_address,omitempty"`
	DeliveryInstructions string       `json:"delivery_instructions,omitempty"`
	DeliveryType         string       `json:"delivery_type,omitempty"`
	ExpectedDelivery     ScheduleSlot `json:"expected_delivery_ts,omitempty"`
	ExpectedPickUp       ScheduleSlot `json:"expected_pick_up_ts,omitempty"`
	ExternalId           string       `json:"external_id,omitempty"`
	InsuredCurrency      string       `json:"insured_currency,omitempty"`
	InsuredValue         float64      `json:"insured_value,omitempty"`
	Parcels              []Parcel     `json:"parcels,omitempty"`
	PickUpAddress        Address      `json:"pick_up_address,omitempty"`
	PickUpInstructions   string       `json:"pick_up_instructions,omitempty"`
	ServiceType          string       `json:"service_type,omitempty"`
	UndeliverableAddress *Address     `json:"undeliverable_address,omitempty"`
}

type LabelResponse struct {
	Content []byte
}
