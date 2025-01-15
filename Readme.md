![screenshot](assets/product-screenshot.png)
## Getting Started

### Getting binary
If you have `go` installed, you can build from source
```sh
go build -o gitContributionGraph cmd/git_contribution_graph/main.go
```
Or download from releses page.

### Github API access token
You can create new access token at https://github.com/settings/tokens, `read:user` scope is enough.

### Dependencies
This version needs `git`, `bash`(if you have git bash that is also ok) and `sed`. Plan for future versions is to only need `git`.

## Usage

### Github Contributions Graph
```sh
gitContributionGraph -token=$GH_CONTRIBUTION_KEY -username=ScaryFrogg
```
By default graph shows current year from beggining until present time. If you want to specify period you can do that by providing `from` and `to` flags using ISO-8601 format.
```sh
gitContributionGraph -token=$GH_CONTRIBUTION_KEY -username=ScaryFrogg -from=2023-01-01T00:00:00Z -to=2023-12-31T23:59:00Z 
```
> **_NOTE:_**  GitHub API doesn't accept period longer than 1 year.

### Local Git reposetories
This is the default mode and you don't need to provide any options, just call binary and graph will be printed if you are inside git repository.
Options `from` and `to` are optional and work the same as in github mode.
Other examples of usages:
- if you want to call graph on demand often you can set alias in your .bashrc or equivalent for other shells:
```bash
alias gcg='gitContributionGraph'
```
- if you want to be able to change dirs with this command and print result after you can do something like
```bash
function cd {
    builtin cd "$@" && gitContributionGraph
}
```
now every time when you cd into repo folder graph will get printed.
