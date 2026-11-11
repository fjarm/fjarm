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
  * If the user references an existing issue (e.g., `#123`), inspect it using `gh issue view 123` and inspect its plan artifacts using `gh issue artifact list 123 --type plan`.
  * View the content of an existing plan using: `gh issue artifact view <issue-number> <artifact-number>`.
  * Otherwise, create a new issue using `gh issue create`.
* **Planning via Issue Artifacts**:
  * Create an issue artifact of type `plan` named `<issue-number>-short-description-implementation-plan` if one does not already exist on the target issue using:
    `gh issue artifact create <issue-number> --type plan --name "<issue-number>-short-description-implementation-plan" --body-file <file-path>` (Note: pass "-" as the argument to `--body-file` to read from standard input).
* **Progress Tracking**:
  * When using a Markdown checklist to list tasks, use `- [x]` to indicate a completed task.
  * Use `- [/]` to indicate a blocked task.
  * Use `- [ ]` to indicate an incomplete task.
  * Retrieve the `<artifact-number>` from `gh issue artifact list <issue-number>`, and update the checklist as tasks are completed using:
    `gh issue artifact edit <issue-number> <artifact-number> --body-file <file-path>` (Note: pass "-" as the argument to `--body-file` to read from standard input).
* **Amendments & Pivots**:
  * Append an amendment section to the plan artifact detailing the reason for pivot, adjusted goals, and revised checklist (see layout below).
  * Post a brief issue comment alerting reviewers to the change in direction:
    `gh issue comment <issue-number> --body "..."`
* **Fallback**:
  * If the repository or CLI environment does not have issue artifacts enabled, place and maintain the plan directly in the GitHub Issue body.
  * If GitHub is unavailable after retrying, temporarily track the plan in a local Markdown file under `./docs/plans/` named `<issue-number>-short-description-implementation-plan.md`.

## GitHub Issue and Plan Specifications
* **Issue Title**: Follow conventional prefix format: `<type>: <short-kebab-or-sentence-description>` (e.g., `feat: add workout rest timer` or `fix: resolve memory leak in routine viewmodel`).
* **Required Issue Body Layout**:
```markdown
## Objective
A few (1–5) sentences summarizing the ultimate goal.

## Research & Context
High-level summary of findings (including any web links to relevant research), affected components, and link to the plan artifact.
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

### 1. Task 1 Implementation Details

### 2. Task 2 Implementation Details

## Verification
Explicit commands and test scenarios to prove success (e.g., unit test commands, UI verification steps).
```

* **Plan Amendment Layout**:
  * When amending an active plan, append the following section to the bottom of the artifact:
```markdown
---

## [Update: YYYY-MM-DD HH:MM] - Reason for Pivot
* **Context**: Explain why the original plan failed or why an update is necessary (e.g., unexpected error, scope change).
* **Objective**: Summarize the adjusted goal.
* **Changes to Tasks**: Specify which original steps are being abandoned, modified, or added.
* **Revised Task Checklist**:
  - [ ] New or modified step 1
  - [ ] New or modified step 2

### 1. New or Modified Step 1 Implementation Details

### 2. New or Modified Step 2 Implementation Details
```

## Pull Request Linking
* Any PR associated with the task must reference the issue in its description (e.g., `Fixes #<issue-number>` or `Closes #<issue-number>`).