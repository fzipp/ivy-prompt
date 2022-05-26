# ivy-prompt

A command-line interface wrapper for Rob Pike's [Ivy](https://github.com/robpike/ivy),
an [APL](https://tryapl.org/)-like big number calculator. It provides a
[Readline](https://tiswww.case.edu/php/chet/readline/rltop.html)-style input
prompt with an input history and tab-completion.

The original `ivy` command interacts with the world via _standard input/output_,
which integrates well with a Unix or
[Plan 9](https://p9f.org) environment or a  text editor like
[Acme](https://research.swtch.com/acme).
However, if you prefer or are used to a mode of interaction more akin to Bash,
Haskell's GHCi or Python's REPL with an interactive line editor, then this
project is for you.

Even though Ivy is described by its author as "a plaything" it should not
be underestimated.
[Watch Russ Cox' videos](https://www.youtube.com/playlist?list=PLrwpzH1_9ufMLOB6BAdzO08Qx-9jHGfGg)
of solutions for [Advent of Code 2021](https://adventofcode.com/2021)
using Ivy for a demonstration of its capabilities, or take the built-in
tour with the `)demo` command.

## Installation

Install Ivy itself, if you haven't already done so:

```
go install robpike.io/ivy@latest
```

The `ivy` binary should be in the `PATH`.

Install this wrapper:

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

## License

This project is free and open source software licensed under the
[BSD 3-Clause License](LICENSE).
