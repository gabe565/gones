package main

import (
	"archive/zip"
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
)

const BaseUrl = "https://datomatic.no-intro.org/index.php"
const System = "45"

var postFields = map[string]string{
	"system_selection": "45",
	"sys_list_order":   "2",
	"naming":           "0",
	"hash":             "0",
	"fill_parents":     "1",
	"hdl_xml_":         "Prepare",
}

func NewDownloader() (*Downloader, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	generate := &Downloader{
		client: &http.Client{
			Jar: jar,
		},
	}

	return generate, nil
}

type Downloader struct {
	client *http.Client
}

func (g *Downloader) Url() string {
	u, err := url.Parse(BaseUrl)
	if err != nil {
		log.Panic(err)
	}
	q := u.Query()
	q.Set("page", "download")
	q.Set("s", System)
	q.Set("op", "xml")
	u.RawQuery = q.Encode()
	return u.String()
}

func (g *Downloader) Run() error {
	if err := g.initSession(); err != nil {
		return err
	}

	url, err := g.prepareDownload()
	if err != nil {
		return err
	}

	res, err := g.download(url)
	if err != nil {
		return err
	}

	f, err := g.unzip(res.Body)
	if err != nil {
		return err
	}
	defer func(f io.ReadCloser) {
		_ = f.Close()
	}(f)

	log.WithField("file", "nes.xml").Info("Write to file")
	out, err := os.Create("../nointro/nes.xml")
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

var ErrInvalidResponse = errors.New("invalid response")

func (g *Downloader) initSession() error {
	// Get session cookie
	log.WithField("url", BaseUrl).Info("Initialize session")
	res, err := g.client.Head(BaseUrl)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return ErrInvalidResponse
	}

	return nil
}

var ErrMissingLocation = errors.New("missing location")

func (g *Downloader) prepareDownload() (string, error) {
	// Get form params to request a download
	body, contentType, err := g.buildPostForm(postFields)
	if err != nil {
		return "", err
	}

	// Disable following redirects since we want the destination URL
	defer func(checkRedirect func(*http.Request, []*http.Request) error) {
		g.client.CheckRedirect = checkRedirect
	}(g.client.CheckRedirect)

	g.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	// Request download
	log.WithField("url", g.Url()).Info("Request download")
	res, err := g.client.Post(g.Url(), contentType, body)
	if err != nil {
		return "", err
	}

	// Drain response body
	_, _ = io.Copy(io.Discard, res.Body)
	_ = res.Body.Close()

	if res.StatusCode != http.StatusFound {
		return "", ErrInvalidResponse
	}

	location := res.Header.Get("Location")
	if location == "" {
		return "", ErrMissingLocation
	}

	url, err := res.Request.URL.Parse(location)
	if err != nil {
		return "", err
	}

	return url.String(), nil
}

func (g *Downloader) download(url string) (*http.Response, error) {
	body, contentType, err := g.buildPostForm(map[string]string{
		"lazy_mode": "Download",
	})
	if err != nil {
		return nil, err
	}

	log.WithField("url", url).Info("Begin download")
	res, err := g.client.Post(url, contentType, body)
	if err != nil {
		return res, err
	}

	if res.StatusCode != http.StatusOK {
		return res, ErrInvalidResponse
	}

	return res, nil
}

var ErrNoFiles = errors.New("no files")

func (g *Downloader) unzip(src io.ReadCloser) (io.ReadCloser, error) {
	b, err := io.ReadAll(src)
	if err != nil {
		return nil, err
	}
	_ = src.Close()

	zipr, err := zip.NewReader(bytes.NewReader(b), int64(len(b)))
	if err != nil {
		return nil, err
	}

	if len(zipr.File) == 0 {
		return nil, ErrNoFiles
	}

	log.WithField("file", zipr.File[0].Name).Info("Begin unzip")
	f, err := zipr.File[0].Open()
	if err != nil {
		return nil, err
	}

	return f, nil
}
