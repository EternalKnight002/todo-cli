````markdown
# 📝 todo-cli

A simple **command-line Todo application** written in [Go](https://go.dev/).  
Tasks are stored locally in a JSON file (`~/.todo/tasks.json` by default, or use the `TODO_FILE` environment variable to override).  

🔗 **Repository:** [EternalKnight002/todo-cli](https://github.com/EternalKnight002/todo-cli)

---

## ✨ Features
- Add tasks with a short description
- List all tasks
- Mark tasks as done
- Edit tasks
- Remove tasks
- Clear all tasks
- Persistent storage in JSON
- Cross-platform (Linux, macOS, Windows)

---

## ⚡ Installation

### Clone & build
```bash
git clone https://github.com/EternalKnight002/todo-cli.git
cd todo-cli
go build -o todo
````

This will produce a binary named `todo` (or `todo.exe` on Windows).

---

## 🚀 Usage

### Add a task

```bash
./todo add "Buy groceries"
```

### List tasks

```bash
./todo list
```

Example output:

```
1) [ ] Buy groceries
2) [x] Finish blog post
    completed: 2025-09-21 17:45
```

### Mark a task done

```bash
./todo do 1
```

### Edit a task

```bash
./todo edit 2 "Finish blog post and publish on GitHub"
```

### Remove a task

```bash
./todo rm 1
```

### Clear all tasks

```bash
./todo clear
```

---

## ⚙️ Storage

By default, tasks are saved to:

* **Linux/macOS:** `~/.todo/tasks.json`
* **Windows:** `%USERPROFILE%\.todo\tasks.json`

You can override this location by setting an environment variable:

```bash
export TODO_FILE=./tasks.json
```

---

## 🛠️ Development

Run locally without building:

```bash
go run main.go add "Test task"
go run main.go list
```

Run tests (once you add them):

```bash
go test ./...
```

---

## 📦 Roadmap

* [ ] Add flags (e.g. `--all`, `--done`, `--pending`)
* [ ] Replace JSON with BoltDB for more robust storage
* [ ] Add interactive TUI with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
* [ ] Package with [GoReleaser](https://goreleaser.com/) for cross-platform binaries

---

## 👤 Author

**Aman** ([@EternalKnight002](https://github.com/EternalKnight002))

---

## 📜 License

This project is licensed under the GNU General Public License v3.0 See the [LICENSE](LICENSE) file for details.
````




