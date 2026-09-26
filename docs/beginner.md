# Hermes for people new to Go

## What you built

Go is a language. `go build` turns your code into **one program file**. That file is `bin/hermes`. It is not an npm package. Other people do not need Node to run it.

```text
source code  →  go build  →  a file named hermes  →  you run that file
```

## Why `hermes doctor` said command not found

Your terminal has a list of folders it searches when you type a name. That list is called **PATH**.

- `./bin/hermes` means: run the file in *this* folder. Always works.
- `hermes` means: search PATH. If the file is not in one of those folders, zsh says `command not found`.

You did not fail. You just have not copied the file onto PATH yet.

## Put Hermes on this Mac (once)

From the Hermes folder:

```bash
make install
```

That copies `bin/hermes` to `~/.local/bin/hermes`.

Then **close the terminal and open a new one**, or run:

```bash
export PATH="$HOME/.local/bin:$PATH"
hermes doctor
```

To make that permanent, add the `export PATH=...` line to `~/.zshrc`.

You do **not** publish to npm. A friend can:

```bash
curl -fsSL https://tryhermes.pages.dev/install | sh
```

## Test without installing globally

```bash
cd /Users/mac/Documents/GitHub/Hermes
./bin/hermes doctor
```

Same program. The `./` is the only difference.
