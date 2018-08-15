package main

import "net/http"
import "fmt"
import "io/ioutil"
import "strings"
//import "math/rand"
//import "time"
/*import "database/sql" */
import "strconv"
import "os"


func main() {
    InsertIssuesInDatabase("justinsb")    
}
    
func GetNumberOfPagesOfIssues(user string) int {

    contents := SendRequest( GetGithubStartPath(user) )

    startanchor := "class=\"pagination\""
    endanchor := "class=\"next_page\""

    pagination := strings.Split(strings.Split(contents, startanchor)[1] , endanchor)[0]

    paginationsplit := strings.Split(pagination[:len(pagination)-8], ">")
    endpage, _ := strconv.Atoi(paginationsplit[len(paginationsplit)-1])
   
    fmt.Println(endpage) 
    
    return endpage

    }

func InsertIssuesInDatabase(user string) {
    numpages := GetNumberOfPagesOfIssues(user)

    for i:=1; i<=numpages; i++ {
      contents := SendRequest( GetGithubIssuePageStartPath(user, string(i)) )
      fmt.Println("page:",i,"/",numpages, GetGithubIssuePageStartPath(user, string(i)))
      InsertIssuesFromPage(contents, user)
    }
}

func InsertIssuesFromPage(contents string, user string) {
   
   startanchor := "li id=\"issue_"
   endanchor := "\""
   contentssplit := strings.Split(contents, startanchor)
   fmt.Println("Got",len(contentssplit),"issues")
   fmt.Println(contentssplit[0])

   os.Exit(1)

   for i:=1; i<len(contentssplit); i++ {
     fmt.Println("issue",i)
      number := strings.Split(contentssplit[i], endanchor)[0]
     InsertIssueInDataBase(number, user) 
   }
}

func InsertIssueInDataBase(number string, user string) {
    contents := SendRequest( GetIssueByNumber(number) )
    
    startanchor := "<div class=\"edit-comment-hide\">"
    endanchor := "<div class=\"comment-reactions  js-reactions-container \">"

    issuetext := strings.Split(strings.Split(contents, startanchor)[1] , endanchor)[0]

    fmt.Println("hey here is the issuetext ", issuetext)

/* database, _ := sql.Open("sqlite3", "./nraboy.db")
    statement, _ := database.Prepare("CREATE TABLE IF NOT EXISTS people (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
    statement.Exec()
    statement, _ = database.Prepare("INSERT INTO people (firstname, lastname) VALUES (?, ?)")
    statement.Exec("Nic", "Raboy")
    rows, _ := database.Query("SELECT id, firstname, lastname FROM people")
    var id int
    var firstname string
    var lastname string
    for rows.Next() {
        rows.Scan(&id, &firstname, &lastname)
        fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname)
*/
}
func SendRequest(url string) string {

    response, _ := http.Get(url)
    bodybytes, _ := ioutil.ReadAll(response.Body)
    return string(bodybytes)
    }

func GetGithubStartPath(user string) string {
    return "https://github.com/kubernetes/kubernetes/issues/created_by/"+user
}

func GetGithubIssuePageStartPath(user string, pagenum string) string {
    return "https://github.com/kubernetes/kubernetes/issues/created_by/"+user+"?page="+strconv.Itoa(1)
//+"&q=is%3Aopen+is%3Aissue+author%3A"+user
}

func GetIssueByNumber(number string) string { 
    return "https://github.com/kubernetes/kubernetes/issues/"+number
}

   /* https://github.com/kubernetes/kubernetes/issues/created_by/justinsb?page=2&q=is%3Aopen+is%3Aissue+author%3Ajustinsb
    */


