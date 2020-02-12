package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/jaswdr/faker"
	"github.com/jessevdk/go-flags"
	_ "github.com/motemen/go-loghttp/global"
	"github.com/trilogy-group/randombugz/fbapi"
	"log"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var Opts struct {
	Instance    string `long:"fb-instance" env:"FB_INSTANCE" description:"FB instance to populate"`
	Token       string `long:"fb-token" env:"FB_TOKEN" description:"FB api token"`
	FileInput   string `short:"f" long:"file-input" description:"Input csv file original;elapsed;milestone;user'"`
	EmailPrefix string `long:"email-prefix" description:"Used for creating new users. New user will have email: <email-prefix>+<random-number>@gmail.com"`
}

type CaseInput struct {
	Original  int
	Elapsed   int
	Milestone int
	User      int
}

type Person struct {
	ixPerson string
	token    string
	email    string
}

const layoutFB string = "2006-01-02 15:04:05Z"

func main() {
	if _, err := flags.Parse(&Opts); err != nil {
		os.Exit(1)
	}
	initRestClient()
	rand.Seed(time.Now().UnixNano())
	fmt.Printf("Read %s", Opts.FileInput)
	csvFile, _ := os.Open(Opts.FileInput)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.Comma = ';'
	cases := readCaseInput(reader)

	numberOfCases := len(cases)
	numberOfMilestones := 1
	numberOfUsers := 1
	for i := range cases {
		numberOfMilestones = max(cases[i].Milestone, numberOfMilestones)
		numberOfUsers = max(cases[i].User, numberOfUsers)
	}

	ixUsers := make(map[int]Person)
	for i := 0; i < numberOfUsers; i++ {
		userUniqueSuffix := strconv.Itoa(rand.Int() % 1_000_000)

		email := Opts.EmailPrefix + "+" + userUniqueSuffix + "@gmail.com"
		fullName := "estimator " + userUniqueSuffix
		ixPerson := fbapi.CreateUser(fbapi.JSON{
			"token":     Opts.Token,
			"sEmail":    email,
			"sFullName": fullName,
			"nType":     "1", //Admin
			"sPassword": "123123",
		})
		ixToken := fbapi.GetToken(fbapi.JSON{
			"email":    email,
			"password": "123123",
		})
		ixUsers[i+1] = Person{
			ixPerson: ixPerson,
			token:    ixToken,
			email:    email,
		}
		fmt.Printf("Created user %s, email: %s, token: %s\n", ixPerson, email, ixToken)
	}

	fmt.Printf("Create %d cases\n", numberOfCases)
	ixProject, projectTitle := createNewProject(ixUsers[1])

	ixFixFors := make(map[int]int)
	for i := 0; i < numberOfMilestones; i++ {
		ixFixFors[i+1] = createNewMilestone(ixProject, ixUsers[1])
	}

	for i := 0; i < numberOfCases; i++ {
		generateCase(i, cases[i], ixProject, ixFixFors, ixUsers)
	}
	fmt.Printf("All done for project %s\n", projectTitle)
}

func readCaseInput(reader *csv.Reader) []CaseInput {
	var cases []CaseInput
	lines, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	for i, line := range lines {
		if i == 0 {
			continue
		}
		original, err := strconv.Atoi(line[0])
		if err != nil {
			log.Fatal(err)
		}
		elapsed, err := strconv.Atoi(line[1])
		if err != nil {
			log.Fatal(err)
		}
		milestone, err := strconv.Atoi(line[2])
		if err != nil {
			log.Fatal(err)
		}
		user, err := strconv.Atoi(line[3])
		if err != nil {
			log.Fatal(err)
		}
		cases = append(cases, CaseInput{
			Original:  original,
			Elapsed:   elapsed,
			Milestone: milestone,
			User:      user,
		})
	}
	return cases
}

func createNewProject(person Person) (int, string) {
	fakeNames := faker.New()
	uniqueSuffix := strconv.Itoa(rand.Int() % 100)
	randomTitle := fakeNames.Lorem().Word() + uniqueSuffix
	fmt.Printf("Create '%s' project\n", randomTitle)
	ixProject := fbapi.CreateProject(fbapi.JSON{
		"token":                  person.token,
		"sProject":               randomTitle,
		"ixPersonPrimaryContact": person.ixPerson,
	})
	return ixProject, randomTitle
}

func createNewMilestone(ixProject int, person Person) int {
	fakeNames := faker.New()
	uniqueSuffix := strconv.Itoa(rand.Int() % 100)
	randomTitle := fakeNames.Payment().CreditCardType() + uniqueSuffix
	ixProjectStr := strconv.Itoa(ixProject)
	fmt.Printf("Create '%s' milestone\n", randomTitle)
	return fbapi.CreateMilestone(fbapi.JSON{
		"token":       person.token,
		"sFixFor":     randomTitle,
		"ixProject":   ixProjectStr,
		"fAssignable": "1",
	})
}

func generateCase(i int, caseInput CaseInput, ixProject int, ixFixFors map[int]int, ixUsers map[int]Person) {
	ixFixFor := ixFixFors[caseInput.Milestone]
	person := ixUsers[caseInput.User]
	title := "Case " + strconv.Itoa(i)
	originalEstimation := caseInput.Original
	elapsed := caseInput.Elapsed
	fmt.Printf("Create case %s\n", title)
	ixBug := fbapi.CreateCase(fbapi.JSON{
		"token":              person.token,
		"sTitle":             title,
		"hrsCurrEst":         strconv.Itoa(originalEstimation),
		"ixProject":          strconv.Itoa(ixProject),
		"ixFixFor":           strconv.Itoa(ixFixFor),
		"ixPersonAssignedTo": person.ixPerson,
	})

	if elapsed > 0 {
		fmt.Printf("Resolve case %d\n", ixBug)
		createTimeInterval(i, person.ixPerson, ixBug, elapsed, person.token)
		fbapi.ResolveCase(fbapi.JSON{
			"token": person.token,
			"ixBug": ixBug,
		})
		fmt.Printf("Close case %d\n", ixBug)
		fbapi.CloseCase(fbapi.JSON{
			"token": person.token,
			"ixBug": ixBug,
		})
	}
}

func createTimeInterval(i int, ixPerson string, ixBug int, elapsed int, ixToken string) {
	t := time.Date(2019, 4, 1, 0, 0, 0, 0, time.UTC)
	dtStart := t.AddDate(0, 0, i)
	dtEnd := dtStart.Add(time.Hour * time.Duration(elapsed))

	fbapi.CreateTimeInterval(fbapi.JSON{
		"token":    ixToken,
		"ixPerson": ixPerson,
		"ixBug":    strconv.Itoa(ixBug),
		"dtStart":  dtStart.Format(layoutFB),
		"dtEnd":    dtEnd.Format(layoutFB),
	})
}

func initRestClient() {
	fbapi.RestClient = resty.
		New().
		SetHostURL(Opts.Instance+"/api").
		SetHeader("Content-Type", "application/json")
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
