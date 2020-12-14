package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

func GetHelloWorldServer(ctx context.Context, wg *sync.WaitGroup) func() error {
	return func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`Hello, world!`))
			if err != nil {
				return
			}
		})
		server := &http.Server{Addr: ":7000", Handler: mux}
		errChan := make(chan error, 1)

		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(shutCtx); err != nil {
				errChan <- fmt.Errorf("关闭 hello world 服务出现错误: %w", err)
			}
			fmt.Println("hello world 服务已关闭")
			close(errChan)
			wg.Done()
		}()

		fmt.Println("启动 hello world 服务...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("启动 hello world 服务出现错误: %w", err)
		}
		fmt.Println("关闭 hello world 服务...")
		err := <-errChan
		wg.Wait()
		return err
	}
}

func GetHelloNameServer(ctx context.Context, wg *sync.WaitGroup) func() error {
	return func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			params := r.URL.Query()
			name := params.Get("name")

			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(fmt.Sprintf("hello, %s ~", name)))
			if err != nil {
				return
			}
		})
		server := http.Server{Addr: ":8000", Handler: mux}
		errChan := make(chan error, 1)

		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(shutCtx); err != nil {
				errChan <- fmt.Errorf("关闭 hello name 服务出现错误: %w", err)
			}
			fmt.Println("hello name 服务已关闭")
			close(errChan)
			wg.Done()
		}()

		fmt.Println("启动 hello name 服务...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("启动 hello name 服务出现错误: %w", err)
		}
		fmt.Println("关闭 hello name 服务...")
		err := <-errChan
		wg.Wait()
		return err
	}

}

func GetEchoServer(ctx context.Context, wg *sync.WaitGroup) func() error {
	return func() error {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := io.Copy(w, r.Body)
			if err != nil {
				return
			}
		})
		server := http.Server{Addr: ":9000", Handler: mux}
		errChan := make(chan error, 1)

		go func() {
			<-ctx.Done()
			shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			if err := server.Shutdown(shutCtx); err != nil {
				errChan <- fmt.Errorf("关闭 echo 服务出现错误: %w", err)
			}
			fmt.Println("echo 服务已关闭")
			close(errChan)
			wg.Done()
		}()

		fmt.Println("启动 echo 服务...")
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			return fmt.Errorf("启动 echo 服务出现错误: %w", err)
		}
		fmt.Println("关闭 echo 服务...")
		err := <-errChan
		wg.Wait()
		return err
	}
}
