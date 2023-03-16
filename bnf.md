# bnf


Set literals
name{value otherValue nested{set}}

Anonymous set
{literals can be anonymous}

Everything is a function. Functions can be called with arguments
expr < function

An async function which can run forever. It is passed the expression on the left which it can call any number of times
consumer > expr

When an input is given to a function map, the inputs name is read, and the function with the same name is called with the sets children.
functionMap(switchName{expr} otherSwitchName{expr expr})

Consumers can be given functionMaps, similar to a switch statement in functionality
consumer > functionMap

Everything always flows left to right, so stacking multiple expressions, functionmaps, etc. is possible.
consumer > functionMap expr functionMap

The function of a set is the identity function.


<program>		::= <expression> | <expression> <program>

<expression>	::= <set> | <function_call> | <function_map> | <identifier>

<function_call>	::= <identifier> | <identifier> <function_call>

<set>			::= <identifier> '{' <set_list> '}' | '{' <set_list> '}'

<set_list>		::= <identifier> | <set> | <identifier> <set_list> | <set> <set_list> 

<identifier>	::= <string>

<declaration>	::= <identifier> '<' <expression>

<stream>		::= <identifier> '>' <expression>

<function_map>	::= <identifier> '(' <set_list> ')'
