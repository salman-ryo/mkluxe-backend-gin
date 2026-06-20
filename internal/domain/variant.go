package domain

// Variant represents a specific type of a production variant, such as size or color.

type Variant struct {
	SKU            string            `bson:"sku" json:"sku"`
	Price          float64           `bson:"price" json:"price"`
	CompareAtPrice float64           `bson:"compare_at_price,omitempty" json:"compare_at_price,omitempty"`
	Stock          int               `bson:"stock" json:"stock"`
	IsDefault      bool              `bson:"is_default" json:"is_default"`
	Attributes     map[string]string `bson:"attributes,omitempty" json:"attributes,omitempty"`
}
