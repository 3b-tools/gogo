//go:build ignore

// this file is inherently broken, so we ignore it for the purposes of testing

package testdata

import "fmt"

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
