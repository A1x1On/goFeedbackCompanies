package helper

import(
	"github.com/pkg/errors"
)

func IfErr(err error) {
	if err != nil {
		panic(err)
	}
}

func IfError(err error, msg string) {
	if err != nil {
		panic(errors.Wrap(err, msg))
	}
}

// func TError(err errorm i interface{}, msg string) {
// 	if err != nil {

// 	switch v := r.(type) {
// 		case ParseJson:
// 			fmt.Println("Stringer:", v)
// 		default:
// 			fmt.Println("Unknown")
// 		}
// 	}
// }