# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## next

### Added
- Tab-completion for `)lib`, `sys "prec"`, `sys "trace"`, `sys "write"`

## [0.5.0] - 2024-12-13

### Added
- Tab-completion for `part`, `rand`, `where`, `)clear`, `)last`
- Tab-completion for `)debug` settings
- Tab-completion for control flow keywords
  `:while`, `:if`, `:elif`, `:else`, `:end`, `:ret`
- File name completion for `)get` and `)save` without quotes

## [0.4.0] - 2024-12-12

### Added
- Tab-completion for `mdiv`, `mix`, `opdelete`, `print`, `split`,
  `sys "read"`, `)timezone`

### Changed
- The project now requires Go >= 1.22.0

## [0.3.0] - 2023-08-06

### Added
- Tab-completion for `box`, `conj`, `first`, `inv`, `sys` functions

### Fixed
- Don't create config directory on history load

## [0.2.0] - 2023-07-11

### Added
- Tab-completion for `count`, `flatten`, `intersect`, `union`, `unique`, `var`
- Tab-completion for user-defined variables
- File name completion for `)get` and `)save`

### Changed
- The project now requires Go >= 1.20.0

## [0.1.0] - 2022-05-26

### Added
- Initial release

[0.4.0]: https://github.com/fzipp/ivy-prompt/compare/v0.3.0...v0.4.0
[0.3.0]: https://github.com/fzipp/ivy-prompt/compare/v0.2.0...v0.3.0
[0.2.0]: https://github.com/fzipp/ivy-prompt/compare/v0.1.0...v0.2.0
[0.1.0]: https://github.com/fzipp/ivy-prompt/releases/tag/v0.1.0
