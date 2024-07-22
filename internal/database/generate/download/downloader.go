package main

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

const BaseURL = "https://datomatic.no-intro.org/index.php"

func NewDownloader(systemName string) (*Downloader, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	generate := &Downloader{
		client: &http.Client{
			Jar: jar,
		},

		SystemName: systemName,
	}

	return generate, nil
}

type Downloader struct {
	client *http.Client

	SystemName string
	SystemID   string
}

func (g *Downloader) URL() string {
	u, err := url.Parse(BaseURL)
	if err != nil {
		log.Fatal().Err(err).Str("url", BaseURL).Msg("Failed to parse base url")
	}
	q := u.Query()
	q.Set("page", "download")
	if g.SystemID != "" {
		q.Set("s", g.SystemID)
	}
	q.Set("op", "xml")
	u.RawQuery = q.Encode()
	return u.String()
}

func (g *Downloader) Run(ctx context.Context) error {
	var err error
	g.SystemID, err = g.getSystemID(ctx, g.SystemName)
	if err != nil {
		return err
	}

	postParams, err := g.getFormParams(ctx)
	if err != nil {
		return err
	}

	url, name, value, err := g.prepareDownload(ctx, postParams)
	if err != nil {
		return err
	}

	f, err := g.download(ctx, url, name, value)
	if err != nil {
		return err
	}
	defer func(f io.ReadCloser) {
		_ = f.Close()
	}(f)

	log.Info().Str("file", "nes.xml").Msg("Write to file")
	out, err := os.Create("internal/database/nointro/nes.xml")
	if err != nil {
		return err
	}

	if _, err := io.Copy(out, f); err != nil {
		return err
	}

	return out.Close()
}

func (g *Downloader) buildPostForm(fields map[string]string) (io.Reader, string, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	for k, v := range fields {
		if err := bodyWriter.WriteField(k, v); err != nil {
			return bodyBuf, bodyWriter.FormDataContentType(), err
		}
	}

	if err := bodyWriter.Close(); err != nil {
		return bodyBuf, bodyWriter.FormDataContentType(), err
	}

	return bodyBuf, bodyWriter.FormDataContentType(), nil
}

var ErrSystemNotFound = errors.New("could not find system ID")

var ErrNoForm = errors.New(`could not find form"`)

func (g *Downloader) getSystemID(ctx context.Context, systemName string) (string, error) {
	log.Info().
		Str("url", g.URL()).
		Str("systemName", g.SystemName).
		Msg("Get system ID")

	res, err := g.request(ctx, http.MethodGet, g.URL(), "", nil, http.StatusOK)
	if err != nil {
		return "", err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", err
	}

	sel := doc.Find(`form[name="main_form"]`)
	if sel.Length() == 0 {
		return "", ErrNoForm
	}

	node := sel.Eq(0)

	// Iterate over selects
	var system string
	node.Find("select").Each(func(_ int, s *goquery.Selection) {
		if name, exists := s.Attr("name"); exists && name == "system_selection" {
			s.Find("option").Each(func(_ int, s *goquery.Selection) {
				if strings.TrimSpace(s.Text()) == systemName {
					if value, exists := s.Attr("value"); exists {
						system = value
					}
				}
			})
		}
	})

	if system == "" {
		return "", fmt.Errorf("%w: %s", ErrSystemNotFound, systemName)
	}

	log.Info().Str("id", system).Msg("Got system ID")
	return system, nil
}

