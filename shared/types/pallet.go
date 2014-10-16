package types

type Pallet struct {
	// Name of the pallet
	Name string `json:"name" binding:"required"`
	// Git url of the pallet
	Url string `json:"url" binding:"required"`
}
