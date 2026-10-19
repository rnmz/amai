<div align="center">
  <img src="readme/logo.png" alt="logo" width="400">
</div>

# 🌸 Amai

**A lightweight blog designed for simplicity.**

[**View Gallery**](./readme/gallery) | [**Report Bug**](https://github.com/rnmz/amai/issues) | [**Writing Guide**](./docs/)

---

## ✨ Features

  * **Markdown First:** Write your content in Markdown and let Amai handle the rendering.
  * **Extended Articles Capabilities:** Built-in math formulas, text markers, zero-dependency footnotes, and interactive code blocks.
  * **High Compatibility:** Extensive support for standard Markdown specifications.
  * **Blazing Fast:** Powered by Go and Vue for optimal performance.

---

## 🎨 Themes


  * **Classic:** [Gallery](./readme/gallery/classic/) | [Frontend](./amai-classic/)

---

## 🛠️ Tech Stack

Amai is built using a robust and modern set of tools:

| Component              | Technology                               |
|------------------------|------------------------------------------|
| **Backend Framework**  | [Gin Gonic](https://gin-gonic.com/)      |
| **Database Driver**    | [sqlite3](https://pkg.go.dev/modernc.org/sqlite)    |
| **SQL Toolkit**        | [sqlx](https://github.com/jmoiron/sqlx/) |
| **Frontend Framework** | [Vue](https://vuejs.org/)                |
| **Markdown Parser**    | [Marked.js](https://marked.js.org/)      |
| **LaTeX Parser**       | [KaTeX.js](https://katex.org/)           |

---

## 📝 Markdown & Articles Features

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
  * 🇯🇵 **Japanese** (100%, AI translation)

---

## 📋 ToDo List
- [ ] Custom Frontend Themes (Classic, Cyberpunk, Y2K)
- [ ] Add more languages (i18n)
- [ ] Content filters
- [ ] Setup panel (TUI) & Docker setup
- [ ] Advanced layout options for PDF export
