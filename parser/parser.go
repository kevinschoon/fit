package parser

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Parser interface {
	Parse(string) (float64, error)
}

// Parsers loads Parser types from the array of strings
// Example:
// 1,Time,2006-01-02
// ^-^----^----index(int),name(string),format(string)
func ParsersFromArgs(args []string) (map[int]Parser, error) {
	parsers := make(map[int]Parser)
	for _, arg := range args {
		split := strings.Split(arg, ",")
		if len(split) < 3 {
			return nil, fmt.Errorf("Bad parser opts: %s", arg)
		}
		index, err := strconv.ParseInt(split[0], 0, 64)
		if err != nil {
			return nil, fmt.Errorf("Bad parser opts: %s", arg)
		}
		switch split[1] {
		case "Time":
			if len(split) != 3 {
				return nil, fmt.Errorf("Bad parser opts: %s", arg)
			}
			parsers[int(index)] = TimeParser{Format: split[2]}
		default:
			return nil, fmt.Errorf("Unknown parser: %s", split[1])
		}
	}
	return parsers, nil
}

type TimeParser struct {
	Format string
}

func (t TimeParser) Parse(v string) (float64, error) {
	parsed, err := time.Parse(t.Format, v)
	if err != nil {
		return 0.0, err
	}
	return float64(parsed.Unix()), nil
}
