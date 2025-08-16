package main

import (
	"fmt"
	"sort"
	"strconv"
)

func main() {

	// s5 := []int{1, 2, 3, 4, 5}
	// s5 = append(s5[:2], s5[3:]...)
	// fmt.Printf("%v", s5)
	// s5 = s5[:1]
	// fmt.Printf("%v", s5)
	// s5 = s5[:3]
	// fmt.Printf("%v", s5)
	// s5 = s5[1:3]
	// fmt.Printf("%v", s5)

	// a := 2
	// var p *int = &a
	// fmt.Println(a, p, *p, &p)

	//两数之和
	// nums := []int{2, 7, 11, 15}
	// target := 9
	// target := 6
	// nums := []int{3, 2, 4}
	// res := twoSum(nums, target)
	// fmt.Println(res)

	//7、合并区间
	// arr := [][]int{
	// 	{4, 5},
	// 	{1, 4},
	// }
	// merged := merge(arr)
	// fmt.Println(merged)

	//6、删除有序数组中的重复项
	nums := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	size := removeDuplicates(nums)
	fmt.Println(size)

	//5、加一
	// digits := []int{1, 2, 9}
	// plusOne(digits)

	//4、最长公共前缀
	// strs := []string{"flower", "flow", "flight"}
	// str := longestCommonPrefix(strs)
	// fmt.Println(str)

	//3、 有效的括号
	// a := isValid("(())")
	// fmt.Println(a)

	// 2、判断一个数是否是回文数
	// isPalindrome(123321)

	//1、只出现一次的数字
	// nums := []int{4, 1, 2, 1, 2}
	// res := singleNumber(nums)
	// fmt.Println(res)
}

/*
两数之和
给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数
*/
func twoSum(nums []int, target int) []int {
	mp := make(map[int]int, len(nums))

	for i, v := range nums {
		mp[v] = i + 1
	}
	fmt.Println(mp)
	arr := []int{}
	for i, v := range nums {
		a := mp[target-v]
		if a != 0 && i != a-1 {
			arr = append(arr, i)
			arr = append(arr, a-1)
			break
		}
	}
	return arr
}

/*
合并区间：以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。
请你合并所有重叠的区间，并返回一个不重叠的区间数组，该数组需恰好覆盖输入中的所有区间。
可以先对区间数组按照区间的起始位置进行排序，然后使用一个切片来存储合并后的区间，遍历排序后的区间数组，
将当前区间与切片中最后一个区间进行比较，如果有重叠，则合并区间；如果没有重叠，则将当前区间添加到切片中。
*/
func merge(intervals [][]int) [][]int {
	if len(intervals) == 0 {
		return intervals
	}
	// 先按区间起始位置排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	// 初始化结果切片，放入第一个区间
	merged := [][]int{intervals[0]}

	for _, current := range intervals[1:] {
		// 获取结果中最后一个区间
		last := merged[len(merged)-1]
		// 如果当前区间的起始位置小于等于结果中最后区间的结束位置，说明有重叠
		if current[0] <= last[1] {
			// 合并区间：更新结果中最后区间的结束位置为两者较大值
			if current[1] > last[1] {
				last[1] = current[1]
			}
		} else {
			// 没有重叠直接添加进新切片中
			merged = append(merged, current)
		}
	}

	return merged
}

/*
删除有序数组中的重复项：给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。可以使用双指针法，
一个慢指针 i 用于记录不重复元素的位置，一个快指针 j 用于遍历数组，当 nums[i] 与 nums[j] 不相等时，将 nums[j] 赋值给 nums[i + 1]，并将 i 后移一位。
*/
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	//不重复元素的位置
	i := 0
	for j := 1; j < len(nums); j++ {
		if nums[i] != nums[j] {
			i++
			nums[i] = nums[j]
		}
	}
	return i + 1
}

