package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/sync/errgroup"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}
}

func run(ctx context.Context) error {
	s := &http.Server{
		Addr: ":80",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}
	//戻り値にエラーが含まれるゴルーチン間の並列処理を実装
	eg, ctx := errgroup.WithContext(ctx)
	//別ゴルーチンで無記名関数でHTTPサーバーを起動
	eg.Go(func() error {
		//19行目のhttp.Serverを経由してサーバーを起動
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	//チャネルから通知を受信した場合、run()を終了
	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	//別ゴルーチン処理を待機
	return eg.Wait()
}
