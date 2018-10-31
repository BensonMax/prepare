package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		//定义一个cmdduixiang
		cmd *exec.Cmd
		err error
	)
	//通过函数获得cmd 对象
	cmd = exec.Command("C:\\cygwin64\\bin\\bash.exe", "-c", "echo 1")

	err = cmd.Run()

	fmt.Println(err)
}
