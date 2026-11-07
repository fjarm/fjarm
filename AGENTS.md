# Agent Execution Guidelines

## Context Inheritance (Cascading Rules)
You operate in a hierarchical environment. Before executing any task in a specific subdirectory, you MUST check if that directory or its parents contain a local `AGENTS.md` file.
* **Rule:** Merge local `AGENTS.md` instructions with these global instructions.
* **Conflict Resolution:** Local/subfolder instructions always override global instructions for code implementation details.
* **Submodule Context:** If your plan involves a subfolder (e.g., `/android`), your plan must explicitly state that it is adopting the rules found in `/android/AGENTS.md`.

## Existing Submodule Context
* [android/AGENTS.md](android/AGENTS.md)
* [proto/AGENTS.md](proto/AGENTS.md)

## Core Behavior Boundaries
* **Plan via GitHub Issues**: Before making any code modifications (creating, updating, or deleting code, build, or configuration files), you MUST have an active GitHub Issue that serves as the plan of record.
* **Scope**: Tasks that are purely conversational, investigative, or read-only do not require an issue. Codebase exploration, reading files, and running non-destructive CLI commands to research and establish the plan are permitted before the GitHub Issue is created.
* **Existing vs. New Issues**:
  * If the user references an existing issue (e.g., `#123`), inspect it using `gh issue view 123` and update its plan if necessary.
  * Otherwise, create a new issue using `gh issue create`.
* **Planning via Issue Artifacts**:
  * Create an issue artifact of type `plan` named `<issue-number>-implementation-plan` if one does not already exist on the target issue using: `gh issue artifact create <issue-number> --type plan --name "implementation-plan" --file <path-or-stdin>`
* **Progress Tracking**:
  * When using a Markdown checklist to list tasks, use `- [x]` to indicate a completed task.
  * Use `- [/]` to indicate a blocked task.
  * Use `- [ ]` to indicate an incomplete task.
  * Update the artifact's checklist as tasks are worked on using: `gh issue artifact edit <issue-number> <artifact-id>`
* **Amendments & Pivots**:
  * Update the plan artifact with an appended amendment section detailing the reason for pivot, adjusted goals, and revised checklist using the same format described in the GitHub Issue and Plan Specifications header.
  * Post a brief issue comment alerting reviewers to the change in direction.
* **Fallback**:
  * If the repository or CLI environment does not have issue artifacts enabled, place and maintain the plan directly in the GitHub Issue body.
  * If after retrying, GitHub is unavailable, temporarily use a locally stored and edited markdown file to track the plan.

## GitHub Issue and Plan Specifications
* **Issue Title**: Follow conventional prefix format: `<type>: <short-kebab-or-sentence-description>` (e.g., `feat: add workout rest timer` or `fix: resolve memory leak in routine viewmodel`).
* **Required Issue Body Layout**:
```markdown
## Objective
A few (1–5) sentences summarizing the ultimate goal.

## Research & Context
Findings, relevant files, architectural constraints and/or diagrams, and reference links.
```
* **Required Issue Plan Artifact Layout**:
  * Structure the `plan` artifact with the following layout:
```markdown
## Objective
A few (1–5) sentences summarizing the ultimate goal.

## Research & Context
Findings, relevant files, architectural constraints and/or diagrams, and reference links.

## Tasks Checklist
- [ ] Task 1
- [ ] Task 2

### 1. Task 1 implementation details

### 2. Task 2 implementaiton details

## Verification
Explicit commands and test scenarios to prove success (e.g., unit test commands, UI verification steps).
```

## Pull Request Linking
* Any PR associated with the task must reference the issue in its description (e.g., `Fixes #<issue-number>` or `Closes #<issue-number>`).
