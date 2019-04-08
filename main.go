package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// The Field type is simply a slice of field data. This comes
// from the engine on a turn-by turn basis, and should NEVER be
// updated by the ai.
type Field []rune

// X,Y position in relation to a Field
type Position struct {
	X int
	Y int
}

// GameState type keeps track of all the game variables that
// get set by the game engine.
type GameState struct {
	TimeBank    int      			// maximum time in ms that bot can have
	TimePerMove int      			// Time in ms that is added to timeBank each move
	PlayerNames []string 			// slice of player names, 0 index is player1
	MyBot       string   			// the name our bot has been assigned
	MyId        int      			// our bot id, which is also the index
	EnemyId     int      			// EnemyId, inverse of MyId
	FieldWidth  int      			// Width of the playing field
	FieldHeight int      			// Height of the playing field
	GameRound   int      			// Current round of the game
	GameField   map[string]Field    // Game field data
	ThisPieceType string 			// The current piece type
	ThisPiecePosition Position 		// The position of the current piece
	NextPieceType string			// The type of piece that will spawn next
	RowPoints   map[string]int      // Amount of points a player has
	Combo		map[string]int		// Height of the current combo for a player
	Skips		map[string]int		// The amount of skips a player has
}

// update settings from key value pair
func (f *GameState) UpdateSetting(key string, value string) {
	switch key {
	case "timebank":
		f.TimeBank, _ = strconv.Atoi(value)
	case "time_per_move":
		f.TimePerMove, _ = strconv.Atoi(value)
	case "player_names":
		f.PlayerNames = strings.Split(value, ",")
		f.GameField = make(map[string]Field, 2)
		f.RowPoints = make(map[string]int, 2)
		f.Combo = make(map[string]int, 2)
		f.Skips = make(map[string]int, 2)
	case "your_bot":
		f.MyBot = value
	case "your_botid":
		f.MyId, _ = strconv.Atoi(value)
		f.EnemyId = f.MyId ^ 1
	case "field_width":
		f.FieldWidth, _ = strconv.Atoi(value)
	case "field_height":
		f.FieldHeight, _ = strconv.Atoi(value)
	default:
		_, _ = io.LogError("Invalid key '%s'", key)
	}
}

func (f *GameState) UpdateGame(values[]string){
	if values[0]=="game"{
		switch values[1] {
		case "round":
			f.GameRound, _ = strconv.Atoi(values[2])
		case "this_piece_type":
			f.ThisPieceType = values[2]
		case "next_piece_type":
			f.NextPieceType = values[2]
		case "this_piece_position":
			tpp := strings.Split(values[2], ",")
			x,_ := strconv.Atoi(tpp[0])
			y, _ := strconv.Atoi(tpp[1])
			f.ThisPiecePosition = Position{x,y}
		}
	} else {
		switch values[1] {
		case "row_points":
			f.RowPoints[values[0]],_ = strconv.Atoi(values[2])
		case "combo":
			f.Combo[values[0]],_ = strconv.Atoi(values[2])
		case "skips":
			f.Skips[values[0]],_ = strconv.Atoi(values[2])
		case "field":
			newFieldString := strings.ReplaceAll(values[2], ",", "")
			f.GameField[values[0]] = Field(newFieldString)
		}
	}
}

func (f *GameState) EnemyName() string {
	return f.PlayerNames[f.MyId^1]
}

// a way to get the value a position in MY field.
func (f *GameState) MyFieldValue(x,y int) (rune, bool){
	if x < 0 || x >= f.FieldWidth {
		// x is outside of field bounds
		return 0, false
	}
	if y < 0 || y >= f.FieldHeight {
		// y is outside of field bounds
		return 0, false
	}
	index := (y * f.FieldWidth) + x
	return f.GameField[f.MyBot][index], true
}

// method of easily getting the value of a position in the ENEMY grid.
func (f *GameState) EnemyFieldValue(x,y int)(rune, bool){
	if x < 0 || x >= f.FieldWidth {
		// x is outside of field bounds
		return 0, false
	}
	if y < 0 || y >= f.FieldHeight {
		// y is outside of field bounds
		return 0, false
	}
	index := (y * f.FieldWidth) + x
	return f.GameField[f.EnemyName()][index], true
}

// BotAI is the interface an AI must implement
type BotAi interface {
	fmt.Stringer
	// Get move should return the a comma delimited list of moves that should be performed by the current piece.
	//
	// gs: The current game state. Nearly all the information that you need is in this object.
	// t: is the time (ms) allow to make the next move. Do not exceed this time.
	GetMove(gs *GameState, t int) string
}

var io = &BotIO // get io pointer
var file = flag.String("f","", "file to use as input for testing")

func main() {
	flag.Parse()
	if len(*file) > 0 {
		// redirect file to input
		f,err := os.Open(*file)
		if err == nil{
			io.in = bufio.NewReader(f)
		}else {
			_, _ = io.LogError("Unable to open file %s: %s", *file, err.Error())
		}
	}
	rand.Seed(time.Now().Unix())

	// todo - replace the following line with getting your bot
	bot := GetRandomBot()

	runBot(bot)
}

// The actual game loop. runBot reads input from the Riddles.io game engine,and
// requests
func runBot(ai BotAi){
	gs := &GameState{
		FieldWidth:  7,
		FieldHeight: 6,
		PlayerNames: []string{"", ""},
	}
	var text string
	var textParts []string
	var err error
	running := true
	for running {
		text, err = io.ReadLine()
		if err != nil {
			_, _ = io.LogError("Error getting input: %s", err)
			continue
		}
		if len(text) == 0 {
			continue
		}
		textParts = strings.Split(text, " ")
		switch textParts[0] {
		case "action":
			if textParts[1] == "move" {
				var turnTime, _ = strconv.Atoi(textParts[2])
				_, _ = io.WriteLine("%s", ai.GetMove(gs, turnTime))
			}
		case "update":
			gs.UpdateGame(textParts[1:])
		case "settings":
			gs.UpdateSetting(textParts[1], textParts[2])
		case "quit", "end":
			running = false
		}
	}
}


/////////////
// AI CODE //
///////////////////////////////////////////////////////////////////////////////////////
// Below is the the code for a bot that plays random moves
// todo -- Replace this bot with one that plays smarter. Good Luck!
type RandomBot string
var moveList = []string{"down", "left", "right", "turnleft", "turnright"}
func (r *RandomBot) GetMove(gs *GameState, t int) string{
	var buffer bytes.Buffer
	for i := 0; i < rand.Intn(10); i++ {
		buffer.WriteString(fmt.Sprintf("%s,",moveList[rand.Intn(len(moveList))]))
	}
	buffer.WriteString("drop")
	return buffer.String()
}

func (r *RandomBot) String() string{
	return string(*r)
}

func GetRandomBot() BotAi{
	r := RandomBot("RandomBot")
	var b BotAi = &r
	return b
}