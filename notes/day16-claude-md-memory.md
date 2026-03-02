# Day 16 – CC: CLAUDE.md and Memory

## Concepts

### The CLAUDE.md hierarchy
Claude Code loads instruction files from multiple locations, merging them from lowest to highest priority:

```
managed policy (Anthropic/admin)   ← highest, cannot be overridden
  └── project CLAUDE.md            ← .claude/CLAUDE.md or root CLAUDE.md (checked into git)
        └── user ~/.claude/CLAUDE.md  ← applies to every project you work in
              └── CLAUDE.local.md  ← personal project overrides (gitignored)
```

Rules from higher-priority files override lower ones. A managed policy (e.g., from an enterprise admin) always wins.

**Practical rule**: put shared team conventions in `CLAUDE.md` (committed to git), personal preferences in `~/.claude/CLAUDE.md`, and short-lived overrides in `CLAUDE.local.md`.

### Auto memory
Claude Code can write notes to a per-project memory file:

```
~/.claude/projects/<encoded-project-path>/memory/MEMORY.md
```

- The first **200 lines** are injected into the system prompt at the start of every session.
- Claude writes to this file when it discovers patterns, mistakes, or preferences worth remembering.
- You can view and edit it directly or via `/memory` in a session.
- Disable it with `autoMemory: false` in `~/.claude/settings.json`.

### .claude/rules/ — modular instructions
Instead of one giant CLAUDE.md you can split instructions into topic files:

```
.claude/rules/
├── code-style.md        # formatting preferences
├── git-workflow.md      # commit/PR conventions
└── api-conventions.md   # project-specific API patterns
```

Optional front-matter lets you scope a rule to specific file paths:

```yaml
---
paths:
  - "handlers/**/*.go"
---
Always check for errors from c.Bind() before using the bound struct.
```

Rules without `paths:` apply globally.

### CLAUDE.local.md
A personal project-level override file. Claude Code gitignores it automatically. Use it for:
- Personal shortcuts or reminders not meant for the team
- Temporary instructions during a focused refactor session
- Overriding a team rule for your local workflow

---

## Key takeaways
- Four layers: managed → project → user → local
- `~/.claude/CLAUDE.md` = your global preferences across every repo
- `CLAUDE.md` in the repo root = team-shared project instructions
- `CLAUDE.local.md` = personal local overrides (gitignored)
- Auto memory persists observations across sessions (first 200 lines auto-loaded)
- `.claude/rules/` lets you scope instructions to specific file paths

---

## Community research: what people put in global CLAUDE.md

Researched real dotfiles repos and blog posts. Here's what the community consistently does.

### Notable real-world examples

