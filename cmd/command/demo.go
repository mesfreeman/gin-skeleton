package command

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	year  int
	month int
)

// dateCmd represents the date command
var dateCmd = &cobra.Command{
	Use:   "date",
	Short: "这是一个日历查询小工具",
	Run: func(cmd *cobra.Command, args []string) {
		// 判断是否合法
		if year < 1000 && year > 9999 {
			_, _ = fmt.Fprintf(os.Stderr, "invalid year should in [1000, 9999], actual:%d\n", year)
			os.Exit(1)
		}
		if month < 1 && year > 12 {
			_, _ = fmt.Fprintf(os.Stderr, "invalid month should in [1, 12], actual:%d\n", month)
		}

		// 调用要执行的程序
		showCalendar()
	},
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(dateCmd)

	// 命令标示
	dateCmd.PersistentFlags().IntVarP(&year, "year", "y", 0, "year to show (should in [1000, 9999]")
	dateCmd.PersistentFlags().IntVarP(&month, "month", "m", 0, "month to show (should in [1, 12]")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// 日历显示
func showCalendar() {
	now := time.Now()
	showYear := year
	if showYear == 0 { // 默认使用今年
		showYear = int(now.Year())
	}

	showMonth := time.Month(month)
	if showMonth == 0 {
		showMonth = now.Month()
	}

	// 打印星期
	weekdays := []string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}
	for _, weekday := range weekdays {
		fmt.Printf("%5s", weekday)
	}
	fmt.Println()

	showTime := time.Date(showYear, showMonth, 1, 0, 0, 0, 0, now.Location())
	for {
		startWd := showTime.Weekday()
		fmt.Printf("%s", strings.Repeat(" ", int(startWd)*5))

		for ; startWd <= time.Saturday; startWd++ {
			fmt.Printf("%5d", showTime.Day())
			showTime = showTime.Add(time.Hour * 24)
			if showTime.Month() != showMonth {
				fmt.Println()
				return
			}
		}
		fmt.Println()
	}
}
