package paack

type Parcel struct {
	Barcode          string `json:"barcode"`
	Height           int    `json:"height"`
	Length           int    `json:"length"`
	Type             string `json:"type"`
	VolumetricWeight int    `json:"volumetric_weight"`
	Weight           int    `json:"weight"`
	Width            int    `json:"width"`
}

type Parcels struct {
	Parcels []Parcel `json:"parcels"`
}

type ReplaceParcelsResponse struct {
	Success bool `json:"success"`
}
