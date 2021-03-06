package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Booklist struct {
	Booklist struct {
		Booklist struct {
			Name         string   `json:"name"`
			BookIds      []string `json:"bookIds"`
			Description  string   `json:"description"`
			BooklistId   string   `json:"booklistId"`
			TotalCount   int      `json:"totalCount"`
			CollectCount int      `json:"collectCount"`
		} `json:"booklist"`
	} `json:"booklist"`
}

type Booklists struct {
	Booklists []struct {
		BooklistId   string `json:"booklistId"`
		CollectCount int    `json:"collectCount"`
		Name         string `json:"name"`
		Description  string `json:"description"`
		TotalCount   int    `json:"totalCount"`
	} `json:"booklists"`
}

type BookInfo struct {
	BookId         string `json:"bookId"`
	Title          string `json:"title"`
	Author         string `json:"author"`
	Cover          string `json:"cover"`
	Category       string `json:"category"`
	Intro          string `json:"intro"`
	Publisher      string `json:"publisher"`
	PublishTime    string `json:"publishTime"`
	BookSize       int    `json:"bookSize"`
	Star           int    `json:"star"`
	NewRatingCount int    `json:"newRatingCount"`
}

type BestBookMarkItem struct {
	Bookid     string `json:"bookId"`
	Uservid    string `json:"userVid"`
	Bookmarkid string `json:"bookmarkId"`
	Chapteruid int    `json:"chapterUid"`
	Range      string `json:"range"`
	Marktext   string `json:"markText"`
	Totalcount int    `json:"totalCount"`
}

type BestBookMarks struct {
	Totalcount string             `json:"totalCount"`
	Items      []BestBookMarkItem `json:"items"`
	Chapters   []struct {
		Bookid     string `json:"bookId"`
		Chapteruid int    `json:"chapterUid"`
		Chapteridx int    `json:"chapterIdx"`
		Title      string `json:"title"`
	} `json:"chapters"`
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	getBookLists()
}

func getBookLists() {
	result, err := getAndSaveResponse("booklists", "https://i.weread.qq.com/market/booklists?count=200&type=0")
	check(err)
	var booklists Booklists
	json.Unmarshal([]byte(result), &booklists)

	f, err := os.Create("README.md")
	check(err)
	defer f.Close()
	fmt.Fprintln(f, "# ??????????????????????????????")

	for index, booklist := range booklists.Booklists {
		log.Println(index, booklist.Name, booklist.CollectCount, booklist.TotalCount)
		fmt.Fprintln(f, "\n###", strconv.Itoa(index+1), booklist.Name)
		fmt.Fprintln(f, "\n"+booklist.Description)

		bookIds, err := getBookList(booklist.BooklistId)
		check(err)
		fmt.Fprintln(f, "\n<details>\n<summary>"+strconv.Itoa(booklist.TotalCount)+" books, "+strconv.Itoa(booklist.CollectCount)+" likes"+"</summary>")
		for index, bookId := range bookIds {
			if index > 30 {
				break
			}
			bookInfo := getBookInfo(bookId)
			booklist.Name = strings.TrimSpace(booklist.Name)
			booklist.Name = strings.ReplaceAll(booklist.Name, " ", "%20")
			booklist.Name = strings.ReplaceAll(booklist.Name, "/", "-")

			bookInfo.Title = strings.ReplaceAll(bookInfo.Title, "/", "-")
			fmt.Fprintf(f, "\n1. [%s](books/%s/%s.md)", bookInfo.Title, booklist.Name, strings.ReplaceAll(bookInfo.Title, " ", "%20"))

			bestMark := getBestBookMarks(bookInfo, booklist.Name)
			if bestMark.Totalcount > 300 {
				fmt.Fprintln(f, "\n\t> "+strconv.Itoa(bestMark.Totalcount)+" "+bestMark.Marktext)
			}
		}
		fmt.Fprintln(f, "\n</details>")
		// os.Exit(1)
	}
}

func getBookInfo(bookId string) BookInfo {
	result, err := getAndSaveResponse(bookId+"-bookinfo", "https://i.weread.qq.com/book/info?bookId="+bookId)
	check(err)
	var bookInfo BookInfo
	json.Unmarshal([]byte(result), &bookInfo)
	check(err)
	return bookInfo
}

