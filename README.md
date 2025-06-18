<p align="left">
    <img src="assets/banner.png" alt="Quill Banner" width="50%">
</p>

![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/ThePat02/quill/go.yml)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/thepat02/quill)
![GitHub License](https://img.shields.io/github/license/thepat02/quill)


Quill is a simple, lightweight and easy-to-use scripting language designed for interactive fiction and branching dialog written in Go. Dialog is written in a natural, human-readable format, making it accessible for writers and developers alike.

> [!IMPORTANT]
> This project is still in development and is missing some features. If you encounter any issues, please open an issue on GitHub.

```python
MrsBennet: "My dear Mr. Bennet, have you heard that Netherfield Park is let at last?"

CHOICE {
    "Yes." { MrBennet: "I think I have!" } [lie, noncanonical],
    "No." {
        MrBennet: "I have not!"
        GOTO news
    }
}

MrsBennet: "Oh, you haven't! I am quiet sure of it!"
MrBennet: "Well, you caught me, but I cannot believe it is true."

LABEL news

MrsBennet: "But it is, for Mrs. Long has just been here, and she told me all about it."
MrsBennet: "Do you not want to know who has taken it?"
MrBennet: "You want to tell me, and I have no objection to hearing it."

END
``` 
For a comprehensive example showcasing all available syntax features, please refer to [syntax.q](/examples/syntax.q)!

## Usage
If you are on Linux, you can grab the binary from the latest workflow artifacts. On windows you need to build it from source, which is straightforward.

You can run Quill in your terminal with the following command. Make sure you have the Quill binary in your PATH or specify the path to the binary directly.
```bash
quill [flags] <file>
```

### Building
To build Quill from source on Linux, you need to have Go installed. Then, run the following command:

```bash
go build -o quill ./cmd/quill
```

## Utilities
- VS Code Extension: https://github.com/ThePat02/quill-vscode
- Linter (Use the `-p` flag to only parse the file without executing it.)