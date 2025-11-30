package handler

import (
	"errors"
	"strings"
)

func GetMetricName(path string) (string, error) {
	params := strings.Split(path, "/")
	if len(params) == 0 || params[0] == "" {
		return "", errors.New("there is no metric name")
	}
	return params[0], nil
}
