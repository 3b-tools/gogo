package testdata

import "fmt"
import "github.com/morganhein/gogo"

func BasicDescription(ctx gogo.Context) error {
	ctx.SetDescription("set a description")
	return nil
}

func BasicSetHelp(ctx gogo.Context) error {
	ctx.SetHelp("is this a thing?")
	return nil
}

func BasicArgument(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1)
	return nil
}

func BasicDescriptionArgument(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).Description("describe what this argument does")
	return nil
}

func BasicCtxChained(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.SetDescription("set a description, this can use any go code to set the value").
		SetHelp("is this a thing?")

	fmt.Println(var1, var2)

	return nil
}

func BasicArgumentChained(ctx gogo.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).
		Description("describe what this argument does").
		AllowedValues("1", "2", "3").
		RestrictedValues("4", "5", "6").
		Default("1")

	fmt.Println(var1, var2)

	return nil
}

// the below are not valid GoGo functions

// not valid because the first parameter is not a gogo.Context
func NoGoGoCtx() error {
	return nil
}

// not valid because the first parameter is not a gogo.Context
func GoGoCtxInWrongPos(var1 string, ctx gogo.Context, var2 bool) error {
	return nil
}

// not valid because we can only return an error
func WrongReturnType(ctx gogo.Context, var1 string, var2 bool) string {
	return ""
}

/*
Other requirements we should test:
arguments must be of type string, bool, int, or float64
*/
