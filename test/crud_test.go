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
		defer func() { _ = rsp2.Body.Close() }()
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

		res, err := io.ReadAll(rsp.Body)
		require.NoError(t, err)
		assert.Equal(t, "tea", gjson.GetBytes(res, "name").String())
		assert.Equal(t, "drink", gjson.GetBytes(res, "type").String())
		assert.Equal(t, "Hot Drinks", gjson.GetBytes(res, "drinkType").String())
		assert.True(t, gjson.GetBytes(res, "uuid").Exists())
		assert.NotZero(t, gjson.GetBytes(res, "createdAt").Time())
		assert.NotZero(t, gjson.GetBytes(res, "updatedAt").Time())
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

	t.Run("No name", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		req := bytes.NewBufferString(`{"title": "tea", "drinkType": "Hot Drinks"}`)
		rsp, err := http.Post(srv.URL+"/drinks", "application/json", req)
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()
		assert.Equal(t, http.StatusBadRequest, rsp.StatusCode)
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
		defer func() { _ = rsp2.Body.Close() }()
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

func TestUpdate(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		reqBody := bytes.NewBufferString(`{"name": "tea", "drinkType": "Sleep Drinks"}`)
		rsp, err := http.Post(srv.URL+"/drinks", "application/json", reqBody)
		require.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()

		reqBody2 := bytes.NewBufferString(`{"name": "tea", "drinkType": "Hot Drinks"}`)
		req2, err := http.NewRequest(http.MethodPut, srv.URL+"/drinks/tea", reqBody2)
		assert.NoError(t, err)

		rsp2, err := http.DefaultClient.Do(req2)
		assert.NoError(t, err)
		defer func() { _ = rsp2.Body.Close() }()

		assert.Equal(t, http.StatusOK, rsp2.StatusCode)

		res2, err := io.ReadAll(rsp2.Body)
		require.NoError(t, err)
		assert.Equal(t, "tea", gjson.GetBytes(res2, "name").String())
		assert.Equal(t, "drink", gjson.GetBytes(res2, "type").String())
		assert.Equal(t, "Hot Drinks", gjson.GetBytes(res2, "drinkType").String())
		assert.True(t, gjson.GetBytes(res2, "uuid").Exists())
		assert.NotZero(t, gjson.GetBytes(res2, "createdAt").Time())
		assert.NotZero(t, gjson.GetBytes(res2, "updatedAt").Time())

		rsp3, err := http.Get(srv.URL + "/drinks/tea")
		assert.NoError(t, err)
		defer func() { _ = rsp3.Body.Close() }()
		assert.Equal(t, http.StatusOK, rsp3.StatusCode)

		res3, err := io.ReadAll(rsp3.Body)
		require.NoError(t, err)

		assert.JSONEq(t, string(res2), string(res3))
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		reqBody := bytes.NewBufferString(`{"name": "tea", "drinkType": "Hot Drinks"}`)
		req, err := http.NewRequest(http.MethodPut, srv.URL+"/drinks/tea", reqBody)
		assert.NoError(t, err)

		rsp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()

		assert.Equal(t, http.StatusNotFound, rsp.StatusCode)
	})

	t.Run("Invalid kind", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		reqBody := bytes.NewBufferString(`{"name": "tea", "drinkType": "Hot Drinks"}`)
		req, err := http.NewRequest(http.MethodPut, srv.URL+"/drink/tea", reqBody)
		assert.NoError(t, err)

		rsp, err := http.DefaultClient.Do(req)
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

		reqBody := bytes.NewBufferString(`{"name": "herbal tea", "drinkType": "Hot Drinks"}`)
		req, err := http.NewRequest(http.MethodPut, srv.URL+"/drinks/herbal tea", reqBody)
		assert.NoError(t, err)

		rsp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()

		assert.Equal(t, http.StatusBadRequest, rsp.StatusCode)
	})
}

