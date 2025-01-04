package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
)

type DataSholat struct {
	Tanggal string `json:"tanggal"`
	Shubuh  string `json:"shubuh"`
	Dzuhur  string `json:"dzuhur"`
	Ashr    string `json:"ashr"`
	Magrib  string `json:"magrib"`
	Isya    string `json:"isya"`
}

func main() {
	s := gocron.NewScheduler(time.UTC)

	_, err := s.Every(1).Minutes().Do(func() {
		// fmt.Println("loop started")
		loopFunction()
	})
	if err != nil {
		panic(err)
	}

	s.StartAsync()

	select {}
}

func loopFunction() {

	jsonFile, err := os.Open("jadwal_sholat.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()

	bytes, err := io.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}

	dataJadwal := []DataSholat{}

	err = json.Unmarshal(bytes, &dataJadwal)
	if err != nil {
		panic(err)
	}

	mapSholat := make(map[string]DataSholat)

	for _, data := range dataJadwal {
		mapSholat[data.Tanggal] = data
	}

	timeNow := time.Now()
	dateString := timeNow.Format("2006-01-02")

	if data, ok := mapSholat[dateString]; ok {
		// loop data inside of the struct
		mapDataSholat := toMap(data)
		timeHour := timeNow.Format("15:04")
		hourNow, _ := time.Parse("15:04", timeHour)
		for key, dataSholat := range mapDataSholat {
			diff := dataSholat.Sub(hourNow)
			oneMinute, _ := time.ParseDuration("1m")

			// checking if the diff time is around 10 minutes
			if diff >= oneMinute*9 && diff < oneMinute*10 && diff > 0 {
				fmt.Println("going to notify slack, ", diff, " - ", key, " - ", dataSholat.Format("15:04"))
				err = callSlack(key, dataSholat.Format("15:04"))
				if err != nil {
					fmt.Println("ERROR CALL SLACK", err)
				}
			}
		}
	}
}

func toMap(data DataSholat) map[string]time.Time {
	result := make(map[string]time.Time)
	layout := "15:04"
	timeHour, _ := time.Parse(layout, data.Shubuh)
	result["Shubuh"] = timeHour
	timeHour, _ = time.Parse(layout, data.Dzuhur)
	result["Dzuhur"] = timeHour
	timeHour, _ = time.Parse(layout, data.Ashr)
	result["Ashr"] = timeHour
	timeHour, _ = time.Parse(layout, data.Magrib)
	result["Magrib"] = timeHour
	timeHour, _ = time.Parse(layout, data.Isya)
	result["Isya"] = timeHour

	return result
}

type DataCallSlack struct {
	Name string `json:"name"`
	Time string `json:"time"`
}

func callSlack(sholat, waktu string) error {
	url := ""

	dataCallSlack := DataCallSlack{
		Name: sholat,
		Time: waktu,
	}

	dataByte, err := json.Marshal(dataCallSlack)
	if err != nil {
		return err
	}

	client := &http.Client{}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(dataByte))
	if err != nil {
		return err
	}

	resp, err := client.Do(request)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("got error")
	}

	return nil
}
