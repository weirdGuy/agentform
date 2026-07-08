# CLAUDE.md — kastor

Kastor is "Terraform for AI agents": a declarative HCL spec compiled to agent frameworks or reconciled against hosted platforms. **SPEC.md is the source of truth** — read it before making design decisions. If code and SPEC.md conflict, flag it; don't silently diverge.

## What this is

- Go CLI (`kastor`) that parses `.agent`, `.tool`, `.prompt`, and `kastor.hcl` project files
- Two execution paths: `kastor build` (codegen → LangGraph first) and `kastor plan/apply` (platform reconciler → OpenAI Assistants first)
- Non-goals for v0: being a runtime, executing agents, eval harnesses

## Architecture

```
cmd/kastor/         CLI entrypoint (cobra)
internal/
  parser/           HCL decode (hashicorp/hcl/v2) → AST
  schema/           typed config structs, validation
  module/           directory walk → symbol table, cross-file reference resolution
  graph/            DAG construction, cycle detection, topo sort
  build/            codegen engine + per-target generators (build/langgraph/)
  provider/         platform reconcilers (provider/openai/)
  state/            state file read/write, locking, diff
```

## Commands

```bash
go build ./...                 # build everything
go test ./...                  # run all tests
go test ./internal/parser/     # test one package
go vet ./...                   # static checks
gofmt -l .                     # formatting check (must be clean)
```

## Conventions

- Go 1.22+, standard library first; approved deps: cobra, hashicorp/hcl/v2, go-cmp (tests)
- All packages under `internal/` except `cmd/`; no public API surface in v0
- Errors: wrap with `fmt.Errorf("context: %w", err)`; every user-facing diagnostic states what was found, what was expected, and where — file:line plus block address (e.g. `agent.weather: unknown reference model.fastt`)
- Table-driven tests; fixtures live in `testdata/` per package (valid + invalid HCL samples)
- Every parser/validation feature needs at least one negative test (bad input → expected diagnostic)
- Providers implement the common interface: `Read / Create / Update / Delete / Diff`
- Keep codegen output deterministic — same spec must always produce byte-identical output (needed for testing and CI diffs)

## Domain rules to enforce (from SPEC.md)

- Agent owns model + IO contract; prompts are pure templates with `requires` variables
- Every prompt variable must be satisfiable from the agent's inputs/outputs → else compile error
- References (`agent.x.output.y`, `model.x`, `tool.x`, `prompt.x`) build the DAG; cycles are a compile error
- `depends_on` is the explicit fallback only — never infer data flow from it
- A tool has exactly one `source` block, `kind` ∈ mcp | http | builtin | runtime | script

## Workflow

- Small PRs mapped to GitHub issues; reference issue number in commits
- Never commit directly to main; always branch (`feat/<issue>-<slug>`) + PR
- `kastor validate` must stay fast — it runs on every save in editor integrations later
- When adding a block field: update schema struct → validation → parser test fixtures → SPEC.md if it's a design change
- Before claiming a milestone or feature is code complete, attempt its acceptance command (e.g. `kastor validate` / `kastor build` on the examples) and confirm the output — passing tests alone don't count