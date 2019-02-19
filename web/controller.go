package web

import (
	"encoding/json"
	"gonote/framework/context"
	"gonote/framework/logger"
	"io/ioutil"
	"runtime/debug"
)

type Note struct {
	Data []string
	Dict map[string]string
}

func Create(ctx *context.Context) {
	jsonStr, err := ioutil.ReadAll(ctx.Body)
	if err != nil {
		ctx.ResponseWriter.Write([]byte("an error occurred"))
		ctx.ResponseWriter.WriteHeader(500)
		logger.Errorf(err.Error())
		logger.Errorf(string(debug.Stack()))
		return
	}
	var note Note
	json.Unmarshal(jsonStr, &note)
	note.Dict[""] = "1"
}