/*
给定一个表示 大整数 的整数数组 digits，其中 digits[i] 是整数的第 i 位数字。这些数字按从左到右，从最高位到最低位排列。这个大整数不包含任何前导 0。
将大整数加 1，并返回结果的数字数组。
*/
func plusOne(digits []int) []int {
	result := 0
	for _, v := range digits {
		result = result*10 + v
	}
	fmt.Println(result)
	result += 1
	fmt.Println(result)

	str := strconv.Itoa(result)
	fmt.Println(str)
	num := []int{}
	for _, v := range str {
		a := string(v)
		i, _ := strconv.Atoi(a)
		num = append(num, i)
	}
	fmt.Println(num)
	return num
}

/*
编写一个函数来查找字符串数组中的最长公共前缀。
如果不存在公共前缀，返回空字符串 ""。
*/
func longestCommonPrefix(strs []string) string {
	fmt.Println(strs)
	if len(strs) == 0 {
		return ""
	}
	first := strs[0]
	fmt.Println(first)
	for i := 0; i < len(first); i++ {
		char := first[i]
		// fmt.Println(reflect.TypeOf(char))
		fmt.Println(string(char))
		for j := 0; j < len(strs); j++ {
			// 如果其余字符串长度不足或字符不匹配
			fmt.Println(strs[j][i])
			if i >= len(strs[j]) || strs[j][i] != char {
				return first[:i] // 返回前缀部分
			}
		}
	}
	// 如果所有字符都匹配，返回整个第一个字符串
	return first
}

/*
给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串 s ，判断字符串是否有效。
有效字符串需满足：
左括号必须用相同类型的右括号闭合。
左括号必须以正确的顺序闭合。
每个右括号都有一个对应的相同类型的左括号。
*/
func isValid(s string) bool {
	// 1. 如果字符串为空，则认为有效
	if len(s) == 0 {
		return true
	}

	// 2. 创建一个映射，存储右括号对应的左括号
	// 这样在遇到右括号时，可以快速查找其匹配的左括号
	mapper := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	// 3. 使用切片模拟栈。rune类型用于存储字符。
	stack := []rune{}

	// 4. 遍历字符串中的每一个字符 (使用 range 确保正确处理多字节字符，尽管这里只有单字节的括号)
	for _, char := range s {
		switch char {
		case '(', '{', '[': // 如果是左括号，压入栈
			stack = append(stack, char)
		case ')', '}', ']': // 如果是右括号
			// 检查栈是否为空，如果为空说明没有匹配的左括号
			if len(stack) == 0 {
				return false
			}
			// 弹出栈顶元素（即最近的左括号）
			top := stack[len(stack)-1]
			stack = stack[:len(stack)-1] // 模拟弹出

			// 检查弹出的左括号是否与当前右括号匹配
			if mapper[char] != top {
				return false // 不匹配，无效
			}
		default:
			// 理论上题目只包含这六种字符，如果包含其他字符，可以根据需求选择返回false或忽略。
			// 这里我们假设输入只包含题目规定的字符，无需额外处理。
		}
	}

	// 5. 遍历结束后，如果栈为空，说明所有括号都已正确匹配
	return len(stack) == 0
}

/*
给你一个整数 x ，如果 x 是一个回文整数，返回 true ；否则，返回 false 。
回文数是指正序（从左向右）和倒序（从右向左）读都是一样的整数。
例如，121 是回文，而 123 不是。
*/
func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	str := strconv.Itoa(x)
	fmt.Println(str)
	runes := []rune(str)
	fmt.Println(runes)
	r := make([]rune, len(runes))
	a := 0
	for i := len(runes) - 1; i >= 0; i-- {
		fmt.Printf("runes[%v] = %v \n", i, runes[i])
		r[a] = runes[i]
		a++
	}
	// fmt.Println(str == string(r))
	return str == string(r)
}

/*
136. 只出现一次的数字：给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。
可以使用 for 循环遍历数组，结合 if 条件判断和 map 数据结构来解决，例如通过 map 记录每个元素出现的次数，然后再遍历 map 找到出现次数为1的元素。
*/
func singleNumber(nums []int) int {
	res := make(map[int]int)
	for _, v := range nums {
		res[v] = res[v] + 1
	}
	for k, v := range res {
		if v == 1 {
			return k
		}
	}
	return 0
}
