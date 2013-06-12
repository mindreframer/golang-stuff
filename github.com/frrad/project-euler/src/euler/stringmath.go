package euler

import (
	"strconv"
)

func StringFastFracAdd(num1 string, den1 string, num2 string, den2 string) (string, string) {
	return StringSum(StringProd(num1, den2), StringProd(num2, den1)), StringProd(den1, den2)
}

func StringDigitSum(str string) (total int) {
	for i := 0; i < len(str); i++ {
		is, _ := strconv.Atoi(str[i : i+1])
		total += is
	}
	return
}

func StringSum(string1 string, string2 string) string {
	length1 := int64(len(string1))
	length2 := int64(len(string2))
	length := 1 + Max(length1, length2)
	string1 = StringReverse(string1)
	string2 = StringReverse(string2)

	sum := make([]int, length)
	str1 := make([]int, length)
	str2 := make([]int, length)

	for i := int64(0); i < length; i++ {
		a := 0
		b := 0
		sum[i] = 0

		if i < length1 {
			a, _ = strconv.Atoi(string(string1[i]))

		}
		str1[i] = a

		if i < length2 {
			b, _ = strconv.Atoi(string(string2[i]))

		}

		str2[i] = b

	}

	for i := int64(0); i < length-1; i++ {
		total := str1[i] + str2[i] + sum[i]
		sum[i] = total % 10
		sum[i+1] = (total - total%10) / 10
	}

	answer := makestring(sum)

	answer = StringReverse(answer)
	answer = StringTrim(answer)

	return answer
}

func times10tothe(n int, x string) string {

	for i := 0; i < n; i++ {
		x = x + "0"
	}

	return x
}

func makelist(s string) []int {
	list := make([]int, len(s))
	for i := 0; i < len(s); i++ {
		list[i], _ = strconv.Atoi(s[i : i+1])
	}
	return list
}

func makestring(table []int) string {
	answer := ""
	for i := 0; i < len(table); i++ {
		answer += strconv.Itoa(table[i])
	}
	return answer
}

func StringProd(x string, y string) (product string) {

	if len(x) < len(y) {
		return StringProd(y, x)
	}

	product = "0"

	if len(y) == 1 {
		if y == "0" || x == "0" {
			return "0"
		}
		if y == "1" {
			return x
		}

		xlist := makelist(StringReverse(x))
		prodlist := make([]int, len(x)+1)
		y, _ := strconv.Atoi(y)
		for i := 0; i < len(x); i++ {
			prodlist[i] += (xlist[i] * y) % 10

			if prodlist[i] >= 10 {
				prodlist[i+1] += prodlist[i] / 10
				prodlist[i] = prodlist[i] % 10
			}

			prodlist[i+1] += (xlist[i] * y) / 10

		}

		return StringTrim(StringReverse(makestring(prodlist)))

	}

	for i := 0; i < len(y); i++ {
		product = StringSum(product, times10tothe(i, StringProd(x, y[len(y)-i-1:len(y)-i])))
	}

	return
}

func StringExp(a string, b int64) string {
	if b == 0 {
		return "1"
	}
	if b == 1 {
		return a
	}
	if b%2 == 0 {
		temp := StringExp(a, b/2)
		return StringProd(temp, temp)
	}
	return StringProd(a, StringExp(a, b-1))
}

//Removes 0-padding on Integer Strings
func StringTrim(a string) string {
	if a == "0" {
		return a
	}
	place := 0
	for i := 0; i < len(a) && string(a[i]) == "0"; i++ {
		place = i + 1
	}

	output := ""

	for i := place; i < len(a); i++ {
		output += string(a[i])
	}

	return output
}
