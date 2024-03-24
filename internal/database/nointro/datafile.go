package nointro

import "encoding/xml"

// http://www.logiqx.com/Dats/datafile.dtd

type Datafile struct {
	XMLName xml.Name `xml:"datafile"`
	Headers []Header `xml:"header"`
	Games   []Game   `xml:"game"`
}

type Header struct {
	XMLName     xml.Name `xml:"header"`
	Name        string   `xml:"name"`
	Description string   `xml:"description"`
	Version     string   `xml:"version"`
	Date        string   `xml:"date"`
	Author      string   `xml:"author"`
	URL         string   `xml:"url"`
}

type Game struct {
	XMLName     xml.Name `xml:"game"`
	Name        string   `xml:"name,attr"`
	Description string   `xml:"description"`
	Roms        []Rom    `xml:"rom"`
}

type Release struct {
	XMLName xml.Name `xml:"release"`
	Name    string   `xml:"name,attr"`
	Region  string   `xml:"region,attr"`
}

type Rom struct {
	XMLName xml.Name `xml:"rom"`
	Name    string   `xml:"name,attr"`
	Size    int      `xml:"size,attr"`
	CRC     string   `xml:"crc,attr"`
	MD5     string   `xml:"md5,attr"`
	SHA1    string   `xml:"sha1,attr"`
}
