package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"cli_cmd/calculator"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "calc",
		Usage:   "一个简单的命令行计算器",
		Version: "1.0.0",
		Authors: []*cli.Author{
			{Name: "Calculator CLI Team"},
		},

		// 全局标志
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "interactive",
				Aliases: []string{"i"},
				Usage:   "启动交互模式",
				Value:   false,
			},
			&cli.BoolFlag{
				Name:    "verbose",
				Aliases: []string{"V"}, // 改为大写V避免与默认version标志冲突
				Usage:   "显示详细输出",
				Value:   false,
			},
		},

		// 默认动作 - 处理单个表达式或启动交互模式
		Action: func(c *cli.Context) error {
			interactive := c.Bool("interactive")
			verbose := c.Bool("verbose")

			// 获取命令行参数（表达式）
			args := c.Args().Slice()

			if interactive || len(args) == 0 {
				return runInteractiveMode(verbose)
			}

			// 将所有参数连接成一个表达式
			expression := strings.Join(args, "")
			return calculateAndPrint(expression, verbose)
		},

		// 子命令
		Commands: []*cli.Command{
			{
				Name:      "eval",
				Aliases:   []string{"e"},
				Usage:     "计算单个表达式",
				ArgsUsage: "EXPRESSION",
				Action: func(c *cli.Context) error {
					verbose := c.Bool("verbose") // 从全局标志获取
					args := c.Args().Slice()

					if len(args) == 0 {
						return fmt.Errorf("请提供一个数学表达式")
					}

					expression := strings.Join(args, "")
					return calculateAndPrint(expression, verbose)
				},
			},
			{
				Name:    "interactive",
				Aliases: []string{"repl", "shell"},
				Usage:   "启动交互式计算器",
				Action: func(c *cli.Context) error {
					verbose := c.Bool("verbose")
					return runInteractiveMode(verbose)
				},
			},
			{
				Name:  "test",
				Usage: "运行内置测试用例",
				Action: func(c *cli.Context) error {
					return runBuiltinTests()
				},
			},
		},

		// 使用错误处理
		OnUsageError: func(c *cli.Context, err error, isSubcommand bool) error {
			fmt.Fprintf(c.App.ErrWriter, "错误: %v\n\n", err)
			if !isSubcommand {
				cli.ShowAppHelp(c)
			}
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "应用程序错误: %v\n", err)
		os.Exit(1)
	}
}

// calculateAndPrint 计算表达式并打印结果
func calculateAndPrint(expression string, verbose bool) error {
	if verbose {
		fmt.Printf("正在计算表达式: %s\n", expression)
	}

	result, err := calculator.Calculate(expression)
	if err != nil {
		return fmt.Errorf("计算错误: %v", err)
	}

	if verbose {
		fmt.Printf("表达式: %s\n", expression)
		fmt.Printf("结果: %s\n", calculator.FormatResult(result))
	} else {
		fmt.Println(calculator.FormatResult(result))
	}

	return nil
}

// runInteractiveMode 运行交互模式
func runInteractiveMode(verbose bool) error {
	fmt.Println("🧮 欢迎使用命令行计算器!")
	fmt.Println("支持的操作: +, -, *, /, %, ()")
	fmt.Println("输入 'help' 查看帮助，'exit' 或 'quit' 退出")
	fmt.Println("示例: 1+2*3, (1+2)*3, 10%3")
	fmt.Println(strings.Repeat("-", 50))

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("calc> ")

		input, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("读取输入失败: %v", err)
		}

		input = strings.TrimSpace(input)

		// 处理特殊命令
		switch strings.ToLower(input) {
		case "exit", "quit", "q":
			fmt.Println("再见! 👋")
			return nil
		case "help", "h":
			printHelp()
			continue
		case "clear", "cls":
			// 清屏
			fmt.Print("\033[H\033[2J")
			continue
		case "":
			continue
		}

		// 计算表达式
		if verbose {
			fmt.Printf("正在计算: %s\n", input)
		}

		result, err := calculator.Calculate(input)
		if err != nil {
			fmt.Printf("❌ 错误: %v\n", err)
			continue
		}

		formattedResult := calculator.FormatResult(result)
		fmt.Printf("📊 %s = %s\n", input, formattedResult)
	}
}

// printHelp 打印帮助信息
func printHelp() {
	fmt.Println("\n📖 计算器帮助:")
	fmt.Println("  支持的运算符:")
	fmt.Println("    +  : 加法 (例: 1+2)")
	fmt.Println("    -  : 减法 (例: 5-3)")
	fmt.Println("    *  : 乘法 (例: 2*3)")
	fmt.Println("    /  : 除法 (例: 8/2)")
	fmt.Println("    %  : 取模 (例: 10%3)")
	fmt.Println("    () : 括号 (例: (1+2)*3)")
	fmt.Println("\n  特殊命令:")
	fmt.Println("    help  : 显示此帮助")
	fmt.Println("    clear : 清屏")
	fmt.Println("    exit  : 退出程序")
	fmt.Println("\n  运算优先级:")
	fmt.Println("    1. 括号 ()")
	fmt.Println("    2. 乘法 * 除法 / 取模 %")
	fmt.Println("    3. 加法 + 减法 -")
	fmt.Println()
}

// runBuiltinTests 运行内置测试
func runBuiltinTests() error {
	fmt.Println("🧪 运行内置测试用例...")

	testCases := []struct {
		expression string
		expected   string
		desc       string
	}{
		{"1+2", "3", "基本加法"},
		{"1+2*2", "5", "运算优先级"},
		{"2*3+1", "7", "乘法优先级"},
		{"(1+2)*3", "9", "括号优先级"},
		{"10/2", "5", "除法"},
		{"10%3", "1", "取模运算"},
		{"5-3", "2", "减法"},
		{"-5+3", "-2", "负数"},
		{"2*3*4", "24", "连续乘法"},
		{"100/10/2", "5", "连续除法"},
		{"1+2+3+4", "10", "连续加法"},
		{"10-3-2", "5", "连续减法"},
		{"2*(3+4)", "14", "括号内加法"},
		{"(10-6)/2", "2", "括号内减法"},
		{"3.5+1.5", "5", "小数计算"},
		{"7.5/2.5", "3", "小数除法"},
	}

	passed := 0
	total := len(testCases)

	fmt.Printf("总共 %d 个测试用例:\n\n", total)

	for i, tc := range testCases {
		result, err := calculator.Calculate(tc.expression)
		if err != nil {
			fmt.Printf("❌ 测试 %d: %s - 错误: %v\n", i+1, tc.desc, err)
			continue
		}

		formattedResult := calculator.FormatResult(result)
		if formattedResult == tc.expected {
			fmt.Printf("✅ 测试 %d: %s - %s = %s\n", i+1, tc.desc, tc.expression, formattedResult)
			passed++
		} else {
			fmt.Printf("❌ 测试 %d: %s - %s = %s (期望: %s)\n",
				i+1, tc.desc, tc.expression, formattedResult, tc.expected)
		}
	}

	fmt.Printf("\n📊 测试结果: %d/%d 通过", passed, total)
	if passed == total {
		fmt.Println(" 🎉 所有测试通过!")
	} else {
		fmt.Printf(" ⚠️  有 %d 个测试失败\n", total-passed)
	}

	return nil
}
