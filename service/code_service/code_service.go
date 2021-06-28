package code_service

import (
	"club/dao"
	appError "club/model/error"
	"strings"
)

var (
	codemodel dao.Code
)

func Code(types []string) (*map[string][]dao.Code, error) {
	var newCodes map[string][]dao.Code
	newCodes = make(map[string][]dao.Code)

	codes, err := codemodel.GetByTypes(types)
	if err != nil {
		return nil, err
	}
	c := *codes

	if codes != nil {
		var t string
		codeSlice := []dao.Code{}
		t = c[0].Type

		for _, code := range c {
			if code.Type != t {
				newCodes[t] = codeSlice
				t = code.Type
				codeSlice = []dao.Code{}
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
