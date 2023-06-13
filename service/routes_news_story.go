package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/darron/ff/core"
	"github.com/labstack/echo/v4"
	"gopkg.in/guregu/null.v4"

	"github.com/microcosm-cc/bluemonday"
)

func (s HTTPService) GetNewsStory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusNotFound, "id must not be blank")
	}
	ns, err := s.conf.NewsStoryRepository.Find(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	if ns.ID == "" || ns.RecordID == "" {
		return c.JSON(http.StatusNotFound, "that id does not exist")
	}
	return c.JSON(http.StatusOK, ns)
}

func (s HTTPService) CreateNewsStory(c echo.Context) error {
	ns := &core.NewsStory{}
	if err := c.Bind(ns); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}
	if ns.RecordID == "" {
		return c.JSON(http.StatusBadRequest, errors.New("must supply record ID with request"))
	}
	id, err := s.conf.NewsStoryRepository.Store(ns)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	return c.JSON(http.StatusCreated, id)
}

func (s HTTPService) DownloadNewsStory(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusNotFound, "id must not be blank")
	}
	ns, err := s.conf.NewsStoryRepository.Find(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// If news story exists - let's:
	// Download the HTML.
	text, err := getNewsText(ns.URL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	ns.BodyText = null.NewString(text, true)
	_, err = s.conf.NewsStoryRepository.Store(ns)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusNoContent, ns.ID)
}

func (s HTTPService) DownloadAllNewsStories(c echo.Context) error {
	// Get all Records.
	records, err := s.conf.RecordRepository.GetAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	ns := null.NewString("", false)
	// Let's loop through all the stories.
	for _, record := range records {
		for _, story := range record.NewsStories {
			if story.BodyText == ns {
				fmt.Println("Getting ", story.URL)
				text, err := getNewsText(story.URL)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, err.Error())
				}
				story.BodyText = null.NewString(text, true)
				_, err = s.conf.NewsStoryRepository.Store(&story)
				if err != nil {
					return c.JSON(http.StatusInternalServerError, err.Error())
				}
			}
		}
	}
	return c.JSON(http.StatusNoContent, "got 'em all")
}

func getNewsText(u string) (string, error) {
	resp, err := http.Get(u)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	// Strip all tags and clean it up a bit.
	// TODO: This still needs a ton of work.
	p := bluemonday.StrictPolicy()
	text := p.Sanitize(string(html))
	return cleanupDownload(text), nil
}

func cleanupDownload(text string) string {
	// Let's get rid of ALL the extra lines.
	re := regexp.MustCompile("(?m)^\\s*$[\r\n]*")
	text = strings.Trim(re.ReplaceAllString(text, ""), "\r\n")
	// Let's get rid of all leading and trailing spaces.
	var newLines []string
	for _, line := range strings.Split(text, "\n") {
		newLines = append(newLines, strings.TrimSpace(line))
	}
	text = strings.Join(newLines, "\n")
	return text
}
