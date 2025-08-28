package calculator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestBasicOperations 测试基本运算
func TestBasicOperations(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		desc       string
	}{
		{"1+2", 3, "加法"},
		{"5-3", 2, "减法"},
		{"2*3", 6, "乘法"},
		{"8/2", 4, "除法"},
		{"10%3", 1, "取模"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, err := Calculate(test.expression)
			assert.NoError(t, err, "计算 %s 时不应该出错", test.expression)
			assert.Equal(t, test.expected, result, "表达式 %s 的结果应该是 %f", test.expression, test.expected)
		})
	}
}

// TestOperatorPrecedence 测试运算符优先级
func TestOperatorPrecedence(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		desc       string
	}{
		{"1+2*2", 5, "乘法优先于加法"},
		{"2*3+1", 7, "乘法优先于加法（右侧）"},
		{"10-6/2", 7, "除法优先于减法"},
		{"2+3*4-1", 13, "混合运算优先级"},
		{"10%3+2", 3, "取模优先于加法"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, err := Calculate(test.expression)
			assert.NoError(t, err, "计算 %s 时不应该出错", test.expression)
			assert.Equal(t, test.expected, result, "表达式 %s 的结果应该是 %f", test.expression, test.expected)
		})
	}
}

// TestParentheses 测试括号
func TestParentheses(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		desc       string
	}{
		{"(1+2)*3", 9, "简单括号"},
		{"2*(3+4)", 14, "括号改变优先级"},
		{"(10-6)/2", 2, "括号内减法"},
		{"((1+2)*3)+1", 10, "嵌套括号"},
		{"(5+3)*(2-1)", 8, "两个括号相乘"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, err := Calculate(test.expression)
			assert.NoError(t, err, "计算 %s 时不应该出错", test.expression)
			assert.Equal(t, test.expected, result, "表达式 %s 的结果应该是 %f", test.expression, test.expected)
		})
	}
}

// TestFloatingPoint 测试浮点数运算
func TestFloatingPoint(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		desc       string
	}{
		{"3.5+1.5", 5.0, "小数加法"},
		{"7.5/2.5", 3.0, "小数除法"},
		{"2.5*4", 10.0, "小数乘法"},
		{"10.5-0.5", 10.0, "小数减法"},
		{"0.1+0.2", 0.3, "小数精度"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, err := Calculate(test.expression)
			assert.NoError(t, err, "计算 %s 时不应该出错", test.expression)
			assert.InDelta(t, test.expected, result, 0.0001, "表达式 %s 的结果应该接近 %f", test.expression, test.expected)
		})
	}
}

// TestNegativeNumbers 测试负数
func TestNegativeNumbers(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		desc       string
	}{
		{"-5", -5, "单个负数"},
		{"-5+3", -2, "负数加正数"},
		{"5+-3", 2, "正数加负数"},
		{"-2*3", -6, "负数乘法"},
		{"(-5+3)*2", -4, "括号内负数"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, err := Calculate(test.expression)
			assert.NoError(t, err, "计算 %s 时不应该出错", test.expression)
			assert.Equal(t, test.expected, result, "表达式 %s 的结果应该是 %f", test.expression, test.expected)
		})
	}
}

// TestComplexExpressions 测试复杂表达式
func TestComplexExpressions(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		desc       string
	}{
		{"2*3*4", 24, "连续乘法"},
		{"100/10/2", 5, "连续除法"},
		{"1+2+3+4", 10, "连续加法"},
		{"10-3-2", 5, "连续减法"},
		{"2*3+4*5", 26, "多个乘法表达式"},
		{"(2+3)*(4-1)", 15, "复合括号"},
		{"10*(2+3)/5", 10, "复杂混合运算"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, err := Calculate(test.expression)
			assert.NoError(t, err, "计算 %s 时不应该出错", test.expression)
			assert.Equal(t, test.expected, result, "表达式 %s 的结果应该是 %f", test.expression, test.expected)
		})
	}
}

// TestErrorCases 测试错误情况
func TestErrorCases(t *testing.T) {
	tests := []struct {
		expression string
		desc       string
	}{
		{"", "空表达式"},
		{"1+", "不完整表达式"},
		{"*2", "开头是操作符"},
		{"(1+2", "未闭合括号"},
		{"1+2)", "额外的右括号"},
		{"1/0", "除零错误"},
		{"1%0", "取模零错误"},
		{"abc", "无效字符"},
		{"1+2*", "结尾操作符"},
		{"1*/2", "非法操作符组合"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			_, err := Calculate(test.expression)
			assert.Error(t, err, "表达式 %s 应该产生错误", test.expression)
		})
	}
}

// TestFormatResult 测试结果格式化
func TestFormatResult(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
		desc     string
	}{
		{3.0, "3", "整数应该显示为整数"},
		{3.14159, "3.14159", "小数显示小数部分"},
		{0.0, "0", "零显示为整数"},
		{-5.0, "-5", "负整数"},
		{-3.14, "-3.14", "负小数"},
		{1000000, "1000000", "大整数显示为整数"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result := FormatResult(test.input)
			assert.Equal(t, test.expected, result, "格式化 %f 应该得到 %s", test.input, test.expected)
		})
	}
}

// TestWhitespace 测试空白字符处理
func TestWhitespace(t *testing.T) {
	tests := []struct {
		expression string
		expected   float64
		desc       string
	}{
		{" 1 + 2 ", 3, "空格应该被忽略"},
		{"1+2", 3, "无空格"},
		{"\t1+2\t", 3, "制表符"},
		{"1 * 2 + 3", 5, "混合空格"},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			result, err := Calculate(test.expression)
			assert.NoError(t, err, "计算 %s 时不应该出错", test.expression)
			assert.Equal(t, test.expected, result, "表达式 %s 的结果应该是 %f", test.expression, test.expected)
		})
	}
}

// BenchmarkCalculate 性能基准测试
func BenchmarkCalculate(b *testing.B) {
	expressions := []string{
		"1+2",
		"1+2*3",
		"(1+2)*3+4",
		"10/2+3*4-1",
		"((1+2)*3+4)/5",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		expr := expressions[i%len(expressions)]
		_, _ = Calculate(expr)
	}
}
