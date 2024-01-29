package main

import (
	"strconv"
	"strings"
)

func IntSliceToString(slice []int) string {
	var sb strings.Builder
	for _, val := range slice {
		sb.WriteString(strconv.Itoa(val))
		sb.WriteString(" ")
	}
	return strings.TrimRight(sb.String(), " ")
}
