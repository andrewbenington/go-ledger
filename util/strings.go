package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/unicode/norm"
)

var (
	monthDayPattern    = regexp.MustCompile(`(\d\d)/(\d\d)`)
	moneyPrefixPattern = regexp.MustCompile(` *\+ *| *\$ *`)
)

func NormalizeUnicode(s string) string {
	return norm.NFC.String(s)
}

func ExtractDateFromTitle(title string) (month int, day int) {
	submatches := monthDayPattern.FindStringSubmatch(title)
	if len(submatches) == 3 {
		month, _ = strconv.Atoi(submatches[1])
		day, _ = strconv.Atoi(submatches[2])
	}
	return
}

func ParseMoneyAmount(s string) (float64, error) {
	trimmed := moneyPrefixPattern.ReplaceAllString(s, "")
	trimmed = strings.TrimSpace(trimmed)
	if trimmed == "" {
		return 0, nil
	}
	f, err := strconv.ParseFloat(trimmed, 64)
	if err != nil {
		fmt.Println(trimmed)
	}
	return f, err
}