| Source | Style | Emphasis |
|--------|-------|----------|
| [harperreed/dotfiles](https://github.com/harperreed/dotfiles/blob/master/.claude/CLAUDE.md) | Personal, informal | Relationship, TDD, decision framework, git hooks |
| [trailofbits/claude-code-config](https://github.com/trailofbits/claude-code-config) | Professional, security-focused | Language toolchains, hard limits, supply chain |
| [ctoth gist](https://gist.github.com/ctoth/d8e629209ff1d9748185b9830fa4e79f) | Epistemological | Agentic reasoning, verification, irreversibility |
| [joecotellese.com](https://joecotellese.com/posts/claude-code-global-configuration/) | J.A.R.V.I.S. persona | Decision framework, TDD, language-specific |

### The 8 most common categories

**1. Communication style**
```markdown
- Be direct. No hedging. No filler phrases like "certainly!" or "great question!"
- Push back when you disagree — cite the reason
- Keep explanations proportional to the complexity of the change
```

**2. 3-tier decision framework** (very popular pattern)
```markdown
- Autonomous (no need to ask): fix typos, linting, formatting, obvious bugs
- Discuss first: multi-file changes, new features, API modifications
- Always ask: rewrites, security changes, schema changes, anything irreversible
```

**3. Code quality hard rules**
```markdown
- Prefer simple, readable, maintainable over clever
- Match existing file style over external standards
- No commented-out dead code
- Zero warnings policy
- Fast-fail with actionable error messages
```

**4. Git discipline**
```markdown
- NEVER use --no-verify when committing
- Imperative mood commit messages, ≤72 chars
- Feature branches with PRs; never push to main directly
- Pre-commit hook failures must be fixed, never bypassed
```

**5. Testing mandate**
```markdown
- TDD: write failing test first, confirm failure, then implement
- Test behavior over implementation
- Mock only external dependencies (slow, non-deterministic, external services)
```

**6. Language-specific toolchains**
Where you set per-language tool preferences (e.g. "use uv for Python", "use go fmt for Go"). Useful if you work in multiple languages — set once globally rather than repeating in every project.

**7. Agentic reasoning rules** (from ctoth's gist — more advanced)
```markdown
- Maximum 3 actions before verifying reality matches predictions
- "I don't know" is valid output — preferable to confident confabulation
- Cannot explain why something exists? Don't modify it (Chesterton's Fence)
- Defer to user when: intent is ambiguous, state is unexpected, action is irreversible
```

**8. Preferred CLI tools**
```markdown
- Use ast-grep over grep for code analysis
- Use fd over find
- Use rg for text search
```

### What people consistently avoid in the global file
- Project-specific context (stack, routes, DB schema) — belongs in `./CLAUDE.md`
- Sensitive data (API keys, tokens)
- Very long documentation — global file should be under ~200 lines; use `@imports` for longer content
- Anything team-specific — the global file is personal and not shared

### Key insight
The most useful global files are the ones people **update over time** as they notice friction — not a perfect template written once upfront. Start minimal (5–10 lines), grow it.

---

## Q&A

**Q:** When running `/memory`, got a settings error: `hooks: Expected array, but received undefined`
**A:** The hooks format changed in a recent Claude Code update. Old format had hook properties (like `prompt`) directly in the array item. New format wraps them in a nested `hooks` array with explicit `type`:
```json
// Old (broken)
"PreCompact": [{ "prompt": "..." }]

// New (correct)
"PreCompact": [{ "hooks": [{ "type": "prompt", "prompt": "..." }] }]
```
The outer item can have an optional `matcher` (for tool-specific hooks like PostToolUse). The inner `hooks` array contains the actual hook definitions with `type` + type-specific fields.

---

## Q&A

**Q:** When running `/memory`, got a settings error: `hooks: Expected array, but received undefined`
**A:** The hooks format changed in a recent Claude Code update. Old format had hook properties (like `prompt`) directly in the array item. New format wraps them in a nested `hooks` array with an explicit `type`:
```json
// Old (broken)
"PreCompact": [{ "prompt": "..." }]

// New (correct)
"PreCompact": [{ "hooks": [{ "type": "prompt", "prompt": "..." }] }]
```
The outer item can optionally have a `matcher` (used for tool-specific hooks like PostToolUse). The inner `hooks` array contains the actual hook definitions with `type` + type-specific fields.

**Q:** What are the JSONL files and folders in `~/.claude/projects/<project>/`?
**A:** That directory is Claude Code's full local database for the project — not just memory. Structure:
- `<uuid>.jsonl` — one file per conversation session (JSON Lines: one message per line)
- `<uuid>/subagents/agent-<id>.jsonl` — subagent conversations spawned via the Task tool, nested under the parent session
- `sessions-index.json` — lightweight index for listing recent sessions quickly (powers `/resume`)
- `memory/MEMORY.md` — the auto-memory file; first 200 lines injected into every new session

**Q:** Is there a way to access old conversations in Claude Code?
**A:** Yes. Two ways: (1) `claude --resume <session-id>` from the terminal; (2) type `/resume` inside a session to get a list of recent sessions for the current project and pick one. The session ID is the UUID in the JSONL filename.

**Q:** What are subagents?
**A:** Separate Claude instances spawned by the main session via the Task tool. Benefits: (1) parallelism — run independent tasks simultaneously; (2) context isolation — subagent output doesn't bloat the main conversation. Their full conversation is stored in `subagents/agent-<id>.jsonl` nested under the parent session folder.

**Q:** When running `/init` in a subdirectory, it updated the parent project's CLAUDE.md instead of creating a new one — why?
**A:** `/init` walks up the directory tree first. If it finds an existing CLAUDE.md in a parent, it updates that file rather than creating a duplicate. It only creates a new file if no CLAUDE.md exists anywhere in the tree above. To get a project-specific one, create an empty CLAUDE.md manually in the project root first — then `/init` will update that local file next time.

**Q:** Claude is reading parent project files when running in a subdirectory — is that expected?
**A:** Yes. Claude Code walks up the directory tree at startup and loads every CLAUDE.md it finds (e.g. `~/github/CLAUDE.md`, `~/github/go/CLAUDE.md`, `~/github/go/helloworld/CLAUDE.md`). All are merged, with deeper/closer files taking higher priority. Same concept as `.gitignore`. This means you can put shared instructions for all Go projects in `~/github/go/CLAUDE.md` without repeating them in every project.

---

## Exercises checklist
- [x] Open/create `~/.claude/CLAUDE.md` and add 2–3 global personal preferences
- [x] Improve one section of the project `CLAUDE.md` (updated "Current progress" line)
- [x] Run `/memory` in a session and explore what's been saved
- [x] Run `/init` in a fresh directory and observe the generated output
- [x] Be able to explain the precedence order verbally
