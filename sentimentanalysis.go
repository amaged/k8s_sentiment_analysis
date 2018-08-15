
package main

import "net/http"

import "fmt"
import "io/ioutil"
import "strings"

//import "math/rand"
//import "time"
import "database/sql"
import "strconv"

import _ "github.com/mattn/go-sqlite3"

//import "os"

//global database
var database *sql.DB

func main() {

	databaseres, _ := sql.Open("sqlite3", "issues.db")
	database = databaseres
	//statement, _ := databaseres.Prepare("CREATE TABLE IF NOT EXISTS issues (id INTEGER PRIMARY KEY, username TEXT, title TEXT, message TEXT)")
	//statement.Exec()

	//given username, add information to database
	ParseAndInsertIssues("justinsb")
}

//given username, insert title, message, username into database
func ParseAndInsertIssues(user string) {

	//get ending page number
	paginationcontents := SendRequest(GetGithubStartPath(user))
	pagination := strings.Split(strings.Split(paginationcontents, "class=\"pagination\"")[1], "class=\"next_page\"")[0]
	paginationsplit := strings.Split(pagination[:len(pagination)-8], ">")
	numpages, _ := strconv.Atoi(paginationsplit[len(paginationsplit)-1])

	//get each page of issue
	for i := 1; i <= numpages; i++ {
		//get list of issue numbers list is formatted like:  id="issue_47"
		issuenumberlist := strings.Split(SendRequest(GetGithubIssuePageStartPath(user, strconv.Itoa(i))), "li id=\"issue_")
		//fmt.Println("page:",i,"/",numpages, GetGithubIssuePageStartPath(user, strconv.Itoa(i)))
		//fmt.Println("Got",len(contentssplit),"issues")
		//fmt.Println(contentssplit[0])

		//for every issue number on the current page
		for j := 1; j < len(issuenumberlist); j++ {
			//trim off part after number, and get full issue html
			issuenumber := strings.Split(issuenumberlist[j], "\"")[0]
			issuecontents := SendRequest(GetIssueByNumber(issuenumber))
			//InsertIssueInDataBase(number, user)

			issuetext := strings.Split(strings.Split(issuecontents, "<div class=\"edit-comment-hide\">")[1], "<div class=\"comment-reactions  js-reactions-container \">")[0]
			issuetextparsed_1 := strings.Split(issuetext, "</p>")

			for k := 0; k < len(issuetextparsed_1)-1; k++ {
				issuetextparsed_2 := strings.Split(issuetextparsed_1[k], "<p>")
				if(len(issuetextparsed_2)>1) {
					InsertIssue(user, issuetextparsed_2[1])
					fmt.Println(i,j,k)
				}
			}
		}
	}
}

func InsertIssue(username string, message string) {
	//fmt.Println("InsertIssue title:", title, "username:", username, "message:", message)

	database, _ := sql.Open("sqlite3", "issues.db")
	statement, _ := database.Prepare("INSERT INTO issues (username, message) VALUES (?, ?)")
	statement.Exec(username, message)
}

/*rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
  var id int
  var firstname string
  var lastname string
  for rows.Next() {
      rows.Scan(&id, &firstname, &lastname)
      fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
*/

func SendRequest(url string) string {
	response, _ := http.Get(url)
	bodybytes, _ := ioutil.ReadAll(response.Body)
	return string(bodybytes)
}

func GetGithubStartPath(user string) string {
	return "https://github.com/kubernetes/kubernetes/issues/created_by/" + user
}

func GetGithubIssuePageStartPath(user string, pagenum string) string {
	return "https://github.com/kubernetes/kubernetes/issues/created_by/" + user + "?page=" + pagenum
	//+"&q=is%3Aopen+is%3Aissue+author%3A"+user
}

func GetIssueByNumber(number string) string {
	return "https://github.com/kubernetes/kubernetes/issues/" + number
}

/* https://github.com/kubernetes/kubernetes/issues/created_by/justinsb?page=2&q=is%3Aopen+is%3Aissue+author%3Ajustinsb
 */
