# dex

Experimental compositional programming language. Wip.

## Demonstration

We will define the functions in the parent language.

```go
func printName(arg dex.Node) dex.Set {
	result := arg.Eval(nil)
	fmt.Println(result.Name())
	return result
}

// Define a new parser with no scope, such that a new one is generated.
parser := dex.NewParser(nil)

// Add our printName function to the parser's scope as "print".
parser.Scope.Set("print", dex.NewFnWrapper(printName))
```

This function can now be accessed from dex parsed with the parser we just created.

```go
// This dex expression first creates a map literal with name "helloWorld", which is
// then immediately passed to the following expression, in this case our print function.
code := "helloWorld{} print"

// Evaluate the code.
parser.Run(code)
// Output: helloWorld
```


## Syntax

### Literals

Map literals are defined with `{}`, with the keys separated by spaces. Each key is also a map.
`example{ keyOne keyTwo }`
Maps can be anonymous.
`{ keyOne keyTwo }`
Map literals can be nested.
`nested{ keyOne{ nestedKeyOne } keyTwo }`

### Variables

Variables are defined with `<` operator.
`ourVariable < mapExample{ hello world }`
If an identifier is used without it having `{}`, it is assumed to be a scoped expression (function or map).
`helloWorld{} print`

### Function maps

Function maps are defined with `()`, with the keys separated by spaces. Each key in a function map must be followed by `{}` which needs to contain at least one scoped identifier. Function maps work like a switch statement, where the name of the node given as argument is used as the key in the map, and then the node is passed to the related expression and the value returned.

`ourFnMap( doSomething{aFunction anotherOne} doSomethingElse{something important} )`
Passing a node to this would get the nodes name, and see if its name was either `doSomething` or `doSomethingElse`. Their respective expressions would then be evaluated with the node as argument, and the value returned.

As function maps are expression, they can also be applied to variables.
`fnMapVariable < ourFnMap( doSomething{aFunction anotherOne} doSomethingElse{something important} )`

### Streams

Streams are similar to pipes, where the left hand side is used as the source, and the right hand side is evaluated with the argument. Streams are defined with `>` operator. Streams are asynchronous.

This would wait for someStream to give inputs, which would then be passed to the function map, and the result would be printed.
`someStream > fnMapVariable print`


### Putting it all together

We will assume that we have an http server function which formats requests to some wanted format, the path in this case, calling the stream on the right, and then writing the result to the response writer. We will use this to implement a simple http API.

```
aMapOfValues < { valueOne valueTwo nestedValue{ valueFour } }
helloWorld < helloWorld{}
api < apiMap( /{ helloWorld } /list/{ aMapOfValues } )
serveHTTP > print api
```
This program logs the request path, looks if the path is found in the functionMap, and then returns the result of the functionMap's expression.

