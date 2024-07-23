package utils

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"reflect"
)

func PrintStruct(s interface{}) {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	if t.Kind() == reflect.Struct {
		msg := ""
		for i := 0; i < v.NumField(); i++ {
			field := t.Field(i)
			value := v.Field(i).Interface()
			msg += fmt.Sprintf("%s: %v\n", field.Name, value)
		}
		log.Info().Msgf(msg)
	} else {
		log.Info().Msgf("Provided value is not a struct")
	}
}
