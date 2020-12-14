package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"golang.org/x/sync/errgroup"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(3)

	ctx, cancel := context.WithCancel(context.Background())

	eg, egCtx := errgroup.WithContext(context.Background())
	eg.Go(GetHelloWorldServer(ctx, &wg))
	eg.Go(GetHelloNameServer(ctx, &wg))
	eg.Go(GetEchoServer(ctx, &wg))

	go func() {
		<-egCtx.Done()
		cancel()
	}()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)
		<-signals
		cancel()
	}()

	if err := eg.Wait(); err != nil {
		fmt.Printf("在 goroutines 中的服务出现错误 : %s\n", err)
		os.Exit(1)
	}
	fmt.Println("所有服务已正常关闭！")
}
