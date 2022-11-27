package user

import (
	"fmt"
	"testing"
)

//func TestLogin(t *testing.T) {
//	token, err := GetOneCCNUToken()
//	fmt.Println(token, err)
//}
func TestJsonToStruct(t *testing.T) {
	_, err := Login("2021213990", "Yyjiwtb")
	fmt.Println(err)
}
