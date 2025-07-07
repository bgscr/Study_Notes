package main

import "fmt"

func main() {
	strs := []string{"flower", "flow", "flight"}
	fmt.Println(longestCommonPrefix(strs))

}

func longestCommonPrefix(strs []string) string {
	result := []byte{}
	for i := 0; i < len(strs[0]); i++ {
		s := strs[0][i]
		for j := 1; j < len(strs); j++ {
			if i >= len(strs[j]) || s != strs[j][i] {
				return string(result)
			}
		}
		result = append(result, s)
	}
	return string(result)
}
