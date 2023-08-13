package scanner

import (
	. "github.com/keatonmartin/golox/token"
	"strconv"
)

type parseError struct {
	Line    int
	Message string
}

type Scanner struct {
	source  []byte
	tokens  []Token
	start   int // start is the beginning of the current lexeme being scanned (index into source)
	current int // current is the index of the current char being considered (index into source)
	line    int // line in source being scanned
	Errs    []parseError
}

// NewScanner returns a new scanner from source
func NewScanner(source []byte) Scanner {
	return Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
		Errs:    []parseError{},
	}
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) addToken(t Type, literal interface{}) {
	s.tokens = append(s.tokens, NewToken(t, s.source[s.start:s.current], literal, s.line))
}

func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	return s.tokens
}

// advance consumes the current character
func (s *Scanner) advance() uint8 {
	ret := s.source[s.current]
	s.current++
	return ret
}

// match consumes the current character if the current character is expected,
// otherwise leaving the current character alone
func (s *Scanner) match(expected uint8) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

// peek returns the current character without consuming it
func (s *Scanner) peek() uint8 {
	if s.isAtEnd() {
		return '\000'
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() uint8 {
	if s.current+1 >= len(s.source) {
		return '\000'
	}
	return s.source[s.current+1]
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case ',':
		s.addToken(COMMA, nil)
	case '.':
		s.addToken(DOT, nil)
	case '-':
		s.addToken(MINUS, nil)
	case '+':
		s.addToken(PLUS, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '*':
		s.addToken(STAR, nil)
	case '?':
		s.addToken(QUESTION, nil)
	case ':':
		s.addToken(COLON, nil)
	case '!':
		if s.match('=') {
			s.addToken(BANG_EQUAL, nil)
		} else {
			s.addToken(BANG, nil)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL, nil)
		} else {
			s.addToken(EQUAL, nil)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL, nil)
		} else {
			s.addToken(LESS, nil)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL, nil)
		} else {
			s.addToken(GREATER, nil)
		}
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(SLASH, nil)
		}
	case '"': // parsing string literal
		for s.peek() != '"' && !s.isAtEnd() {
			if s.peek() == '\n' {
				s.line++
			}
			s.advance()
		}
		if s.isAtEnd() {
			s.Errs = append(s.Errs, parseError{
				Line:    s.line,
				Message: "Unterminated string",
			})
		}
		s.advance() // consume closing "
		s.addToken(STRING, s.source[s.start+1:s.current-1])
	case ' ', '\t', '\r':
	case '\n':
		s.line++
	default:
		// parse number literal
		if isDigit(c) {
			for isDigit(s.peek()) { // consume all numbers before fractional part
				s.advance()
			}
			// if a dot is found, expect fractional part
			if s.peek() == '.' && isDigit(s.peekNext()) {
				s.advance()
				for isDigit(s.peek()) {
					s.advance()
				}
			}
			// not handling error, shouldn't expect one.
			num, _ := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
			s.addToken(NUMBER, num)
		} else if isAlpha(c) {
			for isAlphaNumeric(s.peek()) {
				s.advance()
			}
			text := string(s.source[s.start:s.current])
			t, ok := Keywords[text]
			if !ok { // if not a keyword, type is identifier
				t = IDENTIFIER
			}
			s.addToken(t, nil)
		} else {
			s.Errs = append(s.Errs, parseError{
				Line:    s.line,
				Message: "Unexpected character",
			})
		}
	}
}

func isDigit(c uint8) bool {
	return c >= '0' && c <= '9'
}

func isAlpha(c uint8) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c == '_')
}

func isAlphaNumeric(c uint8) bool {
	return isDigit(c) || isAlpha(c)
}
