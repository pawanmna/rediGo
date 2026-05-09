package core

import "errors"

func evalPing(args []string) ([]byte, error) {
	if len(args) >= 2 {
		return nil, errors.New("the number of arguments cannot be > 1 for PING")
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
		return []byte{}, errors.New("unknown command")
	}
}
