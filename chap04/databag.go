package main

import (
	"context"
	"fmt"
)

func main() {
	ProcessRequest("jane", "Super secret value!")
}

type ctxKey int

const (
	_ ctxKey = iota
	ctxUserID
	ctxAuthToken
)

func UserID(ctx context.Context) string {
	return ctx.Value(ctxUserID).(string)
}

func AuthToken(ctx context.Context) string {
	return ctx.Value(ctxAuthToken).(string)
}

func ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf("handling response for %v (%v)\n", UserID(ctx), AuthToken(ctx))
}
