package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// https://apis.scrimba.com/bored/documentation

// http://www.boredapi.com/api/activity/

// {
// 	"activity": "Learn Express.js",
// 	"accessibility": 0.25,
// 	"type": "education",
// 	"participants": 1,
// 	"price": 0.1,
// 	"link": "https://expressjs.com/",
// 	"key": "3943506"
// }

type boringResponse struct {
	Activity      string  `json:"activity"`
	Accessibility float64 `json:"accessibility"`
	Type          string  `json:"type"`
	Participants  int     `json:"participants"`
	Price         float64 `json:"price"`
	Link          string  `json:"link"`
	Key           string  `json:"key"`
}

func main() {

	ctx := context.Background()
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://apis.scrimba.com/bored/api/activity", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return
	}

	var boringResp boringResponse
	if err := json.NewDecoder(resp.Body).Decode(&boringResp); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Printf("%+v\n", boringResp)
}
