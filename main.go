package main

import (
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/jezard/weathervane/conf"
	"io/ioutil"
	"encoding/json"
	"time"
)

type Wxobs struct {
	SiteRep struct {
		Wx struct {
			Param []struct {
				Name string `json:"name"`
				Units string `json:"units"`
				NAMING_FAILED string `json:"$"`
			} `json:"Param"`
		} `json:"Wx"`
		DV struct {
			DataDate time.Time `json:"dataDate"`
			Type string `json:"type"`
			Location []struct {
				I string `json:"i"`
				Lat string `json:"lat"`
				Lon string `json:"lon"`
				Name string `json:"name"`
				Country string `json:"country"`
				Continent string `json:"continent"`
				Elevation string `json:"elevation"`
				Period struct {
					Type string `json:"type"`
					Value string `json:"value"`
					Rep struct {
						D string `json:"D"`
						H string `json:"H"`
						P string `json:"P"`
						S string `json:"S"`
						T string `json:"T"`
						V string `json:"V"`
						W string `json:"W"`
						Pt string `json:"Pt"`
						Dp string `json:"Dp"`
						NAMING_FAILED string `json:"$"`
					} `json:"Rep"`
				} `json:"Period"`
			} `json:"Location"`
		} `json:"DV"`
	} `json:"SiteRep"`
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil) //local: http://192.168.2.100:8080/
}
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	conf := conf.Get()

	
	url := "http://datapoint.metoffice.gov.uk/public/data/val/wxobs/all/json/all?res=hourly&time=2016-04-14T15Z&key=" + conf.MOApiKey;//date needs manually updating for this stub, re-read the sections about all this (Met Office docs)

	resp, err := http.Get(url)
	if err != nil{
		fmt.Printf("%v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	    panic(err.Error())
	}

	var data Wxobs

	json.Unmarshal(body, &data)

	for i, obs := range data.SiteRep.DV.Location{
		fmt.Printf("%d %s %s\n", i, obs.Name,  obs.Elevation)//get any location attribute 
	}
}
