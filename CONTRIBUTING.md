Hi Emacs, this is -*- mode: gfm -*-

# Contributing to gohacks

## Table of Contents

  * [Contributing to gohacks](#contributing-to-verbuild)
    * [Code of Conduct](#code-of-conduct)
    * [Do I have to read this?](#do-i-have-to-read-this)
    * [What should I know before I get started?](#what-should-i-know-before-i-get-started)
      * [Requirements](#requirements)
        * [CMake](#cmake)
          * [Windows](#windows)
          * [macOS](#macos)
          * [GNU/Linux](#gnulinux)
          * [Visual Studio 2017](#visual-studio-2017)
        * [boost](#boost)
          * [GNU/Linux](#gnulinux-1)
          * [Windows](#windows-1)
    * [How can I contribute?](#how-can-i-contribute)
      * [Reporting bugs](#reporting-bugs)
      * [Feature requests](#feature-requests)
      * [Pull requests](#pull-requests)
    * [Style guides](#style-guides)
      * [Git commit messages](#git-commit-messages)
      * [C   style guide](#c-style-guide)
      * [Documentation style guide](#documentation-style-guide)
    * [Additional notes](#additional-notes)
      * [Issue and PR labels](#issue-and-pr-labels)
        * [Issues](#issues)
        * [Pull request labels](#pull-request-labels)

## Code of Conduct

This project and everyone participating in it is coverned by
the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md).  By
participating, you are expected to uphold this code.  Please report
unacceptable behaviour
to [asmodai@gmail.com](mailto:asmodai@gmail.com).

## Do I have to read this?
> **note:** Please do *not* file an issue to ask a question.  You will
> get faster results by emailing the author.

If you have a general question about this software, please feel free
to email the author.

## What should I know before I get started?

### Requirements

#### Go version
gohacks requires [Go 1.18](https://go.dev/) (or newer) due to the use
of generics.

## How can I contribute?

### Reporting bugs

If you find a bug, please use
the [issue tracker](https://github.com/Asmodai/gohacks/issues) to
report a bug.

Please give as much information as possible, including any relevant
stack traces.

### Feature requests

This project is intended to be a generalised library of gnarly hacks,
but I am open to new ideas.

If you have an idea for an enhancement, please open an issue with the
*enhancement* label.

Other requests will be evaluated on a case-by-case version.

### Pull requests

Please ensure to consider the following when opening a pull request:
 * Do *not* include issue numbers in the PR title.
 * Ensure your code is formatted with `go fmt`.
 * Follow the [documentation style guide](#documentation-style-guide).
 * Ensure you have written tests.
 * End **all** files with a newline.
 
## Style guides
 
### Git commit messages
 * Use the present tense ("Add feature", not "Added feature")
 * Use the imperative mood ("Move thing to..." not "Moves thing to...")
 * Limit the first line to 50 characters or less.
 * Reference issues and pull requests *after* the first line.
 * Install the git hooks in `git-hooks`.
 * Do **not** use emoji in either the title or message.
 
### Go style guide

Please ensure you have formatted your code with `go fmt` before submitting.

### Documentation style guide

> It isn't news that developers don't like documenting their code. But you have
> good reason not to. And if you are documenting code, try to stop! It's not too
> late.

Please try to follow the Go style of documentation, and only document
for API, or where clarity is required.

## Additional notes

### Issue and PR labels

This section lists the labels we use to help us track and manage issues and pull
requests.

#### Issues

Label name | description
-----------|------------
`enhancement` | Feature requests.
`bug` | Confirmed bugs or reports that are likely to be bugs.
`question` | Questions, although the issue tracker is not the place for these.
`feedback` | General feedback.
`more-information-needed` | More information is needs to be collected.
`blocked` | Issues blocked on other issues.
`duplicate` | Issues which are duplicates of other issues.
`wontfix` | An issue that will not be fixed.
`invalid` | Issues that are not valid.

#### Pull request labels

Label name | Description
-----------|------------
`work-in-progress` | PRs that are being worked on with more changes to follow.
`needs-review` | PRs that need code review and final approval.
`under-review` | PRs that are currently being reviewed.
`needs-changes` | PRs that require changes after code review.
`needs-testing` | PRs that require manual testing.
