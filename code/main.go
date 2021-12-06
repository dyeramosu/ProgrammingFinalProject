/*
Deepika Yeramosu
Programming for Scientists
Final Project
*/

package main
import (
  "fmt"
  "gifhelper"
)

func main() {
  fmt.Println("Modeling Enzyme Kinetics with Cellular Automaton")

  numParticles:= 300 //number of substrate molecules
  numEnz := 100 //number of enzyme molecules
  numInhib := 50 //number of inhibitor molecules
  size:= 201

  k1 := 6.67 //k1
  k2 := 1.0 //k-1
  k3 := 1.0 //k2
  kI := 6.67 //kI
  kArray := []float64{k1, k2, k3, kI}

  initialBoard := InitializeBoard(size, numParticles, numEnz, numInhib)

  fmt.Println("Running simulation")

  //numSteps is number of boards that will be created - 1
  numSteps:= 1000

  boards:= Diffuse(initialBoard, numSteps, numParticles, numEnz, numInhib, kArray)

  fmt.Println("Simulation run, drawing images.")
	imglist := DrawBoards(boards, 10)

	//converting images to a GIF
	fmt.Println("Boards drawn to images! Now, convert to animated GIF.")
	gifhelper.ImagesToGIF(imglist, "Final")
	fmt.Println("Success! GIF produced.")
}
