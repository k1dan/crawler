package storage

type (
	File struct {
		FileName string
		Item     Item
	}
	Item struct {
		Title      string `json:"title"`
		Condition  string `json:"condition"`
		Price      string `json:"price"`
		ProductURL string `json:"product_url"`
	}
)
