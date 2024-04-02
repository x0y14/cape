package parse

import (
	"fmt"
	"strconv"
	"unicode"
)

var userInput []rune
var curtPos int
var compositeOpSymbols []string
var singleOpSymbols []string

func init() {
	singleOpSymbols = []string{
		"(", ")", "[", "]", "{", "}",
		".", ",", ":", ";",
		"+", "-", "*", "/", "%",
		">", "<",
		"=", "!",
	}
	compositeOpSymbols = []string{
		"==", "!=", ">=", "<=",
		//"+=", "-=", "*=", "/=", "%=",
		"&&", "||",
	}
}

func startWith(s string) bool {
	runes := []rune(s)
	for i := 0; i < len(runes); i++ {
		if len(userInput) <= curtPos+i || userInput[curtPos+i] != runes[i] {
			return false
		}
	}
	return true
}

func isIdentRune(r rune) bool {
	return ('a' <= r && r <= 'z') ||
		('A' <= r && r <= 'Z') ||
		('0' <= r && r <= '9') ||
		('_' == r || '.' == r)
}

func isEof() bool {
	return curtPos >= len(userInput)
}

func consumeComment() string {
	curtPos += 2
	var s string
	for !isEof() {
		if userInput[curtPos] == '\n' {
			break
		}
		s += string(userInput[curtPos])
		curtPos++
	}
	return s
}

func consumeIdent() string {
	var s string
	for !isEof() {
		if !isIdentRune(userInput[curtPos]) {
			break
		}
		s += string(userInput[curtPos])
		curtPos++
	}
	return s
}

func consumeString() string {
	var s string
	// "
	curtPos++

	for !isEof() {
		if userInput[curtPos] == '"' {
			break
		}
		// escaped double quotation
		if userInput[curtPos] == '\\' && userInput[curtPos+1] == '"' {
			s += "\""
			curtPos += 2
			continue
		}
		// newline
		if userInput[curtPos] == '\\' && userInput[curtPos+1] == 'n' {
			s += "\n"
			curtPos += 2
			continue
		}
		// tab
		if userInput[curtPos] == '\\' && userInput[curtPos+1] == 't' {
			s += "\t"
			curtPos += 2
			continue
		}
		// escaped single quotation
		if userInput[curtPos] == '\\' && userInput[curtPos+1] == '\'' {
			s += "'"
			curtPos += 2
			continue
		}
		// escaped? slash
		if userInput[curtPos] == '\\' && userInput[curtPos+1] == '\\' {
			s += "\\"
			curtPos += 2
			continue
		}

		s += string(userInput[curtPos])
		curtPos++
	}

	// "
	curtPos++
	return s
}

func consumeNumber() (string, bool) {
	isFloat := false
	var s string
	for !isEof() {
		if unicode.IsDigit(userInput[curtPos]) {
			s += string(userInput[curtPos])
			curtPos++
			continue
		} else if userInput[curtPos] == '.' {
			// ポイントの次が、数字じゃなければ強制終了
			if len(userInput) <= curtPos+1 ||
				!unicode.IsDigit(userInput[curtPos+1]) {
				break
			}
			s += string(userInput[curtPos])
			curtPos++
			isFloat = true
			continue
		} else {
			break
		}
	}
	return s, isFloat
}

func consumeWhite() string {
	var s string
	for !isEof() {
		if userInput[curtPos] == ' ' || userInput[curtPos] == '\t' {
			s += string(userInput[curtPos])
			curtPos++
		} else {
			break
		}
	}
	return s
}

func tokenize(input string) (*Token, error) {
	userInput = []rune(input)
	curtPos = 0
	var head Token
	cur := &head
Loop:
	for !isEof() {
		//log.Printf("%v\n", string(userInput[curtPos]))
		// white
		if userInput[curtPos] == ' ' || userInput[curtPos] == '\t' {
			_ = consumeWhite()
			continue
		}

		// newline
		if userInput[curtPos] == '\n' || userInput[curtPos] == '\r' {
			curtPos++
			continue
		}

		// comment
		if userInput[curtPos] == '/' && userInput[curtPos+1] == '/' {
			_ = consumeComment()
			continue
		}

		// symbol
		for _, r := range append(compositeOpSymbols, singleOpSymbols...) {
			if startWith(r) {
				tok := newSymbolToken(r)
				cur.Next = tok
				cur = tok
				curtPos += len(r)
				continue Loop
			}
		}

		if isIdentRune(userInput[curtPos]) && !unicode.IsDigit(userInput[curtPos]) {
			id := consumeIdent()
			tok := newToken(Ident, id, 0, 0)
			cur.Next = tok
			cur = tok
			continue
		}

		// string
		if userInput[curtPos] == '"' {
			s := consumeString()
			tok := newLiteralToken(s)
			cur.Next = tok
			cur = tok
			continue
		}

		// number
		if unicode.IsDigit(userInput[curtPos]) {
			numS, isFloat := consumeNumber()
			if isFloat {
				n, err := strconv.ParseFloat(numS, 64)
				if err != nil {
					return nil, err
				}
				tok := newLiteralToken(n)
				cur.Next = tok
				cur = tok
				continue
			} else {
				n, err := strconv.ParseInt(numS, 10, 0)
				if err != nil {
					return nil, err
				}
				tok := newLiteralToken(int(n))
				cur.Next = tok
				cur = tok
				continue
			}
		}
		return nil, fmt.Errorf("unexpected charactor: %v", userInput[curtPos])
	}
	tok := newToken(Eof, "", 0, 0)
	cur.Next = tok
	cur = tok
	return head.Next, nil
}
