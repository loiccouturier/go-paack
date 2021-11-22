package paack

type Address struct {
	City     string `json:"city"`
	Country  string `json:"country"`
	County   string `json:"county,omitempty"`
	Line1    string `json:"line1"`
	Line2    string `json:"line2"`
	PostCode string `json:"post_code"`
}
