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

	type aka Event
	return json.Unmarshal(data, (*aka)(e))
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
	var jsonBlob = []byte(`[
	{"resource_type":"payment","action":"confirmed","data":{"amount":100}},
	{"resource_type":"customer","action":"created","data":{"name":"john"}}
]`)

	var events []Event
	err := json.Unmarshal(jsonBlob, &events)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v", events)
	// output: [{ResourceType:payment Action:confirmed Data:{Amount:100}} {ResourceType:customer Action:created Data:{Name:john}}]
}
