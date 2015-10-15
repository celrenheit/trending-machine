package hubspider

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/celrenheit/spider"
	"github.com/celrenheit/spider/spiderutils"
)

const (
	baseURL = "https://github.com/trending"
	// specificLanguageFmt = "https://github.com/trending?l=%s"
)

type HubSpider struct {
	URL    string
	Parent bool
	Lang   string
	DB     *mgo.Database
}

var langs []string

func New(db *mgo.Database, languages []string) *HubSpider {
	langs = languages
	return &HubSpider{
		URL:    baseURL,
		Parent: true,
		Lang:   "__all__",
		DB:     db,
	}
}

func NewLanguageSpider(lang string) *HubSpider {
	var u string // url
	if lang == "__all__" {
		u = baseURL
	} else {
		u = baseURL + "?l=" + url.QueryEscape(strings.Replace(strings.Replace(strings.ToLower(lang), " ", "-", -1), "#", "sharp", -1))
	}

	return &HubSpider{
		URL:    u,
		Parent: false,
		Lang:   lang,
	}
}

func (h *HubSpider) Setup(parent *spider.Context) (*spider.Context, error) {
	ctx, err := spiderutils.NewHTTPContext("GET", h.URL, nil)
	if parent != nil {
		ctx.SetParent(parent)
	}
	return ctx, err
}

func (h *HubSpider) Spin(ctx *spider.Context) error {
	if h.Parent {
		snapshots, err := FindSnapshotByTime(h.DB, time.Now())
		if err == nil && snapshots != nil {
			fmt.Println("Checked today snapshot has already been saved.")
			return nil
		}
	}

	fmt.Println("Fetching ", h.Lang, "at: ", h.URL)
	if _, err := ctx.DoRequest(); err != nil {
		return err
	}
	fmt.Println("Status:", ctx.Response().Status, "for", h.Lang)

	doc, err := ctx.HTMLParser()
	if err != nil {
		return err
	}

	repos, err := ParseHTMLToRepos(doc)
	if err != nil {
		return err
	}

	ctx.Set("lang", h.Lang)
	ctx.Set("repos", repos)

	if !h.Parent {
		return nil
	}

	// Code to find all languages
	// langs := make([]langWithURL, 0)
	//
	// doc.Find(".select-menu .select-menu-list .select-menu-item[role=\"menuitem\"]").Each(func(i int, s *goquery.Selection) {
	// 	a := s.Find("a.select-menu-item-text")
	// 	lang := a.Text()
	// 	url, _ := a.Attr("href")
	// 	langs = append(langs, langWithURL{lang, url})
	// })

	allDoneChan := make(chan struct{})
	errChan := make(chan error)
	// Fetch for each languages
	var wg sync.WaitGroup
	for i, l := range langs {
		ls := NewLanguageSpider(l)
		newCtx, err := ls.Setup(ctx)
		if err != nil {
			return err
		}
		wg.Add(1)
		go func(ls spider.Spider, newCtx *spider.Context, i int) {
			defer wg.Done()
			time.Sleep(time.Duration(i) * 800 * time.Millisecond)
			if err := ls.Spin(newCtx); err != nil {
				errChan <- err
			}
		}(ls, newCtx, i)

	}

	go func() {
		wg.Wait()
		close(allDoneChan)
	}()

	select {
	case err := <-errChan:
		return err
	case <-allDoneChan:
		return h.SaveAll(ctx)
	}
}

func (h *HubSpider) SaveAll(ctx *spider.Context) error {
	langs := make(map[string]Repos)
	langs[ctx.Get("lang").(string)] = ctx.Get("repos").(Repos)
	for _, c := range ctx.Children {
		langs[c.Get("lang").(string)] = c.Get("repos").(Repos)
	}
	s := NewSnapshot(langs)
	if err := s.Save(h.DB); err != nil {
		return err
	}

	log.Println("Saved", len(langs), "languages.")
	return nil
}
