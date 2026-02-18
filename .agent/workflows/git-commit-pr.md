---
description: Git commit, push, and open a pull request
agent: build
model: openai/gpt-5.1-codex-mini
---

Use the `git-commit` skill to commit changes, push the new branch to the remote repository, and open a pull request.

Before creating the pull request body, detect repository pull request templates and follow them exactly when present:
- `pull_request_template.md`
- `docs/pull_request_template.md`
- `.github/pull_request_template.md`
- `.github/PULL_REQUEST_TEMPLATE/*.md`
- `PULL_REQUEST_TEMPLATE/*.md`
- `docs/PULL_REQUEST_TEMPLATE/*.md`

If a template is found, preserve the template structure and fill every required section in the PR body.
If multiple templates are found, choose the best match for `<user-request>` and include the selected template path in the PR body.

<user-request>
$ARGUMENTS
</user-request>
