package main

import (
	"apishoppingcart_client/src/pb/shoppingcart"
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	//abrindo a conexao
	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("error on get connection. error: ", err)
	}
	defer conn.Close()

	client := shoppingcart.NewShoppingCartServiceClient(conn)
	stream, err := client.AddItem(context.Background())
	if err != nil {
		log.Fatalln("error on get channel to stream. error: ", err)
	}

	waitch := make(chan struct{})
	go func() {
		for {
			response, err := stream.Recv()
			if err == io.EOF {
				close(waitch)
				return
			}
			if err != nil {
				log.Fatalln("error on recv. error: ", err)
			}
			fmt.Printf("response: %+v\n", response)
		}
	}()

	items := []shoppingcart.AddProduct{
		{ProductId: 1, Quantity: 2, PriceUnit: 5.0},
		{ProductId: 2, Quantity: 7, PriceUnit: 12.00},
		{ProductId: 3, Quantity: 17, PriceUnit: 2.50},
	}

	for _, v := range items {
		if err := stream.Send(&v); err != nil {
			log.Fatalln("error on send. error: ", err)
		}
		fmt.Printf("-> send: %+v\n", v)
	}
	stream.CloseSend()
	<-waitch
}
