package jio

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	REmail      string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	RCreditCard string = "^(?:4[0-9]{12}(?:[0-9]{3})?|5[1-5][0-9]{14}|(222[1-9]|22[3-9][0-9]|2[3-6][0-9]{2}|27[01][0-9]|2720)[0-9]{12}|6(?:011|5[0-9][0-9])[0-9]{12}|3[47][0-9]{13}|3(?:0[0-5]|[68][0-9])[0-9]{11}|(?:2131|1800|35\\d{3})\\d{11}|6[27][0-9]{14})$"
)

// String Generates a schema object that matches string data type
func String() *StringSchema {
	return &StringSchema{
		rules: make([]func(*Context), 0, 3),
	}
}

var _ Schema = new(StringSchema)

// StringSchema match string data type
type StringSchema struct {
	baseSchema

	required *bool
	rules    []func(*Context)
}

// SetPriority same as AnySchema.SetPriority
func (s *StringSchema) SetPriority(priority int) *StringSchema {
	s.priority = priority
	return s
}

// PrependTransform same as AnySchema.PrependTransform
func (s *StringSchema) PrependTransform(f func(*Context)) *StringSchema {
	s.rules = append([]func(*Context){f}, s.rules...)
	return s
}

// Transform same as AnySchema.Transform
func (s *StringSchema) Transform(f func(*Context)) *StringSchema {
	s.rules = append(s.rules, f)
	return s
}

// Required same as AnySchema.Required
func (s *StringSchema) Required() *StringSchema {
	s.required = boolPtr(true)
	return s.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Abort(fmt.Errorf("field `%s` is required", ctx.FieldPath()))
		}
	})
}

// Optional same as AnySchema.Optional
func (s *StringSchema) Optional() *StringSchema {
	s.required = boolPtr(false)
	return s.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Skip()
		}
	})
}

// Default same as AnySchema.Default
func (s *StringSchema) Default(value string) *StringSchema {
	s.required = boolPtr(false)
	return s.PrependTransform(func(ctx *Context) {
		if ctx.Value == nil {
			ctx.Value = value
		}
	})
}

// Set same as AnySchema.Set
func (s *StringSchema) Set(value string) *StringSchema {
	return s.Transform(func(ctx *Context) {
		ctx.Value = value
	})
}

// Equal same as AnySchema.Equal
func (s *StringSchema) Equal(value string) *StringSchema {
	return s.Check(func(ctxValue string) error {
		if value != ctxValue {
			return fmt.Errorf("is not %v", value)
		}
		return nil
	})
}

// When same as AnySchema.When
func (s *StringSchema) When(refPath string, condition interface{}, then Schema) *StringSchema {
	return s.Transform(func(ctx *Context) { s.when(ctx, refPath, condition, then) })
}

// Check use the provided function to validate the value of the key.
// Throws an error when the value is not string.
func (s *StringSchema) Check(f func(string) error) *StringSchema {
	return s.Transform(func(ctx *Context) {
		ctxValue, ok := ctx.Value.(string)
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not string", ctx.FieldPath(), ctx.Value))
			return
		}
		if err := f(ctxValue); err != nil {
			ctx.Abort(fmt.Errorf("field `%s` value %v %s", ctx.FieldPath(), ctx.Value, err.Error()))
		}
	})
}

// Valid same as AnySchema.Valid
func (s *StringSchema) Valid(values ...string) *StringSchema {
	return s.Check(func(ctxValue string) error {
		var isValid bool
		for _, v := range values {
			if v == ctxValue {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("not in %v", values)
		}
		return nil
	})
}

// Min check if the length of this string is greater than or equal to the provided length.
func (s *StringSchema) Min(min int) *StringSchema {
	return s.Check(func(ctxValue string) error {
		if len(ctxValue) < min {
			return fmt.Errorf("length less than %d", min)
		}
		return nil
	})
}

// Max check if the length of this string is less than or equal to the provided length.
func (s *StringSchema) Max(max int) *StringSchema {
	return s.Check(func(ctxValue string) error {
		if len(ctxValue) > max {
			return fmt.Errorf("length exceeded %d", max)
		}
		return nil
	})
}

// Length check if the length of this string is equal to the provided length.
func (s *StringSchema) Length(length int) *StringSchema {
	return s.Check(func(ctxValue string) error {
		if len(ctxValue) != length {
			return fmt.Errorf("length not equal to %d", length)
		}
		return nil
	})
}

// Regex check if the value is matched the regex.
func (s *StringSchema) Regex(regex string) *StringSchema {
	re := regexp.MustCompile(regex)
	return s.Check(func(ctxValue string) error {
		if !re.MatchString(ctxValue) {
			return fmt.Errorf("not match with %s", regex)
		}
		return nil
	})
}

// Alphanum check if the string value to only contain a-z, A-Z, and 0-9
func (s *StringSchema) Alphanum() *StringSchema {
	return s.Regex(`^[a-zA-Z0-9]+$`)
}

// Email check if the string value is a valid email format
func (s *StringSchema) Email() *StringSchema {
	return s.Check(func(ctxValue string) error {

		re := regexp.MustCompile(REmail)

		if !re.MatchString(ctxValue) {
			return fmt.Errorf("is not a valid email")
		}

		return nil
	})
}

// CreditCard check if the string value is a valid credit card format
func (s *StringSchema) CreditCard() *StringSchema {
	return s.Check(func(ctxValue string) error {

		re := regexp.MustCompile(RCreditCard)

		if !re.MatchString(ctxValue) {
			return fmt.Errorf("is not a valid credit card number")
		}

		return nil
	})
}

// Token check if the string value to only contain a-z, A-Z, 0-9, and underscore _
func (s *StringSchema) Token() *StringSchema {
	return s.Regex(`^\w+$`)
}

// Convert use the provided function to convert the value of the key.
// Throws an error when the value is not string.
func (s *StringSchema) Convert(f func(string) string) *StringSchema {
	return s.Transform(func(ctx *Context) {
		ctxValue, ok := ctx.Value.(string)
		if !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not string", ctx.FieldPath(), ctx.Value))
			return
		}
		ctx.Value = f(ctxValue)
	})
}

// Lowercase convert the string value to lowercase.
func (s *StringSchema) Lowercase() *StringSchema {
	return s.Convert(strings.ToLower)

}

// Uppercase convert the string value to uppercase.
func (s *StringSchema) Uppercase() *StringSchema {
	return s.Convert(strings.ToUpper)
}

// Trim  emoves whitespace from both sides of the string value.
func (s *StringSchema) Trim() *StringSchema {
	return s.Convert(strings.TrimSpace)
}

// Validate same as AnySchema.Validate
func (s *StringSchema) Validate(ctx *Context) {
	if s.required == nil {
		s.Optional()
	}
	for _, rule := range s.rules {
		rule(ctx)
		if ctx.skip {
			return
		}
	}
	if ctx.Err == nil {
		if _, ok := (ctx.Value).(string); !ok {
			ctx.Abort(fmt.Errorf("field `%s` value %v is not string", ctx.FieldPath(), ctx.Value))
		}
	}
}
