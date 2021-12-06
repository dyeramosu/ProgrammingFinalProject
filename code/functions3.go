package main
import (
  "fmt"
  "math/rand"
  "os"
  "log"
)

/*
SIMPLE
-add substrate molecules randomly
-add enzyme molecules randomly
-add inhibitor molecules randomly
-allow molecules to move randomly throughout space

-if substrate/inhibitor and enzyme come into contact
  -create complex

-ES complex should dissociate into enzyme and product

MICHAELIS-MENTEN
-E + S --> ES with probability p1         p1(n) = ∆t*k1*n       n = #S
-ES --> E + S with probability p2         p2 = ∆t*k2
-ES --> E + P with probability p3         p3 = ∆t*k3
-E + I --> EI with probability pI         pI = ∆t*k1*n1         n1 = #I

my starting conditions: ∆t = 0.01

-p1 = 0.01*k1*numParticles
-p2 = 0.01*k2
-p3 = 0.01*k3
*/

/*
COLOR KEY:
enzyme = 2 (red)
substrates = 4 (white)
products = 5 (yellow)
ES complex = 7 (pink)
inhibitor = 3 (blue)
IE complex = 6 (green)
*/

/*
DIRECTION KEY:
northwest = board[r-1][c-1]       0
north = board[r-1][c]             1
northeast = board[r-1][c+1]       2
east = board[r][c+1]              3
southeast = board[r+1][c+1]       4
south = board[r+1][c]             5
southwest = board[r+1][c-1]       6
west = board[r][c-1]              7
*/

//this is the main function - calls UpdateBoard as well as updates final map
func Diffuse(initialBoard Board, numSteps int, numParticles, numEnz1, numInhib int, kArray []float64) []Board {
  boards := make([]Board, numSteps+1)
  boards[0] = initialBoard

  for i := 1; i<=numSteps; i++ {
    boards[i] = UpdateBoard(boards[i-1], kArray, numParticles, numInhib)
  }

  myMap := MakeMap(numSteps)

  for j:= 0; j<= numSteps; j++ {
  	for r:= 0; r<CountRows(initialBoard); r++ {
  		for c:= 0; c<CountCols(initialBoard); c++ {
  			if boards[j][r][c] == 5 {
  				myMap["ProductCount"][j]++
  				myMap["TotalCount"][j]++
  			} else if boards[j][r][c] == 4 {
  				myMap["SubstrateCount"][j]++
  				myMap["TotalCount"][j]++
  			} else if boards[j][r][c] == 2 {
  				myMap["EnzymeCount"][j]++
  				myMap["TotalCount"][j]++
  			} else if boards[j][r][c] == 7 {
  				myMap["ComplexCount"][j]++
  				myMap["TotalCount"][j]++
  			} else if boards[j][r][c] == 3 {
          myMap["InhibCount"][j]++
        } else if boards[j][r][c] == 6 {
          myMap["InhibComplexCount"][j]++
        }
  		}
  	}
  }
  WriteMaptoFile(myMap)

  return boards
}

//InitializeBoard initializes a board struct with numParticles substrates and numEnz enzymes
  //all molecules are placed randomly throughout the board to start
func InitializeBoard(size int, numParticles, numEnz, numInhib int) Board {
  var b Board
  b = make(Board, size)
  for r:= range b {
    b[r] = make([]int, size)
  }

  for s:= 0; s< size; s++ {
    for j := 0; j<size; j++ {
      b[s][j] = 0
    }
  }

  //add substrate molecules
  for i:= 0; i<numParticles; i++ {
    var c OrderedPair
    c.x = rand.Intn(size)
    c.y = rand.Intn(size)

    //if this spot is already taken...
    if b[c.x][c.y] != 0 {
      for b[c.x][c.y] != 0 {
        c.x = rand.Intn(size)
        c.y = rand.Intn(size)
      }
    }
    b[c.x][c.y] = 4
  }

  //add enzyme molecules
  for j:= 0; j<numEnz; j++ {
    var d OrderedPair
    d.x = rand.Intn(size)
    d.y = rand.Intn(size)

    //if this spot is already taken...
    if b[d.x][d.y] != 0 {
      for b[d.x][d.y] != 0 {
        d.x = rand.Intn(size)
        d.y = rand.Intn(size)
      }
    }
    b[d.x][d.y] = 2
  }

  for k:= 0; k<numInhib; k++ {
    var e OrderedPair
    e.x = rand.Intn(size)
    e.y = rand.Intn(size)

    //if this spot is already taken...
    if b[e.x][e.y] != 0 {
      for b[e.x][e.y] != 0 {
        e.x = rand.Intn(size)
        e.y = rand.Intn(size)
      }
    }
    b[e.x][e.y] = 3
  }

  return b
}

