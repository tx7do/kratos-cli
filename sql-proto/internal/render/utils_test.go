package render

import (
	"fmt"
	"testing"
)

func TestSnakeToPascal(t *testing.T) {
	fmt.Println(snakeToPascal("user_id"))
	fmt.Println(snakeToPascal("id"))
	fmt.Println(snakeToPascalPlus("id"))
	fmt.Println(snakeToPascalPlus("user_id"))
	fmt.Println(snakeToPascalPlus("identity_provider_id"))
}

func TestMakeEntSetNillableFunc(t *testing.T) {
	fmt.Println(makeEntSetNillableFunc("id"))
	fmt.Println(makeEntSetNillableFunc("user_id"))
	fmt.Println(makeEntSetNillableFunc("identity_provider_id"))
	fmt.Println(makeEntSetNillableFunc("last_login_time"))

	fmt.Println(makeEntSetNillableFuncWithTransfer("last_login_time", "timeutil.TimestamppbToTime"))
	fmt.Println(makeEntSetNillableFuncWithTransfer("last_login_time", "timeutil.StringTimeToTime"))
}

func TestRemoveTableCommentSuffix(t *testing.T) {
	str := "产品表"
	fmt.Println(RemoveTableCommentSuffix(str)) // 输出: 产品

	str2 := "product table"
	fmt.Println(RemoveTableCommentSuffix(str2)) // 输出: product
}
