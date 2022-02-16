package paack

type Order struct {
	OrderId                   string       `json:"order_id,omitempty"`
	TrackingId                string       `json:"tracking_id,omitempty"`
	Currency                  string       `json:"cod_currency,omitempty"`
	Amount                    float32      `json:"cod_value,omitempty"`
	Customer                  Customer     `json:"customer"`
	DeliveryAddress           Address      `json:"delivery_address"`
	DeliveryInstructions      string       `json:"delivery_instructions,omitempty"`
	DeliveryType              string       `json:"delivery_type"`
	ExpectedDelivery          ScheduleSlot `json:"expected_delivery_ts"`
	ExpectedPickUp            ScheduleSlot `json:"expected_pick_up_ts"`
	ExternalId                string       `json:"external_id"`
	InsuredCurrency           string       `json:"insured_currency,omitempty"`
	InsuredAmount             float32      `json:"insured_value,omitempty"`
	Parcels                   []Parcel     `json:"parcels"`
	PickUpAddress             Address      `json:"pick_up_address"`
	PickUpInstructions        string       `json:"pick_up_instructions,omitempty"`
	ServiceType               string       `json:"service_type"`
	UndeliverableAddress      *Address     `json:"undeliverable_address,omitempty"`
	UndeliverableInstructions string       `json:"undeliverable_instructions,omitempty"`
	SaleNumber                string       `json:"sale_number,omitempty"`
	OrderDetails              []Field      `json:"order_details,omitempty"`
}

type OrderResponse struct {
	Success struct {
		TrackingId string `json:"tracking_id"`
	} `json:"success"`
}

type CancelResponse struct {
	Success bool `json:"success"`
}

type UpdateResponse struct {
	Success bool `json:"success"`
}

type GetResponse struct {
	Success Order `json:"success"`
}
