# Syntax Reference (for Authors)

### Basic Formatting

* `**bold text**`
* `*italic*`
* `~~strikethrough~~`
* `==highlighted text==` — green highlight
* `\- highlighted text -\` — same, but red

Neither highlight syntax is part of standard Markdown; both are implemented via custom extensions in `marked-render.ts` (`extension mark` / `markRed`). Nested formatting is supported inside both (`==**bold inside highlight**==` works, and `\- **bold inside** -\` works as well).

Unlike `==...==`, the red highlight safely allows regular hyphens within the text — `\- 5-10% growth per quarter -\` will work correctly because the opening and closing sequences (`\-` and `-\`) are distinct enough not to be confused with a single hyphen. Surrounding spaces are optional (`\-text-\` works too), but including them is recommended for source readability.

---

### Headings

```markdown
# H1
## H2
### H3
#### H4
##### H5
###### H6

```

Every heading automatically receives an `id` (slug) and an anchor icon `#` that appears on hover, allowing users to copy a direct link to the section.

#### Centering

Placing a `$` before the hash marks aligns the heading to the center instead of the default left edge:

```markdown
# Regular heading — left-aligned
$## Centered level 2 heading

```

Works for all levels (`$#`…`$######`) and supports nested formatting (`**bold**`, `code`, etc.) identically to regular headings.

The only difference in rendering is that the anchor icon on centered headings is placed inline before the text rather than absolute-positioned to the left (with centered text, a fixed left offset would shift unpredictably depending on the heading length).

---

### Links and Dividers

```markdown
[link text](https://example.com)

---

```

---

### Lists

```markdown
- item
  - nested item
- item

1. one
2. two

```

**Task Lists (Checkboxes):**

```markdown
- [x] completed
- [ ] incomplete

```

Checkbox styles are implemented using `:has()` selectors and do not depend on whether `marked` adds `task-list-item` classes — they work in any case.

---

### Blockquotes

```markdown
> quote text
>
> second paragraph of the same quote

```

---

### Code

**Inline:** `code`

**Syntax-highlighted block** (`highlight.js`, auto-detects language if omitted):

```markdown
```javascript
const x = 1;

```

```

Every code block automatically receives a language badge and a **Copy** button.

---

### Keyboard Keys

```html
<kbd>Ctrl</kbd> + <kbd>S</kbd>

```

---

### Tables

```markdown
| Column | Center | Right-aligned |
|:---|:---:|---:|
| a | b | c |

```

* Columns aligned to the right automatically use a monospaced font with tabular figures (`font-variant-numeric: tabular-nums`) — convenient for financial data.
* **Color formatting for positive/negative/neutral values** — short syntax `[pos=...]` / `[neg=...]` / `[neut=...]`:

```markdown
| Asset | Yield |
|:---|---:|
| Stocks | [pos=+18.4%] |
| Bonds | [neg=−1.2%] |
| Deposit | [neut=0.0%] |

```

Works beyond tables as well — in any inline context (paragraphs, lists). The value inside `[]` cannot contain the `]` character.

**Horizontal scroll on mobile:** wide tables will not break the layout or compress columns — instead, a horizontal scrollbar appears inside the table itself. Implemented via `renderer.table` in `marked-render.ts`, which wraps the output in `<div class="table-wrapper">`, alongside `white-space: nowrap` on `th`/`td` in CSS.

Due to `nowrap`, lengthy text content in a cell (not just numbers) will not wrap. This is intended behavior for tables with short values (numbers, categories); if a table requires multi-line text descriptions in a specific column, override `white-space: normal` targetedly for that instance.

---

### Formulas (KaTeX)

**Inline** — recommended method via `\(...\)`:

```markdown
Formula \(E = mc^2\) right in the text.

```

**Display (Block)** — via `$$...$$`:

```markdown
$$NPV = \sum_{t=1}^{n} \frac{CF_t}{(1+r)^t} - C_0$$

```

Single `$...$` is also supported (`katexDollarExtension`) with a heuristic: if a dollar sign is immediately followed by a digit, it is treated as a monetary amount (`$100`) rather than a formula and is left untouched. This handles most real-world scenarios (`$CF_t$`, `$x$` are formulas; `$100`, `$1,234` are amounts), though the heuristic is not bulletproof (for instance, a formula starting with a digit like `$2x+1$` will not be recognized).

Using `\(...\)` by default is recommended as it carries no such restrictions; `$...$` is available but secondary.

---

### Charts (Chart.js)

A code block with the language set to `chart` (or `chartjs`) containing a Chart.js JSON configuration:

```markdown
​```chart
{
  "type": "bar",
  "data": {
    "labels": ["Jan", "Feb", "Mar"],
    "datasets": [{ "label": "Sales", "data": [10, 20, 15] }]
  }
}
​```

```

The configuration is a raw `ChartConfiguration` object from Chart.js (`type` / `data` / `options`) without any custom wrapper: anything valid in Chart.js itself is valid here. All standard types work — `bar`, `line`, `pie`, `doughnut`, `radar`, `polarArea`, `scatter`, `bubble` — as well as combinations via `options` (stacked bars, mixed datasets, etc.), as `chart.js/auto` is imported with all controllers pre-registered.

**Positioning** — passed as a second parameter via `:` in the block's language tag:

