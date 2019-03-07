package web

import (
	"encoding/json"
	"gonote/pkg/context"
	"net/http"
)

func Create(ctx *context.Context) {
	var note struct {
		Data []interface{}
		Dict map[string]interface{}
	}
	ctx.Input.ReadBody()
	json.Unmarshal(ctx.Input.RawContent, &note)
	var name string
	var num int
	v := note.Dict["name"]
	switch v.(type) {
	case string:
		name = v.(string)
	default:
		ctx.Abort(context.HttpError{Status: http.StatusBadRequest, Message: []byte("400 BadRequest")})
	}

	v = note.Dict["num"]
	switch v.(type) {
	case int:
		num = v.(int)
	case uint:
		num = v.(int)
	case float64:
		num = int(v.(float64))
	default:
		ctx.Abort(context.HttpError{Status: http.StatusBadRequest, Message: []byte("400 BadRequest")})
	}

	err := insert(name, num)
	if err != nil {
		panic(err.Error())
	}
	data := struct {
		Data string `json:"data"`
	}{"已添加"}
	message, _ := json.Marshal(data)
	ctx.Output.AppendContent(message)
	ctx.Output.SetStatus(200)
	ctx.Output.ServeJson()
}

func Update(ctx *context.Context) {
	data := struct {
		Name *string
		Num  *int
	}{}
	uuid := ctx.Input.Args["uuid"].(string)
	json.Unmarshal(ctx.Input.RawContent, &data)
	if data.Name == nil || data.Num == nil {
		ctx.Abort(context.HttpError{Status: http.StatusBadRequest, Message: []byte("400 BadRequest")})
	}

	update(*data.Name, *data.Num, uuid)
	result := struct {
		Data string `json:"data"`
	}{"已修改"}
	message, _ := json.Marshal(result)
	ctx.Output.AppendContent(message)
	ctx.Output.ServeJson()
}

func Delete(ctx *context.Context) {
	uuid := ctx.Input.Args["uuid"].(string)
	if uuid == "" {
		ctx.Abort(context.HttpError{Status: http.StatusBadRequest, Message: []byte("400 BadRequest")})
	}
	delete(uuid)
	result := struct {
		Data string `json:"data"`
	}{"已删除"}
	message, _ := json.Marshal(result)
	ctx.Output.AppendContent(message)
	ctx.Output.ServeJson()
}

func Get(ctx *context.Context) {
	noteList := get()
	message, _ := json.Marshal(noteList)
	ctx.Output.AppendContent(message)
	ctx.Output.ServeJson()
}
