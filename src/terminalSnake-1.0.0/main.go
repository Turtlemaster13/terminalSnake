package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/nsf/termbox-go"
)

func main() {
	snakeDataPath := "/usr/share/terminalSnake.csv"
	var mapSize int
	var snakeSpeed int
	var easyMode string
	//arguments
	flag.IntVar(&mapSize, "s", 15, "the map size")
	flag.IntVar(&snakeSpeed, "v", 8, "snake velocity")
	flag.StringVar(&easyMode, "e", "n", "easy mode (y/n)")
	flag.Parse()

	//setUp
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()

	eventQueue := make(chan termbox.Event)
	go func() {
		for {
			eventQueue <- termbox.PollEvent()
		}
	}()

	termbox.SetInputMode(termbox.InputEsc)
	//termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	displayTopScores(snakeDataPath, mapSize)
	//termbox.Flush()
	gameMap := genMap(mapSize)
	score := 0
	drawMap(gameMap)
	showScore(score, mapSize)

	snake := snake{1, 0, []location{}, location{int(mapSize/2 - 1), int(mapSize/2 - 1)}}
	drawSnake(snake)
	apple := placeApple(mapSize)
	termbox.Flush()

	//framerate
	snakeSpeedTmp := 0

	directionChanged := false
	frameRate := 120
	frameTime := time.Second / time.Duration(frameRate)
	startTime := time.Now()

	fEasyMode := "normal mode"
	if easyMode == "y" {
		fEasyMode = "easy mode"
	}
	defer saveScoresToFile(snakeDataPath, &score, mapSize, snakeSpeed, fEasyMode)
	//game loop
	for {
		elapsedTime := time.Since(startTime)
		select {
		case ev := <-eventQueue:
			switch ev.Type {
			case termbox.EventKey:
				if ev.Key == termbox.KeyEsc || ev.Ch == 'q' {
					return
				} else if ev.Key == termbox.KeyArrowRight && snake.direction != 3 && !directionChanged {
					snake.direction = 1
					directionChanged = true
				} else if ev.Key == termbox.KeyArrowDown && snake.direction != 0 && !directionChanged {
					snake.direction = 2
					directionChanged = true
				} else if ev.Key == termbox.KeyArrowLeft && snake.direction != 1 && !directionChanged {
					snake.direction = 3
					directionChanged = true
				} else if ev.Key == termbox.KeyArrowUp && snake.direction != 2 && !directionChanged {
					snake.direction = 0
					directionChanged = true
				}
			}
		default:
		}
		if elapsedTime >= frameTime { //keep frame rate the same
			startTime = time.Now()
			drawMap(gameMap)
			updateApple(apple)
			drawSnake(snake)
			showScore(score, mapSize)
			//displayTopScores("snakeGame.csv", mapSize)
			//actully do stuff with the snake
			gameOver := false
			if snakeSpeedTmp == snakeSpeed {
				snake, gameOver, apple, score = moveSnake(snake, mapSize, apple, easyMode, score)
				directionChanged = false
				snakeSpeedTmp = 0
			} else {
				snakeSpeedTmp++
			}
			termbox.Flush()
			if gameOver {
				return
			}

		} else {
			sleepTime := frameTime - elapsedTime
			time.Sleep(sleepTime)
		}
	}
}

func drawMap(gameMap [][]int) {
	for y, mapRow := range gameMap {
		for x := range mapRow {
			termbox.SetBg(x*2+2, y+1, termbox.ColorBlue)
			termbox.SetBg(x*2+1, y+1, termbox.ColorBlue)
		}
	}
}

func drawSnake(snake snake) {
	termbox.SetBg(snake.headPosition.x*2, snake.headPosition.y, termbox.ColorGreen)
	termbox.SetBg(snake.headPosition.x*2-1, snake.headPosition.y, termbox.ColorGreen)
	for _, loc := range snake.tail {
		termbox.SetBg(loc.x*2, loc.y, termbox.ColorGreen)
		termbox.SetBg(loc.x*2-1, loc.y, termbox.ColorGreen)
	}
}

