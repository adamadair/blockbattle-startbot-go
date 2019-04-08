package main

// Riddles Game Engine communicator
//
// This can be tweaked to accept input from files and other sources
// which can be useful when testing.

import (
"bufio"
"fmt"
"os"
"strings"
)

type botIO struct {
	out *bufio.Writer
	log *bufio.Writer
	in *bufio.Reader
}

// BotIO is used for all communication with the Riddles game
// engine. This piece is reusable because all Riddles.io competitions
// communicate with the bot in the same way
var BotIO botIO

func init() {
	BotIO.in = bufio.NewReader(os.Stdin)
	BotIO.out = bufio.NewWriter(os.Stdout)
	BotIO.log = bufio.NewWriter(os.Stderr)
}

// get a line of input from the engine
func (b botIO) ReadLine() (line string, err error){
	line, err = b.in.ReadString('\n')
	if err == nil{
		// On windows system there will still be extemporaneous bytes on the end.
		line = strings.TrimSpace(line)
	}
	return
}

// write a line of output to the engine
func (b botIO) WriteLine(msg string, ctx ...interface{}) (i int, e error){
	// [ go note] : to pass multiple return values you must surround them with parenthesis
	i, e = b.out.WriteString(fmt.Sprintf("%s\n", fmt.Sprintf(msg, ctx...)))
	if e == nil{
		e = b.out.Flush()
	}
	return // implicit return with named return variables
}

// Write an INFO message to the log
// These messages show up in the game log.
func (b botIO) LogInfo(msg string, ctx ...interface{}) (i int, e error){
	i, e = b.log.WriteString(fmt.Sprintf("INFO: %s\n", fmt.Sprintf(msg, ctx...)))
	if e == nil{
		e = b.log.Flush()
	}
	return
}

// Write an ERROR message to the log
// These messages show up in the game log.
func (b botIO) LogError(msg string, ctx ...interface{}) (i int, e error){
	s := fmt.Sprintf(msg, ctx...)
	i, e = b.log.WriteString(fmt.Sprintf("ERROR: %s\n", s))
	if e == nil{
		e = b.log.Flush()
	}
	return
}