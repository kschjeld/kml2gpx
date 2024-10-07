package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	gpx "github.com/twpayne/go-gpx"
	"os"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 2 {
		panic("expecting exactly one argument")
	}

	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}

	var doc KML

	if err := xml.Unmarshal(b, &doc); err != nil {
		panic(err)
	}

	coords := doc.Document.Placemark.LineString.CoordinateList()
	fmt.Printf("Read: %v\n# coordinates: %d\n", doc.Document.Placemark.Name, len(coords))
	fmt.Printf("First coord: %d %f %f\n", coords[0].Pos, coords[0].Lon, coords[0].Lat)

	g := &gpx.GPX{
		Version: "1.0",
		Creator: "Kml2Gpx - http://www.edgeworks.no",
	}

	track := &gpx.TrkType{
		Name:   doc.Document.Placemark.Name,
		TrkSeg: nil,
	}
	g.Trk = append(g.Trk, track)
	seg := &gpx.TrkSegType{
		TrkPt:      nil,
		Extensions: nil,
	}
	track.TrkSeg = append(track.TrkSeg, seg)

	for _, c := range coords {
		wpt := &gpx.WptType{
			Lat: c.Lat,
			Lon: c.Lon,
		}
		seg.TrkPt = append(seg.TrkPt, wpt)
	}

	f, err := os.Create(fmt.Sprintf("%s.%s", os.Args[1], "gpx"))
	if err != nil {
		panic(err)
	}
	fw := bufio.NewWriter(f)
	defer f.Close()

	if _, err := fw.WriteString(xml.Header); err != nil {
		panic(err)
	}

	if err := g.WriteIndent(fw, "", " "); err != nil {
		panic(err)
	}
	fw.Flush()

}

type KML struct {
	XMLName  xml.Name `xml:"http://earth.google.com/kml/2.2 kml"`
	Document Document `xml:"Document"`
}

type Document struct {
	Placemark Placemark `xml:"Placemark"`
}

type Placemark struct {
	Name       string     `xml:"name"`
	LineString LineString `xml:"LineString"`
}

type LineString struct {
	Coordinates string `xml:"coordinates"`
}

type Coordinate struct {
	Pos int
	Lat float64
	Lon float64
}

func (l LineString) CoordinateList() []Coordinate {
	var coords []Coordinate
	for _, c := range strings.Split(l.Coordinates, " ") {
		s := strings.Split(c, ",")

		p, err := strconv.Atoi(s[2])
		if err != nil {
			panic(err)
		}
		lon, err := strconv.ParseFloat(s[0], 64)
		if err != nil {
			panic(err)
		}
		lat, err := strconv.ParseFloat(s[1], 64)
		if err != nil {
			panic(err)
		}

		coord := Coordinate{
			Pos: p,
			Lon: lon,
			Lat: lat,
		}
		coords = append(coords, coord)
	}
	return coords
}
