// Nairaland users data and web scraper

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"golang.org/x/net/html"
	"strings"
	"encoding/csv"
	"os"
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
	//
	for i, val := range users {
		link = val.ProfileLink
		page, err = http.Get(link)
		logError(err)
		pagetext, err = ioutil.ReadAll(page.Body)
		logError(err)
		text = string(pagetext)
		doc, err = html.Parse(strings.NewReader(text))
		logError(err)
		generateProfileData(doc, i)
	}
	//fmt.Println(users)
	//data := [][]string{{"a", "b", "c"}, {"1", "2", "3"}, {"x", "y", "z"}}
	data := structToSlice(users)
	file, err := os.Create("first-go-csv.csv")
	logError(err)
	writer := csv.NewWriter(file)
	writer.WriteAll(data)
	file.Close()
}

var users []User
var user_profile_link string
var user_name string

type User struct {
	ProfileLink string
	Name string
	Location string
	TimeRegistered string
	LastSeen string
}

func generateUsersData(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" && len(node.Attr) >= 2 {
		if node.Attr[1].Key == "class" && node.Attr[1].Val == "user" {
			user_profile_link = "https://www.nairaland.com" + node.Attr[0].Val
			user_name = node.FirstChild.Data
			users = append(users, User{user_profile_link, user_name, "", "", ""})
		}
	}
	for i := node.FirstChild; i != nil; i = i.NextSibling {
		generateUsersData(i)
	}
}

func generateProfileData(node *html.Node, ind int) {
	if node.Type == html.ElementNode && node.Data == "b" && node.FirstChild.Data == "Location" {
		fmt.Println(node.NextSibling.Data)
		users[ind].Location = node.NextSibling.Data[2:]
	} else if node.Type == html.ElementNode && node.Data == "b" && node.FirstChild.Data == "Time registered" {
		fmt.Println(node.NextSibling.Data)
		users[ind].TimeRegistered = node.NextSibling.Data[2:]
	} else if node.Type == html.ElementNode && node.Data == "b" && node.FirstChild.Data == "Last seen" {
		if node.NextSibling.NextSibling.NextSibling != nil {
			fmt.Println(ind, users[ind].Name, node.NextSibling.Data + node.NextSibling.NextSibling.FirstChild.Data +
				node.NextSibling.NextSibling.NextSibling.Data +
				node.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild.Data)
			users[ind].LastSeen = node.NextSibling.Data +
				node.NextSibling.NextSibling.FirstChild.Data +
				node.NextSibling.NextSibling.NextSibling.Data +
				node.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild.Data
		} else {
			fmt.Println(ind, users[ind].Name, node.NextSibling.Data + node.NextSibling.NextSibling.FirstChild.Data)
			users[ind].LastSeen = node.NextSibling.Data + node.NextSibling.NextSibling.FirstChild.Data
		}
	}
	for i := node.FirstChild; i != nil; i = i.NextSibling {
		generateProfileData(i, ind)
	}
}

func logError(err error) {
	if err != nil {
		log.Fatal("Error encountered: ", err)
	}
}

func structToSlice(sliceOfStructs []User) [][]string {
	slice := [][]string{{"name", "profile_link", "location", "time_registered", "last_seen"}}
	for _, i := range sliceOfStructs {
		slice = append(slice, []string{i.Name, i.ProfileLink, i.Location, i.TimeRegistered, i.LastSeen})
	}
	return slice
}

















