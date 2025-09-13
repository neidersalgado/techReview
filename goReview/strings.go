package main

import (
	"fmt"
	"strings"
)

func IdentifyPrefixPostfix(userID, email string) bool {
	return strings.HasPrefix(email, userID) || strings.HasSuffix(email, userID)
}

func ContainsEducative(email string) bool {
	return strings.Contains(email, "educative")
}

func MaskUserName(email string) string {
	mailAt := email[strings.Index(email, "@"):]
	userName := email[:strings.Index(email, "@")]
	newStr := strings.Repeat("*", len(userName)-2)
	return string(email[0]) + newStr + string(email[len(email)-1]) + mailAt
}

func main() {
	// Test your functions here
	fmt.Println(IdentifyPrefixPostfix(".io", "evangeline@educative.io")) // true
	fmt.Println(IdentifyPrefixPostfix("UID", "UID-0123"))                // true
	fmt.Println(IdentifyPrefixPostfix("UID", "evangeline@educative.io")) // false
	fmt.Println(ContainsEducative("evangeline@educative.io"))            // true
	fmt.Println(MaskUserName("evangeline@educative.io"))                 // e******e@educative.io
	// 123
}
