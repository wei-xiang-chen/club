package code_service

import (
	"club/model"
	appError "club/pojo/error"
	"strings"
)

var (
	codemodel model.Code
)

func Code(types []string) (*map[string][]model.Code, error) {
	var newCodes map[string][]model.Code
	newCodes = make(map[string][]model.Code)

	codes, err := codemodel.GetByTypes(types)
	if err != nil {
		return nil, err
	}
	c := *codes

	if codes != nil {
		var t string
		codeSlice := []model.Code{}
		t = c[0].Type

		for _, code := range c {
			if code.Type != t {
				newCodes[t] = codeSlice
				t = code.Type
				codeSlice = []model.Code{}
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
