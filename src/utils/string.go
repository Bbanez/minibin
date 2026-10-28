package utils

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode"
)

var N_CHARS = [10]string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
}
var S_CHARS = [62]string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
}
var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
var nameRegex = regexp.MustCompile(`[^\p{L}\s]`)

func RandomAlphabeticString(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteString(
			S_CHARS[GetRandomInt(0, len(S_CHARS)-26)],
		)
	}
	return sb.String()
}

func RandomNumberString(length int) string {
	var sb strings.Builder
	for range length {
		sb.WriteString(
			N_CHARS[GetRandomInt(0, len(N_CHARS)-1)],
		)
	}
	return sb.String()
}

func RandomString(length int) string {
	bytes := make([]byte, (length*6+7)/8)
	rand.Read(bytes)
	return base64.RawURLEncoding.EncodeToString(bytes)[:length]
}

func Sha512(input string) string {
	hash := sha512.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func Sha256(input string) string {
	hash := sha256.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func Sha1(input string) string {
	hash := sha1.New()
	hash.Write([]byte(input))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func HmacSha256(key string, data string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func HmacSha256B(key []byte, data []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func IsValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func ToSlug(s string) string {
	s = strings.ToLower(s)
	s = regexp.MustCompile(`[^\p{L}\d]+`).ReplaceAllString(s, "-")
	s = regexp.MustCompile(`^-+|-+$`).ReplaceAllString(s, "")
	return s
}

func ToCamelCase(s string) string {
	s = strings.ToLower(s[0:1]) + s[1:]
	var i = 0
	for i < len(s) {
		shouldSplit := false
		switch s[i] {
		case '_', ' ', '-':
			shouldSplit = true
		}
		if shouldSplit {
			s = s[:i] + strings.ToUpper(s[i+1:i+2]) + s[i+2:]
			i += 2
		} else {
			i++
		}
	}
	return s
}

func ToPascalCase(s string) string {
	s = strings.ToUpper(s[0:1]) + s[1:]
	var i = 0
	for i < len(s) {
		shouldSplit := false
		switch s[i] {
		case '_', ' ', '-':
			shouldSplit = true
		}
		if shouldSplit {
			s = s[:i] + strings.ToUpper(s[i+1:i+2]) + s[i+2:]
			i += 2
		} else {
			i++
		}
	}
	return s
}

func ToSnakeCase(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	s = regexp.MustCompile(`[^a-z0-9_]+`).ReplaceAllString(s, "")
	return s
}

func ToKebabCase(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "-")
	s = strings.ReplaceAll(s, "_", "-")
	s = regexp.MustCompile(`[^a-z0-9-]+`).ReplaceAllString(s, "")
	return s
}

func ToUpperSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) && i > 0 && (unicode.IsLower(rune(s[i-1])) || unicode.IsDigit(rune(s[i-1]))) {
			result = append(result, '_')
		}
		if r == '-' || r == ' ' {
			result = append(result, '_')
		} else if r != '_' {
			result = append(result, unicode.ToUpper(r))
		} else {
			result = append(result, '_')
		}
	}
	return strings.Trim(strings.ReplaceAll(string(result), "__", "_"), "_")
}

func TextBetween(s string, begin string, end string) *string {
	src := NewString(s)
	startIndex := src.IndexOf(begin, 0)
	if startIndex == -1 {
		return nil
	}
	startIndex += len(begin)
	endIndex := src.IndexOf(end, startIndex)
	if endIndex == -1 {
		return nil
	}
	result := s[startIndex:endIndex]
	return &result
}

func AllTextBetween(s string, begin string, end string) []string {
	src := NewString(s)
	output := []string{}
	index := struct {
		begin int
		end   int
	}{
		begin: src.IndexOf(begin, 0),
		end:   0,
	}
	if index.begin == -1 {
		return output
	}
	index.begin += len(begin)
	index.end = src.IndexOf(end, index.begin)
	if index.end == -1 {
		return output
	}
	output = append(output, s[index.begin:index.end])
	index.end += len(end)
	for true {
		index.begin = src.IndexOf(begin, index.end)
		if index.begin == -1 {
			break
		}
		index.begin += len(begin)
		index.end = src.IndexOf(end, index.begin)
		if index.end == -1 {
			break
		}
		output = append(output, s[index.begin:index.end])
		index.end += len(end)
	}
	return output
}

type StringSlice struct {
	s []string
}

func NewStringSlice(s []string) *StringSlice {
	return &StringSlice{s: s}
}

func (ss *StringSlice) Append(s string) *StringSlice {
	ss.s = append(ss.s, s)
	return ss
}

func (ss *StringSlice) AppendMany(s []string) *StringSlice {
	ss.s = append(ss.s, s...)
	return ss
}

func (ss *StringSlice) Get() *[]string {
	return &ss.s
}

func KeepUniqueStrings(s []string) []string {
	m := map[string]bool{}
	for _, str := range s {
		m[str] = true
	}
	out := []string{}
	for str := range m {
		out = append(out, str)
	}
	return out
}

// ----------------------------------------------------------------------------
// --- String Object
// ----------------------------------------------------------------------------

type String struct {
	*string
}

func NewString(s string) String {
	return String{string: &s}
}

func (s *String) IndexOf(searchStr string, positionOffs int) int {
	str := *s.string
	if positionOffs < 0 || positionOffs >= len(str) {
		return -1
	}
	idx := strings.Index(str[positionOffs:], searchStr)
	if idx == -1 {
		return -1
	}
	return idx + positionOffs
}

func (s *String) Substring(start int, end int) String {
	str := *s.string
	if start < 0 || start >= len(str) {
		return String{string: new(string)}
	}
	if end < 0 || end > len(str) {
		end = len(str)
	}
	substr := str[start:end]
	return String{string: &substr}
}
