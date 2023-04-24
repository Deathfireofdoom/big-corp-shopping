package main

import (
	"cart-service/internal/cart_service"
	"context"
)

func main() {
	ctx := context.Background()
	cart_service.Service.Run(ctx)
}