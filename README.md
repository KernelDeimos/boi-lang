Boi lang
--------

```
boi! say "Hello, boi!" boi
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
Currently, the only features implemented are calling functions and
setting/accessing variables.

### Statement
Every statement begins with a keyword (ex: `boi!`) and ends with the
statement terminator, which is `boi`.

### Function Call
Every function call begins with the `boi!` keyword, followed by a valid
identifier, then a list of tokens.

### List of built-in functions
#### `say` function
Example:
```
boi! say "Hello, Boi!" boi
```
Output:
```
Hello, Boi!
```
### `set` function
Example:
```
boi! set subject "Boi!" boi
boi! say "Hello, " boi:subject boi
```
Output
```
Hello, Boi!
```

## Terms & Syntax
| Term | Description |
| ---- | ----------- |
| Identifier | A valid Boi-lang identifier is an alphanumeric string starting with a letter. |
| String | A string can be `"in double-quotes with \"escaped quotes\""`, or `outside\ quotes\ with\ escaped\ spaces`. |
| Token | A token in Boi-lang refers to an input value, which is a string or variable. |
