// Nairaland users data and web scraper

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"log"
	"golang.org/x/net/html"
	"strings"
	//"bytes"
	//"io"
	//"text/template"
	//"strconv"
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
var user User
var user_profile_link string
var user_name string
//var ind = 0

type User struct {
	ProfileLink string
	Name string
	//Location string
	//TimeRegistered string
	//LastSeen string
}

func generateUsersData(node *html.Node) {
	fmt.Println("loopout")
	if node.Type == html.ElementNode && node.Data == "a" && len(node.Attr) == 2 {
		//fmt.Println(node.Attr[1])
		if node.Attr[1].Key == "class" && node.Attr[1].Val == "user" {
			fmt.Println("loop")
			user_profile_link = node.Attr[0].Val
			user_name = node.FirstChild.Data
			user = User{user_profile_link, user_name}
			users = append(users, user)
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



















