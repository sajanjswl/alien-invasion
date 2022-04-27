package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	Map           string
	AlienCount    int
	CityCount     int
	ArrayOfAlines []*Alien
	ArrayOfCities []*CityNode
)

// map to store cities distinctively
var CityHashMap = make(map[string]*CityNode)

type Alien struct {
	Id          int
	Movecount   int
	IsDestroyed bool
}

type CityNode struct {
	Id                 int
	Name               string
	North              *CityNode
	South              *CityNode
	East               *CityNode
	West               *CityNode
	VisitingAlines     []*Alien
	IsDestroyed        bool
	DestroyPrintSwitch bool
}

func main() {

	flag.StringVar(&Map, "map", "./map.txt", "Map to traverse")
	flag.IntVar(&AlienCount, "alien-count", 1, "nuber of alien that arives on earth")
	flag.Parse()

	fileReader()
	unleashAliens(AlienCount)
	printMap()

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
		// reads cities line by line and the split into words for processing
		cities := strings.Split(scanner.Text(), " ")
		node := &CityNode{}

		for i, name := range cities {
			if i == 0 {
				_, ok := CityHashMap[name]
				if !ok {
					addNodes(node, name)
				}
			}
			// checks if the next city is in north of the exisiting and
			// adds it to the north of the existing city
			if strings.Contains(name, "north=") {
				cityNode, ok := CityHashMap[strings.TrimLeft(name, "north=")]
				if !ok {
					node.North = &CityNode{}
					addNodes(node.North, strings.TrimLeft(name, "north="))
				} else {
					node.North = cityNode
				}
			}

			if strings.Contains(name, "south=") {
				cityNode, ok := CityHashMap[strings.TrimLeft(name, "south=")]
				if !ok {
					node.South = &CityNode{}
					addNodes(node.South, strings.TrimLeft(name, "south="))
				} else {
					node.South = cityNode
				}

			}

			if strings.Contains(name, "east=") {
				cityNode, ok := CityHashMap[strings.TrimLeft(name, "east=")]
				if !ok {
					node.East = &CityNode{}
					addNodes(node.East, strings.TrimLeft(name, "east="))
				} else {
					node.East = cityNode
				}
			}

			if strings.Contains(name, "west=") {
				cityNode, ok := CityHashMap[strings.TrimLeft(name, "west=")]
				if !ok {
					node.West = &CityNode{}
					addNodes(node.West, strings.TrimLeft(name, "west="))
				} else {
					node.West = cityNode
				}

			}

		}

	}

	file.Close()

}

// add the nodes to the city based on direction
func addNodes(node *CityNode, name string) {
	CityCount = CityCount + 1
	node.Id = CityCount
	node.Name = name
	ArrayOfCities = append(ArrayOfCities, node)
	CityHashMap[name] = node

}

// releases the alien to random cities initally
func unleashAliens(alienCount int) {
	var wg sync.WaitGroup
	for i := 0; i < alienCount; i++ {
		rand.Seed(time.Now().UnixNano())
		randomCity := rand.Intn((len(ArrayOfCities)-1)-0+1) + 0
		alien := &Alien{Id: i}

		var travellableCity []*CityNode
		travellableCity = append(travellableCity, ArrayOfCities[randomCity])
		wg.Add(1)
		go traverse(alien, travellableCity, &wg)

	}
	wg.Wait()
}

func traverse(alien *Alien, travellableCity []*CityNode, wg *sync.WaitGroup) {
	if len(travellableCity) == 0 {
		defer wg.Done()
		return
	}
	alien.Movecount = alien.Movecount + 1
	currentCity := travellableCity[0]
	currentCity.VisitingAlines = append(currentCity.VisitingAlines, alien)

	if alien.Movecount > 100000 {
		fmt.Printf("Printing alien that has travesed 10000 alienID %d  cound %d \n", alien.Id, alien.Movecount)
		defer wg.Done()
		return
	}
	if currentCity.East != nil && !currentCity.East.IsDestroyed {
		travellableCity = append(travellableCity, currentCity.East)
	}

	if currentCity.West != nil && !currentCity.West.IsDestroyed {
		travellableCity = append(travellableCity, currentCity.West)
	}
	if currentCity.North != nil && !currentCity.North.IsDestroyed {
		travellableCity = append(travellableCity, currentCity.North)
	}

	if currentCity.South != nil && !currentCity.South.IsDestroyed {
		travellableCity = append(travellableCity, currentCity.South)
	}

	if len(currentCity.VisitingAlines) > 1 && !currentCity.IsDestroyed {
		currentCity.IsDestroyed = true
		if !currentCity.DestroyPrintSwitch {
			fmt.Printf("%s has been destroyed by alien %d and alien %d! \n\n", currentCity.Name, currentCity.VisitingAlines[0].Id, currentCity.VisitingAlines[1].Id)
			currentCity.DestroyPrintSwitch = true
		}
		for _, al := range currentCity.VisitingAlines {
			al.IsDestroyed = true
		}
		defer wg.Done()
		return
	}
	if len(travellableCity) <= 1 || alien.IsDestroyed {
		defer wg.Done()
		return
	}
	travellableCity = travellableCity[1:]

	rand.Seed(time.Now().UnixNano())

	randomMove := rand.Intn((len(travellableCity)-1)-0+1) + 0

	travellableCity[0], travellableCity[randomMove] = travellableCity[randomMove], travellableCity[0]

	traverse(alien, travellableCity, wg)

}

func printMap() {

	f, err := os.Create("./output.txt")
	if err != nil {
		log.Fatalf("Failed to write to file")
	}

	for i, v := range ArrayOfCities {
		var city string

		if v.IsDestroyed {
			continue
		}
		if !v.IsDestroyed && i == 0 {
			city = city + v.Name
		}

		if v.East != nil && !v.East.IsDestroyed {
			city = city + " " + "east=" + v.East.Name
		}
		if v.West != nil && !v.West.IsDestroyed {
			city = city + " " + "west=" + v.West.Name
		}
		if v.North != nil && !v.North.IsDestroyed {
			city = city + " " + "north=" + v.North.Name
		}

		if v.South != nil && !v.South.IsDestroyed {
			city = city + " " + "south=" + v.South.Name
		}
		if len(city) != 0 {
			fmt.Println(city)
		}
		_, err := f.WriteString(city)
		if err != nil {
			log.Fatal("Falied to write to output file")
		}

	}
}
