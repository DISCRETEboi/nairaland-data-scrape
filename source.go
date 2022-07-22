// Nairaland users data and web scraper

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"golang.org/x/net/html"
	"strings"
)

func main() {
	var link string
	fmt.Print("Enter a thread link >>> ")
	fmt.Scanf("%s", &link)
	page, err := http.Get(link)
	logError(err)
	pagetext, err := ioutil.ReadAll(page.Body)
	logError(err)
	text := string(pagetext)
	doc, err := html.Parse(strings.NewReader(text))
	logError(err)
	generateUsersData(doc)
	page.Body.Close()
	fmt.Println(users)
}

var users []User
var user_profile_link string
var user_name string

type User struct {
	ProfileLink string
	Name string
	//Location string
	//TimeRegistered string
	//LastSeen string
}

func generateUsersData(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" && len(node.Attr) >= 2 {
		if node.Attr[1].Key == "class" && node.Attr[1].Val == "user" {
			user_profile_link = "https://www.nairaland.com" + node.Attr[0].Val
			user_name = node.FirstChild.Data
			users = append(users, User{user_profile_link, user_name})
		}
	}
	for i := node.FirstChild; i != nil; i = i.NextSibling {
		generateUsersData(i)
	}
}

func logError(err error) {
	if err != nil {
		log.Fatal("Error encountered: ", err)
	}
}



















