# ivy-prompt

A line editor interface wrapper for Rob Pike's [Ivy](https://robpike.io/ivy),
an interpreter for an [APL](https://tryapl.org/)-like language.
It provides a
[Readline](https://tiswww.case.edu/php/chet/readline/rltop.html)-style input
prompt with input history and tab-completion.

The original `ivy` command interacts with the world via _standard input/output_,
integrating well with a Unix or
[Plan 9](https://p9f.org) environment or a text editor like
[Acme](https://research.swtch.com/acme).
However, if you prefer or are used to a mode of interaction more akin to Bash,
Haskell's GHCi, or Python's REPL with an interactive line editor, then this
project is for you.

Although Ivy is described by its creator as "a plaything," it should not
be underestimated.
[Watch Russ Cox' videos](https://www.youtube.com/playlist?list=PLrwpzH1_9ufMLOB6BAdzO08Qx-9jHGfGg)
of solutions for [Advent of Code 2021](https://adventofcode.com/2021)
using Ivy for a demonstration of its capabilities or take the built-in
tour with the `)demo` command.

## Installation

First install Ivy itself, if you haven't already done so:

```
go install robpike.io/ivy@latest    # or @master if you want to use
                                    # the development version of Ivy.
```

The `ivy` binary needs to be in the `PATH` for ivy-prompt to locate it.

Then install this wrapper:

```
go install github.com/fzipp/ivy-prompt@latest
```

Run it:

```
ivy-prompt
```

## Line Editing

See the [`liner` documentation](https://github.com/peterh/liner#line-editing) for
a complete table of supported keystrokes / actions. Here's a small selection:

| Keystroke | Action                            |
|-----------|-----------------------------------|
| Tab       | Next completion                   |
| Up        | Previous match from history       |
| Down      | Next match from history           |
| Ctrl-L    | Clear screen (line is unmodified) |
| Ctrl-D    | (if line is empty) quit           |

## Tab-Completion

Tab-completion works not only with built-in identifiers
but also with user-defined variables and operators.

It is context-aware to some extent.
For example, it will provide different completion options
after `)help `, `)get "`, `)save "`, or `sys "`
compared to the usual context.

## Input History

The input history is preserved across sessions
in a file named `ivy_history`
within the `ivy` subdirectory,
located within the user configuration directory.
The specific location of this directory
depends on the operating system:

| Operating system | Config directory                         |
|------------------|------------------------------------------|
| Linux/Unix       | `$XDG_CONFIG_HOME/ivy/`                  |
| macOS            | `$HOME/Library/Application Support/ivy/` |
| Windows          | `%AppData%\ivy\`                         |

## License

This project is free and open source software licensed under the
[BSD 3-Clause License](LICENSE).
