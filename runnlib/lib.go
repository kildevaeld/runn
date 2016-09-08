package runnlib

import "github.com/asaskevich/govalidator"

func init() {
	govalidator.TagMap["notempty"] = govalidator.Validator(func(str string) bool {
		return len(str) > 0
	})
}
