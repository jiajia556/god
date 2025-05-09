package service

import (
	"fmt"
	"os"
)

type NumInt interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

func OutputFatal(msg ...any) {
	fmt.Println(msg...)
	os.Exit(0)
}

func InputNum[T comparable](msg string) (num T, err error) {
	fmt.Println(msg)
	_, err = fmt.Scanf("%d", &num)
	if err != nil {
		return
	}
	return
}

func InputStr(msg string) (str string, err error) {
	fmt.Println(msg)
	_, err = fmt.Scanf("%s", &str)
	if err != nil {
		return
	}
	return
}
