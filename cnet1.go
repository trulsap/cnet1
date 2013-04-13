package main

import (
	"time"
	"math/rand"
	"fmt"
)


var nodes = 15
var minDegree = 3
var maxDegree = 5
var mainc = make(chan bool)
var allchans = make([]chan []int, nodes)
var outs = make([][]int, nodes)
var outc = make([][]chan []int , nodes)

func randBoard() {
	for n := 0; n < nodes; n++ {
		degree := minDegree + rand.Intn(maxDegree - minDegree + 1)
		outs[n] = make([]int, degree)
		for i := 0; i < degree; i++ {			
			x := rand.Intn(nodes)
			for {
				bad := x == n
				for j := 0; j < i; j++ {
					if outs[n][j] == x {
						bad = true
					}					
				}
				outs[n][i] = x
				if !bad { break }	
				x = rand.Intn(nodes)
			}
		}
	}
}

func initOutC() {
	for n := 0; n < nodes; n++ {
		outc[n] = make([]chan []int, len(outs[n]))
		for i := 0; i < len(outc[n]); i++ {
			outc[n][i] = allchans[outs[n][i]]
		}
	}
}

func nodeFunct(name int, myc chan []int, isee []chan []int) {
	for {
		x := <- myc
		fmt.Printf("%d just got %d\n", name, x)
		if x[0] == name {
			fmt.Println("I got it!")
			mainc <- true
		} else {
			seenBefore := false
			for i := 0; i < len(x); i++ {
				seenBefore = seenBefore || x[i] == name
			}
			if !seenBefore {
				conc := make([]int, 1 + len(x))
				for i := 0; i < len(x); i++ {
					conc[i] = x[i]
				}
				conc[len(x)] = name

				for m := 0; m < len(isee); m++ {
					isee[m] <- conc
				}
			} 
		}
	}
}

func printBoard() {
	for n := 0; n < len(outs); n++ {
		fmt.Printf("%d : %d", n, outs[n][0])
		for i := 1; i < len(outs[n]); i++ {
			fmt.Printf(", %d", outs[n][i])
		}
		fmt.Printf("\n")
	}
}

func main() {
	rand.Seed( time.Now().UTC().UnixNano())

	for i := 0; i < len(allchans); i++ {
		allchans[i] = make(chan []int)
	}

	randBoard()
	printBoard()
	initOutC()

	for n := 0; n < nodes ; n++ {
		go nodeFunct(n, allchans[n], outc[n])
	}

	msg := []int{nodes - 1}
	allchans[0] <- msg

	<- mainc
}
