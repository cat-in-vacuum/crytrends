package providers

import (
	"context"
	"github.com/rs/xid"
)

const (
	// ключи контекста
	ctxIDKey = "id_key"
)

// TODO ...тут должен быть конвертор в общий тип ответа для всех  поставщиков
type OHLCAssets struct {
}

// будем использовать для передачи вспомогательных данных внутри процессов
type Context struct {
	context.Context
}

func (ctx *Context) SetID() {
	ctx.Context = context.WithValue(context.Background(), ctxIDKey, xid.New().String())
}

func (ctx *Context) GetID() string {
	if id, ok := ctx.Context.Value(ctxIDKey).(string); ok {
		return id
	}
	ctx.SetID()
	return ctx.GetID()
}
