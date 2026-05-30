package core

import (
	"errors"
	"fmt"
	"io"
	"strconv"
	"time"
)

var RESP_NIL = []byte("$-1\r\n")

func EvalAndRespond(cmd *RedisCmd, c io.ReadWriter) error {
	switch cmd.Cmd {
	case "PING":
		return evalPING(cmd.Args, c)
	case "SET":
		return evalSET(cmd.Args, c)
	case "GET":
		return evalGET(cmd.Args, c)
	case "TTL":
		return evalTTL(cmd.Args, c)
	default:
		return fmt.Errorf("ERR unknown command '%s'", cmd.Cmd)
	}
}

func evalPING(args []string, c io.ReadWriter) error {
	var b []byte
	if len(args) >= 2 {
		return errors.New("ERR wrong number of arguments for 'PING' command")
	}

	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	_, err := c.Write(b)
	return err
}

func evalSET(args []string, c io.ReadWriter) error {
	if(len(args) <= 1){
		return errors.New("(Error) ERR wrong number of arguments for 'set' command.")
	}
	var key, value string
	var exDurationMs int64 = -1 
	key, value = args[0] , args[1]
	for i := 2; i < len(args); i++{
		switch args[i]{
		case "EX","ex":
			i++
			if i == len(args){
				return errors.New("(error) ERR syntax error.")
			}
			exDurationsSec, err := strconv.ParseInt(args[3],10,64)
			if err != nil{
				return errors.New("(error) ERR value is not an integer or out of range.")
			}
			exDurationMs = exDurationsSec * 1000
		default:
			return errors.New("(error) ERR syntax error.")
		}
	}

	Put(key,NewObj(value,exDurationMs))
	c.Write([]byte("+OK\r\n"))
	return nil
}

func evalGET(args []string, c io.ReadWriter) error {
	if(len(args) != 1){
		return errors.New("(Error) ERR wrong number of arguments for 'get' command.")
	}
	key := args[0]
	obj := Get(key)
	if obj == nil {
		c.Write([]byte(RESP_NIL))
		return nil
	}

	if obj.ExpireAt != -1 && obj.ExpireAt < time.Now().UnixMilli() {
		c.Write([]byte(RESP_NIL))
		return nil
	}

	c.Write(Encode(obj.value.(string), false))
	return nil
}

func evalTTL(args []string, c io.ReadWriter) error {
	if(len(args) != 1){
		return errors.New("(Error) ERR wrong number of arguments for 'ttl' command.")
	}
	var key string = args[0]
	obj := Get(key)

	if obj == nil {
		c.Write([]byte(":-2\r\n"))
		return nil
	}
	
	if obj.ExpireAt == -1 {
		c.Write([]byte(":-1\r\n"))
		return nil
	}

	durationMs := obj.ExpireAt - time.Now().UnixMilli()
	if durationMs < 0 {
		c.Write([]byte(":-2\r\n"))
		return nil
	}
	c.Write(Encode(int64(durationMs/1000), false))
	return nil
}