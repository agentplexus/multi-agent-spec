# GitHub PR Review with Claude Code

This guide explains how to have Claude Code review GitHub Pull Requests and post review comments directly to the PR.

## Prerequisites

### 1. GitHub CLI (gh)

The GitHub CLI is required for Claude Code to interact with GitHub PRs.

**Installation:**

```bash
# macOS
brew install gh

# Ubuntu/Debian
sudo apt install gh

# Windows
winget install GitHub.cli
```

**Verify installation:**

```bash
gh --version
```

### 2. GitHub Authentication

Authenticate the GitHub CLI with your account:

```bash
gh auth login
```

Follow the prompts to authenticate via browser or token.

### 3. GitHub Token Permissions

When authenticating, the token needs specific scopes to review PRs:

| Scope | Purpose | Required |
|-------|---------|----------|
| `repo` | Full access to private repositories | Yes (for private repos) |
| `public_repo` | Access to public repositories only | Yes (for public repos only) |
| `write:discussion` | Create/edit PR comments and reviews | Yes |

**For fine-grained tokens (recommended):**

| Permission | Access Level | Purpose |
|------------|--------------|---------|
| Pull requests | Read and write | View PR details, post reviews |
| Contents | Read | Read PR diff and changed files |
| Metadata | Read | Repository metadata |

**Check current authentication:**

```bash
gh auth status
```

**Refresh authentication with required scopes:**

```bash
gh auth refresh -s repo,write:discussion
```

## Reviewing a PR

### Method 1: Direct Review Command

Ask Claude Code to review a PR by number or URL:

```
Review PR #123 and post a code review comment
```

or

```
Review https://github.com/owner/repo/pull/123 and comment on it
```

### Method 2: Step-by-Step

1. **View PR details:**

   ```bash
   gh pr view 123 --json title,body,author,files,commits
   ```

2. **View the diff:**

   ```bash
   gh pr diff 123
   ```

3. **Post a review:**

   ```bash
   gh pr review 123 --approve --body "LGTM! Code looks good."
   ```

   or with comments:

   ```bash
   gh pr review 123 --comment --body "Review comments here..."
   ```

   or request changes:

   ```bash
   gh pr review 123 --request-changes --body "Changes needed..."
   ```

## Review Actions

| Action | Command | Description |
|--------|---------|-------------|
| Approve | `gh pr review N --approve` | Approve the PR |
| Comment | `gh pr review N --comment` | Add review without approval |
| Request Changes | `gh pr review N --request-changes` | Request modifications |
| Add Line Comment | `gh api` (see below) | Comment on specific lines |

### Adding Line-Level Comments

For inline comments on specific lines:

```bash
gh api repos/{owner}/{repo}/pulls/{pr}/comments \
  -f body="Comment text" \
  -f commit_id="SHA" \
  -f path="path/to/file.go" \
  -F line=42 \
  -f side="RIGHT"
```

## Example Review Workflow

When you ask Claude Code to review a PR, it will:

1. **Fetch PR metadata** — Title, description, author, branch info
2. **Fetch the diff** — All changed files and their modifications
3. **Analyze changes** — Security issues, code quality, best practices
4. **Generate review** — Summary with specific findings
5. **Post review** — Submit the review to GitHub

## Troubleshooting

### "Resource not accessible by integration"

Your token lacks required permissions. Re-authenticate:

```bash
gh auth refresh -s repo,write:discussion
```

### "Must have admin rights to Repository"

For organization repos, ensure:

- You have write access to the repository
- Organization allows third-party access (if using OAuth app)

### "Not Found" Error

- Verify the PR number/URL is correct
- Check you have access to the repository
- For private repos, ensure `repo` scope is granted

## Security Considerations

- Claude Code only reads PR content for review purposes
- Review comments are posted under your GitHub identity
- Tokens are stored locally by the `gh` CLI
- Never share tokens or authenticate Claude Code with admin-level tokens

## Related Commands

```bash
# List open PRs
gh pr list

# Check out a PR locally
gh pr checkout 123

# View PR checks/status
gh pr checks 123

# Merge a PR (after approval)
gh pr merge 123
```
