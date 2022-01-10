package test

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nasermirzaei89/core/internal/repository/memory"
	"github.com/nasermirzaei89/core/internal/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
)

func TestList(t *testing.T) {
	t.Parallel()

	t.Run("Empty", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		// first list
		rsp, err := http.Get(srv.URL + "/drinks")
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()
		assert.Equal(t, http.StatusOK, rsp.StatusCode)

		emptyResults := `{"items": []}`

		res, err := io.ReadAll(rsp.Body)
		require.NoError(t, err)

		assert.JSONEq(t, emptyResults, string(res))
	})

	t.Run("After first create", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		firstItemReq := bytes.NewBufferString(`{"name": "tea", "drinkType": "Hot Drinks"}`)
		rsp, err := http.Post(srv.URL+"/drinks", "application/json", firstItemReq)
		require.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()

		res, err := io.ReadAll(rsp.Body)
		require.NoError(t, err)

		rsp2, err := http.Get(srv.URL + "/drinks")
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()
		assert.Equal(t, http.StatusOK, rsp2.StatusCode)

		res2, err := io.ReadAll(rsp2.Body)
		require.NoError(t, err)

		assert.EqualValues(t, 1, gjson.GetBytes(res2, "items.#").Int())
		assert.JSONEq(t, string(res), gjson.GetBytes(res2, "items.0").String())
	})

	t.Run("Invalid type", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		rsp, err := http.Get(srv.URL + "/drink")
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, rsp.StatusCode)
	})
}

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		req := bytes.NewBufferString(`{"name": "tea", "drinkType": "Hot Drinks"}`)
		rsp, err := http.Post(srv.URL+"/drinks", "application/json", req)
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()
		assert.Equal(t, http.StatusCreated, rsp.StatusCode)

		res2, err := io.ReadAll(rsp.Body)
		require.NoError(t, err)
		assert.Equal(t, "tea", gjson.GetBytes(res2, "name").String())
		assert.Equal(t, "drink", gjson.GetBytes(res2, "type").String())
		assert.Equal(t, "Hot Drinks", gjson.GetBytes(res2, "drinkType").String())
		assert.True(t, gjson.GetBytes(res2, "uuid").Exists())
		assert.NotZero(t, gjson.GetBytes(res2, "createdAt").Time())
		assert.NotZero(t, gjson.GetBytes(res2, "updatedAt").Time())
	})

	t.Run("Invalid type", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		req := bytes.NewBufferString(`{"name": "tea", "drinkType": "Hot Drinks"}`)
		rsp, err := http.Post(srv.URL+"/drink", "application/json", req)
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, rsp.StatusCode)
	})

	t.Run("Invalid name", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		req := bytes.NewBufferString(`{"name": "herbal tea", "drinkType": "Hot Drinks"}`)
		rsp, err := http.Post(srv.URL+"/drinks", "application/json", req)
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, rsp.StatusCode)
	})

	t.Run("Duplicate name", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		req := bytes.NewBufferString(`{"name": "tea", "drinkType": "Hot Drinks"}`)
		rsp, err := http.Post(srv.URL+"/drinks", "application/json", req)
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()
		assert.Equal(t, http.StatusCreated, rsp.StatusCode)

		req2 := bytes.NewBufferString(`{"name": "tea", "drinkType": "Hot Drinks"}`)
		rsp2, err := http.Post(srv.URL+"/drinks", "application/json", req2)
		assert.NoError(t, err)
		defer func() { _ = rsp2.Body.Close() }()
		assert.Equal(t, http.StatusConflict, rsp2.StatusCode)
	})
}

func TestRead(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		req := bytes.NewBufferString(`{"name": "tea", "drinkType": "Hot Drinks"}`)
		rsp, err := http.Post(srv.URL+"/drinks", "application/json", req)
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()

		res, err := io.ReadAll(rsp.Body)
		require.NoError(t, err)

		rsp2, err := http.Get(srv.URL + "/drinks/tea")
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()
		assert.Equal(t, http.StatusOK, rsp2.StatusCode)

		res2, err := io.ReadAll(rsp2.Body)
		require.NoError(t, err)

		assert.JSONEq(t, string(res), string(res2))
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		rsp, err := http.Get(srv.URL + "/drinks/tea")
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()
		assert.Equal(t, http.StatusNotFound, rsp.StatusCode)
	})

	t.Run("Invalid type", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		rsp, err := http.Get(srv.URL + "/drink/tea")
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, rsp.StatusCode)
	})

	t.Run("Invalid name", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		rsp, err := http.Get(srv.URL + "/drinks/herbal tea")
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, rsp.StatusCode)
	})
}
