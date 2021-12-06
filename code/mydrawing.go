package main

import (
	"image"
)

//DrawBoards takes in a slice of boards and returns a list of images
func DrawBoards(boards []Board, cellWidth int) []image.Image {
	numGenerations := len(boards) //the number of boards we input are the number of generations
  imageList := make([]image.Image, numGenerations)
	for i := range boards { //for each board inputted
		imageList[i] = DrawBoard(boards[i], cellWidth) //create an image
	}
	return imageList
}

//DrawBoard takes in a single board and outputs an image for that board
func DrawBoard(board Board, cellWidth int) image.Image {
	height := len(board) * cellWidth
	width := len(board[0]) * cellWidth
	c := CreateNewPalettedCanvas(width, height, nil)

	// declare colors
  black:= MakeColor(0, 0, 0)
	blue := MakeColor(0, 0, 255)
	red := MakeColor(255, 0, 0)
  white := MakeColor(255, 255, 255)
  gray := MakeColor(85, 85, 85)
	yellow := MakeColor(255, 255, 0)
	green := MakeColor(0, 255, 0)
	pink := MakeColor(255, 105, 180)


	// fill in colored squares
	for i :=0; i < CountRows(board); i++ {
		for j:= 0; j < CountCols(board); j++ {
			if board[i][j] == 3 {
				c.SetFillColor(blue)
			} else if board[i][j] == 2 {
				c.SetFillColor(red)
			} else if board[i][j] == 1 {
				c.SetFillColor(gray)
			} else if board[i][j] == 4 {
        c.SetFillColor(white)
      } else if board[i][j] == 0 {
        c.SetFillColor(black)
      } else if board[i][j] == 5 {
				c.SetFillColor(yellow)
			} else if board[i][j] == 6 {
				c.SetFillColor(green)
			} else if board[i][j] == 7 {
				c.SetFillColor(pink)
			} else {
        panic("error in values in DrawBoard ")
      }
			x := j * cellWidth
			y := i * cellWidth
			c.ClearRect(x, y, x+cellWidth, y+cellWidth) //tells us the four coordinates of the cell i want to color
			c.Fill() //colors the rectangle
		}
	}

	return GetImage(c)
}

func DrawGridLines(pic Canvas, cellWidth int) {
	w, h := pic.Width(), pic.Height()
	// first, draw vertical lines
	for i := 1; i < w/cellWidth; i++ {
		y := i * cellWidth
		pic.MoveTo(0.0, float64(y))
		pic.LineTo(float64(w), float64(y))
	}
	// next, draw horizontal lines
	for j := 1; j < h/cellWidth; j++ {
		x := j * cellWidth
		pic.MoveTo(float64(x), 0.0)
		pic.LineTo(float64(x), float64(h))
	}
	pic.Stroke()
}
