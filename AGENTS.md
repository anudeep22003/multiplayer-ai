## Planning & Decisions
- **Clarify up front:** Before writing code, ask questions to lock down requirements, constraints, and edge cases. Do not guess on ambiguity.
- **Escalate decisive points:** Never make significant architectural, schema, dependency, or contract decisions unilaterally. Present concise options with trade-offs and wait for approval.
- **Plan before acting:** For non-trivial work, outline the approach and files you will touch before editing.

## Code & Architecture
- **Simplest, but no simpler:** Choose the most straightforward solution that solves the current problem. Avoid speculative abstractions and premature optimization.
- **SOLID & modular:** Keep modules small, cohesive, and single-purpose. Depend on clean interfaces/boundaries so components can be swapped or extended without cascading rewrites.
- **Explicit over clever:** Favor readable, idiomatic code with clear naming and obvious data flow over dense one-liners or magic.
- **Fail fast:** Validate inputs at boundaries and handle errors explicitly—never swallow failures silently.

## Git & Commits
- **Atomic, complete commits:** Each commit must represent one logically distinct, self-contained unit of work that compiles and passes checks on its own.
- **Reviewable diffs:** Keep commits small and focused. Do not mix formatting/refactors with behavioral changes.
- **Stay in scope:** Only touch files relevant to the current task. Leave unrelated code alone.

## Verification
- Run the project's typechecker, linter/formatter, and relevant tests before committing or handing off work.

## Project-Specific Rules
<!-- Add stack conventions, commands, or footguns below as the repo grows -->
