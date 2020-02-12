package fbapi

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"log"
	"strconv"
)

type JSON map[string]interface{}

var RestClient *resty.Client

type CaseResponse struct {
	CaseData `json:"data"`
}

type CaseData struct {
	Case `json:"case"`
}

type Case struct {
	IxBug      int      `json:"ixBug"`
	Operations []string `json:"operations"`
}

type ProjectResponse struct {
	ProjectData `json:"data"`
}

type ProjectData struct {
	Project `json:"project"`
}

type Project struct {
	IxProject int `json:"ixProject"`
}

type MilestoneResponse struct {
	MilestoneData `json:"data"`
}

type MilestoneData struct {
	Milestone `json:"fixfor"`
}

type Milestone struct {
	IxFixFor int `json:"ixFixFor"`
}

type PersonResponse struct {
	PersonData `json:"data"`
}

type PersonData struct {
	Person `json:"person"`
}

type Person struct {
	IxPerson int `json:"ixPerson"`
}

type LogonResponse struct {
	LogonData `json:"data"`
}

type LogonData struct {
	Token string `json:"token"`
}

func CreateCase(body JSON) int {
	response, err := RestClient.R().
		SetBody(body).
		Post("/new")
	if err != nil {
		log.Fatal(err)
	}

	var parsed CaseResponse
	err = json.Unmarshal(response.Body(), &parsed)
	if err != nil {
		log.Fatal(err)
	}
	return parsed.CaseData.IxBug
}

func CreateTimeInterval(body JSON) {
	response, err := RestClient.R().
		SetBody(body).
		Post("/newInterval")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(response.Body()))
}

func CreateUser(body JSON) string {
	response, err := RestClient.R().
		SetBody(body).
		Post("/newPerson")
	if err != nil {
		log.Fatal(err)
	}

	var parsed PersonResponse
	err = json.Unmarshal(response.Body(), &parsed)
	if err != nil {
		log.Fatal(err)
	}
	ixPerson := strconv.Itoa(parsed.PersonData.IxPerson)
	if ixPerson == "0" {
		log.Fatal(string(response.Body()))
	}

	return ixPerson
}

func GetToken(body JSON) string {
	response, err := RestClient.R().
		SetBody(body).
		Post("/logon")
	if err != nil {
		log.Fatal(err)
	}

	var parsed LogonResponse
	err = json.Unmarshal(response.Body(), &parsed)
	if err != nil {
		log.Fatal(err)
	}

	return parsed.LogonData.Token
}

func GetCurrentUser(body JSON) int {
	response, err := RestClient.R().
		SetBody(body).
		Post("/viewPerson")
	if err != nil {
		log.Fatal(err)
	}

	var parsed PersonResponse
	err = json.Unmarshal(response.Body(), &parsed)
	if err != nil {
		log.Fatal(err)
	}

	return parsed.PersonData.Person.IxPerson
}

func CreateProject(body JSON) int {
	response, err := RestClient.R().
		SetBody(body).
		Post("/newProject")
	if err != nil {
		log.Fatal(err)
	}

	var parsed ProjectResponse
	err = json.Unmarshal(response.Body(), &parsed)
	if err != nil {
		log.Fatal(err)
	}

	return parsed.ProjectData.Project.IxProject
}

func CreateMilestone(body JSON) int {
	response, err := RestClient.R().
		SetBody(body).
		Post("/newFixFor")
	if err != nil {
		log.Fatal(err)
	}

	var parsed MilestoneResponse
	err = json.Unmarshal(response.Body(), &parsed)
	if err != nil {
		log.Fatal(err)
	}

	return parsed.MilestoneData.Milestone.IxFixFor
}

func toPrettyString(v interface{}) string {
	prettyString, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	return string(prettyString)
}

func printJson(v interface{}) {
	prettyString := toPrettyString(v)
	log.Println(prettyString)
}

func ResolveCase(body JSON) {
	response, err := RestClient.R().
		SetBody(body).
		Post("/resolve")
	if err != nil {
		log.Fatal(err)
	}

	var parsed CaseResponse
	err = json.Unmarshal(response.Body(), &parsed)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseCase(body JSON) {
	response, err := RestClient.R().
		SetBody(body).
		Post("/close")
	if err != nil {
		log.Fatal(err)
	}

	var parsed CaseResponse
	err = json.Unmarshal(response.Body(), &parsed)
	if err != nil {
		log.Fatal(err)
	}
}
