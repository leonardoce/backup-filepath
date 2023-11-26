# Go project structure

This is a possible structure for a Go project exposing web services and a
frontend application.

It is influenced by https://github.com/golang-standards/project-layout

## Development workflow

To start the application you can use:

```shell
make
```

The following command will execute the linters locally, and is designed to be used before filing a pull request:

```shell
make check
```

## How to create your application from this template

Use the script in `scripts/rename.sh` to rename it to your project name:

```shell
./scripts/rename.sh github.com/myname/myproject myproject
```

## How to configure the release workflow

The release workflow need a personal access token to work. The automatic
`GITHUB_TOKEN` access token doesn't allow to start workflow from other
workflows, and this is needed for the image building process.

You can generate a personal access token following [the GitHub
documentation][1].
IMPORTANT: you'll need to use a **classic** personal access token as the
*fine-grained personal access token* don't work correctly. The token should have
permission to read and write from the repository and to manage the workflows.

Once you have your token (it starts with `ghp_`) you need to create a GitHub
action secret named `REPO_PAT`, whose value is the token itself. Do that
following the instructions in [the GitHub
documentation][2].

[1]: https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token#creating-a-personal-access-token-classic
[2]: https://docs.github.com/en/actions/security-guides/encrypted-secrets#creating-encrypted-secrets-for-a-repository