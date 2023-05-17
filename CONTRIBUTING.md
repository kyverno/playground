# Contributing Guidelines

We welcome all contributions, suggestions, and feedback, so please do not hesitate to reach out!

Before you contribute, please take a moment to review and agree to abide by our community [Code of Conduct](/CODE_OF_CONDUCT.md).

- [Contributing Guidelines](#contributing-guidelines)
  - [Engage with us](#engage-with-us)
  - [Ways you can contribute](#ways-you-can-contribute)
    - [1. Report issues](#1-report-issues)
    - [2. Fix or Improve Documentation](#2-fix-or-improve-documentation)
    - [3. Submit Pull Requests](#3-submit-pull-requests)
      - [How to Create a PR](#how-to-create-a-pr)
      - [Developer Certificate of Origin (DCO) Sign off](#developer-certificate-of-origin-dco-sign-off)
  - [Release Processes](#release-processes)

## Engage with us

The Kyverno website has the most updated information on [how to engage with the Kyverno community](https://kyverno.io/community/) including its maintainers and contributors. There are three classes of contributors possible: Contributor, Code Owner, and Maintainer. Please see the [Contributing section on the website](https://kyverno.io/community/#contributing) for the requirements and privileges afforded to each.

Join our community meetings to learn more about Kyverno and engage with other contributors.

## Ways you can contribute

### 1. Report issues

Issues to Kyverno help improve the project in multiple ways including the following:

- Report potential bugs
- Request a feature
- Request a sample policy

### 2. Fix or Improve Documentation

The [Kyverno website](https://kyverno.io), like the main Kyverno codebase, is stored in its own [git repo](https://github.com/kyverno/website). To get started with contributions to the documentation, [follow the guide](https://github.com/kyverno/website#contributing) on that repository.

### 3. Submit Pull Requests

[Pull requests](https://docs.github.com/en/github/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/about-pull-requests) (PRs) allow you to contribute back the changes you've made on your side enabling others in the community to benefit from your hard work. They are the main source by which all changes are made to this project and are a standard piece of GitHub operational flows.

New contributors may easily view all [open issues labeled as good first issues](https://github.com/kyverno/playground/issues?q=is%3Aissue+is%3Aopen+label%3A%22good+first+issue%22) allowing you to get started in an approachable manner.

<!-- Once you wish to get started contributing to the code base, please refer to our [running in development mode](https://github.com/kyverno/kyverno/wiki/Running-in-development-mode) for local setup guide. -->

#### How to Create a PR

Head over to the project repository on GitHub and click the **"Fork"** button. With the forked copy, you can try new ideas and implement changes to the project.

1. **Clone the repository to your device:**

Get the link of your forked repository, paste it in your device terminal and clone it using the command.

```sh
git clone https://hostname/YOUR-USERNAME/YOUR-REPOSITORY
```

2. **Create a branch:**

Create a new brach and navigate to the branch using this command.

```sh
git checkout -b <new-branch>
```

Great, it's time to start hacking! You can now go ahead to make all the changes you want.

3. **Stage, Commit, and Push changes:**

Now that we have implemented the required changes, use the command below to stage the changes and commit them.

```sh
git add .
```

```sh
git commit -s -m "Commit message"
```

The `-s` signifies that you have signed off the commit.

Go ahead and push your changes to GitHub using this command.

```sh
git push
```

#### Developer Certificate of Origin (DCO) Sign off

For contributors to certify that they wrote or otherwise have the right to submit the code they are contributing to the project, we are requiring everyone to acknowledge this by signing their work which indicates you agree to the DCO found [here](https://developercertificate.org/).

To sign your work, just add a line like this at the end of your commit message:

```sh
Signed-off-by: Random J Developer <random@developer.example.org>
```

This can easily be done with the `-s` command line option to append this automatically to your commit message.

```sh
git commit -s -m 'This is my commit message'
```

## Release Processes

Review the [release process](./docs/RELEASE.md).
