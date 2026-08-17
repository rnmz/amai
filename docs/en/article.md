# Syntax Guide (for Authors)

### Basic Formatting

* `**bold text**`
* `*italic*`
* `~~strikethrough~~`
* `{{mark:color text}}` — highlight (background color)

The highlight feature is not part of standard Markdown; it is implemented as a custom extension in `marked-render.ts` (`extension mark`). Specify the color immediately after `mark:` with no space before the colon.

There are two ways to set the color:

**1. Curated Palette** — `yellow`, `red`, `green`, `blue`, `orange`, `purple`. Each has a carefully matched background and text color pair in the stylesheet (for example, `red` matches the text shade of `[neg=...]` in tables). Use these colors when you need semantic highlighting (error, success, warning, etc.):

```markdown
{{mark:yellow important text}}
{{mark:red critical value}}
{{mark:green increased by 12%}}

```

`yellow` is the "default" highlight without an additional CSS class (similar to a classic `<mark>` tag). In the `amai-classic` theme, it renders as a light green shade.

**2. Any Custom CSS Color** — named (`aquamarine`, `tomato`, `cornflowerblue`, etc.) or hex (`#7fffd4`) — whenever you need a color outside the curated palette:

```markdown
{{mark:aquamarine text}}
{{mark:#ff9900 text}}

```

⚠️ **Important difference:** Custom colors change **only the background**; the text color remains standard (matching the main paragraph text) as contrast is not automatically adjusted for dark backgrounds. If you pick a color that is too dark or too saturated, the text might become illegible. Stick to light/pastel shades, or rely on the curated palette from option 1 for guaranteed readability.

If an unparsed color is provided (typos, invalid CSS strings, etc.), it silently falls back to the default yellow highlight so that content is never lost.

Nested formatting (bold text, links, and even inline code wrapped in single backticks) is fully supported inside both color methods:

```markdown
{{mark:red **important**: see `npm run build` and the [documentation](https://example.com)}}
{{mark:aquamarine supports **bold text** and [links](https://example.com)}}

```

⚠️ **Limitation:** The closing sequence is always the **first** occurrence of `}}` after the opening `{{mark:`. If the inner text contains `}}` (e.g., raw JSON or templates with double curly braces), the highlight will close prematurely. This rarely affects standard post content.

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

Every heading automatically receives a unique `id` (slug) and an anchor icon `#` on hover, allowing users to copy direct links to sections.

#### Centered Headings

Prefixing the hash symbols with `$` centers the heading alignment instead of the default left-aligned text:

```markdown
# Regular heading — left-aligned
$## Centered level-2 heading

```

This works across all levels (`$#`…`$######`) and supports nested formatting (`**bold**`, `code`, etc.) just like standard headings.

The only structural difference is that centered headings place the anchor icon inline before the text instead of using absolute positioning on the left side (since absolute left positioning breaks visually depending on heading length).

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
- [ ] uncompleted

```

Checkbox styles rely on CSS `:has()` selectors and work universally, regardless of whether `marked` appends `task-list-item` classes.

---

### Blockquotes

```markdown
> Blockquote text
>
> Second paragraph in the same blockquote

```

---

### Code

**Inline:** `code`

**Syntax-highlighted code block** (`highlight.js`, auto-detects language if omitted):

```markdown
```javascript
const x = 1;

```

```

Each code block automatically displays a language badge alongside a **Copy** button.

---

### Keybindings

```html
<kbd>Ctrl</kbd> + <kbd>S</kbd>

```

---

### Tables

```markdown
| Column | Centered | Right-aligned |
|:---|:---:|---:|
| a | b | c |

```

* Right-aligned columns automatically use a monospaced font with tabular numbers (`font-variant-numeric: tabular-nums`), optimized for financial data.
* **Semantic status indicators for positive/negative/neutral values** — shorthand syntax: `[pos=...]` / `[neg=...]` / `[neut=...]`:

```markdown
| Asset | Returns |
|:---|---:|
| Stocks | [pos=+18.4%] |
| Bonds | [neg=−1.2%] |
| Savings | [neut=0.0%] |

```

This shorthand works inline across all contexts (paragraphs, lists, etc.). Note that values inside `[]` cannot contain a closing bracket `]`.

**Mobile Horizontal Scroll:** Wide tables feature an inline container scrollbar rather than compressing column widths or breaking layout responsiveness. This is handled via `renderer.table` in `marked-render.ts`, which wraps table outputs inside a `<div class="table-wrapper">` combined with CSS `white-space: nowrap` on `th`/`td` elements.

Due to `white-space: nowrap`, text content within cells will not wrap to a new line. While ideal for short numeric or categorical data, you should override this rule locally (`white-space: normal`) for columns containing lengthy multi-line text descriptions.

---

### Math Formulas (KaTeX)

**Inline** — recommended approach using `\(...\)`:

```markdown
Inline formula \(E = mc^2\) inside a sentence.

```

**Block level** — using `$$...$$`:

```markdown
$$
NPV = \sum_{t=1}^{n} \frac{CF_t}{(1+r)^t} - C_0
$$

