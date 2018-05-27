package beater

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/elastic/beats/libbeat/logp"
)

type QueryResponse struct {
	Size       int
	Start      int
	Limit      int
	IsLastPage bool
	Values     []Event
}

type Event struct {
	Id           int
	Timestamp    string
	Author       Author
	EventType    string
	Entities     []Entity
	IpAddress    string
	EventMessage string
	Entries      []Entry
}

type Author struct {
	Id   int
	Name string
	Type string
}

type Entity struct {
	Id      int
	Name    string
	Type    string
	Subtype string
	Primary bool
}

type Entry struct {
	PropertyName string
	OldValue     string
	NewValue     string
}

func QueryAuditLog(bt *Crowdbeat,
	start_datetime string,
	end_datetime string) []Event {
	var events []Event
	var queryResponse QueryResponse
	page := 0
	for {
		bodyBytes := MakeRequest(
			bt,
			fmt.Sprint(
				`/rest/admin/latest/auditlog/query?start=`, page),
			"POST",
			fmt.Sprint(`{"beforeOrOn": "`, end_datetime, `",`,
				`"onOrAfter": "`, start_datetime, `"}`),
		)

		unmarshalError := json.Unmarshal(bodyBytes, &queryResponse)
		if unmarshalError != nil {
			logp.Err(unmarshalError.Error())
		}

		events = append(events, queryResponse.Values...)
		page++

		if queryResponse.IsLastPage {
			logp.Info("Found unpaginated results with %v events",
				queryResponse.Size)
			break
		} else {
			logp.Info("Found paginated results with %v events",
				queryResponse.Size)
		}
	}
	return events
}

func MakeRequest(bt *Crowdbeat,
	url string,
	method string,
	payload string) []byte {
	client := &http.Client{}
	request, responseError := http.NewRequest(
		method, bt.config.CrowdUrl+url, bytes.NewBuffer([]byte(payload)))
	request.SetBasicAuth(bt.config.CrowdUsername, bt.config.CrowdPassword)
	request.Header.Add("Content-Type", "application/json")
	response, responseError := client.Do(request)
	if responseError != nil {
		logp.Err(responseError.Error())
	}
	defer response.Body.Close()
	bodyBytes, readError := ioutil.ReadAll(response.Body)
	if readError != nil {
		logp.Err(readError.Error())
	}
	return bodyBytes
}
