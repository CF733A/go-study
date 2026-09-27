package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	var res string
	res = contextWithTimeout(context.Background(), 1*time.Second, 2*time.Second)
	fmt.Println(res)
	res = contextWithTimeout(context.Background(), 2*time.Second, 1*time.Second)
	fmt.Println(res)
}

func contextWithTimeout(ctx context.Context, contextTimeout time.Duration, timeAfter time.Duration) string {
	deadline := time.Now().Add(timeAfter)
	deadlineCtx, deadlineCancel := context.WithDeadline(ctx, deadline)
	defer deadlineCancel()

	timeoutCtx, timeoutCancel := context.WithTimeout(ctx, contextTimeout)
	defer timeoutCancel()

	select {
	case <-deadlineCtx.Done():
		return "превышено время ожидания"
	case <-timeoutCtx.Done():
		return "превышено время ожидания контекста"
	}
}