func TestPatch(t *testing.T) {
	t.Parallel()

	t.Run("Valid JSON Patch", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		reqBody := bytes.NewBufferString(`{"name": "tea", "drinkType": "Sleep Drinks", "alcoholPercentage": 10}`)
		rsp, err := http.Post(srv.URL+"/drinks", "application/json", reqBody)
		require.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()

		reqBody2 := bytes.NewBufferString(`[
{"op": "replace", "path": "/drinkType", "value": "Hot Drinks"},
{"op": "remove", "path": "/alcoholPercentage"}
]`)
		req2, err := http.NewRequest(http.MethodPatch, srv.URL+"/drinks/tea", reqBody2)
		assert.NoError(t, err)

		req2.Header.Set("Content-Type", "application/json-patch+json")

		rsp2, err := http.DefaultClient.Do(req2)
		assert.NoError(t, err)
		defer func() { _ = rsp2.Body.Close() }()

		assert.Equal(t, http.StatusOK, rsp2.StatusCode)

		res2, err := io.ReadAll(rsp2.Body)
		require.NoError(t, err)
		assert.Equal(t, "tea", gjson.GetBytes(res2, "name").String())
		assert.Equal(t, "drink", gjson.GetBytes(res2, "type").String())
		assert.Equal(t, "Hot Drinks", gjson.GetBytes(res2, "drinkType").String())
		assert.False(t, gjson.GetBytes(res2, "alcoholPercentage").Exists())
		assert.True(t, gjson.GetBytes(res2, "uuid").Exists())
		assert.NotZero(t, gjson.GetBytes(res2, "createdAt").Time())
		assert.NotZero(t, gjson.GetBytes(res2, "updatedAt").Time())

		rsp3, err := http.Get(srv.URL + "/drinks/tea")
		assert.NoError(t, err)
		defer func() { _ = rsp3.Body.Close() }()
		assert.Equal(t, http.StatusOK, rsp3.StatusCode)

		res3, err := io.ReadAll(rsp3.Body)
		require.NoError(t, err)

		assert.JSONEq(t, string(res2), string(res3))
	})

	t.Run("Valid Merge Patch", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		reqBody := bytes.NewBufferString(`{"name": "tea", "drinkType": "Sleep Drinks", "alcoholPercentage": 10}`)
		rsp, err := http.Post(srv.URL+"/drinks", "application/json", reqBody)
		require.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()

		reqBody2 := bytes.NewBufferString(`{"drinkType": "Hot Drinks", "alcoholPercentage": null}`)
		req2, err := http.NewRequest(http.MethodPatch, srv.URL+"/drinks/tea", reqBody2)
		assert.NoError(t, err)

		req2.Header.Set("Content-Type", "application/merge-patch+json")

		rsp2, err := http.DefaultClient.Do(req2)
		assert.NoError(t, err)
		defer func() { _ = rsp2.Body.Close() }()

		assert.Equal(t, http.StatusOK, rsp2.StatusCode)

		res2, err := io.ReadAll(rsp2.Body)
		require.NoError(t, err)
		assert.Equal(t, "tea", gjson.GetBytes(res2, "name").String())
		assert.Equal(t, "drink", gjson.GetBytes(res2, "type").String())
		assert.Equal(t, "Hot Drinks", gjson.GetBytes(res2, "drinkType").String())
		assert.False(t, gjson.GetBytes(res2, "alcoholPercentage").Exists())
		assert.True(t, gjson.GetBytes(res2, "uuid").Exists())
		assert.NotZero(t, gjson.GetBytes(res2, "createdAt").Time())
		assert.NotZero(t, gjson.GetBytes(res2, "updatedAt").Time())

		rsp3, err := http.Get(srv.URL + "/drinks/tea")
		assert.NoError(t, err)
		defer func() { _ = rsp3.Body.Close() }()
		assert.Equal(t, http.StatusOK, rsp3.StatusCode)

		res3, err := io.ReadAll(rsp3.Body)
		require.NoError(t, err)

		assert.JSONEq(t, string(res2), string(res3))
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		reqBody := bytes.NewBufferString(`{"drinkType": "Hot Drinks", "alcoholPercentage": null}`)
		req, err := http.NewRequest(http.MethodPatch, srv.URL+"/drinks/tea", reqBody)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/merge-patch+json")

		rsp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()

		assert.Equal(t, http.StatusNotFound, rsp.StatusCode)
	})

	t.Run("Invalid kind", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		reqBody := bytes.NewBufferString(`{"drinkType": "Hot Drinks", "alcoholPercentage": null}`)
		req, err := http.NewRequest(http.MethodPatch, srv.URL+"/drink/tea", reqBody)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/merge-patch+json")

		rsp, err := http.DefaultClient.Do(req)
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

		reqBody := bytes.NewBufferString(`{"drinkType": "Hot Drinks", "alcoholPercentage": null}`)
		req, err := http.NewRequest(http.MethodPatch, srv.URL+"/drinks/herbal tea", reqBody)
		assert.NoError(t, err)

		req.Header.Set("Content-Type", "application/merge-patch+json")

		rsp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()

		assert.Equal(t, http.StatusBadRequest, rsp.StatusCode)
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		reqBody := bytes.NewBufferString(`{"name": "tea", "drinkType": "Sleep Drinks"}`)
		rsp, err := http.Post(srv.URL+"/drinks", "application/json", reqBody)
		require.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()

		req2, err := http.NewRequest(http.MethodDelete, srv.URL+"/drinks/tea", nil)
		assert.NoError(t, err)

		rsp2, err := http.DefaultClient.Do(req2)
		assert.NoError(t, err)
		defer func() { _ = rsp2.Body.Close() }()

		assert.Equal(t, http.StatusNoContent, rsp2.StatusCode)

		rsp3, err := http.Get(srv.URL + "/drinks/tea")
		assert.NoError(t, err)
		defer func() { _ = rsp3.Body.Close() }()
		assert.Equal(t, http.StatusNotFound, rsp3.StatusCode)
	})

	t.Run("Not Found", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		req, err := http.NewRequest(http.MethodDelete, srv.URL+"/drinks/tea", nil)
		assert.NoError(t, err)

		rsp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()

		assert.Equal(t, http.StatusNotFound, rsp.StatusCode)
	})

	t.Run("Invalid kind", func(t *testing.T) {
		t.Parallel()

		repo := memory.NewItemRepository()

		h := transport.New(repo)

		srv := httptest.NewServer(h)
		defer srv.Close()

		req, err := http.NewRequest(http.MethodDelete, srv.URL+"/drink/tea", nil)
		assert.NoError(t, err)

		rsp, err := http.DefaultClient.Do(req)
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

		req, err := http.NewRequest(http.MethodDelete, srv.URL+"/drinks/herbal tea", nil)
		assert.NoError(t, err)

		rsp, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		defer func() { _ = rsp.Body.Close() }()

		assert.Equal(t, http.StatusBadRequest, rsp.StatusCode)
	})
}
