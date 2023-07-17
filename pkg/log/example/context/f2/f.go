package main

import (
	"context"

	"github.com/wangweihong/eazycloud/pkg/log"
)

func main() {
	// save key/value pair in context
	ctx := log.WithFieldPair(context.Background(), "ip", "10.30.100.111")

	// save logger in context
	ctx = log.F(ctx).WithContext(ctx)

	// every log will carry ip:10.30.100.111
	log.FromContext(ctx).Info("bbb")
	log.FromContext(ctx).Info("okok", log.String("name", "bbb"))
}
