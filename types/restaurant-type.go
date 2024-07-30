package types

type Restaurant struct {
	Title       string
	Image       string
	Time        string
	Pickup      bool
	Delivery    bool
	Rating      float64
	RatingCount int
	Menu        []string
}

type Rating struct {
	Rating float64
}
