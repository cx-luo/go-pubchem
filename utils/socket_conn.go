// Package utils coding=utf-8
// @Project : go-pubchem
// @Time    : 2024/1/11 11:03
// @Author  : chengxiang.luo
// @Email   : chengxiang.luo@foxmail.com
// @File    : socket_conn.go
// @Software: GoLand
package utils

import (
	"errors"
	"github.com/go-resty/resty/v2"
)

func SendDataByApi(data interface{}, apiUrl string) ([]byte, error) {
	client := resty.New()
	resp, err := client.R().
		SetBody(data).
		SetHeader("Content-Type", "application/json").
		Post(apiUrl)

	if err != nil || resp.StatusCode() != 200 {
		return resp.Body(), errors.New(resp.Status())
	}
	return resp.Body(), nil
}
