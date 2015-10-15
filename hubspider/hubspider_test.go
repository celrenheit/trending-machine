package hubspider

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"gopkg.in/mgo.v2"
)

func TestCorrectURL(t *testing.T) {
	var tests = []struct {
		lang, expected string
	}{
		{"Objective-C++", "https://github.com/trending?l=objective-c%2B%2B"},
		{"Cap'n Proto", "https://github.com/trending?l=cap%27n-proto"},
		{"API Blueprint", "https://github.com/trending?l=api-blueprint"},
		{"DIGITAL Command Language", "https://github.com/trending?l=digital-command-language"},
		{"C#", "https://github.com/trending?l=csharp"},
	}

	for _, tt := range tests {
		actual := NewLanguageSpider(tt.lang)
		if actual.URL != tt.expected {
			t.Errorf("URL(%s): expected %s, actual %s", tt.lang, tt.expected, actual.URL)
		}
	}
}

func TestImport(t *testing.T) {
	session, _ := mgo.Dial(":27017")
	db := session.DB("ghtrending")

	saveAllDir(db, "github-trending/")
}

func saveAllDir(db *mgo.Database, path string) {

	list, _ := ioutil.ReadDir(path)
	for _, f := range list {
		if f.IsDir() && strings.HasPrefix(f.Name(), "2014") {
			fmt.Println(f.Name())
			saveAllDir(db, filepath.Join(path, f.Name()))
		} else if strings.HasSuffix(f.Name(), ".md") && (strings.HasPrefix(f.Name(), "2015") || strings.HasPrefix(f.Name(), "2014")) {
			file, err := os.Open(filepath.Join(path, f.Name()))
			if err != nil {
				fmt.Println(err)
				return
			}
			defer file.Close()

			contents, err := ioutil.ReadAll(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			unsafe := blackfriday.MarkdownBasic(contents)
			html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
			// fmt.Println(string(html))
			doc, err := goquery.NewDocumentFromReader(bytes.NewReader(html))
			langs := make(map[string]Repos)
			doc.Find("h4").Each(func(i int, hs *goquery.Selection) {
				var repos Repos

				hs.NextUntil("h4").Find("li").Each(func(i int, s *goquery.Selection) {
					repo := &Repo{}
					repo.URL, _ = s.Find("a").Attr("href")
					repo.Name = strings.Replace(strings.Replace(s.Find("a").Text(), " ", "", -1), "\n", "", -1)
					// repo.Description =
					parts := strings.Split(s.Text(), ":")
					repo.Description = strings.TrimSpace(strings.Replace(strings.Join(parts[1:], ""), "\n", "", -1))

					img, ok := s.Find("img").Attr("src")
					if ok {
						repo.BuiltBy = Contributors{Contributor{
							Avatar: img,
							// Username: "$$unknown$$",
						}}
					}
					repos = append(repos, repo)
				})
				fmt.Println("len:", len(repos))
				langs[strings.Title(hs.Text())] = repos
			})
			_, filename := filepath.Split(file.Name())
			ymd := strings.Split(strings.TrimSuffix(filename, ".md"), "-")
			s := &Snapshot{
				Date:      fmt.Sprintf("%s-%s-%s", ymd[2], ymd[1], ymd[0]),
				Languages: langs,
			}
			s.Save(db)
		}
	}
}
