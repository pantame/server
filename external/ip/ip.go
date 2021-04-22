package ip

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pantame/server/entities"
	"github.com/pantame/server/models/ip"
	"net/http"
	"time"
)

type responseIPApi struct {
	Status        string  `json:"status"`
	Message       string  `json:"message"`
	ContinentCode string  `json:"continentCode"`
	Country       string  `json:"country"`
	CountryCode   string  `json:"countryCode"`
	Region        string  `json:"region"`
	RegionName    string  `json:"regionName"`
	City          string  `json:"city"`
	District      string  `json:"district"`
	Zip           string  `json:"zip"`
	Lat           float64 `json:"lat"`
	Lon           float64 `json:"lon"`
	TimeZone      string  `json:"timezone"`
	Currency      string  `json:"currency"`
	ISP           string  `json:"isp"`
	ORG           string  `json:"org"`
	Mobile        bool    `json:"mobile"`
	Proxy         bool    `json:"proxy"` // Proxy, VPN or Tor exit address
	Hosting       bool    `json:"hosting"`
}

func Query(IP string) (*entities.IPData, error) {
	ipData, err := ip.GetIP("ip_date", ip.GenerateIPDate(IP))
	if err == nil {
		return ipData, errors.New("alreadyExists")
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,continentCode,country,countryCode,region,regionName,city,district,zip,lat,lon,timezone,currency,isp,org,mobile,proxy,hosting,query", IP), nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result responseIPApi

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	if result.Status != "success" {
		return nil, errors.New(result.Message)
	}

	return &entities.IPData{
		ID:            0,
		IP:            IP,
		IPDate:        IP + "-" + time.Now().Format("2006-01-02"),
		Date:          time.Now().Format("2006-01-02"),
		ContinentCode: result.ContinentCode,
		Country:       result.Country,
		CountryCode:   result.CountryCode,
		Region:        result.Region,
		RegionName:    result.RegionName,
		City:          result.City,
		District:      result.District,
		Zip:           result.Zip,
		Lat:           result.Lat,
		Lon:           result.Lon,
		TimeZone:      result.TimeZone,
		Currency:      result.Currency,
		ISP:           result.ISP,
		ORG:           result.ORG,
		Mobile:        result.Mobile,
		Proxy:         result.Proxy,
		Hosting:       result.Hosting,
		Register:      time.Now().Unix(),
	}, nil
}
