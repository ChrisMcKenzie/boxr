package types

type Pallet struct {
	// Name of the pallet
	Name string `db:"name" json:"name" binding:"required"`
	// Git url of the pallet
	Url string `db:"url" json:"url" binding:"required"`
}
