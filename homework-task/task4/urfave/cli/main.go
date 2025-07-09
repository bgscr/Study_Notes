package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/urfave/cli/v3"
)

func main() {
	cmd := &cli.Command{
		Commands: []*cli.Command{
			{
				Name:    "+",
				Aliases: []string{"add", "加法", "加"},
				Usage:   "加法：calc + 数字1 数字2",
				Action: func(c context.Context, cmd *cli.Command) error {
					return calc(c, cmd)
				},
			},
			{
				Name:    "-",
				Aliases: []string{"sub", "减法", "减"},
				Usage:   "减法：calc - 数字1 数字2",
				Action: func(c context.Context, cmd *cli.Command) error {
					return calc(c, cmd)
				},
			},
			{
				Name:    "*",
				Aliases: []string{"mul", "乘法", "乘"},
				Usage:   "乘法：calc * 数字1 数字2",
				Action: func(c context.Context, cmd *cli.Command) error {
					return calc(c, cmd)
				},
			},
			{
				Name:    "/",
				Aliases: []string{"div", "除法", "除"},
				Usage:   "除法：calc / 数字1 数字2",
				Action: func(c context.Context, cmd *cli.Command) error {
					return calc(c, cmd)
				},
			},
			{
				Name:    "%",
				Aliases: []string{"mod", "取模", "取余"},
				Usage:   "取模：calc % 整数1 整数2",
				Action: func(c context.Context, cmd *cli.Command) error {
					return calc(c, cmd)
				},
			},
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

// 统一处理运算逻辑
func calc(_ context.Context, cmd *cli.Command) error {
	if cmd.NArg() != 2 {
		return fmt.Errorf("错误：需要两个参数\n示例：calc %v 5 3", cmd.Name)
	}

	x, err := strconv.ParseFloat(cmd.Args().Get(0), 64)
	if err != nil {
		return fmt.Errorf("第一个参数无效：%v", err)
	}
	y, err := strconv.ParseFloat(cmd.Args().Get(1), 64)
	if err != nil {
		return fmt.Errorf("第二个参数无效：%v", err)
	}

	switch cmd.Name {
	case "+":
		fmt.Printf("结果：%f\n", x+y)
	case "-":
		fmt.Printf("结果：%f\n", x-y)
	case "*":
		fmt.Printf("结果：%f\n", x*y)
	case "/":
		if y == 0 {
			return fmt.Errorf("除数不能为零")
		}
		fmt.Printf("结果：%f\n", x/y)
	case "%":
		// 取模需要整数
		a := int(x)
		b := int(y)
		if b == 0 {
			return fmt.Errorf("取模的除数不能为零")
		}
		fmt.Printf("结果：%d\n", a%b)
	default:
		return fmt.Errorf("不支持的操作符")
	}
	return nil
}
