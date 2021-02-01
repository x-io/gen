package gen

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/x-io/gen/core"
)

func performRequest(r http.Handler, method, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestRouteRawPath(t *testing.T) {
	route := New()
	// route.UseRawPath = true

	route.Post("/project/:name/build/:num", func(c *core.Context) {
		name := c.Params("name")
		num := c.Params("num")

		assert.Equal(t, name, c.Params("name"))
		assert.Equal(t, num, c.Params("num"))

		assert.Equal(t, "Some/Other/Project", name)
		assert.Equal(t, "222", num)
	})

	w := performRequest(route, "POST", "/project/Some%2FOther%2FProject/build/222")
	assert.Equal(t, http.StatusOK, w.Code)
}