```markdown
​```chart:left    → chart on the left, text wraps around on the right
​```chart:right   → chart on the right, text wraps around on the left
​```chart:float   → narrow wrapping (280px, "magazine" style for small inline illustrations)
​```chart:full    → full width, breaks the text column
​```chart         → no modifier — full-width block element without text wrapping

```

```markdown
Text before chart...

​```chart:right
{ "type": "line", "data": { ... } }
​```

Text after — wraps around the chart automatically via CSS float, requiring no manual `<div>` wrappers.

```

**Format Limitations:**

* The config is pure JSON, meaning **functions and callbacks are not supported** (`options.plugins.tooltip.callbacks`, custom scale formatters on `scales.ticks.callback`, etc.) — JSON cannot serialize functions. Only declarative Chart.js configurations are available.
* If the JSON is invalid or Chart.js fails to initialize, a block with the `chart-error` class displaying "Invalid chart config" is shown, preventing the post itself from breaking.
* Multiple `chart:left`/`chart:right` blocks placed back-to-back without intervening text will align side-by-side due to float behavior rather than stacking vertically.
* `chart:full` applies `clear: both`, properly clearing preceding floated charts without requiring manual intervention.
* In the admin preview, charts re-render with a debounce on input, but only while the "Preview" tab is active. Typing a long post with many charts may cause a slight delay before they appear when switching to preview mode.

---

### Footnotes

```markdown
Text with a footnote[^label].

[^label]: Footnote explanation.

```

* The label (`label`) can be any unique identifier within the document.
* Definitions can be placed anywhere in the document; ordinal numbers are assigned sequentially based on where the reference appears in the text.
* Referencing a label without a corresponding definition leaves `[^label]` as plain text.
* A definition without a matching reference in the text is silently ignored.

---

### Disclaimers

Three equivalent syntax options are supported — choose based on convenience.

**1. Block syntax via fence** (recommended for standalone warning paragraphs):

```markdown
​```disclaimer
Disclaimer text.
​```

```

**2. Inline syntax within a sentence** — uses the same prefix inside single backticks:

```markdown
Regular text, and here is a `disclaimer important warning` right in the middle of a paragraph.

```

Differs from standard inline code only by the `disclaimer ` prefix (with a trailing space) — a function call like ``disclaimer()`` will not trigger it due to the absence of the space.

**3. Raw HTML** — for complete layout control over the block:

```html
<div class="disclaimer">
Disclaimer text.
</div>

```

Both new variants (fence and inline) do not parse Markdown inside — text is escaped and output as-is, matching how `marked` handles raw HTML block content. If formatting like `**bold**` or links is required inside a disclaimer, use the raw HTML variant.

---

### Images

```markdown
![alt text](url)

```

---

### Raw HTML

Since only admins can publish posts, `marked` does not sanitize HTML — custom `<div>` tags, inline styles, etc., can be inserted directly into `.md` files whenever Markdown capabilities fall short.

---

## Known Limitations

* **Formulas vs. Dollars:** A single `$...$` is not parsed as a formula if immediately followed by a digit (protecting `$100`, `$1,234`). Consequently, formulas starting with a digit (e.g., `$2x+1$`) are not recognized and require `\(...\)`. Edge cases like `$100$150` (two `$` signs with no space) might theoretically trigger false positives, though this rarely occurs in practice.
* **Tables and `white-space: nowrap`:** Table cells do not wrap text (enabling reliable horizontal scrolling on mobile instead of column crushing). Tables containing long text descriptions in cells will extend the scroll area rather than wrapping onto multiple lines.
* **Footnotes inside code blocks:** The `processFootnotes` preprocessor is unaware of ````` block boundaries, meaning literal `[^1]` text inside code examples will still be replaced if matched.
* **Charts — Declarative config only:** The JSON config in `chart` blocks does not support Chart.js functions/callbacks (custom tooltips, axis formatters, etc.) — only what can be expressed with static data.
* **Disclaimers (fence/inline) — No nested Markdown:** Text renders as-is without processing `**bold**`, links, etc. Raw HTML is required for rich formatting inside disclaimers.
* **`:has()` Selectors:** Required for task checkboxes (table wrappers remain unaffected). Supported natively in Chrome, Edge, and Safari; Firefox requires version 121+ (late 2023). On older browser versions, checkboxes will render, but the list item may retain its native bullet point.
* **Bundle Size:** `highlight.js` is included with full language support (~190 languages), and `chart.js/auto` includes all registered chart types. If bundle weight becomes a priority, consider transitioning to `highlight.js/lib/core` with manual language registration and `chart.js` (without `/auto`) with targeted type registration.

---

## Extension Points

Future features can be added via the same pattern using `marked.use({ extensions: [...] })` in `marked-render.ts`:

* **Mermaid** — Diagrams and flowcharts (````mermaid` blocks) for sequence diagrams, architecture flows, and structural charts that fall outside Chart.js data types.
* **Code Line Numbers** — For posts discussing specific lines of code.
* **Diff Highlighting** — `highlight.js` supports the `diff` language out of the box; rendering styles just need verification.
* **Functions/Callbacks in Chart Configs** — If custom formatting (e.g., tooltip formatters) becomes necessary, a `chart-preset` system can be introduced to reference predefined `options` stored in code via short keys instead of passing full JSON configs in Markdown.