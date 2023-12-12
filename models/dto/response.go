package dto

import (
	"encoding/json"
	"log"
	"notes/pkg"
)

type Response struct {
	Result string          `json:"result"`
	Data   json.RawMessage `json:"data,omitempty"`
	Error  string          `json:"error,omitempty"`
}

func (resp *Response) GetJson() (byteResp []byte, err error) {
	myErr := pkg.NewMyError("package pkg: func GetJson()")
	if resp.Data == nil {
		resp.Data = json.RawMessage(`{}`)
	}
	byteResp, err = json.Marshal(resp)
	if err != nil {
		e := myErr.Wrap(err, "")
		log.Println(e.Error())
		return byteResp, e
	}
	return byteResp, nil
}
