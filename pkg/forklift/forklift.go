package forklift

import "github.com/Secret-Ironman/boxr/pkg/api"

type Forklift interface {
	Deploy(api.Pallet, Shelves) (shelf string, err error)
}
