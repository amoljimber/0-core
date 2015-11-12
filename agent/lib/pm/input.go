package pm

import (
	"bufio"
	"github.com/Jumpscale/agent2/agent/lib/pm/core"
	"io"
	"log"
	"regexp"
	"strconv"
	"strings"
)

var pmMsgPattern, _ = regexp.Compile("^(\\d+)(:{2,3})(.*)$")

type streamConsumer struct {
	cmd    *core.Cmd
	reader io.Reader
	level  int
	Signal chan int
}

func newStreamConsumer(cmd *core.Cmd, reader io.Reader, level int) *streamConsumer {
	return &streamConsumer{
		cmd:    cmd,
		reader: reader,
		level:  level,
		Signal: make(chan int),
	}
}

// read input until the end (or closed)
// process all messages as speced x:: or x:::
// other messages that has no level are assumed of level consumer.level
func (consumer *streamConsumer) consume(handler MessageHandler) {
	reader := bufio.NewReader(consumer.reader)
	var level int
	var message string
	var multiline = false

	for {
		line, err := reader.ReadString('\n')

		if err != nil && err != io.EOF {
			log.Println(err)
			return
		}

		line = strings.TrimRight(line, "\n")

		if line != "" {
			if !multiline {
				matches := pmMsgPattern.FindStringSubmatch(line)
				if matches == nil {
					//use default level.
					handler(&Message{
						Cmd:     consumer.cmd,
						Level:   consumer.level,
						Message: line,
					})
				} else {
					l, _ := strconv.ParseInt(matches[1], 10, 0)
					level = int(l)
					message = matches[3]

					if matches[2] == ":::" {
						multiline = true
					} else {
						//single line message
						handler(&Message{
							Cmd:     consumer.cmd,
							Level:   level,
							Message: message,
						})
					}
				}
			} else {
				/*
				   A known issue is that if stream was closed (EOF) before
				   we receive the ::: termination of multiline string. We discard
				   the uncomplete multiline string message.
				*/
				if line == ":::" {
					multiline = false
					//flush message
					handler(&Message{
						Cmd:     consumer.cmd,
						Level:   level,
						Message: message,
					})
				} else {
					message += "\n" + line
				}
			}
		}

		if err == io.EOF {
			consumer.Signal <- 1
			close(consumer.Signal)
			return
		}
	}
}

func (consumer *streamConsumer) Consume(handler MessageHandler) {
	go consumer.consume(handler)
}
