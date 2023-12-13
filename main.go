package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"pgregory.net/rand"
	//"math/rand"
)

const samples = 9999999
const doorCount = 5
const openDoors = 3

var switchWin = 0
var noSwitchWin = 0
var wg sync.WaitGroup
var mu1 sync.Mutex
var mu2 sync.Mutex

func main() {
	var start = time.Now()

	fmt.Printf("Number of CPUs: %v\n", runtime.NumCPU())
	fmt.Printf("Number of Samples: %v\n", samples)
	fmt.Printf("Number of Doors: %v\n", doorCount)
	fmt.Printf("Number of Doors Opened: %v\n", openDoors)

	wg.Add(samples * 2)

	for i := 0; i < samples; i++ {
		montyHall(true)
		montyHall(false)
		//go montyHall(true)
		//go montyHall(false)
	}

	wg.Wait()
	fmt.Printf("\nSwitch Wins: %v\n", switchWin)
	fmt.Printf("Switch Win Percentage: %v %%\n", float64(switchWin)/float64(samples)*100)
	fmt.Printf("\nNo Switch Wins: %v\n", noSwitchWin)
	fmt.Printf("No Switch Win Percentage: %v %%\n", float64(noSwitchWin)/float64(samples)*100)

	var end = time.Now()
	fmt.Printf("\nFinished in %v s", end.Sub(start).Seconds())
}

func montyHall(sw bool) {
	defer wg.Done()

	var r = rand.New()
	//var r = rand.New(rand.NewSource(time.Now().UnixNano()))

	var opened [doorCount]bool
	var correct = r.Intn(doorCount)
	var chosen = 0

	for i := 0; i < openDoors; i++ {
		var o = 0
		for {
			o = r.Intn(doorCount)
			if o != correct && o != chosen && !opened[o] {
				break
			}
		}
		opened[o] = true
	}

	if sw {
		var prevChosen = chosen
		for {
			chosen = r.Intn(doorCount)
			if chosen != prevChosen && !opened[chosen] {
				break
			}
		}
	}

	if chosen == correct {
		if sw {
			mu1.Lock()
			switchWin++
			mu1.Unlock()
		} else {
			mu2.Lock()
			noSwitchWin++
			mu2.Unlock()
		}
	}
}
