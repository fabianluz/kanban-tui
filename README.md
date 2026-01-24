# Terminal Kanban Board (TUI)

A lightweight, keyboard-centric Kanban board that runs entirely in your terminal. Built with Go and the [Charm](https://charm.sh/) ecosystem (`Bubble Tea`, `Lipgloss`).

PLACEHOLDER OF A SCREENSHOT

## 🚀 Features

* **Keyboard First Navigation:** Vim-like movement (`h`, `j`, `k`, `l`) to traverse columns and tasks.
* **Kanban Workflow:** Move tasks between *To Do*, *In Progress*, and *Done* with a single keystroke.
* **Local Persistence:** Automatically saves your board state to a local JSON file (`board.json`).
* **CRUD Operations:** Create, Read, Update, and Delete tasks directly from the TUI.
* **Export/Backup:** Backup your board data to a separate JSON file with a simple command.
* **Responsive UI:** Automatically adapts column sizes to fit your terminal window.

## 🛠️ Tech Stack

* **Language:** Go (Golang)
* **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea) (The Elm Architecture in Go)
* **Styling:** [Lipgloss](https://github.com/charmbracelet/lipgloss) (CSS for the terminal)
* **Input Handling:** [Bubbles](https://github.com/charmbracelet/bubbles)

## 📦 Installation

Ensure you have [Go installed](https://go.dev/dl/) (version 1.18+ recommended).

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/YOUR_USERNAME/kanban-tui.git](https://github.com/YOUR_USERNAME/kanban-tui.git)
    cd kanban-tui
    ```

2.  **Install dependencies:**
    ```bash
    go mod tidy
    ```

3.  **Run the application:**
    ```bash
    go run .
    ```

4.  **(Optional) Build a binary:**
    ```bash
    go build -o kanban
    ./kanban
    ```

## 🎮 Controls

| Key | Action |
| :--- | :--- |
| **Navigation** | |
| `h` / `←` | Move focus to the **Left** column |
| `l` / `→` | Move focus to the **Right** column |
| `k` / `↑` | Move cursor **Up** |
| `j` / `↓` | Move cursor **Down** |
| **Actions** | |
| `Enter` / `Space` | **Move** selected task to the next column |
| `n` | **New** task (opens input box) |
| `e` | **Edit** selected task |
| `d` / `Backspace` | **Delete** selected task |
| `E` (Shift+e) | **Export** board to `backup_kanban.json` |
| `q` / `Ctrl+c` | **Quit** application |

## 💾 Data Persistence

The application saves your tasks to `board.json` in the same directory where the executable is run.
* **Auto-Save:** Happens automatically whenever you add, move, edit, or delete a task.
* **Backup:** Pressing `E` creates a snapshot named `backup_kanban.json`.

## 🤝 Contributing

Contributions are welcome!
1.  Fork the project.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.