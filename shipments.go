package bigcommerce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Shipment struct {
	ID                   int64            `json:"id,omitempty"`
	OrderId              int64            `json:"order_id,omitempty"`
	CustomerId           int64            `json:"customer_id,omitempty"`
	OrderAddressId       int64            `json:"order_address_id,omitempty"`
	DateCreated          string           `json:"date_created,omitempty"`
	TrackingNumber       string           `json:"tracking_number"`
	MerchantShippingCost string           `json:"merchant_shipping_cost"`
	ShippingMethod       string           `json:"shipping_method"`
	Comments             string           `json:"comments"`
	ShippingProvider     string           `json:"shipping_provider"`
	TrackingCarrier      string           `json:"tracking_carrier"`
	BillingAddress       *ShipmentAddress `json:"billing_address,omitempty"`
	ShippingAddress      *ShipmentAddress `json:"shipping_address,omitempty"`
	Items                []ShipmentItem   `json:"items"`
}

type ShipmentAddress struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	Company     string `json:"company"`
	Street1     string `json:"street_1"`
	Street2     string `json:"street_2"`
	City        string `json:"city"`
	State       string `json:"state"`
	Zip         string `json:"zip"`
	Country     string `json:"country"`
	CountryIso2 string `json:"country_iso2"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
}

type ShipmentItem struct {
	OrderProductId int64 `json:"order_product_id"`
	ProductId      int64 `json:"product_id,omitempty"`
	Quantity       int64 `json:"quantity"`
}

// GetOrderShipments retrieves all shipments that belong to a specific order
func (bc *Client) GetOrderShipments(orderId int64, filters map[string]string) ([]Shipment, error) {
	var params []string
	for k, v := range filters {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}
	url := fmt.Sprintf("/v2/orders/%d/shipments?%s", orderId, strings.Join(params, "&"))

	req := bc.getAPIRequest(http.MethodGet, url, nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		if res.StatusCode == http.StatusNoContent {
			return []Shipment{}, nil
		}
		return nil, err
	}

	var shipments []Shipment
	err = json.Unmarshal(body, &shipments)
	if err != nil {
		return nil, err
	}
	return shipments, nil
}

// CreateOrderShipment creates a new shipment belonging to an order.
// If the shipment does not contain all products, bigcommerce will by default tag the order as partially done
func (bc *Client) CreateOrderShipment(orderId int64, shipment Shipment) (*Shipment, error) {
	url := fmt.Sprintf("/v2/orders/%d/shipments", orderId)

	// Make sure shipment doesn't have any fields that are not allowed
	shipment = Shipment{
		OrderAddressId:   shipment.OrderAddressId,
		TrackingNumber:   shipment.TrackingNumber,
		ShippingMethod:   shipment.ShippingMethod,
		Comments:         shipment.Comments,
		ShippingProvider: shipment.ShippingProvider,
		TrackingCarrier:  shipment.TrackingCarrier,
		Items:            shipment.Items,
	}

	reqJSON, err := json.Marshal(shipment)
	if err != nil {
		return nil, err
	}

	req := bc.getAPIRequest(http.MethodPost, url, bytes.NewReader(reqJSON))
	res, err := bc.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := processBody(res)

	if err != nil {
		if res.StatusCode == http.StatusNoContent {
			return &Shipment{}, nil
		}
		return nil, err
	}

	var s *Shipment
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// DeleteOrderShipments deletes ALL shipments belonging to an order
func (bc *Client) DeleteOrderShipments(orderId int64) (bool, error) {
	url := fmt.Sprintf("/v2/orders/%d/shipments", orderId)

	req := bc.getAPIRequest(http.MethodDelete, url, nil)
	_, err := bc.HTTPClient.Do(req)

	if err != nil {
		return false, err
	}

	return true, nil
}

// DeleteOrderShipment deletes a single shipment under an order
func (bc *Client) DeleteOrderShipment(orderId int64, shipmentId int64) (bool, error) {
	url := fmt.Sprintf("/v2/orders/%d/shipments/%d", orderId, shipmentId)

	req := bc.getAPIRequest(http.MethodDelete, url, nil)
	_, err := bc.HTTPClient.Do(req)

	if err != nil {
		return false, err
	}

	return true, nil
}

// GetOrderShipment retrieves a single shipment
func (bc *Client) GetOrderShipment(orderId int64, shipmentId int64) (*Shipment, error) {
	url := fmt.Sprintf("/v2/orders/%d/shipments/%d", orderId, shipmentId)

	req := bc.getAPIRequest(http.MethodGet, url, nil)
	res, err := bc.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := processBody(res)
	if err != nil {
		if res.StatusCode == http.StatusNoContent {
			return &Shipment{}, nil
		}
		return nil, err
	}

	var shipment *Shipment
	err = json.Unmarshal(body, &shipment)
	if err != nil {
		return nil, err
	}
	return shipment, nil
}

// UpdateOrderShipment updates an existing shipment belonging to an order.
// If the shipment does not contain all products, bigcommerce will by default tag the order as partially done
func (bc *Client) UpdateOrderShipment(orderId int64, shipment Shipment) (*Shipment, error) {
	url := fmt.Sprintf("/v2/orders/%d/shipments/%d", orderId, shipment.ID)

	// Make sure shipment doesn't have any fields that are not allowed
	shipment = Shipment{
		OrderAddressId:   shipment.OrderAddressId,
		TrackingNumber:   shipment.TrackingNumber,
		ShippingMethod:   shipment.ShippingMethod,
		Comments:         shipment.Comments,
		ShippingProvider: shipment.ShippingProvider,
		TrackingCarrier:  shipment.TrackingCarrier,
		Items:            shipment.Items,
	}

	reqJSON, err := json.Marshal(shipment)
	if err != nil {
		return nil, err
	}

	req := bc.getAPIRequest(http.MethodPut, url, bytes.NewReader(reqJSON))
	res, err := bc.HTTPClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := processBody(res)

	if err != nil {
		if res.StatusCode == http.StatusNoContent {
			return &Shipment{}, nil
		}
		return nil, err
	}

	var s *Shipment
	err = json.Unmarshal(body, &s)
	if err != nil {
		return nil, err
	}
	return s, nil
}
