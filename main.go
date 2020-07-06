package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const (
	baseURL = "https://github.com"
)

var (
	username = "useraccount"
	password = "password"
)

type App struct {
	Client *http.Client
}

type AuthenticityToken struct {
	Token string
}

type Project struct {
	Name string
}

func (app *App) getToken() AuthenticityToken {
	loginURL := baseURL + "/login"
	client := app.Client

	response, err := client.Get(loginURL)

	if err != nil {
		log.Fatalln("Error fetching response. ", err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	token, _ := document.Find("input[name='authenticity_token']").Attr("value")

	authenticityToken := AuthenticityToken{
		Token: token,
	}

	return authenticityToken
}

func (app *App) login() []Project {
	client := app.Client

	authenticityToken := app.getToken()

	loginURL := baseURL + "/session"
	data := url.Values{
		"authenticity_token": {authenticityToken.Token},
		"login":              {username},
		"password":           {password},
	}

	response, err := client.PostForm(loginURL, data)

	if err != nil {
		log.Fatalln(err)
	}

	defer response.Body.Close()

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal("Error loading HTTP response body. ", err)
	}

	var projects []Project
	document.Find("[data-filterable-for=\"dashboard-repos-filter-left\"]").Find(".width-full").Each(func(i int, s *goquery.Selection) {
		name, _ := s.Find("a").Attr("href")
		project := Project{
			Name: name,
		}

		projects = append(projects, project)

	})

	return projects
}

func main() {
	jar, _ := cookiejar.New(nil)

	app := App{
		Client: &http.Client{Jar: jar},
	}

	projects := app.login()

	for index, project := range projects {
		fmt.Printf("%d: %s\n", index+1, project.Name)
	}
}
