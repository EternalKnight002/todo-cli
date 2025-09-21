// main.go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Task struct {
	ID          int64      `json:"id"`
	Title       string     `json:"title"`
	Done        bool       `json:"done"`
	CreatedAt   time.Time  `json:"created_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type Tasks []Task

func tasksFilePath() (string, error) {
	if p := os.Getenv("TODO_FILE"); p != "" {
		return p, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".todo")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	return filepath.Join(dir, "tasks.json"), nil
}

func loadTasks() (Tasks, error) {
	path, err := tasksFilePath()
	if err != nil {
		return nil, err
	}

	// If file doesn't exist, return empty list
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return Tasks{}, nil
	}

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var ts Tasks
	if err := json.Unmarshal(b, &ts); err != nil {
		// backup the corrupted file so user can inspect
		backup := fmt.Sprintf("%s.broken.%d", path, time.Now().Unix())
		_ = os.WriteFile(backup, b, 0o644) // best-effort
		// Inform user and start fresh
		fmt.Fprintf(os.Stderr, "Warning: tasks file corrupted. Backed up to %s and starting with empty list.\n", backup)
		return Tasks{}, nil
	}
	return ts, nil
}

func saveTasks(ts Tasks) error {
	path, err := tasksFilePath()
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(ts, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	// write temp file
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	// atomic move
	return os.Rename(tmp, path)
}

func nextID(ts Tasks) int64 {
	var max int64
	for _, t := range ts {
		if t.ID > max {
			max = t.ID
		}
	}
	return max + 1
}

func findIndexByID(ts Tasks, id int64) int {
	for i, t := range ts {
		if t.ID == id {
			return i
		}
	}
	return -1
}

func cmdAdd(args []string) error {
	_ = args // silence linter if you don't use args directly here
	if len(args) == 0 {
		return errors.New("usage: todo add <task title>")
	}
	title := strings.Join(args, " ")
	ts, err := loadTasks()
	if err != nil {
		return err
	}
	id := nextID(ts)
	t := Task{ID: id, Title: title, Done: false, CreatedAt: time.Now()}
	ts = append(ts, t)
	if err := saveTasks(ts); err != nil {
		return err
	}
	fmt.Printf("Added %d: %s\n", id, title)
	return nil
}

func cmdList(args []string) error {
	_ = args
	ts, err := loadTasks()
	if err != nil {
		return err
	}
	if len(ts) == 0 {
		fmt.Println("No tasks.")
		return nil
	}
	for _, t := range ts {
		check := " "
		if t.Done {
			check = "x"
		}
		fmt.Printf("%d) [%s] %s\n", t.ID, check, t.Title)
		if t.CompletedAt != nil {
			fmt.Printf("    completed: %s\n", t.CompletedAt.Format("2006-01-02 15:04"))
		}
	}
	return nil
}

func cmdDo(args []string) error {
	_ = args
	if len(args) == 0 {
		return errors.New("usage: todo do <id>")
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}
	ts, err := loadTasks()
	if err != nil {
		return err
	}
	i := findIndexByID(ts, id)
	if i == -1 {
		return fmt.Errorf("task %d not found", id)
	}
	if ts[i].Done {
		fmt.Println("Already completed.")
		return nil
	}
	now := time.Now()
	ts[i].Done = true
	ts[i].CompletedAt = &now
	if err := saveTasks(ts); err != nil {
		return err
	}
	fmt.Printf("Marked %d done\n", id)
	return nil
}

func cmdRemove(args []string) error {
	_ = args
	if len(args) == 0 {
		return errors.New("usage: todo rm <id>")
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}
	ts, err := loadTasks()
	if err != nil {
		return err
	}
	i := findIndexByID(ts, id)
	if i == -1 {
		return fmt.Errorf("task %d not found", id)
	}
	ts = append(ts[:i], ts[i+1:]...)
	if err := saveTasks(ts); err != nil {
		return err
	}
	fmt.Printf("Removed %d\n", id)
	return nil
}

func cmdEdit(args []string) error {
	_ = args
	if len(args) < 2 {
		return errors.New("usage: todo edit <id> <new title>")
	}
	id, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}
	newTitle := strings.Join(args[1:], " ")
	ts, err := loadTasks()
	if err != nil {
		return err
	}
	i := findIndexByID(ts, id)
	if i == -1 {
		return fmt.Errorf("task %d not found", id)
	}
	ts[i].Title = newTitle
	if err := saveTasks(ts); err != nil {
		return err
	}
	fmt.Printf("Updated %d\n", id)
	return nil
}

func cmdClear(args []string) error {
	_ = args
	path, err := tasksFilePath()
	if err != nil {
		return err
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return err
	}
	fmt.Println("All tasks cleared.")
	return nil
}

func usage() {
	fmt.Println(`Usage: todo <command> [args]
Commands:
  add <title>       Add a task
  list              List tasks
  do <id>           Mark task done
  rm <id>           Remove task
  edit <id> <title> Edit task title
  clear             Remove all tasks
  help              Show this help`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}
	cmd := os.Args[1]
	args := os.Args[2:]
	var err error
	switch cmd {
	case "add":
		err = cmdAdd(args)
	case "list":
		err = cmdList(args)
	case "do", "complete":
		err = cmdDo(args)
	case "rm", "remove":
		err = cmdRemove(args)
	case "edit":
		err = cmdEdit(args)
	case "clear":
		err = cmdClear(args)
	case "help":
		usage()
		return
	default:
		fmt.Println("Unknown command:", cmd)
		usage()
		return
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
