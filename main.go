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
		fmt.Printf("Usage: %s [KML file]", os.Args[0])
		os.Exit(1)
	}

	b, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Printf("ERROR: Unable to open input file: %s\n", err)
		os.Exit(1)
	}

	var doc KML
	if err := xml.Unmarshal(b, &doc); err != nil {
		fmt.Printf("ERROR: Unable to parse KML file: %s\n", err)
		os.Exit(1)
	}

	coords, err := doc.Document.Placemark.LineString.CoordinateList()
	if err != nil {
		fmt.Printf("ERROR: Unable to parse coordinates in KML file: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Read: %v, number of coordinates: %d\n", doc.Document.Placemark.Name, len(coords))

	// Convert waypoints
	var wpts []*gpx.WptType
	for i, c := range coords {
		wpt := &gpx.WptType{
			Lat:  c.Lat,
			Lon:  c.Lon,
			Name: fmt.Sprintf("%d", i+1),
		}
		wpts = append(wpts, wpt)
	}

	// Create GPX structure
	g := &gpx.GPX{
		Version: "1.0",
		Creator: "Kml2Gpx - http://www.edgeworks.no",
		Trk: []*gpx.TrkType{{
			Name: doc.Document.Placemark.Name,
			TrkSeg: []*gpx.TrkSegType{{
				TrkPt: wpts,
			}},
		}},
	}

	// Write to output file
	outputFile := fmt.Sprintf("%s.%s", os.Args[1], "gpx")
	f, err := os.Create(outputFile)
	if err != nil {
		fmt.Printf("ERROR: Unable to create output file: %s\n", err)
		os.Exit(1)
	}

	fw := bufio.NewWriter(f)
	defer f.Close()
	if _, err := fw.WriteString(xml.Header); err != nil {
		fmt.Printf("ERROR: Unable to write output file: %s\n", err)
		os.Exit(1)
	}

	if err := g.WriteIndent(fw, "", " "); err != nil {
		fmt.Printf("ERROR: Unable to write output file: %s\n", err)
		os.Exit(1)
	}
	_ = fw.Flush()
	fmt.Printf("Wrote output file: %s\n", outputFile)
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

func (l LineString) CoordinateList() ([]Coordinate, error) {
	var coords []Coordinate
	for _, c := range strings.Split(l.Coordinates, " ") {
		s := strings.Split(c, ",")

		p, err := strconv.Atoi(s[2])
		if err != nil {
			return nil, fmt.Errorf("failed to convert number: %v", err)
		}
		lon, err := strconv.ParseFloat(s[0], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert coordinate: %v", err)
		}
		lat, err := strconv.ParseFloat(s[1], 64)
		if err != nil {
			return nil, fmt.Errorf("failed to convert coordinate: %v", err)
		}

		coord := Coordinate{
			Pos: p,
			Lon: lon,
			Lat: lat,
		}
		coords = append(coords, coord)
	}
	return coords, nil
}
