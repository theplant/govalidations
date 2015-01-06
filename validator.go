package govalidations

import (
	"regexp"
	"strings"
)

type ValueGetter func(object interface{}) interface{}

type Validator func(object interface{}) []*Error

func MessageSwitcher(vd func(object interface{}) string, name string) Validator {
	return func(object interface{}) (r []*Error) {
		message := vd(object)
		if message == "" {
			return
		}
		r = append(r, &Error{
			Name:    name,
			Message: message,
		})
		return
	}
}

func Custom(vd func(object interface{}) bool, name string, message string) Validator {
	return func(object interface{}) (r []*Error) {
		if vd(object) {
			return
		}

		r = append(r, &Error{
			Name:    name,
			Message: message,
		})
		return
	}
}

func DynamicMessage(vd func(object interface{}) (b bool, name string, message string)) Validator {
	return func(object interface{}) (r []*Error) {
		b, n, m := vd(object)
		if b {
			return
		}

		r = append(r, &Error{
			Name:    n,
			Message: m,
		})
		return
	}
}

func Regexp(vg ValueGetter, matcher *regexp.Regexp, name string, message string) Validator {
	return Custom(func(object interface{}) bool {
		value := vg(object).(string)
		return matcher.MatchString(value)
	}, name, message)
}

func Presence(vg ValueGetter, name string, message string) Validator {
	return Custom(func(object interface{}) bool {
		value := vg(object).(string)
		return strings.TrimSpace(value) != ""
	}, name, message)
}

func Limitation(vg ValueGetter, min int, max int, name string, message string) Validator {
	return Custom(func(object interface{}) bool {
		value := vg(object).(string)
		return len(value) >= min && len(value) <= max
	}, name, message)
}

func Prohibition(vg ValueGetter, min int, max int, name string, message string) Validator {
	return Custom(func(object interface{}) bool {
		value := vg(object).(string)
		return len(value) < min || len(value) > max
	}, name, message)
}

var forbiddenStrings []string = []string{
	`<script`,
	`</script>`,
	`<style`,
	`</style>`,
	`<img`,
	`</img>`,
	`<embed`,
	`</embed>`,
	`<object`,
	`</object>`,
	`<video`,
	`</video>`,
	`<audio`,
	`</audio>`,
	`<source`,
	`</source>`,
	`<track`,
	`</track>`,
	`<iframe`,
	`</iframe>`,
	`<frame`,
	`</frame>`,
	`<input`,
	`</input>`,
	`<base`,
	`</base>`,
	`<applet`,
	`</applet>`,
	`<link`,
	`</link>`,
}

func AvoidScriptTag(vg ValueGetter, name string, message string) Validator {
	return Custom(func(object interface{}) bool {
		value := strings.TrimSpace(vg(object).(string))
		if value == "" {
			return true
		}

		for _, str := range forbiddenStrings {
			if strings.Contains(strings.ToLower(value), str) {
				return false
			}
		}

		htmlTagRegexp := regexp.MustCompile(`<\/?\w+[^>]*>`)
		match := htmlTagRegexp.MatchString(value)
		return !match
	}, name, message)
}
