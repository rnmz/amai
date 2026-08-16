# 🌸 Amai

**A lightweight, minimalist blog designed for simplicity.**

[**View Gallery**](./readme/gallery) | [**Report Bug**](https://github.com/rnmz/amai/issues) | [**Writing Guide**](./docs/)

---

## ✨ Features

  * **Markdown First:** Write your content in Markdown and let Amai handle the rendering.
  * **Extended Post Capabilities:** Built-in math formulas, text markers, zero-dependency footnotes, and interactive code blocks.
  * **High Compatibility:** Extensive support for standard Markdown specifications.
  * **Blazing Fast:** Powered by Go and Vue for optimal performance.

---

## 🛠️ Tech Stack

Amai is built using a robust and modern set of tools:

| Component | Technology |
| :--- | :--- |
| **Backend Framework** | [Gin Gonic](https://gin-gonic.com/) |
| **Database Driver** | [lib/pq](https://github.com/lib/pq/) |
| **SQL Toolkit** | [sqlx](https://github.com/jmoiron/sqlx/) |
| **Frontend Framework** | [Vue](https://vuejs.org/) |
| **Markdown Parser** | [Marked.js](https://marked.js.org/) |
| **LaTeX Parser** | [KaTeX.js](https://katex.org/) |

---

## 📝 Markdown & Post Features

Amai ensures consistent formatting across standard specs while offering custom enhancements for rich content creation. Check out the [**Writing Guide**](./docs/) for a full example of what you can do in a post.

### Spec Compliance
  * **Markdown 1.0:** `100%`
  * **CommonMark 0.31:** `98%`
  * **GitHub Flavored Markdown (GFM):** `97%`
  * **KaTeX:** `100%`

### Extended Features & Custom Syntax
  * 📐 **Math & Formulas:** Native KaTeX rendering for inline `\(...\)` and block `$$...$$` equations.
  * 🖍️ **Text Highlighting:** Highlight key text using `==custom mark==` syntax.
  * 📌 **Footnotes:** Lightweight native footnote processing (`[^label]`) without external heavy libraries.
  * 💻 **Enhanced Code Blocks:** Syntax highlighting powered by Highlight.js with automatic language badges and a **Copy** button.
  * 📊 **Smart Financial Tables:** Monospace tabular numbers (`tabular-nums`) for right-aligned columns and custom state styling (`.num-positive`, `.num-negative`).

---

## 🌍 Localization

We currently support the following languages:

  * 🇺🇸 **English** (100%)
  * 🇷🇺 **Russian** (100%)
  * 🇯🇵 **Japanese** (100%)

---

## 📋 ToDo list
- [x] Dark theme
- [ ] Better design
- [ ] More frontend themes
- [ ] Add more languages (i18n)
- [ ] Content filters
- [ ] Setup panel (TUI) & hard way (docker)
- [ ] Advanced layout options for PDF export

---

### 🚀 Getting Started

Easy way: WiP

---

Hard way (docker): WiP

---

Hard way (without docker):

1. **Clone the repo**
    ```bash
    git clone https://github.com/rnmz/amai.git
    ```

2. **Install packages**

    Caddy:
    ```bash
    sudo apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | sudo gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
    curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | sudo tee /etc/apt/sources.list.d/caddy-stable.list
    sudo chmod o+r /usr/share/keyrings/caddy-stable-archive-keyring.gpg
    sudo chmod o+r /etc/apt/sources.list.d/caddy-stable.list
    sudo apt update
    sudo apt install caddy
    ```

    Go (1.26.4):
    ```bash
    sudo rm -rf /usr/local/go
    curl -OL https://go.dev/dl/go1.26.4.linux-amd64.tar.gz
    sudo tar -C /usr/local -xzf go1.26.4.linux-amd64.tar.gz
    echo 'export PATH=$PATH:/usr/local/go/bin:$(go env GOPATH)/bin' >> ~/.bashrc
    source ~/.bashrc
    ```

    Node (nvm):
    ```bash
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.40.6/install.sh | bash
    ```

3. **Install dependencies**
    ```bash
    go mod tidy
    cd ./amai-frontend
    npm install
    ```

4. **Configure your server**

    Create a `.env` file based on `.env.example`:
    ```bash
    cp .env.example .env
    ```

    *Note: All variables are required for the server to start correctly!*

5. **Configure Caddy**

    Point `/etc/caddy/Caddyfile` at your backend (adjust the domain and port to match your `.env`):
    ```caddyfile
    your-domain.com {
        reverse_proxy localhost:443
    }
    ```
    Then reload Caddy:
    ```bash
    sudo systemctl reload caddy
    ```

6. **Build and run**
    ```bash
    go build -o amai .
    ./amai
    cd ./amai-frontend
    npm run build
    ```
