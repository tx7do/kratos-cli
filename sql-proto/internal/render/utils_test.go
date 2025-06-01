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
