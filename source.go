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
	"strconv"
)

func main() {
	var link string
	fmt.Println("Enter the thread link below (to process a default link, just press Enter)")
	fmt.Print("Link >>> ")
	fmt.Scanf("%s", &link)
	if link == "" {
		link = "https://www.nairaland.com/7243961/christian-how-often-pray"
	}
	link0 := link
	var page *http.Response
	var pageTrack *http.Response
	var ppage *http.Response
	var err error
	var pagetext []byte
	var text string
	var doc *html.Node
	x := 1
	for {
		page, err = http.Get(link)
		logError(err)
		if pageTrack == nil {
			// do nothing
		} else if page.Request.URL.Path == pageTrack.Request.URL.Path || x == 10000 {
			break
		}
		pagetext, err = ioutil.ReadAll(page.Body)
		logError(err)
		text = string(pagetext)
		doc, err = html.Parse(strings.NewReader(text))
		logError(err)
		generateUsersData(doc)
		fmt.Println("The processing of the webpage at", page.Request.URL.Path, "was successful!")
		link = link0 + "/" + strconv.Itoa(x)
		pageTrack = page
		x++
	}
	page.Body.Close()
	generateUniqueUsers()
	for i, val := range users {
		fmt.Println("Processing profile", i+1, "with username", users[i].Name, "...")
		link = val.ProfileLink
		ppage, err = http.Get(link)
		if err != nil {
			fmt.Println("Error processing profile at index", i+1, "[", err, "]")
			continue
		}
		pagetext, err = ioutil.ReadAll(ppage.Body)
		logError(err)
		text = string(pagetext)
		doc, err = html.Parse(strings.NewReader(text))
		logError(err)
		generateProfileData(doc, i)
	}
	ppage.Body.Close()
	data := structToSlice(users)
	file, err := os.Create("first-go-csv.csv")
	logError(err)
	writer := csv.NewWriter(file)
	writer.WriteAll(data)
	file.Close()
	fmt.Println("A csv file 'first-go-csv.csv' has been written to the current working directory")
	fmt.Println("DONE! Now, check the csv file :)")
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
	PersonalText string
	Gender string
	Twitter string
	TimeSpentOnline string
}

func generateUsersData(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" && len(node.Attr) >= 2 {
		if node.Attr[1].Key == "class" && node.Attr[1].Val == "user" {
			user_profile_link = "https://www.nairaland.com" + node.Attr[0].Val
			user_name = node.FirstChild.Data
			users = append(users, User{user_profile_link, user_name, "", "", "", "", "", "", ""})
		}
	}
	for i := node.FirstChild; i != nil; i = i.NextSibling {
		generateUsersData(i)
	}
}

func generateProfileData(node *html.Node, ind int) {
	if node.Type == html.ElementNode && node.Data == "b" && node.FirstChild.Data == "Location" {
		users[ind].Location = node.NextSibling.Data[2: ]
	} else if node.Type == html.ElementNode && node.Data == "b" && node.FirstChild.Data == "Time registered" {
		users[ind].TimeRegistered = node.NextSibling.Data[2: ]
	} else if node.Type == html.ElementNode && node.Data == "b" && node.FirstChild.Data == "Last seen" {
		if node.NextSibling.NextSibling.NextSibling != nil {
			users[ind].LastSeen = (node.NextSibling.Data +
				node.NextSibling.NextSibling.FirstChild.Data +
				node.NextSibling.NextSibling.NextSibling.Data +
				node.NextSibling.NextSibling.NextSibling.NextSibling.FirstChild.Data)[2: ]
		} else {
			users[ind].LastSeen = (node.NextSibling.Data + node.NextSibling.NextSibling.FirstChild.Data)[2: ]
		}
	} else if node.Type == html.ElementNode && node.Data == "b" && node.FirstChild.Data == "Personal text" {
		users[ind].PersonalText = node.NextSibling.Data[2: ]
	} else if node.Type == html.ElementNode && node.Data == "b" && node.FirstChild.Data == "Gender" {
		users[ind].Gender = node.NextSibling.Data[2: ]
	} else if node.Type == html.ElementNode && node.Data == "b" && node.FirstChild.Data == "Twitter" {
		users[ind].Twitter = node.NextSibling.Data[2: ]
	} else if node.Type == html.ElementNode && node.Data == "b" && node.FirstChild.Data == "Time spent online" {
		users[ind].TimeSpentOnline = node.NextSibling.Data[2: ]
	}
	for i := node.FirstChild; i != nil; i = i.NextSibling {
		generateProfileData(i, ind)
	}
}

func generateUniqueUsers() {
	unique_users := []User{}
	for _, val := range users {
		if sliceContains(unique_users, val) {
			// do nothing
		} else {
			unique_users = append(unique_users, val)
		}
	}
	users = unique_users
}

func sliceContains(suser []User, user User) bool {
    for _, a := range suser {
        if a.ProfileLink == user.ProfileLink {
            return true
        }
    }
    return false
}

func logError(err error) {
	if err != nil {
		log.Fatal("Error encountered: ", err)
	}
}

func structToSlice(sliceOfStructs []User) [][]string {
	slice := [][]string{{"name", "profile_link", "location", "time_registered", "last_seen", "personal_text", "gender", "twitter", "time_spent_online"}}
	for _, i := range sliceOfStructs {
		slice = append(slice, []string{i.Name, i.ProfileLink, i.Location, i.TimeRegistered, i.LastSeen, i.PersonalText, i.Gender, i.Twitter, i.TimeSpentOnline})
	}
	return slice
}

















