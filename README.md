# GoDSL

A powerful and flexible parser for meta-annotated variables and function calls in Go. This project helps you create domain-specific languages (DSLs) by allowing you to expose Go variables and functions through a simple, custom syntax.

> [!NOTE]  
> Head over to the [GoDSL Examples](https://github.com/toxyl/go-dsl-examples) repository for example code that showcases how GoDSL can be used.

## What is a DSL?

A Domain-Specific Language (DSL) is a programming language designed for a specific purpose or domain, rather than for general-purpose programming. Unlike general-purpose languages like Go or Python, DSLs are tailored to solve problems in a particular domain, making them more expressive and easier to use for that specific purpose. Common examples include:
- SQL for database queries
- HTML for web page structure
- Regular expressions for text pattern matching
- Configuration languages like YAML or TOML

DSLs are particularly useful when you want to:
- Create a simple scripting interface for your Go application
- Build a configuration language
- Create a mathematical or logical expression evaluator
- Expose Go functionality through a custom syntax

## What does GoDSL do?

GoDSL is a tool that generates a parser for a domain-specific language (DSL) from Go source code with special annotations. The parser implementation (located in the `parser` folder) is generic and will be created in your package directory. Note that files with the prefixes `dsl_` and `template_` will be automatically removed when running GoDSL, so you should avoid using these prefixes in your own files.

After copying the base files, GoDSL generates a `dsl_init.go` file containing an `init()` function that loads all annotated functions and variables into the parser. Once this setup is complete, you can begin using your custom language.

## Features

- Parse nested function calls and expressions
- Support named and positional arguments
- Validate and convert int, float, bool, string, and image types
- Handle optional parameters with defaults
- Manage global variables
- Visualize AST trees in debug mode
- Reference script arguments ($1, $2, etc.)
- Process inline comments (# comment #)
- Handle escaped characters in strings and comments
- Flexible expression whitespace
- Detailed error reporting
- Image processing capabilities:
  - Support for RGBA and NRGBA image formats
  - Proper transparency handling with alpha channel support
  - Pre-multiplied alpha support for correct blending
  - Comprehensive test suite for image operations including:
    - Test image generators for validation:
      - Gradient patterns with alpha transitions
      - Checkerboard patterns with alternating transparency
      - Color wheels with radial transparency
      - Noise patterns with random transparency
      - Color bands and edge case test patterns

## Syntax

The parser supports a simple but powerful syntax for function calls and variable management. Here are the key features:

- **Basic Function Calls**: Call functions with arguments separated by spaces, like `functionName(arg1 arg2 "string arg")`
- **Named Arguments**: Use named parameters for more readable function calls, such as `functionName(param1=value1 param2=value2)`
- **Nested Function Calls**: Combine function calls by nesting them, for example `outerFunction(innerFunction(arg1 arg2) arg3)`
- **Variable Assignment**: Create and set variables using the syntax `variableName: value`
- **Argument References**: Reference script arguments using `$1`, `$2`, etc., as in `functionName($1 $2)`
- **Comments**: Add inline comments using the `#` symbol, like `functionName(arg1 # This is a comment # arg2)`. You can escape the `#` character using `\#` if needed.
- **Strings**: Enclose text in `"` characters, like `"hello world"`. You can escape the `"` character using `\"` if needed.

> [!NOTE]  
> The same argument style (named or unnamed) must be used consistently throughout the entire expression, including any nested function calls.  
> These expressions are therefore **invalid**:
> - `functionName(arg1 param2=value2 arg3)` (mixing styles in the same call)
> - `outerFunction(innerFunction(arg1 arg2) param=value)` (different styles in nested calls)

## Type System

The parser supports a set of basic types to handle different kinds of data:

- **Integer**: Whole numbers like `42` for counting and discrete values
- **Float**: Decimal numbers like `3.14` for precise calculations and measurements
- **Boolean**: Logical values `true` and `false` for conditional operations
- **String**: Text values enclosed in `"` characters, like `"hello \"world"`. You can escape the `"` character using `\"` if needed
- **Image**: Image data in RGBA/RGBA64 or NRGBA/NRGBA64 format, supporting 8-bit and 16-bit color depths with full alpha channel transparency

## Error Handling

The parser provides detailed error messages to help you identify and fix issues in your DSL code. Here are the types of errors that will be reported:

- **Syntax Errors**: Missing or mismatched parentheses, unterminated strings or comments
- **Function Errors**: Unknown function names, invalid function calls, or unexpected tokens
- **Argument Errors**: Invalid argument types, invalid named arguments, or type conversion failures
- **Variable Errors**: Invalid variable assignments, reference errors, or missing variables
- **Runtime Errors**: Any errors that occur during the execution of your DSL code

## Parser Flow

The parser processes your DSL code through several distinct stages:

1. **Tokenization**: The input is broken down into individual tokens, carefully handling strings, variables, and special characters
2. **Lexical Analysis**: The sequence of tokens is validated to ensure proper syntax and structure
3. **Parsing**: An abstract syntax tree (AST) is constructed from the validated tokens
4. **Evaluation**: The function calls and variable assignments are executed according to the AST structure

You can find detailed test coverage of these stages in the `parser/pkg_test.go` file.

## Development Workflow

Here are some important requirements and considerations for developing with GoDSL:

- **Package Initialization**: Do not write your own `init()` function for the package, as GoDSL will generate one that loads everything the DSL needs
- **Namespace Organization**: GoDSL uses a pseudo-namespace under `dsl` in the package root to minimize interference with existing files. All DSL-related types will be prefixed with `dsl`
- **File Management**: GoDSL copies source code and templates to your package directory, prefixing them with `dsl_` (source code) and `template_` (templates). Avoid using these prefixes for your own files as they will be removed during subsequent GoDSL runs
- **Language Exposure**: By default, no language features are exposed. You'll need to write your own code to expose the desired functionality
- **Generated Files**: Do not edit the generated files if you plan to use GoDSL to update your language later

### Write annotated Go code

`/src/my-project/example/definition.go`:  

```go
package main

var (
    // @Name:  last
    // @Desc:  last result
    // @Range: -
    // @Unit:  -
    res = 0.0
)

// @Name: add
// @Desc: Adds two numbers
// @Param:      x      -    0..10   0   First number
// @Param:      y      -    0..10   0   Second number
// @Returns:    result -    0..20   0   Sum of the numbers
func add(x, y float64) (result float64, err error) {
    return x + y, nil
}
```

This will define a DSL with one variable and one function, each is annotated so GoDSL can parse out the necessary information for validation and documentation. 

### Definining Functions

When defining functions for your DSL, keep these requirements in mind:

- **Function Location**: Functions must be defined at the package level
- **Parameter Count**: Functions can have any number of parameters
- **Return Values**: Functions must return a pair of values, with the second value being an `error`
- **Supported Types**: The following types are allowed for parameters and returns:
  - `float*` (any float type)
  - `int*` (any integer type)
  - `uint*` (any unsigned integer type)
  - `bool`
  - `string`
  - `*image.RGBA` (8-bit RGBA image type)
  - `*image.NRGBA` (8-bit non-premultiplied RGBA image type)
  - `*image.RGBA64` (16-bit RGBA image type)
  - `*image.NRGBA64` (16-bit non-premultiplied RGBA image type)

Each function must be annotated with the following information:
- **@Name**: The function's name
- **@Desc**: A description of what the function does
- **@Param**: For each parameter, specify:
  - Name
  - Unit (`-` for unitless)
  - Range (`-` for no range)
  - Default value (`-` for no default)
  - Description
- **@Returns**: For the first return value, specify:
  - Name
  - Unit (`-` for unitless)
  - Range (`-` for no range)
  - Default value (`-` for no default)
  - Description

> [!NOTE]  
> While you can annotate the `error` return value, it's recommended to omit it for functions that never return an error to keep the documentation clean. The `error` return is used internally by the parser to determine if a function executed successfully.

### Defining Variables
- Variables **MUST** be defined in a single block
- GoDSL expects these annotations:
    - **@Name** is the name of the variable.  
    - **@Desc** is the description of the variable.      
    - **@Range** is the range of the variable (omit or use `-` for no range).
    - **@Unit** is the unit of the variable (omit or use `-` for no unit).

## Generate the DSL

To get started with your DSL, follow these steps:

1. **Build and Install GoDSL**:
   If you haven't built `go-dsl` yet, you'll need to build and install it first:
   ```bash
   go build -o go-dsl ./app/
   sudo cp go-dsl /usr/local/bin/
   ```

2. **Generate the DSL**:
   Navigate to your project's root directory and run:
   ```bash
   cd /src/my-project
   go-dsl "basic" "Basic Example" "A basic example implementation" "1.0.0" "basic" example/
   ```

   The command takes these arguments:
   - `id`: Language identifier (e.g., `cpp` for C++ or `go` for Golang)
   - `name`: The name of your language
   - `description`: A description of your language
   - `version`: The version of your language
   - `extension`: File extension for your language (e.g., `go` for `*.go` files)
   - `packages`: The packages to scan for annotated functions and variables (each package will generate a separate DSL)

Once complete, your package directories will contain all necessary files for the DSL, including an `init()` function in `dsl_init.go` that prepares the pseudo-namespace and loads all variables and functions. You're now ready to use your DSL!

## Write your main() function

Here's a simple example of how to create a CLI tool that processes your DSL:

```go
package main 

import (
    "fmt"
    "strings"

    "github.com/toxyl/flo"
)

func main() {
    // Execute the DSL script using command line arguments
    r, err := dsl.run(strings.Join(os.Args[1:], " "), true) // true enables debug mode
    if err != nil {
        fmt.Println("\x1b[31mError:\x1b[0m", err)
        return
    }

    // Process the script's result
    var res any
    if r.err != nil {
        fmt.Println("\x1b[31mError:\x1b[0m", r.err)
        return
    }
    res = r.value
    fmt.Printf("\x1b[32mResult:\x1b[0m %v\n\n", res)

    // Example of variable manipulation
    v := dsl.vars.get("gx")
    gx := v.get().(float64)
    fmt.Println("\x1b[32mgx\x1b[0m =", "gx", "*", 10)
    dsl.vars.set("gx", gx*10)
    fmt.Println("\x1b[32mgx\x1b[0m =", v.get())

    // Export documentation
    flo.File("doc.md").StoreString(dsl.docMarkdown())
    flo.File("doc.html").StoreString(dsl.docHTML())
}
```

To run your application:
```bash
cd /src/my-project/
go run ./example/
```

Alternatively, you can use the DSL shell for a more interactive experience:

## The DSL Shell

The DSL shell provides an interactive environment for testing and using your DSL. It can be used both for development and as a standalone application, as demonstrated in several examples. To launch it, simply call `dsl.shell()` in your main function:

```go
package main

func main() {
    dsl.shell() // Launch the interactive shell
}
```

The shell will display a welcome message with basic usage examples and available commands.

### Basic Usage

The shell supports several types of operations:

1. **Function Calls**:
   ```go
   add(1 2)          // Basic function call
   add(x=1 y=2)      // Named arguments
   add(1 sub(2 3))   // Nested function calls
   ```

2. **Variable Management**:
   ```go
   x: 42             // Create and set variable
   y: add(x 10)      // Use variables in expressions
   ```

3. **Variable Inspection**:
   ```go
   x                 // Display value of x
   y                 // Display value of y
   ```

### Available Commands

The shell provides several built-in commands to help you work with your DSL:

- `?` - Display the help screen
- `help` - Show the full documentation
- `debug` - Toggle debug mode (displays AST trees)
- `store` - Save the current variable state
- `restore` - Restore the previous variable state
- `export-md` - Export documentation as Markdown
- `export-html` - Export documentation as HTML
- `export-vscode-extension` - Generate a VSCode extension for your DSL
- `search [term]` - Search documentation for variables or functions
- `exit` or `CTRL+D` - Exit the shell
- `TAB` `TAB` - Show autocomplete suggestions for variables and functions
