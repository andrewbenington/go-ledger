package util

import (
	"encoding/csv"
	"fmt"
	"hash/fnv"
	"regexp"
	"strconv"
	"strings"

	color "github.com/PerformLine/go-stockutil/colorutil"

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
		fmt.Println(err, trimmed)
	}
	return f, err
}

func SplitWordsIgnoreQuotes(s string, delim rune) ([]string, error) {
	r := csv.NewReader(strings.NewReader(s))
	r.Comma = delim
	words, err := r.Read()
	if err != nil {
		return nil, err
	}
	return words, nil
}

func StringToRGB(s string) string {
	hash := fnv.New32a()
	hash.Write([]byte(s))
	hashCode := float64(hash.Sum32())
	hue := hashCode / (2 ^ 32) * 360
	red, green, blue := color.HslToRgb(hue, 0.8, 0.4)

	return fmt.Sprintf("%02X%02X%02X", red, green, blue)
}
