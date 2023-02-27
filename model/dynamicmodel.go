package main

import (
	"encoding/json"
	"fmt"
)

type Event struct {
	ResourceType string `json:"resource_type"`
	Action       string `json:"action"`
	Data         any    `json:"data"`
}

func (e *Event) UnmarshalJSON(data []byte) error {
	var inner struct {
		ResourceType string `json:"resource_type"`
	}
	if err := json.Unmarshal(data, &inner); err != nil {
		return err
	}

	switch inner.ResourceType {
	case "payment":
		e.Data = new(Payment)
	case "customer":
		e.Data = new(Customer)
	}

	type eventAlias Event
	return json.Unmarshal(data, (*eventAlias)(e))
}

type Customer struct {
	Name string `json:"name"`
}

func (c *Customer) String() string {
	return fmt.Sprintf("{Name:%s}", c.Name)
}

type Payment struct {
	Amount int `json:"amount"`
}

func (p *Payment) String() string {
	return fmt.Sprintf("{Amount:%d}", p.Amount)
}

func main() {
	bPaymentEvent := []byte(`{"resource_type":"payment","action":"confirmed","data":{"amount":100}}`)
	var payment Event
	err := json.Unmarshal(bPaymentEvent, &payment)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("payment event: %+v\n", payment)

	bCustomerEvent := []byte(`{"resource_type":"customer","action":"created","data":{"name":"john"}}`)
	var customer Event
	err = json.Unmarshal(bCustomerEvent, &customer)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("customer event: %+v\n", customer)
}