func (g *Downloader) getFormParams(ctx context.Context) (map[string]string, error) {
	postParams := make(map[string]string)

	// Get session cookie
	log.Info().Str("url", g.URL()).Msg("Get form params")
	res, err := g.request(ctx, http.MethodGet, g.URL(), "", nil, http.StatusOK)
	if err != nil {
		return postParams, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return postParams, err
	}

	sel := doc.Find(`form[name="main_form"]`)
	if sel.Length() == 0 {
		return postParams, ErrNoForm
	}

	node := sel.Eq(0)

	// Iterate over selects
	node.Find("select").Each(func(_ int, s *goquery.Selection) {
		if name, exists := s.Attr("name"); exists {
			s.Find("option").Each(func(_ int, s *goquery.Selection) {
				if _, exists := s.Attr("selected"); exists {
					if value, exists := s.Attr("value"); exists {
						postParams[name] = value
					}
				}
			})
		}
	})

	// Iterate over inputs
	node.Find("input").Each(func(_ int, s *goquery.Selection) {
		if inputType, exists := s.Attr("type"); exists {
			if inputType == "radio" {
				// Skip unchecked radio buttons
				if _, exists := s.Attr("checked"); !exists {
					return
				}
			}
		}

		if name, exists := s.Attr("name"); exists {
			if value, exists := s.Attr("value"); exists {
				postParams[name] = value
			}
		}
	})

	log.Info().Interface("params", postParams).Msg("Got form params")
	return postParams, nil
}

var (
	ErrMissingLocation = errors.New("missing location")
	ErrNoButton        = errors.New(`could not find download button`)
)

func (g *Downloader) prepareDownload(ctx context.Context, postFields map[string]string) (string, string, string, error) {
	// Get form params to request a download
	body, contentType, err := g.buildPostForm(postFields)
	if err != nil {
		return "", "", "", err
	}

	// Disable following redirects since we want the destination URL
	defer func(checkRedirect func(*http.Request, []*http.Request) error) {
		g.client.CheckRedirect = checkRedirect
	}(g.client.CheckRedirect)

	g.client.CheckRedirect = func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	}

	// Request download
	log.Info().Str("url", g.URL()).Msg("Request download")
	res, err := g.request(ctx, http.MethodPost, g.URL(), contentType, body, http.StatusFound)
	if err != nil {
		return "", "", "", err
	}
	_, _ = io.Copy(io.Discard, res.Body)
	_ = res.Body.Close()

	location := res.Header.Get("Location")
	if location == "" {
		return "", "", "", ErrMissingLocation
	}

	url, err := res.Request.URL.Parse(location)
	if err != nil {
		return "", "", "", err
	}

	res, err = g.request(ctx, http.MethodGet, url.String(), "", nil, http.StatusOK)
	if err != nil {
		return "", "", "", err
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return "", "", "", err
	}

	sel := doc.Find(`#content form`)
	if sel.Length() == 0 {
		return "", "", "", ErrNoForm
	}
	node := sel.Eq(0)

	// Iterate over selects
	sel = node.Find(`input[type="submit"]`)
	node = sel.Eq(0)

	name, ok := node.Attr("name")
	if !ok {
		return "", "", "", ErrNoButton
	}

	value, ok := node.Attr("value")
	if !ok {
		return "", "", "", ErrNoButton
	}

	// Drain response body
	_, _ = io.Copy(io.Discard, res.Body)
	_ = res.Body.Close()

	return url.String(), name, value, nil
}

var ErrNoFiles = errors.New("no files")

func (g *Downloader) download(ctx context.Context, url, name, value string) (io.ReadCloser, error) {
	body, contentType, err := g.buildPostForm(map[string]string{
		name: value,
	})
	if err != nil {
		return nil, err
	}

	log.Info().Str("url", url).Msg("Begin download")
	res, err := g.request(ctx, http.MethodPost, url, contentType, body, http.StatusOK)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.Copy(io.Discard, res.Body)
		_ = res.Body.Close()
	}()

	log.Info().Str("contentLength", res.Header.Get("Content-Length")).Msg("Began download")

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	_ = res.Body.Close()

	zipr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, err
	}

	if len(zipr.File) == 0 {
		return nil, ErrNoFiles
	}

	log.Info().Str("file", zipr.File[0].Name).Msg("Begin unzip")
	f, err := zipr.File[0].Open()
	if err != nil {
		return nil, err
	}
	return f, nil
}

var ErrInvalidResponse = errors.New("invalid response")

func (g *Downloader) request(ctx context.Context, method, url, contentType string, body io.Reader, expectStatus int) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	res, err := g.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != expectStatus {
		return nil, fmt.Errorf(
			`%w: got %q, expected %q`,
			ErrInvalidResponse,
			res.Status,
			fmt.Sprintf("%d %s", expectStatus, http.StatusText(expectStatus)),
		)
	}
	return res, nil
}
