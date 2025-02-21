// Partially based on https://github.com/google/cel-go/blob/v0.23.2/examples/custom_global_function_test.go
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common/types"
	"github.com/google/cel-go/common/types/ref"
)

func main() {
	env, err := cel.NewEnv(
		cel.Function("httpget",
			cel.Overload("httpget_string_string",
				[]*cel.Type{cel.StringType},
				cel.StringType,
				cel.UnaryBinding(func(u ref.Val) ref.Val {
					url := string(u.(types.String))

					resp, err := http.Get(url)
					if err != nil {
						return types.String(err.Error())
					}
					defer resp.Body.Close()
					body, _ := io.ReadAll(resp.Body)

					return types.String(body)
				},
				),
			),
		),
		cel.Function("println",
			cel.Overload("println_stringarray_string",
				[]*cel.Type{cel.StringType},
				cel.StringType,
				cel.FunctionBinding(func(args ...ref.Val) ref.Val {
					var sb strings.Builder
					for _, v := range args {
						arg := v.(types.String)
						sb.WriteString(string(arg))
					}
					out := sb.String()
					fmt.Println(out)
					return types.String(out)
				},
				),
			),
		),
	)
	if err != nil {
		log.Fatalf("environment creation error: %v\n", err)
	}

	ast, iss := env.Compile(`println(httpget('https://httpbin.org/anything'))`)
	if iss.Err() != nil {
		log.Fatalln(iss.Err())
	}
	prg, err := env.Program(ast)
	if err != nil {
		log.Fatalf("Program creation error: %v\n", err)
	}

	_, _, err = prg.Eval(map[string]any{})
	if err != nil {
		log.Fatalf("Evaluation error: %v\n", err)
	}
}
