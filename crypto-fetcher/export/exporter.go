package export

import "crypto-fetcher/store"

type Exporter interface {
	Export(store *store.PriceStore, filename string) error
}