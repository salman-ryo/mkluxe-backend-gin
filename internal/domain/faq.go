package domain

// FAQ represents a frequently asked question embedded in a product
type FAQ struct {
	Question  string `bson:"question" json:"question"`
	Answer    string `bson:"answer" json:"answer"`
	SortOrder int    `bson:"sort_order" json:"sort_order"`
}
