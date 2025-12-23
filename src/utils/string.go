package utils

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"strings"
	"unicode"
)

var S_CHARS = [62]string{
	"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z",
	"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
}

func RandomAlphabeticString(length int) string {
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteString(
			S_CHARS[GetRandomInt(0, len(S_CHARS)-26)],
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

func ToPascalCase(s string) string {
	s = strings.ToUpper(s[0:1]) + s[1:]
	var i = 0
	for i < len(s) {
		shouldSplit := false
		if s[i] == '_' {
			shouldSplit = true
		} else if s[i] == ' ' {
			shouldSplit = true
		} else if s[i] == '-' {
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

func ToCamelCase(s string) string {
	s = strings.ToLower(s[0:1]) + s[1:]
	var i = 0
	for i < len(s) {
		shouldSplit := false
		if s[i] == '_' {
			shouldSplit = true
		} else if s[i] == ' ' {
			shouldSplit = true
		} else if s[i] == '-' {
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
func StringRef(s string) *string {
	return &s
}
