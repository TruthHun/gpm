package serv

import (
	"fmt"
	"testing"
)

func TestExecCommand(t *testing.T) {
	fmt.Println(execCommand("ping", "www.baidu.com"))
}
