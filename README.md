# vscjsmd
A simple parser from VSCode Snippets to Markdown table for single prefix snippets, useful for *vscode snippets creators*.

#### Install
If you have Go installed:
- `git clone https://github.com/Metalymph/vscjsmd.git`
- inside the root project directory run `go install`

Otherwise, copy the binary from `cmd` directory in your `$PATH` or just use it. There is a version for Linux, Windows and MacOs (arm64, amd64 and universal).

#### Example
Taking a fragment from my lorenzopirro.rust-flash-snippets extension:

```json
{
    "loop": {
        "prefix": "loop",
        "body": [
            "loop {",
            "\t$0",
            "}"
        ],
        "description": "infinite loop"
    },
    "for": {
        "prefix": "for",
        "body": [
            "for ${1:index}, ${2:value} in $0 {",
            "\t",
            "}"
        ],
        "description": "for loop"
    }
}
```

Running in terminal `vscjsmd snippets.code-snippets` creates a `snippets_table.md` file with the following not formatted markdown table:
```
| prefix | description |
| :----- | :---------- |
| for | for loop |
| loop | infinite loop |
```
You can easily copy it on your snippets extension `README.md` file and format it than.