```

Single dollar signs `$...$` are supported (`katexDollarExtension`) using a heuristic: if a number directly follows the opening `$`, it is interpreted as a monetary value (`$100`) rather than a LaTeX expression. This correctly resolves most edge cases (`$CF_t$`, `$x$` render as math; `$100`, `$1,234` render as currency), though edge cases starting math with digits (`$2x+1$`) will not match.

Prefer `\(...\)` as the default standard to avoid parsing quirks; treat `$...$` as an alternative syntax option.

---

### Charts (Chart.js)

Create a code block using the language identifier `chart` (or `chartjs`) containing a valid Chart.js JSON configuration:

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

The payload is a standard Chart.js `ChartConfiguration` object (`type` / `data` / `options`). Any native Chart.js configurations work out-of-the-box (`bar`, `line`, `pie`, `doughnut`, `radar`, `polarArea`, `scatter`, `bubble`, stacked bar configurations, mixed dataset types, etc.) as `chart.js/auto` is loaded with all built-in controllers.

**Layout and Alignment** — set using a colon `:` modifier on the code block language:

```markdown
​```chart:left    → Floats left; text wraps around the right side
​```chart:right   → Floats right; text wraps around the left side
​```chart:float   → Compact float (280px width) for magazine-style inline illustrations
​```chart:full    → Full-width block element; breaks text columns
​```chart         → Default block element; full width without text wrapping

```

```markdown
Text preceding the chart...

​```chart:right
{ "type": "line", "data": { ... } }
​```

Text following the chart automatically wraps around it via CSS float without requiring manual `<div>` wrappers.

```

**Format Limitations:**

* Configurations must be valid JSON: **functions and callbacks are unsupported** (`options.plugins.tooltip.callbacks`, `scales.ticks.callback`, etc.) because JSON cannot serialize functions. Only declarative configurations are supported.
* Invalid JSON or initialization failures render a graceful fallback container with class `chart-error` and the message "Invalid chart config", preventing runtime post crashes.
* Sequential `chart:left`/`chart:right` blocks without intervening text will stack horizontally (standard CSS float behavior).
* `chart:full` applies `clear: both`, correctly clearing any preceding floated charts.
* In the admin preview pane, charts re-render with a debounce on every edit while the "Preview" tab is active. Complex posts with numerous charts may experience slight rendering delays during rapid edits.

---

### Footnotes

```markdown
Text containing a footnote link[^label].

[^label]: Footnote detailed explanation.

```

* Labels are arbitrary strings, but must remain unique within the document context.
* Definitions can appear anywhere in the file; sequential numbering is resolved by order of appearance in text.
* Unmatched labels (`[^label]` without a corresponding definition) render as plain text.
* Unused definitions (without text links) are silently discarded.

---

### Disclaimers

Choose from three supported syntax approaches based on preference:

**1. Fenced Block** (recommended for dedicated warning callouts):

```markdown
​```disclaimer
Disclaimer callout text.
​```

```

**2. Inline Variant** — using backticks with the `disclaimer` prefix:

```markdown
Standard paragraph text, with a `disclaimer inline warning callout` embedded inside.

```

Differentiated from standard inline code by the `disclaimer ` prefix (including trailing space). Standard code snippets like ``disclaimer()`` will not trigger this block as they lack the trailing space.

**3. Raw HTML** — for complete layout control:

```html
<div class="disclaimer">
Disclaimer text.
</div>

```

Fenced and inline variants escape inner content and do not parse nested Markdown. If you require **bold text**, links, or other inline styles inside a disclaimer block, use the raw HTML pattern.

---

### Images

```markdown
![alt text](url)

```

---

### Raw HTML

Admin content is not sanitized by `marked` by default. You can embed arbitrary HTML tags (`<div>`, inline CSS, etc.) directly into `.md` files whenever Markdown capabilities fall short.

---

## Known Limitations

* **Math Expressions vs. Dollar Signs:** Single dollar signs `$...` followed directly by digits are treated as currency (`$100`, `$1,234`). Math formulas starting with a digit (`$2x+1$`) will not match; use `\(...\)` instead. Adjacent dollar amounts (`$100$150`) may trigger a false positive, though this pattern is rare in practice.
* **Table Non-Wrapping (`white-space: nowrap`):** Table cells enforce single-line rendering to preserve mobile scroll behavior. Columns containing lengthy prose will widen table bounds rather than wrapping text lines.
* **Footnotes in Code Blocks:** The `processFootnotes` preprocessor is context-agnostic to backtick fences; literal strings matching `[^1]` inside code blocks will be transformed.
* **Declarative Charts Only:** Chart JSON parameters cannot process JavaScript functions or dynamic callbacks (custom tooltips, tick formatters, etc.).
* **Fenced/Inline Disclaimers:** Native Markdown parsing is disabled inside fenced and inline disclaimers. Use raw HTML if rich formatting is required inside a disclaimer.
* **Custom Highlights Lack Auto-Contrast:** Custom CSS highlight colors (`{{mark:aquamarine text}}`, `#7fffd4`, etc.) alter background colors while preserving default text color. Select light or pastel values to ensure readability, or stick to curated palette keywords (`yellow`, `red`, `green`, `blue`, `orange`, `purple`).
* **`:has()` CSS Selector Requirements:** Form controls like list checkboxes rely on modern browser CSS support (Chrome/Edge/Safari support this natively; Firefox requires v121+, released late 2023). Older browsers render functional elements, but may retain native bullet points.
* **Bundle Overhead:** `highlight.js` bundles language definitions (~190 languages) and `chart.js/auto` includes default controllers. Optimize bundle weight by migrating to `highlight.js/lib/core` alongside targeted `chart.js` modular imports if necessary.

---

## Extension Points

To introduce additional markup extensions, register custom parsers via `marked.use({ extensions: [...] })` in `marked-render.ts`:

* **Mermaid Diagrams** — render flowcharts, sequence diagrams, and architecture maps using ````mermaid` blocks.
* **Code Line Numbers** — display line numbering for code walkthroughs.
* **Diff Highlighting** — leverage `highlight.js` native `diff` language parsing with custom diff view styling.
* **Chart Presets** — define reusable `chart-preset` identifiers containing complex pre-configured JavaScript functions in the codebase, enabling short code block aliases inside Markdown files.