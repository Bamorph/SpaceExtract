package main

import (
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/paulmach/orb"
	"github.com/paulmach/orb/planar"
)

// Global verbose flag
var verbose bool

// PointsToCm converts points to centimeters assuming a default DPI of 72
func PointsToCm(points float64, DPI float64) float64 {
	return points * (25.4 / DPI) / 10
}

// Space struct represents a space with a title and a list of coordinates
type Space struct {
	Title string
	Path  []orb.Point
}

// Area calculates the area of the space using the orb/planar library
func (s Space) Area() float64 {
	polygon := orb.Polygon{s.Path}
	return planar.Area(polygon) // Convert cm² to m²
}

// Perimeter calculates the perimeter of the space
func (s Space) Perimeter() float64 {
	polygon := orb.Polygon{s.Path}
	return planar.Length(polygon) // Convert cm to meters
}

// ParseCoordinates parses the coordinates from the string representation
func ParseCoordinates(pathStr string) []orb.Point {
	coords := []orb.Point{}

	// Remove brackets at the start and end
	pathStr = strings.TrimPrefix(pathStr, "[")
	pathStr = strings.TrimSuffix(pathStr, "]")

	// Split the string into pairs
	pairs := strings.Split(pathStr, "][")

	for _, pair := range pairs {
		points := strings.Fields(pair)
		if len(points) == 2 {
			x, errX := strconv.ParseFloat(points[0], 64)
			y, errY := strconv.ParseFloat(points[1], 64)
			if errX == nil && errY == nil {
				coords = append(coords, orb.Point{PointsToCm(x, 72), PointsToCm(y, 72)})
			} else {
				fmt.Printf("Failed to convert coordinates from: %v\n", pair)
			}
		}
	}

	// Close the polygon if not closed
	if len(coords) > 0 && coords[0] != coords[len(coords)-1] {
		coords = append(coords, coords[0])
	}

	return coords
}

// FindSpacesFromEndObj reads the content of a PDF and extracts the spaces
func FindSpacesFromEndObj(pdfPath string) []Space {
	content, err := ioutil.ReadFile(pdfPath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	// Convert to string and use regex to find matches
	contentStr := string(content)
	pattern := `<</Type/Space/Title\((.*?)\)/Path\[(.*?)\]/C\[.*?\]/CA .*?>>`
	re := regexp.MustCompile(pattern)
	matches := re.FindAllStringSubmatch(contentStr, -1)

	var spaces []Space
	for _, match := range matches {
		title := match[1]
		path := ParseCoordinates(match[2])
		spaces = append(spaces, Space{Title: title, Path: path})
	}

	return spaces
}

// ExportSpacesToCSV exports the space details to a CSV file
func ExportSpacesToCSV(spaces []Space, pdfPath string) {
	dir, pdfFilename := splitFilePath(pdfPath)
	baseName := strings.TrimSuffix(pdfFilename, ".pdf")
	csvPath := fmt.Sprintf("%s/%s.csv", dir, baseName)

	file, err := os.Create(csvPath)
	if err != nil {
		log.Fatalf("Failed to create CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Title", "Area (m2)", "Perimeter (m)"})

	// Write rows
	for _, space := range spaces {
		writer.Write([]string{
			space.Title,
			fmt.Sprintf("%.2f", space.Area()),
			fmt.Sprintf("%.2f", space.Perimeter()),
		})
	}
	if verbose {
		fmt.Printf("Data exported to %s\n", csvPath)
	}
}

// splitFilePath splits a full file path into directory and file name
func splitFilePath(fullPath string) (string, string) {
	pathParts := strings.Split(fullPath, string(os.PathSeparator))
	dir := strings.Join(pathParts[:len(pathParts)-1], string(os.PathSeparator))
	filename := pathParts[len(pathParts)-1]
	return dir, filename
}

func main() {
	// Check if the PDF file path is provided as a command-line argument
	if len(os.Args) < 2 {
		log.Fatal("Please provide a PDF file path as an argument.")
	}

	// Get the PDF file path from the command-line arguments
	pdfPath := os.Args[1]

	// Check for the verbose flag "-v"
	if len(os.Args) > 2 && os.Args[2] == "-v" {
		verbose = true
	}

	spaces := FindSpacesFromEndObj(pdfPath)
	if verbose {
		for _, space := range spaces {
			fmt.Printf("Title: %s\n", space.Title)
			fmt.Printf("Path: %v\n", space.Path)
			fmt.Printf("Area: %.2f square meters\n", space.Area())
			fmt.Printf("Perimeter: %.2f meters\n", space.Perimeter())
			fmt.Println()
		}
	}

	ExportSpacesToCSV(spaces, pdfPath)
}
