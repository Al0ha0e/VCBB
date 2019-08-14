package slave_side

import (
	"fmt"
	"testing"

	"github.com/go-redis/redis"
)

const url = "http://127.0.0.1:8080/hello/test"

func TestPyExecuter(t *testing.T) {
	cli := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	cli.Set("A", "1", 0)
	cli.Set("B", "2", 0)
	exe := NewPyExecuter(url)
	code :=
		`def func():
	print(input[0],input[1],input[0]+input[1])
	return input[0]+input[1]
output=[func()]`
	ansChan := make(chan *executeResult, 1)
	exe.Run(1, []string{"A", "B"}, code, ansChan)
	ans := <-ansChan
	if ans.err != nil {
		fmt.Println(ans.err)
		return
	}
	ansstr, _ := cli.Get(ans.result[0][0]).Result()
	println("ANS", ansstr)
}
