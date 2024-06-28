package main

import (
	"crypto/rand"
    "encoding/hex"
)

func genSessionid () string {
    raw := make([]byte, 32)
    _, err := rand.Read(raw)
    if err != nil {panic(err)}
    return hex.EncodeToString(raw)
}
