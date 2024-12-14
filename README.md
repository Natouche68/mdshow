# MDShow

A CLI tool that allows you to create beautiful presentations from a single Markdown file.

> In fact you can even run this README as a presentation !

---

## Installation

You can either download a premade binary from the releases, or download it with Go (recommanded) :

```bash
go install github.com/Natouche68/mdshow@latest
```

---

## Usage

In the command line, you can run the tool with the following command :

```bash
mdshow [markdown_file] [theme_name]

# Example
mdshow README.md catppuccin
```

---

Each slide is separated by a horizontal line (`---`).

```md
# This is the slide 1

---

## This is the slide 2
```

---

### Themes

The following themes are available :

- `catppuccin` _(default)_
- `ocean`
- `business`
- `natural`

---

### Building presentations

You can also build presentation into HTML to run them without the CLI tool.

```bash
mdshow build [markdown_file] [theme_name]

# Example
mdshow build README.md catppuccin
```

The resulting HTML file will be named as the markdown file with `.html` extension.
