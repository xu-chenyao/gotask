package calculator

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Token 表示表达式中的一个标记
type Token struct {
	Type  TokenType
	Value string
}

// TokenType 标记类型
type TokenType int

const (
	NUMBER TokenType = iota
	PLUS
	MINUS
	MULTIPLY
	DIVIDE
	MODULO
	LPAREN
	RPAREN
	EOF
)

// Lexer 词法分析器
type Lexer struct {
	input    string
	position int
	current  rune
}

// NewLexer 创建新的词法分析器
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input:    strings.ReplaceAll(input, " ", ""), // 移除所有空格
		position: 0,
	}
	if len(l.input) > 0 {
		l.current = rune(l.input[0])
	}
	return l
}

// advance 移动到下一个字符
func (l *Lexer) advance() {
	l.position++
	if l.position >= len(l.input) {
		l.current = 0 // EOF
	} else {
		l.current = rune(l.input[l.position])
	}
}

// skipWhitespace 跳过空白字符
func (l *Lexer) skipWhitespace() {
	for l.current != 0 && unicode.IsSpace(l.current) {
		l.advance()
	}
}

// readNumber 读取数字
func (l *Lexer) readNumber() string {
	start := l.position
	for l.current != 0 && (unicode.IsDigit(l.current) || l.current == '.') {
		l.advance()
	}
	return l.input[start:l.position]
}

// NextToken 获取下一个标记
func (l *Lexer) NextToken() Token {
	for l.current != 0 {
		l.skipWhitespace()

		if unicode.IsDigit(l.current) {
			return Token{NUMBER, l.readNumber()}
		}

		switch l.current {
		case '+':
			l.advance()
			return Token{PLUS, "+"}
		case '-':
			l.advance()
			return Token{MINUS, "-"}
		case '*':
			l.advance()
			return Token{MULTIPLY, "*"}
		case '/':
			l.advance()
			return Token{DIVIDE, "/"}
		case '%':
			l.advance()
			return Token{MODULO, "%"}
		case '(':
			l.advance()
			return Token{LPAREN, "("}
		case ')':
			l.advance()
			return Token{RPAREN, ")"}
		default:
			return Token{EOF, ""}
		}
	}
	return Token{EOF, ""}
}

// Parser 语法分析器
type Parser struct {
	lexer        *Lexer
	currentToken Token
}

// NewParser 创建新的语法分析器
func NewParser(lexer *Lexer) *Parser {
	p := &Parser{lexer: lexer}
	p.currentToken = p.lexer.NextToken()
	return p
}

// eat 消费当前标记
func (p *Parser) eat(tokenType TokenType) error {
	if p.currentToken.Type == tokenType {
		p.currentToken = p.lexer.NextToken()
		return nil
	}
	return fmt.Errorf("期望标记类型 %v，但得到 %v", tokenType, p.currentToken.Type)
}

// factor 解析因子（数字或括号表达式）
func (p *Parser) factor() (float64, error) {
	token := p.currentToken

	if token.Type == NUMBER {
		if err := p.eat(NUMBER); err != nil {
			return 0, err
		}
		value, err := strconv.ParseFloat(token.Value, 64)
		if err != nil {
			return 0, fmt.Errorf("无法解析数字: %s", token.Value)
		}
		return value, nil
	}

	if token.Type == LPAREN {
		if err := p.eat(LPAREN); err != nil {
			return 0, err
		}
		result, err := p.expr()
		if err != nil {
			return 0, err
		}
		if err := p.eat(RPAREN); err != nil {
			return 0, err
		}
		return result, nil
	}

	if token.Type == MINUS {
		if err := p.eat(MINUS); err != nil {
			return 0, err
		}
		result, err := p.factor()
		if err != nil {
			return 0, err
		}
		return -result, nil
	}

	if token.Type == PLUS {
		if err := p.eat(PLUS); err != nil {
			return 0, err
		}
		return p.factor()
	}

	return 0, fmt.Errorf("意外的标记: %s", token.Value)
}

// term 解析项（处理 *, /, % 运算符）
func (p *Parser) term() (float64, error) {
	result, err := p.factor()
	if err != nil {
		return 0, err
	}

	for p.currentToken.Type == MULTIPLY || p.currentToken.Type == DIVIDE || p.currentToken.Type == MODULO {
		token := p.currentToken

		switch token.Type {
		case MULTIPLY:
			if err := p.eat(MULTIPLY); err != nil {
				return 0, err
			}
			right, err := p.factor()
			if err != nil {
				return 0, err
			}
			result = result * right

		case DIVIDE:
			if err := p.eat(DIVIDE); err != nil {
				return 0, err
			}
			right, err := p.factor()
			if err != nil {
				return 0, err
			}
			if right == 0 {
				return 0, fmt.Errorf("除零错误")
			}
			result = result / right

		case MODULO:
			if err := p.eat(MODULO); err != nil {
				return 0, err
			}
			right, err := p.factor()
			if err != nil {
				return 0, err
			}
			if right == 0 {
				return 0, fmt.Errorf("模运算的除数不能为零")
			}
			// Go的浮点数取模运算
			result = float64(int(result) % int(right))
		}
	}

	return result, nil
}

// expr 解析表达式（处理 +, - 运算符）
func (p *Parser) expr() (float64, error) {
	result, err := p.term()
	if err != nil {
		return 0, err
	}

	for p.currentToken.Type == PLUS || p.currentToken.Type == MINUS {
		token := p.currentToken

		switch token.Type {
		case PLUS:
			if err := p.eat(PLUS); err != nil {
				return 0, err
			}
			right, err := p.term()
			if err != nil {
				return 0, err
			}
			result = result + right

		case MINUS:
			if err := p.eat(MINUS); err != nil {
				return 0, err
			}
			right, err := p.term()
			if err != nil {
				return 0, err
			}
			result = result - right
		}
	}

	return result, nil
}

// Calculate 计算数学表达式
func Calculate(expression string) (float64, error) {
	if strings.TrimSpace(expression) == "" {
		return 0, fmt.Errorf("表达式不能为空")
	}

	lexer := NewLexer(expression)
	parser := NewParser(lexer)

	result, err := parser.expr()
	if err != nil {
		return 0, err
	}

	// 检查是否还有未处理的标记
	if parser.currentToken.Type != EOF {
		return 0, fmt.Errorf("表达式解析不完整，剩余: %s", parser.currentToken.Value)
	}

	return result, nil
}

// FormatResult 格式化结果输出
func FormatResult(result float64) string {
	// 如果结果是整数，显示为整数格式
	if result == float64(int64(result)) {
		return fmt.Sprintf("%.0f", result)
	}
	// 否则显示为浮点数，最多6位小数
	return fmt.Sprintf("%.6g", result)
}