func moveSnake(snake snake, mapSize int, apple location, easyMode string, score int) (snake, bool, location, int) {
	gameOver := false
	lastPosition := snake.headPosition
	switch snake.direction {
	case 0:
		snake.headPosition.y--
	case 1:
		snake.headPosition.x++
	case 2:
		snake.headPosition.y++
	case 3:
		snake.headPosition.x--
	}
	snake.tail = append([]location{lastPosition}, snake.tail...)
	if len(snake.tail) > snake.length {
		snake.tail = snake.tail[:len(snake.tail)-1]
	}
	cell := termbox.GetCell((snake.headPosition.x)*2, snake.headPosition.y)

	if cell.Bg == termbox.ColorRed {
		snake.length++
		apple = placeApple(mapSize)
		score++
	} else if cell.Bg == termbox.ColorDefault && easyMode == "y" {
		if snake.direction == 0 {
			snake.headPosition.y = mapSize
		} else if snake.direction == 1 {
			snake.headPosition.x = 1
		} else if snake.direction == 2 {
			snake.headPosition.y = 1
		} else if snake.direction == 3 {
			snake.headPosition.x = mapSize
		}
	} else if cell.Bg == termbox.ColorGreen && easyMode == "y" {

	} else if cell.Bg != termbox.ColorBlue {
		gameOver = true
	}

	return snake, gameOver, apple, score
}

func genMap(deminsions int) [][]int {
	gameMap := make([][]int, deminsions)
	for i := 0; i < deminsions; i++ {
		gameMap[i] = make([]int, deminsions)
	}
	return gameMap
}

func placeApple(mapSize int) location {
	placed := false
	var appleX, appleY int
	for !placed {
		appleX = rand.Intn(mapSize)
		appleY = rand.Intn(mapSize)
		if termbox.GetCell(appleX+2, appleY+1).Bg == termbox.ColorBlue {
			termbox.SetBg(appleX*2+1, appleY+1, termbox.ColorRed)
			termbox.SetBg(appleX*2+2, appleY+1, termbox.ColorRed)
			placed = true
		}
	}
	return location{appleX, appleY}

}
func updateApple(apple location) {
	appleX, appleY := apple.x, apple.y
	termbox.SetBg(appleX*2+1, appleY+1, termbox.ColorRed)
	termbox.SetBg(appleX*2+2, appleY+1, termbox.ColorRed)
}
func showScore(score int, mapSize int) {
	maxScore := (mapSize * mapSize) - 1

	scoreStr := fmt.Sprintf("current score: %d", score)
	maxScoreStr := fmt.Sprintf("max Score: %d", maxScore)
	scoreX := mapSize*2 + 3
	for _, char := range scoreStr {
		termbox.SetCell(scoreX, 5, char, termbox.ColorDefault, termbox.ColorDefault)
		scoreX++
	}
	scoreX = mapSize*2 + 3
	for _, char := range maxScoreStr {
		termbox.SetCell(scoreX, 6, char, termbox.ColorDefault, termbox.ColorDefault)
		scoreX++
	}
}
func saveScoresToFile(saveFile string, score *int, mapSize int, speed int, fEasyMode string) {
	file, err := os.OpenFile(saveFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	newRow := []string{fmt.Sprintf("%d", *score), fmt.Sprintf("%d", mapSize), fmt.Sprintf("%d", speed), fEasyMode}
	err = writer.Write(newRow)
	if err != nil {
		panic(err)
	}
	writer.Flush()
}

func displayTopScores(saveFile string, mapSize int) {
	file, err := os.Open(saveFile)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	scores, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	if len(scores) > 1 {
		scores = scores[1:]
	} else {
		scores = [][]string{}
	}
	sort.Slice(scores, func(i, j int) bool {
		num1, _ := strconv.Atoi(scores[i][0])
		num2, _ := strconv.Atoi(scores[j][0])
		return num1 < num2
	})

	scoreStr := "-----High-Scores-----"
	scoreX := mapSize*2 + 3
	scoreY := 9
	for _, char := range scoreStr {
		termbox.SetCell(scoreX, scoreY, char, termbox.ColorDefault, termbox.ColorDefault)
		scoreX++
	}
	scoreY++
	for _, maxScores := range reverse(scores[len(scores)-min(len(scores), 6):]) {
		scoreX = mapSize*2 + 3
		scoreStr = fmt.Sprintf("score: %s, map size: %s, snake speed: %s, game mode: %s", maxScores[0], maxScores[1], maxScores[2], maxScores[3])

		for _, char := range scoreStr {
			termbox.SetCell(scoreX, scoreY, char, termbox.ColorDefault, termbox.ColorDefault)
			scoreX++
		}
		scoreY++
	}
}

func min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func reverse(oldArray [][]string) [][]string {
	newArray := [][]string{}

	for _, i := range oldArray {
		newArray = append([][]string{i}, newArray...)
	}
	return newArray
}

type snake struct {
	direction    int
	length       int
	tail         []location
	headPosition location
}

type location struct {
	x, y int
}
