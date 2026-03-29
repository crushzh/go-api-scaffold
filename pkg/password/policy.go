package password

import (
	"fmt"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// Policy 密码策略
type Policy struct {
	MinLength      int
	RequireUpper   bool
	RequireNumber  bool
	RequireSpecial bool
}

// DefaultPolicy 默认密码策略（最小长度 8）
func DefaultPolicy() *Policy {
	return &Policy{MinLength: 8}
}

// Validate 验证密码是否符合策略，返回错误列表
func (p *Policy) Validate(pwd string) []string {
	var errs []string
	minLen := p.MinLength
	if minLen <= 0 {
		minLen = 8
	}
	if len(pwd) < minLen {
		errs = append(errs, fmt.Sprintf("密码长度不能少于 %d 位", minLen))
	}
	if p.RequireUpper {
		has := false
		for _, r := range pwd {
			if unicode.IsUpper(r) {
				has = true
				break
			}
		}
		if !has {
			errs = append(errs, "密码必须包含大写字母")
		}
	}
	if p.RequireNumber {
		has := false
		for _, r := range pwd {
			if unicode.IsDigit(r) {
				has = true
				break
			}
		}
		if !has {
			errs = append(errs, "密码必须包含数字")
		}
	}
	if p.RequireSpecial {
		has := false
		for _, r := range pwd {
			if strings.ContainsRune("!@#$%^&*()_+-=[]{}|;':\",./<>?", r) {
				has = true
				break
			}
		}
		if !has {
			errs = append(errs, "密码必须包含特殊字符")
		}
	}
	return errs
}

// ValidateError 将验证结果合并为单一错误
func (p *Policy) ValidateError(pwd string) error {
	errs := p.Validate(pwd)
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("%s", strings.Join(errs, "；"))
}

// Hash 使用 bcrypt 加密密码
func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

// Verify 验证密码是否匹配
func Verify(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
