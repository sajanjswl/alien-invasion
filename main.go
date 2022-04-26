package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var (
	Map           string
	AlienCount    int
	ArrayOfAlines []*Alien
	ArrayOfCities []*CityNode
	CityCount     int
)

type Alien struct {
	Id          int
	Movecount   int
	IsDestroyed bool
}

type CityNode struct {
	Id             int
	Name           string
	North          *CityNode
	South          *CityNode
	East           *CityNode
	West           *CityNode
	VisitingAlines []*Alien
	MinTravelTime  time.Duration
	IsDestroyed    bool
}

func main() {

	flag.StringVar(&Map, "map", "./map.txt", "Map to traverse")
	flag.IntVar(&AlienCount, "alien-count", 42, "nuber of alien that arives on earth")
	flag.Parse()

	// fmt.Println("Map", Map)
	// fmt.Println("number", AlienCount)
	fileReader()

}

// reads input from map
func fileReader() {

	file, err := os.Open(Map)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {

		cities := strings.Split(scanner.Text(), " ")
		node := &CityNode{}

		for i, v := range cities {
			if i == 0 {
				CityCount = CityCount + 1
				node.Id = CityCount
				fmt.Println("Printing city count", CityCount)

				node.Name = v
				node.MinTravelTime = time.Duration((len([]rune(v))))
			}

			if strings.Contains(v, "north=") {
				CityCount = CityCount + 1
				nodeNorth := &CityNode{
					Id:   CityCount,
					Name: strings.TrimLeft(v, "north="),
				}
				ArrayOfCities = append(ArrayOfCities, nodeNorth)
				node.North = nodeNorth
				node.MinTravelTime = time.Duration((len([]rune(v))))
			}

			if strings.Contains(v, "south=") {
				CityCount = CityCount + 1
				nodeSouth := &CityNode{
					Id:   CityCount,
					Name: strings.TrimLeft(v, "south="),
				}
				ArrayOfCities = append(ArrayOfCities, nodeSouth)
				node.South = nodeSouth
				node.MinTravelTime = time.Duration((len([]rune(v))))
			}

			if strings.Contains(v, "east=") {
				CityCount = CityCount + 1
				nodeEast := &CityNode{
					Id:   CityCount,
					Name: strings.TrimLeft(v, "east="),
				}

				ArrayOfCities = append(ArrayOfCities, nodeEast)
				node.East = nodeEast
				node.MinTravelTime = time.Duration((len([]rune(v))))
			}

			if strings.Contains(v, "west=") {
				CityCount = CityCount + 1
				nodeWest := &CityNode{
					Id:   CityCount,
					Name: strings.TrimLeft(v, "west="),
				}
				ArrayOfCities = append(ArrayOfCities, nodeWest)
				node.West = nodeWest
				node.MinTravelTime = time.Duration((len([]rune(v))))
			}

		}

		// fmt.Println("Printing cities", city)

		// txtlines = append(txtlines, scanner.Text())
	}

	file.Close()

	for _, v := range ArrayOfCities {
		fmt.Printf("City %d Name %s time %s \n", v.Id, v.Name, v.MinTravelTime)
	}

	// for _, eachline := range txtlines {
	// 	fmt.Println(eachline)
	// }
}