//UpdateBoard takes board as input and outputs a new board after calling RandomWalk on the input board
func UpdateBoard(b Board, kArray []float64, numParticles, numInhib int) Board {
  //create new board
	newBoard := CopyBoard(b)
  bsize := CountRows(b)

  //iterate through all cells in board and move molecules randomly
  for i:= 0; i<bsize; i++ {
    for j:= 0; j<bsize; j++ {
      if newBoard[i][j] != 0 {
        RandomWalk(newBoard, i, j, kArray, numParticles, numInhib)
      }
    }
  }

  //if you come across a complex, dissociate into P + E or S + E at probability rate
  p2 := kArray[1] * 0.01 //probability of ES --> E + S
  p3 := kArray[2] * 0.01 //probability of ES --> E + P
  //0.01 is change in time

  if kArray[1]>9 || kArray[2]>9  {
    panic("k2 and/or k3 are greater than 9, choose smaller k2 and/or k3 values")
  }

  for i:= 0; i<bsize; i++ {
    for j:= 0; j<bsize; j++ {
      if newBoard[i][j] == 7 {
        r_int:= rand.Intn(100)
        newp3 := int(p3*100)
        for z:= 0; z<newp3; z++ {
          if r_int == z+1 {
            CreateProduct(newBoard, i, j)
          }
        }
      }
    }
  }

  for i:= 0; i<bsize; i++ {
    for j:= 0; j<bsize; j++ {
      if newBoard[i][j] == 7 {
        r_int := rand.Intn(100)
        newp2 := int(p2*100)
        for z:= 0; z<newp2; z++ {
          if r_int == z+1 {
            CreateSubstrate(newBoard, i, j)
          }
        }
      }
    }
  }

  return newBoard
}

//RandomWalk takes the board and position as input and, if at a molecule, randomly moves it
//if the current molecule is a substrate, it checks to see if we are near an enzyme
  //if so, calls CreateComplex
//if the current molecule is an enzyme, it checks to see if we are near a substrate
  //if so, calls CreateComplex
func RandomWalk(b Board, r, c int, kArray []float64, numParticles, numInhib int) {
  diffRate := 1 //number of steps each molecule takes per board update

  for k:= 0; k<diffRate; k++ {
    //pick random direction from 0 to 7
    dir := rand.Intn(8)

    //if direction is not valid, pick random directions until find one that is valid
    if !isValid(r, c, dir, b) {
      for !isValid(r, c, dir, b) {
        dir = rand.Intn(8)
      }
    }
    //now we have a valid direction

    //if we are at a substrate molecule
    if b[r][c] == 4 {
      //set current position to be 0
      b[r][c] = 0
      //if substrate molecule is next to an enzyme
      if AtEnzyme(b, r, c) {
        //create SE complex with probability rate p1
        p1 := kArray[0]*0.01*float64(numParticles)
        if p1 >= 1 {
          CreateComplex(b, r, c, 1)
        }
        return
      } else { //not at enzyme
        //move molecule in random direction
        MoveMolecule(b, r, c, dir, 4)
      }

    //else, if we are an enzyme molecule
    } else if b[r][c] == 2 {
      //set current position to be 0
      b[r][c] = 0
      //if enzyme molecule is next to a substrate
      if AtSubstrate(b, r, c) {
        //create SE complex with probability rate p1
        p1:= kArray[0]*0.01*float64(numParticles)
        if p1 >= 1 {
          CreateComplex(b, r, c, 1)
        }
        return
      } else { //not at substrate
        //move molecule in random direction
        MoveMolecule(b, r, c, dir, 2)
      }

    //else if we are at a product molecule
    } else if b[r][c] == 5 {
      //set current position to be 0
      b[r][c] = 0
      //move molecule in random direction
      MoveMolecule(b, r, c, dir, 5)

    //else if we are at a complex molecule
    } else if b[r][c] == 7 {
      //set current position to be 0
      b[r][c] = 0
      //move molecule in random direction
      MoveMolecule(b, r, c, dir, 7)

    //else if we are at an inhibitor molecule
    } else if b[r][c] == 3 {
      //set current position to be 0
      b[r][c] = 0
      //if inbibitor molecule is next to an enzyme
      if AtEnzyme(b, r, c) {
        //create IE complex with probability rate pI
        pI:= kArray[3]*0.01*float64(numInhib)
        if pI >= 1 {
          CreateBadComplex(b, r, c, 1)
        }
        return
      } else { //not at enzyme
        //move molecule in random direction
        MoveMolecule(b, r, c, dir, 3)
      }

    //else if we are at an IE complex
    } else if b[r][c] == 6 {
      //set current position to be 0
      b[r][c] = 0
      //move molecule in random direction
      MoveMolecule(b, r, c, dir, 6)
    }
  }
}

