package handlers

import (
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rmq/types"
)

type CatalogHandler struct {
	// searchRpc search.SearchClient
	// marketingRpc marketing.MarketingClient
}

func NewCatalogHandler() *CatalogHandler {
	return &CatalogHandler{}
}

func (h *CatalogHandler) HandleProductUpdate(event *types.ProductEvent) error {
	_, ok := event.Data.(*types.ProductData)
	if !ok {
		return zeroerr.ErrInvalidEventData
	}

	// TODO: Update search index
	// TODO: Notify marketing service about new/updated products
	return nil
}
