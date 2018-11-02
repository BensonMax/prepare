package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	err    error
	output []byte
}

func main() {
	//执行一个cmd,让它在一个协程里去执行， "sleep 2;echo hello"
	//1秒的时候，我们杀死cmd
	var (
		ctx        context.Context
		cancelFunc context.CancelFunc
		cmd        *exec.Cmd
		resultChan chan *result
		res        *result
	)
	//创建一个结果队列
	resultChan = make(chan *result, 1000)

	ctx, cancelFunc = context.WithCancel(context.TODO())

	go func() {
		var (
			output []byte
			err    error
		)
		cmd = exec.CommandContext(ctx, "C:\\cygwin64\\bin\\bash.exe", "-c", "sleep 2;echo hello")

		//执行任务，捕获输出
		output, err = cmd.CombinedOutput()

		//把任务输出结果传给main协程
		resultChan <- &result{
			err,
			output,
		}
	}()

	//睡一秒
	time.Sleep(1 * time.Second)

	//取消上下文 —> 导致 "exec.CommandContex" 中断（context.Context 里面由有一个context: chan, cancelFunc(): 关闭 chan）
	//cmd 会监听 ctx.Done() ,然后kill pid，杀死子进程
	cancelFunc()

	//在main协程中等待子协程退出，打印执行结果X'z'x
	res = <-resultChan

	fmt.Println(res.err, string(res.output))

}