//MoveMolecule moves the current position of the molecule to a random position in its neighborhood
  //and then sets the current position to be the new position
func MoveMolecule(b Board, r, c, dir, num int) {
  if dir == 0 {
    b[r-1][c-1] = num
    r = r-1
    c = c-1
  } else if dir == 1 {
    b[r-1][c] = num
    r = r-1
  } else if dir == 2 {
    b[r-1][c+1] = num
    r = r-1
    c = c+1
  } else if dir == 3 {
    b[r][c+1] = num
    c = c+1
  } else if dir == 4 {
    b[r+1][c+1] = num
    r = r+1
    c = c+1
  } else if dir == 5 {
    b[r+1][c] = num
    r = r+1
  } else if dir == 6 {
    b[r+1][c-1] = num
    r = r+1
    c = c-1
  } else if dir == 7 {
    b[r][c-1] = num
    c = c-1
  } else {
    panic("invalid direction given to CreateProduct")
  }
}

//AtEnzyme returns true if an enzyme is in the neighborhood of the substrate molecule
func AtEnzyme(board Board, r, c int) bool {
  //if direction 0 is valid and enzyme is at direction 0
  if isValid1(r, c, 0, board) && board[r-1][c-1] == 2 {
    //set enzyme position to be 0 and return true
    board[r-1][c-1] = 0
    return true

  //if direction 1 is valid and enzyme is at direction 1
  } else if isValid1(r, c, 1, board) && board[r-1][c] == 2 {
    board[r-1][c] = 0
    return true

  //if direction 2 is valid and enzyme is at direction 2
  } else if isValid1(r, c, 2, board) && board[r-1][c+1] == 2 {
    board[r-1][c+1] = 0
    return true

  //if direction 3 is valid and enzyme is at direction 3
  } else if isValid1(r, c, 3, board) && board[r][c+1] == 2 {
    board[r][c+1] = 0
    return true

  //if direction 4 is valid and enzyme is at direction 4
  } else if isValid1(r, c, 4, board) && board[r+1][c+1] == 2 {
    board[r+1][c+1] = 0
    return true

  //if direction 5 is valid and enzyme is at direction 5
  } else if isValid1(r, c, 5, board) && board[r+1][c] == 2 {
    board[r+1][c] = 0
    return true

  //if direction 6 is valid and enzyme is at direction 6
  } else if isValid1(r, c, 6, board) && board[r+1][c-1] == 2 {
      board[r+1][c-1] = 0
      return true

  //if direction 7 is valid and enzyme is at direction 7
  } else if isValid1(r, c, 7, board) && board[r][c-1] == 2 {
    board[r][c-1] = 0
    return true
  }
  return false
}

//AtSubstrate returns true if a substrate is in the neighborhood of the enzyme molecule
func AtSubstrate(board Board, r, c int) bool {
  //if direction 0 is valid and substrate is at direction 0
  if isValid1(r, c, 0, board) && board[r-1][c-1] == 4 {
    //set substrate position to be 0 and return true
    board[r-1][c-1] = 0
    return true

  //if direction 1 is valid and substrate is at direction 1
  } else if isValid1(r, c, 1, board) && board[r-1][c] == 4 {
    board[r-1][c] = 0
    return true

  //if direction 2 is valid and substrate is at direction 2
  } else if isValid1(r, c, 2, board) && board[r-1][c+1] == 4 {
    board[r-1][c+1] = 0
    return true

  //if direction 3 is valid and substrate is at direction 3
  } else if isValid1(r, c, 3, board) && board[r][c+1] == 4 {
    board[r][c+1] = 0
    return true

  //if direction 4 is valid and substrate is at direction 4
  } else if isValid1(r, c, 4, board) && board[r+1][c+1] == 4 {
    board[r+1][c+1] = 0
    return true

  //if direction 5 is valid and substrate is at direction 5
  } else if isValid1(r, c, 5, board) && board[r+1][c] == 4 {
    board[r+1][c] = 0
    return true

  //if direction 6 is valid and substrate is at direction 6
  } else if isValid1(r, c, 6, board) && board[r+1][c-1] == 4 {
      board[r+1][c-1] = 0
      return true

  //if direction 7 is valid and substrate is at direction 7
  } else if isValid1(r, c, 7, board) && board[r][c-1] == 4 {
    board[r][c-1] = 0
    return true
  }
  return false
}

