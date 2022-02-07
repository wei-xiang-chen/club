package code_service

import (
	"club/dal"
	appError "club/models/error"
	"strings"
)

var (
	codemodel dal.Code
)

func Code(types []string) (*map[string][]dal.Code, error) {
	var newCodes map[string][]dal.Code
	newCodes = make(map[string][]dal.Code)

	codes, err := codemodel.GetByTypes(types)
	if err != nil {
		return nil, err
	}
	c := *codes

	if codes != nil {
		var t string
		codeSlice := []dal.Code{}
		t = c[0].Type

		for _, code := range c {
			if code.Type != t {
				newCodes[t] = codeSlice
				t = code.Type
				codeSlice = []dal.Code{}
			}

			codeSlice = append(codeSlice, code)
		}
		newCodes[t] = codeSlice
	}

	return &newCodes, nil
}

func CheckCode(_type string, option *string) error {
	*option = strings.ToUpper(*option)

	count, err := codemodel.CheckCode(_type, *option)
	if err != nil {
		return err
	}

	if *count > 0 {
		return nil
	} else {
		return appError.AppError{Message: "Code is illegal." + "[" + *option + "]"}
	}
}
