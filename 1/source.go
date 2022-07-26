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
	"regexp"
)

func main() {
	var forumLink string
	fmt.Print("Enter forum link >>> ")
	fmt.Scanf("%s", &forumLink)
	if forumLink == "" {
		forumLink = "https://www.nairaland.com/education"
	}
	generateThreadLinks(forumLink)
	var page *http.Response
	var pageTrack *http.Response
	var ppage *http.Response
	var err error
	var pagetext []byte
	var text string
	var doc *html.Node
	var i = 0
	var usersTrack []User
	for j, link := range threads {
		link0 := link
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
	users = generateUniqueUsers(users)
	setDiff(users, usersTrack)
	fmt.Println("Total number of profile collected from thread [", link, "]:", len(sub_users))
	fmt.Println("Total number of profile collected [cumulative]:", len(users))
	usersTrack = users
	for _, val := range sub_users {
		fmt.Println("Processing profile", i+1, "with username", users[i].Name, "...")
		link = val.ProfileLink
		ppage, err = http.Get(link)
		if err != nil {
			fmt.Println("Error processing profile at index", i+1, "[", err, "]")
			i++
			continue
		}
		pagetext, err = ioutil.ReadAll(ppage.Body)
		logError(err)
		text = string(pagetext)
		doc, err = html.Parse(strings.NewReader(text))
		logError(err)
		generateProfileData(doc, i)
		i++
	}
	ppage.Body.Close()
	sub_users = []User{}
	fmt.Println("Thread", j+1, "processed")
	}
	columnsAmend()
	data := structToSlice(users)
	file, err := os.Create("out-data/first-go-csv.csv")
	logError(err)
	writer := csv.NewWriter(file)
	writer.WriteAll(data)
	file.Close()
	fmt.Println("A csv file 'first-go-csv.csv' has been written to the sub-directory 'out-data'")
	fmt.Println("DONE! Now, check the csv file :)")
}

var users []User
var sub_users []User
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
	NoOfPosts string
	NoOfTopics string
	//NoOfFollowing string
}

func generateUsersData(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "a" && len(node.Attr) >= 2 {
		if node.Attr[1].Key == "class" && node.Attr[1].Val == "user" {
			user_profile_link = "https://www.nairaland.com" + node.Attr[0].Val
			user_name = node.FirstChild.Data
			users = append(users, User{user_profile_link, user_name, "", "", "", "", "", "", "", "", ""})
			sub_users = append(sub_users, User{user_profile_link, user_name, "", "", "", "", "", "", "", "", ""})
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
		on := node.NextSibling.NextSibling.NextSibling
		if on != nil {
			if on.NextSibling.NextSibling == nil {
				users[ind].LastSeen = (node.NextSibling.Data +
					node.NextSibling.NextSibling.FirstChild.Data +
					on.Data +
					on.NextSibling.FirstChild.Data)[2: ]
			} else {
				users[ind].LastSeen = (node.NextSibling.Data +
					node.NextSibling.NextSibling.FirstChild.Data +
					on.Data +
					on.NextSibling.FirstChild.Data +
					on.NextSibling.NextSibling.Data +
					on.NextSibling.NextSibling.NextSibling.FirstChild.Data)[2: ]
			}
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
	} else if node.Type == html.ElementNode && node.Data == "a" && node.FirstChild != nil {
		regex1, _ := regexp.Compile("view " + strings.ToLower(users[ind].Name) + "'s posts \\([0-9]+\\)")
		regex2, _ := regexp.Compile("view " + strings.ToLower(users[ind].Name) + "'s topics \\([0-9]+\\)")
		nFDlower := strings.ToLower(node.FirstChild.Data)
		match1 := regex1.MatchString(nFDlower); match2 := regex2.MatchString(nFDlower)
		if match1 == true {
			regex, _ := regexp.Compile("\\([0-9]+\\)")
			match := regex.FindString(nFDlower)
			nop := strings.ReplaceAll(match, ")", ""); nop = strings.ReplaceAll(nop, "(", "")
			users[ind].NoOfPosts = nop
		} else if match2 == true {
			regex, _ := regexp.Compile("\\([0-9]+\\)")
			match := regex.FindString(nFDlower)
			nop := strings.ReplaceAll(match, ")", ""); nop = strings.ReplaceAll(nop, "(", "")
			users[ind].NoOfTopics = nop
		}
	}
	for i := node.FirstChild; i != nil; i = i.NextSibling {
		generateProfileData(i, ind)
	}
}

var threads []string
func generateThreadLinks(forumLink string) {
	page, err := http.Get(forumLink)
	logError(err)
	pagetext, err := ioutil.ReadAll(page.Body)
	logError(err)
	text := string(pagetext)
	doc, err := html.Parse(strings.NewReader(text))
	logError(err)
	parseForumLinks(doc)
}

var next2NodeParse = false
var nextNodeParse = false
func parseForumLinks(node *html.Node) {
	if node.Type == html.ElementNode && node.Data == "img" && node.NextSibling != nil {
		if node.NextSibling.Data == " " {
			next2NodeParse = true
		}
	}
	if node != nil && node.NextSibling != nil {
		if node.Data == " " && node.NextSibling.Data == "b" {
			nextNodeParse = true
		}
	}
	if node.Type == html.ElementNode && node.Data == "b" && node.FirstChild.Data == "a" {
		if next2NodeParse && nextNodeParse {
			threadLink := "https://www.nairaland.com" + node.FirstChild.Attr[0].Val
			threads = append(threads, threadLink)
			next2NodeParse = false; nextNodeParse = false
		}
	}
	for i := node.FirstChild; i != nil; i = i.NextSibling {
		parseForumLinks(i)
	}
}

func generateUniqueUsers(usersSlice []User) []User {
	unique_users := []User{}
	for _, val := range usersSlice {
		if sliceContains(unique_users, val) {
			// do nothing
		} else {
			unique_users = append(unique_users, val)
		}
	}
	//users = unique_users
	return unique_users
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
	slice := [][]string{{"name", "profile_link", "location", "time_registered", "last_seen", "personal_text", "gender",
	"twitter", "time_spent_online", "no_of_posts", "no_of_topics"}}
	for _, i := range sliceOfStructs {
		slice = append(slice, []string{i.Name, i.ProfileLink, i.Location, i.TimeRegistered, i.LastSeen, i.PersonalText, i.Gender,
			i.Twitter, i.TimeSpentOnline, i.NoOfPosts, i.NoOfTopics})
	}
	return slice
}

func columnsAmend() {
	//replace all empty fields of 'no_of_topics' with '0'
	for i, user := range users {
		if user.NoOfTopics == "" {
			users[i].NoOfTopics = "0"
		}
	}
}

func setDiff(major []User, minor []User) {
	var track bool
	sub_users = []User{}
	for _, i := range major {
		track = true
		for _, j := range minor {
			if i == j {
				track = false
				break
			}
		}
		if track == true {
			sub_users = append(sub_users, i)
		}
	}
}