//CreateComplex creates SE complex when S and E are in eachothers neighborhoods
func CreateComplex(b Board, r, c, enz int) {
  //pick a random direction
  dir := rand.Intn(8)

  //if direction is not valid, pick random directions until find one that is
  if !isValid(r, c, dir, b) {
    for !isValid(r, c, dir, b) {
      dir = rand.Intn(8)
    }
  }

  //if direction is ___
    //set cell of that direction to be SE complex
  if dir == 0 {
    b[r-1][c-1] = 7
  } else if dir == 1 {
    b[r-1][c] = 7
  } else if dir == 2 {
    b[r-1][c+1] = 7
  } else if dir == 3 {
    b[r][c+1] = 7
  } else if dir == 4 {
    b[r+1][c+1] = 7
  } else if dir == 5 {
    b[r+1][c] = 7
  } else if dir == 6 {
    b[r+1][c-1] = 7
  } else if dir == 7 {
    b[r][c-1] = 7
  } else {
    panic("invalid direction given to CreateComplex")
  }
}

//CreateBadComplex creates IE complex when I and E are in eachothers neighborhoods
func CreateBadComplex(b Board, r, c, enz int) {
  //pick a random direction
  dir := rand.Intn(8)

  //if direction is not valid, pick random directions until find one that is
  if !isValid(r, c, dir, b) {
    for !isValid(r, c, dir, b) {
      dir = rand.Intn(8)
    }
  }

  //if direction is ___
    //set cell of that direction to be IE complex
  if dir == 0 {
    b[r-1][c-1] = 6
  } else if dir == 1 {
    b[r-1][c] = 6
  } else if dir == 2 {
    b[r-1][c+1] = 6
  } else if dir == 3 {
    b[r][c+1] = 6
  } else if dir == 4 {
    b[r+1][c+1] = 6
  } else if dir == 5 {
    b[r+1][c] = 6
  } else if dir == 6 {
    b[r+1][c-1] = 6
  } else if dir == 7 {
    b[r][c-1] = 6
  } else {
    panic("invalid direction given to CreateComplex")
  }
}

//CreateProduct creates E and P molecules by dissociation of SE complex
func CreateProduct(b Board, r, c int) {
  //pick a random direction
  dir := rand.Intn(8)

  //if direction is not valid, pick random directions until find one that is
  if !isValid(r, c, dir, b) {
    for !isValid(r, c, dir, b) {
      dir = rand.Intn(8)
    }
  }

  //if direction is ____
    //set current cell to be enzyme
    //set direction cell to be product
  if dir == 0 {
    b[r][c] = 2
    b[r-1][c-1] = 5

  } else if dir == 1 {
    b[r][c] = 2
    b[r-1][c] = 5

  } else if dir == 2 {
    b[r][c] = 2
    b[r-1][c+1] = 5

  } else if dir == 3 {
    b[r][c] = 2
    b[r][c+1] = 5

  } else if dir == 4 {
    b[r][c] = 2
    b[r+1][c+1] = 5

  } else if dir == 5 {
    b[r][c] = 2
    b[r+1][c] = 5

  } else if dir == 6 {
    b[r][c] = 2
    b[r+1][c-1] = 5

  } else if dir == 7 {
    b[r][c] = 2
    b[r][c-1] = 5

  } else {
    panic("invalid direction given to CreateProduct")
  }
}

//isValid returns true if direction cell is within board AND if direction cell is empty
func isValid(r, c, direction int, board Board) bool {
  if direction > 7 || direction < 0 {
    panic("invalid direction given")
  }

  //if row == 0 and if direction is north, northeast or northwest
  if r == 0 && (direction == 0 || direction == 1 || direction == 2) {
      return false
    //if c == 0 and if direction is northwest, west or southwest
  } else if c == 0 && (direction == 0 || direction == 6 || direction == 7){
      return false
    //if r is the last row and if direction is southeast, south or southwest
  } else if r == (CountRows(board)-1) && (direction == 4 || direction == 5 || direction == 6){
      return false
    //if c is the last column and if direction is northeast, east or southeast
  } else if c == (CountCols(board)-1) && (direction == 2 || direction == 3 || direction == 4) {
      return false
  } else { //direction is valid
    if direction == 0 {
      if board[r-1][c-1] != 0 {
        return false
      }
    } else if direction == 1 {
      if board[r-1][c] != 0 {
        return false
      }
    } else if direction == 2 {
      if board[r-1][c+1] != 0 {
        return false
      }
    } else if direction == 3 {
      if board[r][c+1] != 0 {
        return false
      }
    } else if direction == 4 {
      if board[r+1][c+1] != 0 {
        return false
      }
    } else if direction == 5 {
      if board[r+1][c] != 0 {
        return false
      }
    } else if direction == 6 {
      if board[r+1][c-1] != 0 {
        return false
      }
    } else if direction == 7 {
      if board[r][c-1] != 0 {
        return false
      }
    }
  }
  return true
}

