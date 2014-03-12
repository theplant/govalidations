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

func Regexp(vg ValueGetter, matcher *regexp.Regexp, name string, message string) Validator {
	return Custom(func(object interface{}) bool {
		value := vg(object).(string)
		return matcher.MatchString(value)
	}, name, message)
}

func Presence(vg ValueGetter, name string, message string) Validator {
	return Custom(func(object interface{}) bool {
		value := vg(object).(string)
		return strings.Trim(value, " 　") != ""
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
	`<script>`,
	`<script`,
	`</script>`,
	`<style>`,
	`<style`,
	`</style>`,
	`<img>`,
	`<img`,
	`</img>`,
	`<embed>`,
	`<embed`,
	`</embed>`,
}

func AvoidNormalHtmlTag(vg ValueGetter, name string, message string) Validator {
	return Custom(func(object interface{}) bool {
		value := vg(object).(string)
		for _, str := range forbiddenStrings {
			if strings.Contains(strings.ToLower(value), str) {
				return false
			}
		}
		return true
	}, name, message)
}

func AvoidScriptTag(vg ValueGetter, name string, message string) Validator {
	return Custom(func(object interface{}) bool {
		value := vg(object).(string)
		htmlTagRegexp := regexp.MustCompile(`(<[^<>]+>)`)
		match := htmlTagRegexp.MatchString(value)
		return !match
	}, name, message)
}
