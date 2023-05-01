package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type ViaCEP struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
}

type Cdn struct {
	Code       string `json:"code"`
	State      string `json:"state"`
	City       string `json:"city"`
	District   string `json:"district"`
	Address    string `json:"address"`
	Status     int    `json:"status"`
	Ok         bool   `json:"ok"`
	StatusText string `json:"statusText"`
}

func main() {

	c1 := make(chan ViaCEP)
	c2 := make(chan Cdn)
	cep := "69060-432" // 69060-432

	go func() {
		resp, error := http.Get(fmt.Sprintf("https://cdn.apicep.com/file/apicep/%s.json", cep))
		if error != nil {
			panic(error)
		}
		defer resp.Body.Close()
		body, error := ioutil.ReadAll(resp.Body)
		if error != nil {
			panic(error)
		}
		var cdn Cdn
		error = json.Unmarshal(body, &cdn)
		if error != nil {
			panic(error)
		}
		fmt.Println(cdn)
		c2 <- cdn
	}()

	go func() {
		resp, error := http.Get("http://viacep.com.br/ws/" + cep + "/json/")
		if error != nil {
			panic(error)
		}
		defer resp.Body.Close()
		body, error := ioutil.ReadAll(resp.Body)
		if error != nil {
			panic(error)
		}
		var viaCEP ViaCEP
		error = json.Unmarshal(body, &viaCEP)
		if error != nil {
			panic(error)
		}
		fmt.Println(viaCEP)
		c1 <- viaCEP
	}()

	select {
	case msg1 := <-c1:
		fmt.Printf("received ViaCEP %v", msg1)

	case msg2 := <-c2:
		fmt.Printf("received Cdn %v", msg2)

	case <-time.After(time.Second * 1):
		println("timeout")
	}

}