//isValid1 returns true if direction cell is within board
func isValid1(r, c, direction int, board Board) bool {
  if direction > 7 || direction < 0 {
    panic("invalid direction given")
  }

  //if row == 0 and if direction is north, northeast or northwest
  if r == 0 && (direction == 0 || direction == 1 || direction == 2) {
      return false
    //if c == 0 and if direction is northwest, west or southwest
  } else if c == 0 && (direction == 0 || direction == 6 || direction == 7){
      return false
    //if r is the last row and if direction is southeast, south or southwest
  } else if r == (CountRows(board)-1) && (direction == 4 || direction == 5 || direction == 6){
      return false
    //if c is the last column and if direction is northeast, east or southeast
  } else if c == (CountCols(board)-1) && (direction == 2 || direction == 3 || direction == 4) {
      return false
  }
  return true
}

//CopyBoard creates a new Board struct and copies all cells from input board to new board
func CopyBoard(b Board) Board {
	var newBoard Board
  bsize := CountRows(b)

  newBoard = make(Board, bsize)
  for r:= range newBoard {
    newBoard[r] = make([]int, bsize)
  }

  for i:= 0; i<bsize; i++ {
    for j:= 0; j<bsize; j++ {
      newBoard[i][j] = b[i][j]
    }
  }

	return newBoard
}

//CountRows returns number of rows in the input board
func CountRows(board Board) int {
  if len(board) == 0 {
    panic("invalid number of rows given")
  }
	return len(board)
}

//CountCols returns number of cols in the input board
func CountCols(board Board) int {
	// assume that we have a rectangular board
	if CountRows(board) == 0 {
		panic("Error: empty board given to CountCols")
	}
	// give # of elements in 0-th row
	return len(board[0])
}

//dissociation of SE back to S + E
func CreateSubstrate(b Board, r, c int) {
  //pick a random direction
  dir := rand.Intn(8)

  //if direction is not valid, pick random directions until find one that is
  if !isValid(r, c, dir, b) {
    for !isValid(r, c, dir, b) {
      dir = rand.Intn(8)
    }
  }

  //if direction is ____
    //set current cell to be enzyme
    //set direction cell to be substrate
  if dir == 0 {
    b[r][c] = 2
    b[r-1][c-1] = 4

  } else if dir == 1 {
    b[r][c] = 2
    b[r-1][c] = 4

  } else if dir == 2 {
    b[r][c] = 2
    b[r-1][c+1] = 4

  } else if dir == 3 {
    b[r][c] = 2
    b[r][c+1] = 4

  } else if dir == 4 {
    b[r][c] = 2
    b[r+1][c+1] = 4

  } else if dir == 5 {
    b[r][c] = 2
    b[r+1][c] = 4

  } else if dir == 6 {
    b[r][c] = 2
    b[r+1][c-1] = 4

  } else if dir == 7 {
    b[r][c] = 2
    b[r][c-1] = 4

  } else {
    panic("invalid direction given to CreateSubstrate")
  }
}

//MakeMap creates a map of strings to int arrays
func MakeMap(numSteps int) map[string][]int {
  var newMap = make(map[string][]int)

  newMap["SubstrateCount"] = make([]int, numSteps+1)
  newMap["EnzymeCount"] = make([]int, numSteps+1)
  newMap["ComplexCount"] = make([]int, numSteps+1)
  newMap["ProductCount"] = make([]int, numSteps+1)
  newMap["TotalCount"] = make([]int, numSteps+1)
  newMap["InhibCount"] = make([]int, numSteps+1)
  newMap["InhibComplexCount"] = make([]int, numSteps+1)

  return newMap
}

//WriteMaptoFile takes final output map and creates file to be used in RStudio
func WriteMaptoFile(myMap map[string][]int) {
	outfile, err := os.Create("finalMap")
	if err != nil {
		log.Fatal("Cannot create file", err)
	}
	defer outfile.Close()

	for key, val := range myMap {
    fmt.Fprint(outfile, key, "\t")
    for i:= 0; i<len(myMap[key]); i++ {
      fmt.Fprint(outfile, val[i], "\t")   //, "\n")
    }
    fmt.Fprint(outfile, "\n")
	}
}
