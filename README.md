# Cata Programming Language

Cata is a small, fast and safe programming language. It does not have a garbage
collector, but uses affine types to ensure memory safety. This both provides
high performance and makes Cata well-suited for embedded and systems
programming. However, in contrast to other similar programming languages like
Rust that also rely on affine types, Cata is designed to be a very small and
simple programming language. This means that it should be possible for a single
person to understand Cata in its entirety. Also, a major goal of Cata is to
compile very quickly.

Currently, Cata is not yet usable. The compiler is written in Go, is not
optimized, generates only C code (which makes compilation slow), and does not
provide good error diagnostics. As soon as the compiler works well enough, it
will be rewritten in Cata itself.

# Roadmap
The current language implementation:
- [x] Basic language features, variables, function calls
- [x] Structs
- [x] Control flow
- [x] Linear types
- [ ] Drop insertion
- [ ] Generics
- [ ] Borrowing
- [ ] Built-in operators
- [ ] Arrays and slices
- [ ] C interop
- [ ] Enums and unions
- [ ] Error handling
- [ ] Struct functions
- [ ] Interfaces
- [ ] Ergonomic improvements

Long-term language goals:
- [ ] The current compiler
- [ ] Standard library
- [ ] Self-hosted compiler
- [ ] Language specification
- [ ] More compilation targets
- [ ] Tooling
