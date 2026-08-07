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
* **Always Plan First**: Before executing any terminal command, code modification, or file creation, you MUST create or update a Markdown planning file.
* **No Inline Execution**: Do not skip the planning phase. If a plan file does not exist for the current task, stop and create it.

## Plan File Specifications
* **File Format**: The file must be exclusively written in standard GitHub Flavored Markdown (GFM).
* **File Naming Convention**: Every plan file name must follow the strict pattern: `PLAN_YYYYMMDD_kebab-case-short-description.md`.
* **Example**: `PLAN_20260727_auth-bug-fix.md`.
* **Time Check**: Before generating a plan, use your closest available system tool (e.g., terminal `date` command) to fetch the exact current timestamp.
* **File Creation**: Write your initial plan to the newly generated timestamped filename using standard GFM.
* **Directory**: Save all plan files inside the `./docs/plans/` directory. Create this directory if it does not exist.

## Required Plan Layout
The top level of the plan Markdown file must explicitly follow this structure:
1. `# Objective` - A few (one to five) sentences summarizing the ultimate goal.
2. `# Research & Context` - Findings, relevant files, relevant links/URLs with summaries, and known constraints.
3. `# Tasks Checklist` - A Markdown checklist (`- [ ]`) breaking down individual, actionable steps.
4. `# Verification` - Explicit commands to run to prove success (e.g., test scripts).

## Plan Amendment Rules (No Overwriting)
* **Never Overwrite**: Once a timestamped plan file is created for a task, you are forbidden from deleting, replacing, or modifying existing content in that file.
* **Append via Sections**: If a plan needs an update, correction, or pivot during execution, append a horizontal separator (`---`) followed by a new Markdown section to the bottom of the file.
* **Timestamp Each Pivot**: The contents of a new amendment section must begin with a level-2 header containing the exact update timestamp and a brief title.

## Format for Plan Updates
* **Plan Update**: When appending an update, use the following exact structure at the bottom of the active file:

---

## [Update: YYYY-MM-DD HH:MM] - Reason for Pivot
* **Context**: Explain why the original plan failed or why an update is necessary (e.g., unexpected error, scope change).
* **Objective**: Summarize the ultimate goal of the new step.
* **Changes to Tasks**: Specify which original steps are being abandoned, modified, or added.
* **Revised Task Checklist**:
    - [ ] New or modified step 1
    - [ ] New or modified step 2

## Status Tracking
* **Complete Task**: When using a Markdown checklist to list tasks, use `- [x]` to indicate a completed task.
* **Blocked Task**: Use `- [/]` to indicate a blocked task.
* **Incomplete Task**: Use `- [ ]` to indicate an incomplete task.
