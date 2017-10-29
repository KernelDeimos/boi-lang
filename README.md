Boi lang
--------

The documentation is lacking, so for now just refer to this code snippet.

```
boi, "Hello, Boi!" boi

boi! set greet "Hello," boi
boi! set subject "Boi!" boi
boi, boi:greet " " boi:subject boi

boi, [int 72] [int 69] [int 76] [int 76] [int 79] boi

boi: tmp [int 0] boi
bloop < boi:tmp [int 10] boi
	boi, "I say this 10 times" boi
	boi: tmp [+ boi:tmp [int 1]] boi
BOI

boi! set a [int 1] boi
boi! set b [int 2] boi
boi, [dec [+ boi:a boi:b]] boi

boi, [+ A [int 3]] boi

boi: prob [dec [IsEven 42]] boi
boi, "There is a " boi:prob "% chance that 42 is an even number" boi

boi? cat true boi
	ONE SCHWIFTY BOI
	boi, "value is " boi:SCHWIFTY boi
BOI
boi, "value is " boi:SCHWIFTY boi

```

### What is Boi-lang?
This is an experimental language developed as a joke.

Boi-lang takes things way outside the box. It features dynamic scoping
and only one variable type (byte array), which Boi-lang lets you pass
around however you please.

My inspiration for this work of "art" can be attributed to a Facebook group with
the following title:
> we post the word "boi" every day until the James Webb telescope is launched

Every line contains the word "boi".

## Why?
Developing a silly programming language is fun, and makes it possible to
explore new ideas without worrying about quality and consistency.

## How?
Every statement begins with a keyword (ex: `boi!`) and ends with the
statement terminator, which is `boi`.

There are different keywords which will explained further:

| keyword | Description |
| ------- | ----------- |
| `boi!`  | Call a function |
| `boi?`  | Call a function, and execute succeeding statements if it returns true |
| `boi,`  | Shorthand to call say function |
| `boi:`  | Shorthand to call set function |

### Using Functions
Every function call begins with the `boi!` keyword, followed by a list of tokens
where the first token is the function name. Note that since variables are tokens,
the function name can be taken from a variable.

#### `say` function
The say function outputs its input arguments to standard out
Example:
```
boi! say "Hello, Boi!" boi
```
Output:
```
Hello, Boi!
```
#### `set` function
The set function takes two parameters - a variable name and value
Example:
```
boi! set subject "Boi!" boi
boi! say "Hello, " boi:subject boi
```
Output
```
Hello, Boi!
```

#### `cat` function
The cat function takes any number of parameters, strings them
together and returns the output so it's available in the
`ret:exit` variable.

### Conditionals
Conditionals distinguish computers from calculators. A language without conditionals
is, well, a calculator. 

Conditionals aren't very useful in Boi-lang yet, but here's an example anyway:
```
boi? cat true boi
    boi! say "the cat function returned true" boi
BOI
```

Note that block statements end with `BOI`.

Also note that "true" is a string. See the "truth semantics" section
below for more information.

## Truth Semantics
Every variable in Boi-lang is an array of bytes. This makes the truth
semantics very simple:

| Situation | Memory (hex) | Result |
| --------- | ------------ | ------ |
| Variable doesn't exist | N/A | false |
| ASCII string 'false' | 66 61 6c 73 65 | false |
| Literal binary value 0 | 00 | false |
| Anything else | any of not the above | true |

## Terms & Syntax
| Term | Description |
| ---- | ----------- |
| Identifier | A valid Boi-lang identifier is any valid string. |
| String | A string can be `"in double-quotes with \"escaped quotes\""`, or `outside\ quotes\ with\ escaped\ spaces`. |
| Token | A token in Boi-lang refers to an input value, which is a string or variable. |
