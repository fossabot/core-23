package main

import (
	"net/http"

	"github.com/nasermirzaei89/core/internal/repository/memory"
	"github.com/nasermirzaei89/core/internal/transport"
	"github.com/nasermirzaei89/env"
	"github.com/pkg/errors"
)

func main() {
	repo := memory.NewMemoryRepository()

	h := transport.New(repo)

	err := http.ListenAndServe(env.GetString("API_ADDRESS", ":80"), h)
	if err != nil {
		panic(errors.Wrap(err, "error on listen and serve http"))
	}
}
