package main

import ."fmt"

func main() {

	var str1 = []byte{'A', 'C', 'd', 'e', 'S', 'A'}

	var str2 []byte
	Printf("\n------------%s------------\n", str1)
	for _, r := range str1 {
		if r > 0x40 && r < 0X5b {
			str2 = append(str2, r+0x20)
		} else {
			str2 = append(str2,r)
		}
	}
	Printf("\n------------%s------------\n", str2)

}
