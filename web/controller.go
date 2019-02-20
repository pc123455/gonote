package web

import (
	"encoding/json"
	"gonote/framework/constant"
	"gonote/framework/context"
	"gonote/framework/logger"
	"io/ioutil"
	"runtime/debug"
)

func Create(ctx *context.Context) {
	jsonStr, err := ioutil.ReadAll(ctx.Body)
	if err != nil {
		logger.Errorf(err.Error())
		logger.Errorf(string(debug.Stack()))
		panic("request body read error")
	}
	var note struct {
		Data []interface{}
		Dict map[string]interface{}
	}
	json.Unmarshal(jsonStr, &note)
	var name string
	var num int
	v := note.Dict["name"]
	switch v.(type) {
	case string:
		name = v.(string)
	default:
		ctx.WriteHeader(constant.BadRequest)
		return
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
		ctx.WriteHeader(constant.BadRequest)
		return
	}

	err = insert(name, num)
	if err != nil {
		panic(err.Error())
	}
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	ctx.ResponseWriter.WriteHeader(constant.Created)
	data := struct {
		Data string `json:"data"`
	}{"已添加"}
	message, _ := json.Marshal(data)
	ctx.ResponseWriter.Write(message)
}

func Update(ctx *context.Context) {
	data := struct {
		Name *string
		Num  *int
	}{}
	uuid := (*ctx.Param)["uuid"].(string)
	jsonStr, err := ioutil.ReadAll(ctx.Body)
	if err != nil {
		logger.Errorf(err.Error())
		logger.Errorf(string(debug.Stack()))
		panic("request body read error")
	}
	json.Unmarshal(jsonStr, &data)
	if data.Name == nil || data.Num == nil {
		ctx.WriteHeader(constant.BadRequest)
		return
	}

	update(*data.Name, *data.Num, uuid)
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	result := struct {
		Data string `json:"data"`
	}{"已修改"}
	message, _ := json.Marshal(result)
	ctx.ResponseWriter.Write(message)
}

func Delete(ctx *context.Context) {
	uuid := (*ctx.Param)["uuid"].(string)
	if uuid == "" {
		ctx.WriteHeader(constant.BadRequest)
		return
	}
	delete(uuid)
	ctx.ResponseWriter.Header().Set("Content-Type", "application/json")
	result := struct {
		Data string `json:"data"`
	}{"已删除"}
	message, _ := json.Marshal(result)
	ctx.ResponseWriter.Write(message)
}

func Get(ctx *context.Context) {

}
