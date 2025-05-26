# tg-chat-summary

`tg-chat-summary` is a Go application that connects to your Telegram account, allows you to select a chat, and then uses a provided Ollama instance to generate a summary of the recent messages in that chat. It features a terminal user interface (TUI) for easy navigation and interaction.

## Features

- Securely connect to Telegram using your API ID and Hash, utilizing opensource telegram client "tdlib"
- List your recent chats.
- Select a chat to view its summary.
- Fetches recent messages from the selected chat.
- Generates a summary of the chat messages using a specified Ollama model.
- Interactive TUI for navigation.

## Prerequisites

Before you begin, ensure you have met the following requirements:

- **Go**: Version 1.18 or higher.
- **TDLib**: The Telegram Database Library.
- **Ollama**: A running Ollama instance with a downloaded model (e.g., Llama3.2). You can find installation instructions for Ollama [here](https://ollama.com/).

## Installation

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd tg-chat-summary
    ```

2.  **Install TDLib:**
    The project includes a script to download and build TDLib.
    ```bash
    ./install-tdlib.sh
    ```
    This script will clone the TDLib repository, build it, and install the necessary libraries into the `td/tdlib` directory within the project. Ensure you have `cmake` and a C++ compiler installed.

3.  **Build the application:**
    ```bash
    go build .
    ```
    This will create an executable named `tg-chat-summary` (or `tg-chat-summary.exe` on Windows).

## Environment Variables

The application requires the following environment variables to be set:

-   `TG_API_ID`: Your Telegram API ID. You can obtain this from [my.telegram.org](https://my.telegram.org/apps).
-   `TG_API_HASH`: Your Telegram API Hash, also obtained from [my.telegram.org](https://my.telegram.org/apps).
-   `OLLAMA_HOST`: The base URL of your running Ollama instance (e.g., `http://localhost:11434`).


Example:
```bash
export TG_API_ID="your_api_id"
export TG_API_HASH="your_api_hash"
export OLLAMA_HOST="http://localhost:11434"
```

## Usage

1.  **Set the environment variables** as described above.
2.  **Run the application:**
    ```bash
    ./tg-chat-summary
    ```

### TUI Navigation

-   **Arrow Up/Down (↑/↓)**: Navigate through the list of chats.
-   **Enter**: Select the highlighted chat to generate and view its summary.
-   **Esc**:
    -   When viewing a chat summary, press `Esc` to return to the chat list.
-   **q**: Quit the application from the chat list view.

## How it Works

1.  The application starts and initializes the TDLib client using your API credentials.
2.  It authenticates with Telegram. You might be prompted for your phone number, code, and 2FA password in the terminal.
3.  Once authenticated, it fetches a list of your recent chats.
4.  The TUI displays the list of chats.
5.  When you select a chat:
    a.  It fetches the recent message history for that chat.
    b.  It constructs a prompt with these messages.
    c.  It sends the prompt to your Ollama instance via the `OLLAMA_HOST`.
    d.  The Ollama model generates a summary.
    e.  The summary is displayed in the TUI.

## Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue for any bugs, features, or improvements.

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.
