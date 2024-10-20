package testdata

import "fmt"
import gogo2 "github.com/morganhein/gogo"

func AliasedDescription(ctx gogo2.Context) error {
	ctx.SetDescription("set a description")
	return nil
}

func AliasedSetHelp(ctx gogo2.Context) error {
	ctx.SetHelp("is this a thing?")
	return nil
}

func AliasedArgument(ctx gogo2.Context, var1 string, var2 bool) error {
	ctx.Argument(var1)
	return nil
}

func AliasedDescriptionArgument(ctx gogo2.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).Description("describe what this argument does")
	return nil
}

func AliasedCtxChained(ctx gogo2.Context, var1 string, var2 bool) error {
	ctx.SetDescription("set a description, this can use any go code to set the value").
		SetHelp("is this a thing?")

	fmt.Println(var1, var2)

	return nil
}

func AliasedArgumentChained(ctx gogo2.Context, var1 string, var2 bool) error {
	ctx.Argument(var1).
		Description("describe what this argument does").
		AllowedValues("1", "2", "3").
		RestrictedValues("4", "5", "6").
		Default("1")

	fmt.Println(var1, var2)

	return nil
}
