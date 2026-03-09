# Bonus – Switching Between Two Claude Accounts

Claude Code has no native profile switcher — but there are two practical approaches depending on where you work.

## Terminal: `CLAUDE_CONFIG_DIR` aliases (recommended)

`CLAUDE_CONFIG_DIR` is a supported env var that tells Claude where to store its config/credentials. Point it at a different directory per account.

Your current account is already stored in `~/.claude` (the default), so you only need a new directory for account 2.

**One-time setup:**
```bash
mkdir ~/.claude-account2
```

Add to `~/.zshrc`:
```bash
alias claude1="claude"
alias claude2="CLAUDE_CONFIG_DIR=~/.claude-account2 claude"
```

```bash
source ~/.zshrc
claude2   # log in with account 2, then /exit
```

After that each alias just works — can run both simultaneously in separate tabs.

## Terminal alternative: `cswap`

A third-party tool that backs up and swaps `~/.claude.json` and OAuth tokens.

```bash
uv tool install claude-swap   # or: pipx install claude-swap

cswap --add-account   # save current logged-in account (repeat for each)
cswap --switch        # rotate to next account
cswap --switch-to user@email.com
cswap --status
cswap --list
```

Restart Claude after switching. Repo: https://github.com/realiti4/claude-swap

## VSCode extension: VS Code Profiles

The extension does **not** respect `CLAUDE_CONFIG_DIR` set in shell profiles. The cleanest isolation is **separate VS Code profiles** — each profile has its own extension storage, so each can be logged into a different Claude account.

`Cmd+Shift+P` → "Profiles: Create Profile" → log into a different account in each profile.

To switch profiles: click the profile name in the bottom-left, or `Cmd+Shift+P` → "Profiles: Switch Profile".

## Tying terminal + VSCode together

Set `CLAUDE_CONFIG_DIR` per VS Code profile so the integrated terminal always matches the extension account.

1. Switch to Profile 2 in VSCode
2. `Cmd+Shift+P` → **"Open User Settings (JSON)"** (opens settings for current profile only)
3. Add:
```json
"terminal.integrated.env.osx": {
  "CLAUDE_CONFIG_DIR": "/Users/mohamedaboelnaga/.claude-account2"
}
```

Profile 1 needs no change — it uses `~/.claude` by default.

## Summary table

| Method | Terminal | VSCode | Simultaneous |
|--------|----------|--------|--------------|
| `CLAUDE_CONFIG_DIR` aliases | Yes | No | Yes |
| `cswap` | Yes | No | No (swap only) |
| VS Code Profiles | No | Yes | Yes (different windows) |

## Check which account is active

```bash
# terminal
claude /status

# or check the config file
cat ~/.claude.json | grep -i email
cat ~/.claude-account2/.claude.json | grep -i email
```

In VSCode: look at the bottom status bar — it shows the logged-in account email.
