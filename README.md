# MDShow

A CLI tool that allows you to create beautiful presentations from a single Markdown file.

> In fact you can even run this README as a presentation !

---

## Installation

You can either download a premade binary from the releases, or download it with Go (recommanded) :

```
go install github.com/Natouche68/mdshow@latest
```

---

## Usage

In the command line, you can run the tool with the following command :

```
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
