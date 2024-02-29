package src

import (
	"crypto/rand"
	"encoding/hex"
)

func GetUID() string {
	b := make([]byte, 8) // 8 bytes will give us 16 hex characters
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func NewCancellable[T map[string]any](inputs T, fetcher func(f func() T) T) T {
	return fetcher(func() T {
		return inputs
	})
}
