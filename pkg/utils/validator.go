package utils

import (
	"unicode"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func validatePasswordComplexity(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	// 基础长度检查 (也可以直接用 tag 写 min=8,max=20，这里做双重保险或逻辑整合)
	if len(password) < 8 || len(password) > 20 {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// 统计满足的种类数量
	count := 0
	if hasUpper {
		count++
	}
	if hasLower {
		count++
	}
	if hasDigit {
		count++
	}
	if hasSpecial {
		count++
	}
	return count >= 3
}

func init() {
	// 获取Gin 使用的 validator 引擎实例
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// 注册名为 "password_complexity" 的规则，对应上面的函数
		// 第一个参数是标签名，第二个参数是函数
		v.RegisterValidation("password_complexity", validatePasswordComplexity)
	}
}