func getBookList(booklistId string) ([]string, error) {
	result, err := getAndSaveResponse(booklistId, "https://i.weread.qq.com/booklist/single?booklistId="+booklistId)
	check(err)
	var booklist Booklist
	json.Unmarshal([]byte(result), &booklist)

	bl := booklist.Booklist.Booklist
	bl.Name = strings.TrimSpace(bl.Name)
	bl.Name = strings.ReplaceAll(bl.Name, "/", "-")
	bl.Name = strings.ReplaceAll(bl.Name, "%20", " ")

	path := "books/" + bl.Name
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0755)
		check(err)
	}
	return bl.BookIds, err
}

func getAndSaveResponse(id, URL string) ([]byte, error) {
	fileName := "temp/" + id + ".json"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		rawCookies, err := ioutil.ReadFile("cookie.txt")
		check(err)

		rawRequest := fmt.Sprintf("GET / HTTP/1.0\r\n%s\r\n\r\n", rawCookies)
		req, _ := http.ReadRequest(bufio.NewReader(strings.NewReader(rawRequest)))
		url, _ := url.Parse(URL)
		jar, _ := cookiejar.New(nil)
		jar.SetCookies(url, req.Cookies())
		client := http.Client{Jar: jar}
		request, _ := http.NewRequest("GET", URL, nil)
		request.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:88.0) Gecko/20100101 Firefox/88.0")
		r, err := client.Do(request)
		check(err)
		defer r.Body.Close()

		body, _ := ioutil.ReadAll(r.Body)
		err = ioutil.WriteFile("temp/"+id+".json", body, 0644)
		check(err)
		log.Println("StatusCode:", r.StatusCode, id, "from:", url, "saved")
		time.Sleep(200 * time.Millisecond)
		return body, err
	}
	body, err := ioutil.ReadFile(fileName)
	if strings.Contains(string(body), "errcode") {
		os.Remove(fileName)
		log.Fatal(string(body))
	}
	return body, err
}

func getBestBookMarks(bookInfo BookInfo, booklistName string) BestBookMarkItem {
	bookId := bookInfo.BookId

	booklistName = strings.ReplaceAll(booklistName, "%20", " ")
	booklistName = strings.ReplaceAll(booklistName, "/", "-")
	f, err := os.Create("books/" + booklistName + "/" + strings.ReplaceAll(bookInfo.Title, "/", "-") + ".md")
	check(err)
	defer f.Close()

	fmt.Fprintln(f, "## "+bookInfo.Title)
	fmt.Fprintln(f, "\n"+bookInfo.Author, " - ", bookInfo.Category)
	fmt.Fprintln(f, "\n> "+bookInfo.Intro)

	result, err := getAndSaveResponse(bookId, "https://i.weread.qq.com/book/bestbookmarks?count=5000&bookId="+bookId)
	check(err)

	var bestBookMarks BestBookMarks
	json.Unmarshal([]byte(result), &bestBookMarks)

	// sort chapters and marks
	sort.Slice(bestBookMarks.Chapters, func(i, j int) bool {
		return bestBookMarks.Chapters[i].Chapteruid < bestBookMarks.Chapters[j].Chapteruid
	})

	sort.Slice(bestBookMarks.Items, func(i, j int) bool {
		range1, _ := strconv.Atoi(strings.Split(bestBookMarks.Items[i].Range, "-")[0])
		range2, _ := strconv.Atoi(strings.Split(bestBookMarks.Items[j].Range, "-")[0])
		return range1 < range2
	})

	if len(bestBookMarks.Items) > 0 {
		bestMark := bestBookMarks.Items[0]
		for _, v := range bestBookMarks.Chapters {
			fmt.Fprintln(f, "\n###", v.Title)
			for i := 0; i < len(bestBookMarks.Items); i++ {
				if bestBookMarks.Items[i].Chapteruid == v.Chapteruid {
					fmt.Fprintln(f, "\n"+bestBookMarks.Items[i].Marktext+" c:"+strconv.Itoa(bestBookMarks.Items[i].Totalcount))
				}
				if bestBookMarks.Items[i].Totalcount > bestMark.Totalcount {
					bestMark = bestBookMarks.Items[i]
				}
			}
		}
		if bestMark.Totalcount > 10 {
			// log.Println(bestMark.Totalcount, bestMark.Marktext)
			return bestMark
		}
	}
	return BestBookMarkItem{}
}
