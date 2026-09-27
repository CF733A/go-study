package main

import (
	"context"
	"time"
)

func main() {
	var res string
	res = contextWithDeadline(context.Background(), 1*time.Second, 2*time.Second)
	println(res)
	res = contextWithDeadline(context.Background(), 2*time.Second, 1*time.Second)
	println(res)
	/* Output:
	context deadline exceeded
	time after exceeded
	*/
}

func contextWithDeadline(ctx context.Context, contextDeadline time.Duration, timeAfter time.Duration) string {
	// var cancel context.CancelFunc
	// defer cancel()

	var res string

	deadline := time.Now().Add(contextDeadline)

	ctx, cancel := context.WithDeadline(ctx, deadline)
	defer cancel()

	ctxTimeOut, canccelTimeOut := context.WithTimeout(ctx, timeAfter)
	defer canccelTimeOut()

	select {
	case <-ctx.Done():
		res = "context deadline exceeded"
	case <-ctxTimeOut.Done():
		res = "time after exceeded"
	}

	return res
}
