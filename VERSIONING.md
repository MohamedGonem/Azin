# Versioning

This project follows Semantic Versioning (SemVer). Version numbers are of the form:
```
MAJOR.MINOR.PATCH[-PRELEASE]
```
For example:
```
0.3.0-alpha.1
0.3.0-beta.2
0.3.0-rc.1
0.3.0
0.3.1
0.4.0
1.0.0
```

The pre-release (`alpha`, `beta`, `rc`) indicate how mature an upcoming release is.

The **major**, **minor** and **patch** numbers describe how compatible a stable release is with the previous one.

## General rules

The goal of SemVer is to communicate compatibility to users. As such, **version numbers describe observable behavior, not the amount of work that went into a release!**

Keep the following rules in mind:

- **NEVER modify an existing release**. If `0.3.0` contains a bug, release `0.3.1`. Never replace or retag `0.3.0`.
- **Do NOT vibe version**. We aren't Linus Torvalds, so a release shouldn't become `2.0.0` just because you got a fairy last night in your dreams and she made you feel like this release is important or contains a lot of work and that's why you should increase the major version.
- **Internal implementation does NOT affect the version**. Rewriting the parser, type checker, optimizer or backend doesn't automatically warrant incrementing the minor version when incrementing the patch version is more suitable relative to the observable user behavior.
- **Always think like an user**. Ask yourself:
    - Can the users do something they couldn't before?
    - Will existing projects stop compiling?
    - Will existing build scripts or tooling break?
- **When in doubt, version according to compatibility, not implementation!**

## Pre-releases

### Alpha

Use **alpha** releases while actively building the next release.

This is where major development happens. Features may be incomplete, unstable or completely redesigned.

Typical alpha work includes:

- Adding new language features.
- Designing or redesigning syntax and semantics.
- Implementing compiler passes.
- Rewriting major components (parser, type checker, optimizer, backend, etc.).
- Experimenting with new ideas.

**Stay in alpha until you've finished building the release.**

If you're still asking yourself:

- Should this syntax change?
- Should this feature exist?
- Should we redesign this API?
- Is this implementation the right approach?

then it's still an alpha.

Examples:
```
0.3.0-alpha.1
0.3.0-alpha.2
0.3.0-alpha.3
```

### Beta

Move to **beta** once the release is **feature-complete**.

Everything planned for the release should already exist. The focus shifts from building new functionality to making existing functionality reliable.

Typical beta work includes:

- Fixing parser, type checker and code generation bugs.
- Fixing optimizer bugs.
- Improving diagnostics.
- Improving optimization quality.
- Improving compile times and memory usage.
- Expanding test coverage.
- Updating documentation.

**Avoid adding new features during beta.**

If somebody proposes a new language feature, compiler flag, backend or standard library module, it should almost always wait for the next release.

Examples:

```
0.3.0-beta.1
0.3.0-beta.2
0.3.0-beta.3
```

### Release candidate (RC)

A **Release Candidate** is a build you believe is ready to ship.

If no serious bugs are found, the final release should simply be the last RC.

Only fix release-blocking issues, such as:

- Internal Compiler Errors (ICEs).
- Incorrect code generation.
- Serious parser or type checker bugs.
- Regressions.
- Packaging or installation problems.

Do not:

- Add new features.
- Change syntax or semantics.
- Perform major refactors.
- Introduce new warnings or compiler options unless required to fix a release blocker.

Examples:
```
0.3.0-rc.1
0.3.0-rc.2
```

If everything looks good, then we release:
```
0.3.0
```

## Stable releases

### Patch

> Patch version Z (x.y.Z | x > 0) MUST be incremented if only backward compatible bug fixes are introduced. A bug fix is defined as an internal change that fixes incorrect behavior.

Increment the **patch** version when improving existing functionality **without adding new capabilities or intentionally breaking compatibility**.

Typical examples:

- Fix parser bugs.
- Fix type checker bugs.
- Fix code generation bugs.
- Fix optimizer bugs.
- Fix Internal Compiler Errors (ICEs).
- Improve optimization quality.
- Improve compile speed.
- Reduce memory usage.
- Improve diagnostics.
- Refactor internal code.
- Add tests.
- Update dependencies.
- Improve documentation.

**Rule of thumb**: Users can't do anything new, they simply get a more correct, faster or more reliable compiler.

Examples:
```
0.3.0 → 0.3.1 
0.3.1 → 0.3.2
```

These are still **patch** releases:

- Rewriting the parser.
- Rewriting the optimizer.
- Rewriting the backend.
- Producing faster or better machine code.
- Fixing the compiler so it correctly rejects programs that were never valid according to the language specification.

The amount of work does not determine the version number.

### Minor

> Minor version Y (x.Y.z | x > 0) MUST be incremented if new, backward compatible functionality is introduced to the public API. It MUST be incremented if any public API functionality is marked as deprecated. It MAY be incremented if substantial new functionality or improvements are introduced within the private code. It MAY include patch level changes. Patch version MUST be reset to 0 when minor version is incremented.

Increment the **minor** version when adding new functionality without breaking existing users.

Typical examples:

- New language features.
- New standard library APIs.
- New compiler flags.
- New warnings.
- New optimization levels.
- New target architectures.
- New code generation backends.
- New output formats.
- Deprecating existing functionality.

**Rule of thumb**: Existing projects continue to compile, but users gain new capabilities.

Examples:
```
0.3.2 → 0.4.0
1.4.7 → 1.5.0
```

Adding features should **not** accidentally break existing code. For example, introducing a new reserved keyword that causes previously valid programs to stop compiling is a breaking change. Since Azin doesn't currently have editions or language versions, that generally requires a major release.

### Major

> Major version X (X.y.z | X > 0) MUST be incremented if any backward incompatible changes are introduced to the public API. It MAY also include minor and patch level changes. Patch and minor versions MUST be reset to 0 when major version is incremented.

Increment the **major** version when intentionally introducing breaking changes.

Typical examples:

- Existing valid programs no longer compile.
- Existing programs change behavior.
- Breaking syntax or semantics.
- Removing language features.
- Removing compiler flags.
- Breaking the public API.
- Changing a documented stable file format.
- Removing previously deprecated functionality.

**Rule of thumb**: Existing users may need to update their code, build scripts or tooling before upgrading.

Examples:

```
1.9.4 → 2.0.0
2.7.1 → 3.0.0
```

Breaking changes should be **intentional**, **documented**, and **announced**. We DO NOT make breaking changes without announcing users!