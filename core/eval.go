package core

import (
	"errors"
	"strings"
)

func evalPing(args []string) ([]byte, error) {
	if len(args) >= 2 {
		return nil, errors.New("ERR wrong number of arguments for PING")
	}

	if len(args) == 0 {
		return Encode("PONG", true)
	} else {
		return Encode(args[0], false)
	}
}

func Eval(cmd *RedisCmd) ([]byte, error) {
	switch cmd.Cmd {
	case "PING":
		return evalPing(cmd.Args)
	default:
		return nil, errors.New(
			"ERR unknown command '" + strings.ToLower(cmd.Cmd) + "'",
		)
	}
}
