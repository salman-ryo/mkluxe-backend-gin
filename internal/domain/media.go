package domain

// Media represents image metadata embedded in a product
type Media struct {
	URL       string `bson:"url" json:"url"`
	AltText   string `bson:"alt_text" json:"alt_text"`
	IsPrimary bool   `bson:"is_primary" json:"is_primary"`
	SortOrder int    `bson:"sort_order" json:"sort_order"`
}
