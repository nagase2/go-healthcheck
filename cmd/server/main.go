package main

import (
	// (一部抜粋)
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	hellopb "example/hello/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type myServer struct {
	hellopb.UnimplementedGreetingServiceServer
}

func (s *myServer) Hello(ctx context.Context, req *hellopb.HelloRequest) (*hellopb.HelloResponse, error) {
	// リクエストからnameフィールドを取り出して
	// "Hello, [名前]!"というレスポンスを返す
	fmt.Println("Helloが呼ばれました..", req.GetName())
	return &hellopb.HelloResponse{
		Message: fmt.Sprintf("🐔Hello, %s!", req.GetName()),
	}, nil
}

// 自作サービス構造体のコンストラクタを定義
func NewMyServer() *myServer {
	return &myServer{}
}

func main() {
	// 1. 8080番portのLisnterを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// 2. gRPCサーバーを作成
	s := grpc.NewServer()

	// 3. gRPCサーバーにGreetingServiceを登録
	// hellopb.RegisterGreetingServiceServer(s, [サーバーに登録するサービス])
	hellopb.RegisterGreetingServiceServer(s, NewMyServer())

	// 4. サーバーリフレクションの設定
	reflection.Register(s)

	// 3. 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// 4.Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
