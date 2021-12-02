package paack

import (
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestClient_CreateOrder(t *testing.T) {

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	orderExternalId := "BC156454-" + strconv.Itoa(r1.Intn(100-5)+5) + "-" + strconv.Itoa(r1.Intn(200-100)+100)

	pickupAddress := Address{
		City:     "BEZONS",
		Country:  "FR",
		Line1:    "1 RUE JEAN CARRASSO",
		PostCode: "95870",
	}

	customerAddress := Address{
		City:     "STRASBOURG",
		Country:  "FR",
		Line1:    "3 RUE DU PARC",
		PostCode: "67000",
	}

	parcel1 := Parcel{
		Barcode:          orderExternalId + "-1",
		Height:           84,
		Length:           33,
		Type:             "standard",
		VolumetricWeight: 3475,
		Weight:           83,
		Width:            22,
	}

	parcel2 := Parcel{
		Barcode:          orderExternalId + "-2",
		Height:           24,
		Length:           43,
		Type:             "standard",
		VolumetricWeight: 3885,
		Weight:           93,
		Width:            82,
	}

	customer := Customer{
		FirstName: "Loic",
		LastName:  "Couturier",
		Address:   customerAddress,
		Language:  "fr",
	}

	today := time.Now()

	order := Order{
		Currency:             "EUR",
		Amount:               40.50,
		Customer:             customer,
		DeliveryAddress:      customerAddress,
		DeliveryInstructions: "",
		DeliveryType:         "direct",
		ExpectedPickUp:       NewScheduleSlot(time.Date(today.Year(), today.Month(), today.Day() + 1, 17, 0, 0, today.Nanosecond(), today.Location()), time.Date(today.Year(), today.Month(), today.Day() + 1, 19, 0, 0, today.Nanosecond(), today.Location())),
		ExpectedDelivery:     NewScheduleSlot(time.Date(today.Year(), today.Month(), today.Day() + 1, 18, 0, 0, today.Nanosecond(), today.Location()), time.Date(today.Year(), today.Month(), today.Day() + 1, 20, 0, 0, today.Nanosecond(), today.Location())),
		ExternalId:           orderExternalId,
		InsuredCurrency:      "EUR",
		InsuredAmount:        40.50,
		Parcels: []Parcel{
			parcel1,
			parcel2,
		},
		PickUpAddress:             pickupAddress,
		PickUpInstructions:        "",
		ServiceType:               "ST2",
		UndeliverableAddress:      &pickupAddress,
		UndeliverableInstructions: "",
	}

	c := NewClient("https://api.staging.paack.app", "https://paack-hq-staging.eu.auth0.com", "https://paack-hq-staging.eu.auth0.com", "https://api.oms.staging.paack.app", "https://ggl-stg-gcp-gw", "https://api.oms.staging.paack.app", "", "")
	r, err := c.CreateOrder(order)

	if err != nil {
		t.Error(err)
	}

	if r != nil {
		if r.Success.TrackingId == "" {
			t.Errorf("No tracking id")
		}
	}

}
