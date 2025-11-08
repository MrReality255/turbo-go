package utils

import (
	"net/http"

	"github.com/kataras/iris/v12"
)

type WebContextHandler struct {
	ctx     iris.Context
	onError func(string, error)
}

func (wch *WebContextHandler) HandleError(hint string, err error) {
	if wch.onError != nil {
		wch.onError(hint, err)
	}
}

func (wch *WebContextHandler) RespondJSON(data interface{}) {
	wch.ctx.StatusCode(http.StatusOK)
	_, err := wch.ctx.Write(ToJSONB(data))
	wch.HandleError("write failed", err)
}

func (wch *WebContextHandler) RespondJSONErr(content any, err error) {
	if err != nil {
		wch.SrvErr(err)
		return
	}
	wch.RespondJSON(content)
}

func (wch *WebContextHandler) ClientErr(ctx iris.Context, err error) {
	ctx.StatusCode(http.StatusBadRequest)
	wch.HandleError("client error", err)
}

func (wch *WebContextHandler) SrvErr(err error) {
	wch.ctx.StatusCode(http.StatusInternalServerError)
	wch.HandleError("server error", err)
}

func (wch *WebContextHandler) GetParam(key string) string {
	return wch.ctx.Params().Get(key)
}

func NewContextHandler(ctx iris.Context, onError func(string, error)) *WebContextHandler {
	return &WebContextHandler{
		ctx:     ctx,
		onError: onError,
	}
}
