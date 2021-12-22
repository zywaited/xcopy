package utils

import "strings"

func ToCameWithSep(s, sep string) string {
	bs := strings.Builder{}
	for {
		m := strings.Index(s, sep)
		if m < 0 {
			break
		}
		bs.WriteString(ToUcFirst(s[:m]))
		s = s[m+len(sep):]
	}
	bs.WriteString(ToUcFirst(s))
	return bs.String()
}

// ToCame 下划线转驼峰
func ToCame(s string) string {
	return ToCameWithSep(s, "_")
}

func ToSnakeFWithSep(s, sep string) string {
	bs := strings.Builder{}
	for i := 0; i < len(s); i++ {
		d := s[i]
		up := d >= 'A' && d <= 'Z'
		if up {
			d += 32
		}
		if i > 0 && up {
			bs.WriteString(sep)
		}
		bs.WriteByte(d)
	}
	return bs.String()
}

// ToSnake 驼峰转下划线
func ToSnake(s string) string {
	return ToSnakeFWithSep(s, "_")
}

// ToLcFirst 首字母小写
func ToLcFirst(s string) string {
	if len(s) > 0 && s[0] >= byte('A') && s[0] <= byte('Z') {
		s = string(s[0]+32) + s[1:]
	}
	return s
}

// ToUcFirst 首字母大写
func ToUcFirst(s string) string {
	if len(s) > 0 && s[0] >= byte('a') && s[0] <= byte('z') {
		s = string(s[0]-32) + s[1:]
	}
	return s
}
