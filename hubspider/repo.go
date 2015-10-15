package hubspider

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Repo struct {
	Name        string       `bson:"name" json:"name"`
	Description string       `bson:"description" json:"description"`
	URL         string       `bson:"url" json:"url"`
	Owner       string       `bson:"owner" json:"owner"`
	BuiltBy     Contributors `bson:"built_by" json:"built_by"`
}
type Repos []*Repo

type Contributor struct {
	Avatar   string `bson:"avatar" json:"avatar"`
	Username string `bson:"username" json:"username"`
}

type Contributors []Contributor

func ParseHTMLToRepos(doc *goquery.Document) (Repos, error) {
	var repos Repos

	doc.Find("li.repo-list-item").Each(func(i int, s *goquery.Selection) {
		repo := &Repo{}
		repoListName := s.Find(".repo-list-name")
		repo.Name = strings.Replace(strings.Replace(repoListName.Text(), " ", "", -1), "\n", "", -1)
		url, _ := repoListName.Find("a").Attr("href")
		repo.URL = "https://github.com" + url
		repo.Owner = repoListName.Find(".prefix").Text()

		repo.Description = strings.Replace(strings.TrimSpace(s.Find("p.repo-list-description").Text()), "\n", "", -1)
		s.Find(".repo-list-meta a img").Each(func(i int, s *goquery.Selection) {
			c := Contributor{}
			c.Username, _ = s.Attr("alt")
			c.Avatar, _ = s.Attr("src")
			c.Username = strings.Replace(c.Username, "@", "", 1)
			repo.BuiltBy = append(repo.BuiltBy, c)
		})
		repos = append(repos, repo)
	})

	return repos, nil
}
