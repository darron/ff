package service

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/darron/ff/core"
	"github.com/labstack/echo/v4"

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
	resp, err := http.Get(ns.URL)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()
	html, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// Strip all tags.
	p := bluemonday.StrictPolicy()
	text := p.Sanitize(string(html))
	fmt.Println(text)
	// TODO: Save it in the DB.
	// ns.BodyText = text
	// nsID, err := s.conf.NewsStoryRepository.Store(ns)
	// if err != nil {
	// 	return c.JSON(http.StatusInternalServerError, err.Error())
	// }
	return c.JSON(http.StatusCreated, id)
}
