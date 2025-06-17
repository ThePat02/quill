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
LABEL start

ELIZABETH: "It is a truth universally acknowledged..." [pride]
DARCY: "That a single man in possession of a good fortune..." [pride]
ELIZABETH: "Must be in want of a wife!"

CHOICE {
    "Ask about Mr. Darcy's fortune" {
        ELIZABETH: "And what is your fortune, sir?"
        DARCY: "My income is considerable."
    },
    "Change the subject" {
        ELIZABETH: "Shall we discuss the weather?"
        GOTO weather
    }
}

LABEL weather
DARCY: "The grounds are quite pleasant today."

CHOICE {
    "Return to start" { GOTO start },
    "Take your leave" { GOTO ending }
}

LABEL ending
ELIZABETH: "Good day, Mr. Darcy."
END
``` 
For a comprehensive example showcasing all available syntax features, please refer to [syntax.q](/examples/syntax.q)!

## Usage
If you are on Linux, you can grab the binary from the latest workflow artifacts. On windows you need to build it from source, which is straightforward.

You can run Quill in your terminal with the following command. Make sure you have the Quill binary in your PATH or specify the path to the binary directly.
```bash
quill [options] <file>
```

### Building
To build Quill from source on Linux, you need to have Go installed. Then, run the following command:

```bash
go build -o quill ./cmd/quill
```