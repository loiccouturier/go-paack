package paack

type Label struct {
	CodCurrency string `json:"cod_currency"`
	CodValue    int    `json:"cod_value"`
	Customer    Customer `json:"customer"`
	DeliveryAddress      Address `json:"delivery_address"`
	DeliveryInstructions string  `json:"delivery_instructions"`
	DeliveryType         string  `json:"delivery_type"`
	ExpectedDelivery     ScheduleSlot `json:"expected_delivery_ts"`
	ExpectedPickUp       ScheduleSlot `json:"expected_pick_up_ts"`
	ExternalId           string   `json:"external_id"`
	InsuredCurrency      string   `json:"insured_currency"`
	InsuredValue         float64  `json:"insured_value"`
	Parcels              []Parcel `json:"parcels"`
	PickUpAddress        Address  `json:"pick_up_address"`
	PickUpInstructions   string   `json:"pick_up_instructions"`
	ServiceType          string   `json:"service_type"`
	UndeliverableAddress Address  `json:"undeliverable_address"`
}
