package examples

import (
	"testing"

	"github.com/PMoneda/http"
)

func TestShouldConnectToPruuWithMock(t *testing.T) {
	http.With(t, func(ctx *http.MockContext) {
		ctx.RegisterMock(&http.ReponseMock{
			ReponseBody: "ok",
		})
		if postToPruu() != "ok" {
			ctx.Fail()
		}
	})
}
