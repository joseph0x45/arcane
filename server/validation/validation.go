package validation

import "regexp"

func IsEmail(str *string) bool {
	re := regexp.MustCompile(`/^((?!\.)[\w-_.]*[^.])(@\w+)(\.\w+(\.\w+)?[^.\W])$/gim;`)
	return re.MatchString(*str)
}
