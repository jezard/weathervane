/****
*
*  Package to create a simple UK weather observation page/site - the aim being to disregard inaccurate forecasts and allow our brains to extrapolate where and when the rain will fall
* 
*  Thanks to Met Office http://www.metoffice.gov.uk/ for their DataPoint feeds
*  Thanks to Matt Holt https://mholt.github.io/json-to-go/ for his JSON-to-go JSON to Go struct generator...
*
****/

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

//capabilities
type Wxcaps struct {
	Resource struct {
		Res string `json:"res"`
		Type string `json:"type"`
		TimeSteps struct {
			TS []time.Time `json:"TS"`
		} `json:"TimeSteps"`
	} `json:"Resource"`
}

//observations
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
	var obs []Wxobs //stores all of our observations
	caps := getCapabilities()

	for _, capability := range caps{
		f := "2006-01-02T15Z"
		s := capability.UTC().Format(f) 
		o := getObservations(s)

		obs = append(obs, o)
	}

	for _, ob := range obs{//test to loop through all observations
		for i, loc := range ob.SiteRep.DV.Location{
			fmt.Printf("%d %s %s\n", i, loc.Name,  loc.Elevation)//get any location attribute 
		}
	}
}

//returns array of available capabilities
func getCapabilities()[]time.Time{
	conf := conf.Get()

	var c []time.Time

	//get the capabilities for the UK observations data feed
	url := "http://datapoint.metoffice.gov.uk/public/data/val/wxobs/all/json/capabilities?res=hourly&key=" + conf.MOApiKey;

	resp, err := http.Get(url)
	if err != nil{
		fmt.Printf("%v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
	    panic(err.Error())
	}

	var data Wxcaps

	json.Unmarshal(body, &data)

	for _, capability := range data.Resource.TimeSteps.TS{
		c = append(c, capability)	
	}
	return c
}

//returns observation data for a single capability
func getObservations(snapshot string)Wxobs{
	conf := conf.Get()

	//get the UK observations data feed
	url := "http://datapoint.metoffice.gov.uk/public/data/val/wxobs/all/json/all?res=hourly&time=" + snapshot + "&key=" + conf.MOApiKey;//date needs manually updating for this stub, re-read the sections about all this (Met Office docs)

	fmt.Printf("%s\n", url)

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

	return data
}



