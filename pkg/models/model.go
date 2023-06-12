package models

type Stock_decrease struct {
	ID           int64 `json:"id" gorm:"primary_key"`
	OrderID      int64 `json:"order_id"`
	ProductRefer int64 `json:"product_id"`
	Quantity     int64 `json:"quantity"`
}

type Product struct {
	Id             int64          `json:"id" gorm:"primary_key"`
	Name           string         `json:"name" `
	Stock          int64          `json:"stock"`
	Price          int64          `json:"price"`
	stock_decrease Stock_decrease `gorm:"foreignKey:ProductRefer"`
}

type Produk struct {
	Id    int64  `json:"id" gorm:"primary_key"`
	Name  string `json:"name" `
	Stock int64  `json:"stock"`
	Price int64  `json:"price"`
}

type Pagination struct {
	Limit int `json:"limit"`
	Page  int `json:"page"`
}
