package bigcommerce

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type InventoryResource struct {
	Inventories []Inventory `json:"data"`
	Meta        Meta        `json:"meta"`
}

type Identity struct {
	Sku       string `json:"sku"`
	VariantID int    `json:"variant_id"`
	ProductID int    `json:"product_id"`
}
type Settings struct {
	SafetyStock      int    `json:"safety_stock"`
	IsInStock        bool   `json:"is_in_stock"`
	WarningLevel     int    `json:"warning_level"`
	BinPickingNumber string `json:"bin_picking_number"`
}
type Inventory struct {
	Identity             Identity `json:"identity"`
	AvailableToSell      int      `json:"available_to_sell"`
	TotalInventoryOnhand int      `json:"total_inventory_onhand"`
	Settings             Settings `json:"settings"`
}

type Links struct {
	Previous string `json:"previous"`
	Current  string `json:"current"`
	Next     string `json:"next"`
}

type Meta struct {
	Pagination Pagination `json:"pagination"`
}

func (bc *Client) GetInventoryForLocation(ID int64, filters map[string]string) (*InventoryResource, error) {
	var params []string
	for k, v := range filters {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}

	url := fmt.Sprintf("/v3/inventory/locations/%d/items", ID) + strings.Join(params, "&")

	req := bc.getAPIRequest(http.MethodGet, url, nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		if res.StatusCode == http.StatusNoContent {
			return &InventoryResource{}, nil
		}
		return nil, err
	}

	var resource InventoryResource
	err = json.Unmarshal(body, &resource)

	if err != nil {
		return nil, err
	}
	return &resource, nil
}
