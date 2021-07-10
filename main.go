package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"github.com/alextotalk/ewesomeProje/models"
	"io"
	"log"
	"strconv"
	"sync"
	"sync/atomic"
)

func main() {
	var shippingSourceAddress = `Kiev`
	var shippingDestinationAddress = `Odessa`
	var cartonBoxDimensions = `[3,4,5]`

	totalJSON := 3000
	amountJSON := 0
	amountXML := 4000

	lowestPrice := getLowestPrise(shippingSourceAddress, shippingDestinationAddress, cartonBoxDimensions, totalJSON, amountJSON, amountXML)
	fmt.Printf("Lowest Prise %v", lowestPrice)
}

func getLowestPrise(a, b, c string, d, e, f int) int {
	type Prices []int
	var m atomic.Value
	m.Store(make(Prices, 3))
	var mu sync.Mutex

	wg := &sync.WaitGroup{}

	wg.Add(1)
	//Api1
	go func(a, b, c string, apiPrice int, wg *sync.WaitGroup) {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		m1 := m.Load().(Prices)

		bodyBytes := []byte(`{"contact address": "` + a + `","warehouse address": "` + b + `","package dimensions": ` + c + `}`)
		//resp, err := http.Post("http://localhost:8000/api1",
		//	"application/json; charset=utf-8", bytes.NewBuffer(bodyBytes))
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//
		//body, err := ioutil.ReadAll(resp.Body)
		ap := strconv.Itoa(apiPrice)
		bodyBytes = []byte(`{"total":` + ap + `}`)
		body := bodyBytes
		//if err != nil {
		//	log.Fatalln(err)
		//}
		d := &models.Total{}
		err := d.UnmarshalJSON(body)
		if err != nil {
			log.Fatalln(err)
		}
		//json.Unmarshal(body, d)
		m1[0] = d.Total
		return
	}(a, b, c, d, wg)

	wg.Add(1)
	//Api2

	go func(a, b, c string, apiPrice int, wg *sync.WaitGroup) {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		m1 := m.Load().(Prices)
		bodyBytes := []byte(`{"consignee": "` + a + `","consignor": "` + b + `","cartons": ` + c + `}`)
		//resp, err := http.Post("http://localhost:8000/api2",
		//	"application/json; charset=utf-8", bytes.NewBuffer(bodyBytes))
		//if err != nil {
		//	log.Fatalln(err)
		//}
		////We Read the response body on the line below.
		//body, err := ioutil.ReadAll(resp.Body)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		ap := strconv.Itoa(apiPrice)
		bodyBytes = []byte(`{"amount":` + ap + `}`)
		body := bodyBytes
		d := &models.Amount{}
		err := d.UnmarshalJSON(body)
		if err != nil {
			log.Fatalln(err)
		}
		//json.Unmarshal(body, d)
		m1[1] = d.Amount
		return
	}(a, b, c, e, wg)

	wg.Add(1)
	//Api3

	go func(a, b, c string, apiPrice int, wg *sync.WaitGroup) {
		defer wg.Done()
		mu.Lock()
		defer mu.Unlock()
		m1 := m.Load().(Prices)
		bodyBytes := []byte(`<?xml version="1.0" encoding="utf-8"?>
	<createInfo>
  		<consignee value="` + a + `"/>
  		<consignor value="` + b + `"/>
  		<cartons value="` + c + `">
	</createInfo>`)
		//resp, err := http.Post("http://localhost:8000/api3",
		//	"application/xml; charset=utf-8", bytes.NewBuffer(bodyBytes))
		//if err != nil {
		//	log.Fatalln(err)
		//}
		////We Read the response body on the line below.
		//body, err := ioutil.ReadAll(resp.Body)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		ap := strconv.Itoa(apiPrice)
		bodyBytes = []byte(`<?xml version="1.0" encoding="utf-8"?><Amount>` + ap + `</Amount>`)
		body := bodyBytes
		input := bytes.NewReader(body)
		decoder := xml.NewDecoder(input)
		amount := 0
		for {
			tok, tokenErr := decoder.Token()
			if tokenErr != nil && tokenErr != io.EOF {
				fmt.Println("error happened", tokenErr)
				break
			} else if tokenErr == io.EOF {
				break
			}
			if tok == nil {
				fmt.Println("t is nil break")
			}

			switch tok := tok.(type) {
			case xml.StartElement:
				if tok.Name.Local == "Amount" {

					if err := decoder.DecodeElement(&amount, &tok); err != nil {
						fmt.Println("error happened", err)
					}
				}
			}
		}
		m1[2] = amount
		return
	}(a, b, c, f, wg)

	wg.Wait()
	sp := m.Load().(Prices)
	for i := 1; i < len(sp); i++ {
		if sp[i] == 0 {
			continue
		}
		if sp[i] < sp[0] {
			sp[0] = sp[i]
		}
	}
	if sp[0] == 0 {
		fmt.Println("Sorry, we cannot offer you a price!")
	}
	return sp[0]
}
