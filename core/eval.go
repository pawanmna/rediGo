package core

import "errors"

func evalPing(args []string) ([]byte, error) {
	if len(args) >= 2 {
		return nil, errors.New("The number of arguments cannot be > 1 in PING")
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
		return evalPing(cmd.Args)
	}
}
