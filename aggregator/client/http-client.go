package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mohamedramadan14/roads-fees-system/types"
)

type Client struct {
	Endpoint string
}

func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
	}
}

func (c *Client) AggregateInvoice(distance types.Distance) error {
	httpc := http.DefaultClient
	b, err := json.Marshal(distance)
	if err != nil {

		return err
	}
	req, err := http.NewRequest(http.MethodPost, c.Endpoint, bytes.NewBuffer(b))
	if err != nil {
		return err
	}
	resp, err := httpc.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to generate invoice, status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	return nil
}
