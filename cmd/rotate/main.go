package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"marketplace-key-rotation/internal/infrai"
	"marketplace-key-rotation/internal/marketplace"
)

func idem() string { b := make([]byte, 8); _, _ = rand.Read(b); return hex.EncodeToString(b) }
func main() {
	c, err := infrai.NewFromEnv()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	var created struct {
		ID  string `json:"key_id"`
		Key string `json:"key_secret"`
	}
	if err := c.CreateKey(ctx, "marketplace", "rotation-staging", []string{"account"}, idem(), &created); err != nil {
		log.Fatal(err)
	}
	var rotated struct {
		Key string `json:"key"`
	}
	if err := c.RotateKey(ctx, created.ID, 24, idem(), &rotated); err != nil {
		log.Fatal(err)
	}
	plan, err := marketplace.PlanRotation([]marketplace.Deployment{{Name: "orders", APIKey: "old"}}, created.Key, rotated.Key)
	if err != nil {
		log.Fatal(err)
	}
	var usage any
	if err := c.Usage(ctx, &usage); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("new key issued once; handoff=%v; usage evidence collected\n", plan.HandoffOrder)
}
