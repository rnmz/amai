import { marked, type Token } from 'marked'
import hljs from 'highlight.js'
import katex from 'katex'
import { Chart, type ChartConfiguration } from 'chart.js/auto'
import { watch } from 'vue'
import { useTheme } from '@/composables/useTheme'

let hljsLinkEl: HTMLLinkElement | null = null

function setHljsTheme(theme: 'light' | 'dark') {
  const href =
    theme === 'dark'
      ? '/styles/github-dark.min.css'
      : '/styles/github.min.css'

  if (!hljsLinkEl) {
    hljsLinkEl = document.createElement('link')
    hljsLinkEl.rel = 'stylesheet'
    document.head.appendChild(hljsLinkEl)
  }
  hljsLinkEl.href = href
}

export function initHljsTheme(): void {
  const { theme } = useTheme()
  watch(theme, (val) => setHljsTheme(val), { immediate: true })
}

function readCssVar(name: string, fallback: string): string {
  if (typeof window === 'undefined') return fallback
  const value = getComputedStyle(document.documentElement).getPropertyValue(name).trim()
  return value || fallback
}

function isDarkVariant(): boolean {
  return document.documentElement.getAttribute('data-theme') === 'dark'
}

function applyChartDefaults() {
  const textColor = readCssVar('--md-text', '#e2f1f8')
  const gridColor = readCssVar('--md-border-light', 'rgba(0, 240, 255, 0.18)')
  const accentColor = readCssVar('--md-accent', '#00f0ff')
  const linkColor = readCssVar('--md-link', '#ff0055')
  const bgColor = readCssVar('--color-bg', '#050811')

  Chart.defaults.font.family = "'Chakra Petch', system-ui, -apple-system, sans-serif"
  Chart.defaults.font.size = 12
  Chart.defaults.font.weight = 'bold'
  Chart.defaults.color = textColor
  Chart.defaults.borderColor = gridColor

  Chart.defaults.plugins.legend.labels.color = textColor
  Chart.defaults.plugins.legend.labels.font = {
    family: "'Orbitron', sans-serif",
    size: 12,
    weight: 'bold',
  }

  Chart.defaults.plugins.tooltip.backgroundColor = bgColor
  Chart.defaults.plugins.tooltip.titleColor = accentColor
  Chart.defaults.plugins.tooltip.bodyColor = textColor
  Chart.defaults.plugins.tooltip.borderColor = linkColor
  Chart.defaults.plugins.tooltip.borderWidth = 1
  Chart.defaults.plugins.tooltip.padding = 10

  Chart.defaults.elements.bar.borderWidth = 2
  Chart.defaults.elements.line.borderWidth = 3
  Chart.defaults.elements.point.radius = 5
  Chart.defaults.elements.point.hoverRadius = 7
}

export function initChartTheme(): void {
  const { theme } = useTheme()

  watch(
    theme,
    () => {
      applyChartDefaults()
      document.querySelectorAll<HTMLCanvasElement>('canvas[data-chart-config]').forEach((canvas) => {
        const existing = chartInstances.get(canvas)
        if (existing) {
          existing.destroy()
          chartInstances.delete(canvas)
        }
      })
      initCharts()
    },
    { immediate: true },
  )
}

export function initMarkdownTheming(): void {
  initHljsTheme()
  initChartTheme()
}

function slugify(text: string): string {
  return text
    .toString()
    .toLowerCase()
    .trim()
    .replace(/[^\p{L}\p{N}\s-]/gu, '')
    .replace(/\s+/g, '-')
}

function escapeHtml(text: string): string {
  return text
    .replace(/&/g, '&amp;')
    .replace(/"/g, '&quot;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
}

const renderer = new marked.Renderer()
const originalTable = marked.Renderer.prototype.table

renderer.table = function (this: typeof renderer, token: Parameters<typeof originalTable>[0]): string {
  const html = originalTable.call(this, token)
  return `<div class="table-wrapper">${html}</div>`
}

function indexOfOrUndefined(src: string, needle: string): number | undefined {
  const idx = src.indexOf(needle)
  return idx === -1 ? undefined : idx
}

const MARK_COLORS = new Set(['yellow', 'red', 'green', 'blue', 'orange', 'purple'])

const NAMED_CSS_COLOR_RE = /^[a-zA-Z]{3,32}$/
const HEX_CSS_COLOR_RE = /^#(?:[0-9a-fA-F]{3}|[0-9a-fA-F]{4}|[0-9a-fA-F]{6}|[0-9a-fA-F]{8})$/

function isSafeCssColor(value: string): boolean {
  return NAMED_CSS_COLOR_RE.test(value) || HEX_CSS_COLOR_RE.test(value)
}

const markExtension = {
  name: 'mark',
  level: 'inline' as const,
  start(src: string) {
    return indexOfOrUndefined(src, '{{mark:')
  },
  tokenizer(this: { lexer: { inlineTokens: (s: string) => Token[] } }, src: string) {
    const rule = /^\{\{mark:([#a-zA-Z0-9]+)[ \t]+([\s\S]+?)\}\}/
    const match = rule.exec(src)
    const color = match?.[1]
    const text = match?.[2]
    if (match && color && text) {
      return {
        type: 'mark',
        raw: match[0],
        color,
        tokens: this.lexer.inlineTokens(text),
      }
    }
    return undefined
  },
  renderer(
    this: { parser: { parseInline: (tokens: Token[]) => string } },
    token: { color: string; tokens: Token[] },
  ) {
    const inner = this.parser.parseInline(token.tokens)

    if (MARK_COLORS.has(token.color)) {
      const classAttr = token.color === 'yellow' ? '' : ` class="mark-${token.color}"`
      return `<mark${classAttr}>${inner}</mark>`
    }

    if (isSafeCssColor(token.color)) {
      return `<mark class="mark-custom" style="--mark-custom-bg: ${token.color}">${inner}</mark>`
    }
    return `<mark>${inner}</mark>`
  },
}

const coloredValueExtension = {
  name: 'coloredValue',
  level: 'inline' as const,
  start(src: string) {
    return src.match(/\[(pos|neg|neut)=/)?.index
  },
  tokenizer(src: string) {
    const rule = /^\[(pos|neg|neut)=([^\]]+)\]/
    const match = rule.exec(src)
    const kind = match?.[1] as 'pos' | 'neg' | 'neut' | undefined
    const value = match?.[2]
    if (match && kind && value) {
      return {
        type: 'coloredValue',
        raw: match[0],
        kind,
        text: value,
      }
    }
    return undefined
  },
  renderer(token: { kind: 'pos' | 'neg' | 'neut'; text: string }) {
    const classMap = {
      pos: 'num-positive',
      neg: 'num-negative',
      neut: 'num-neutral',
    } as const
    return `<span class="${classMap[token.kind]}">${escapeHtml(token.text)}</span>`
  },
}

renderer.code = ({ text, lang }: { text: string; lang?: string; escaped?: boolean }): string => {
  const langToken = (lang || '').split(' ')[0] ?? ''
  const [requestedLang = '', direction] = langToken.split(':')

  if (requestedLang === 'chart' || requestedLang === 'chartjs') {
    const escaped = escapeHtml(text)

    const validDirections = ['left', 'right', 'float', 'full']
    const directionClass =
      direction && validDirections.includes(direction)
        ? ` chart-block-wrapper--${direction}`
        : ''

    return `
<div class="chart-block-wrapper${directionClass}">
  <canvas class="chart-canvas" data-chart-config="${escaped}"></canvas>
</div>`
  }

  if (requestedLang === 'disclaimer') {
    const escaped = escapeHtml(text.trim())
    return `
<div class="disclaimer">
${escaped}
</div>`
  }

  let highlighted: string
  let displayLang: string

  if (requestedLang && hljs.getLanguage(requestedLang)) {
    highlighted = hljs.highlight(text, { language: requestedLang }).value
    displayLang = requestedLang
  } else {
    const auto = hljs.highlightAuto(text)
    highlighted = auto.value
    displayLang = auto.language || 'text'
  }

  return `
<div class="code-block-wrapper">
  <div class="code-block-header">
    <span class="code-lang">${displayLang}</span>
    <button type="button" class="copy-button" data-copy-target>Copy</button>
  </div>
  <pre><code class="hljs language-${displayLang}">${highlighted}</code></pre>
</div>`
}

renderer.codespan = ({ text }: { text: string }): string => {
  const prefix = 'disclaimer '
  if (text.startsWith(prefix)) {
    const content = text.slice(prefix.length).trim()
    return `<span class="disclaimer disclaimer-inline">${escapeHtml(content)}</span>`
  }

  return `<code>${escapeHtml(text)}</code>`
}

renderer.heading = ({ text, depth }: { text: string; depth: number }): string => {
  const slug = slugify(text)
  return `
<h${depth} id="${slug}">
  <a class="heading-anchor" href="#${slug}">#</a>${text}
</h${depth}>`
}

const centeredHeadingExtension = {
  name: 'centeredHeading',
  level: 'block' as const,
  start(src: string) {
    return src.match(/^\$#{1,6}\s/m)?.index
  },
  tokenizer(this: { lexer: { inline: (s: string) => Token[] } }, src: string) {
    const rule = /^\$(#{1,6})[ \t]+(.*)(?:\n|$)/
    const match = rule.exec(src)
    const hashes = match?.[1]
    const text = match?.[2]
    if (match && hashes && text) {
      return {
        type: 'centeredHeading',
        raw: match[0],
        depth: hashes.length,
        text: text.trim(),
        tokens: this.lexer.inline(text.trim()),
      }
    }
    return undefined
  },
  renderer(
    this: { parser: { parseInline: (tokens: Token[]) => string } },
    token: { depth: number; text: string; tokens: Token[] },
  ) {
    const slug = slugify(token.text)
    const html = this.parser.parseInline(token.tokens)
    return `
<h${token.depth} id="${slug}" class="heading-centered"><a class="heading-anchor heading-anchor-inline" href="#${slug}">#</a>${html}</h${token.depth}>`
  },
}

const katexBlockExtension = {
  name: 'katexBlock',
  level: 'block' as const,
  start(src: string) {
    return src.match(/\$\$/)?.index
  },
  tokenizer(src: string) {
    const rule = /^\$\$([\s\S]+?)\$\$/
    const match = rule.exec(src)
    const text = match?.[1]
    if (match && text) {
      return {
        type: 'katexBlock',
        raw: match[0],
        text: text.trim(),
      }
    }
    return undefined
  },
  renderer(token: { text: string }) {
    try {
      return katex.renderToString(token.text, { throwOnError: false, displayMode: true })
    } catch {
      return `<pre>${token.text}</pre>`
    }
  },
}

const katexInlineExtension = {
  name: 'katexInline',
  level: 'inline' as const,
  start(src: string) {
    return src.match(/\\\(/)?.index
  },
  tokenizer(src: string) {
    const rule = /^\\\(([\s\S]+?)\\\)/
    const match = rule.exec(src)
    const text = match?.[1]
    if (match && text) {
      return {
        type: 'katexInline',
        raw: match[0],
        text: text.trim(),
      }
    }
    return undefined
  },
  renderer(token: { text: string }) {
    try {
      return katex.renderToString(token.text, { throwOnError: false, displayMode: false })
    } catch {
      return `<code>${token.text}</code>`
    }
  },
}

const katexDollarExtension = {
  name: 'katexDollar',
  level: 'inline' as const,
  start(src: string) {
    return src.match(/\$[^\d\s]/)?.index
  },
  tokenizer(src: string) {
    const rule = /^\$([^\$\n]+?)\$/
    const match = rule.exec(src)
    const text = match?.[1]
    if (match && text && !/^\d/.test(text)) {
      return {
        type: 'katexDollar',
        raw: match[0],
        text: text.trim(),
      }
    }
    return undefined
  },
  renderer(token: { text: string }) {
    try {
      return katex.renderToString(token.text, { throwOnError: false, displayMode: false })
    } catch {
      return `<code>${token.text}</code>`
    }
  },
}

marked.use({
  renderer,
  extensions: [
    centeredHeadingExtension,
    markExtension,
    coloredValueExtension,
    katexBlockExtension,
    katexInlineExtension,
    katexDollarExtension,
  ],
})

export function initCopyButtons(root: HTMLElement | Document = document): void {
  root.addEventListener('click', async (e: Event) => {
    const target = e.target as HTMLElement
    const btn = target.closest<HTMLButtonElement>('[data-copy-target]')
    if (!btn) return

    const wrapper = btn.closest('.code-block-wrapper')
    const codeEl = wrapper?.querySelector('pre code')
    if (!codeEl) return

    try {
      await navigator.clipboard.writeText(codeEl.textContent ?? '')
      const original = btn.textContent
      btn.textContent = 'Copied!'
      setTimeout(() => {
        btn.textContent = original
      }, 1500)
    } catch (err) {
      console.error('Failed to copy code:', err)
    }
  })
}

const chartInstances = new WeakMap<HTMLCanvasElement, Chart>()

export function initCharts(root: HTMLElement | Document = document): void {
  const isDark = isDarkVariant()

  const fillColors = isDark
    ? ['rgba(0, 240, 255, 0.65)', 'rgba(255, 0, 85, 0.65)', 'rgba(255, 230, 0, 0.65)', 'rgba(0, 255, 102, 0.65)', 'rgba(176, 0, 255, 0.65)']
    : ['rgba(255, 0, 85, 0.75)', 'rgba(0, 163, 173, 0.75)', 'rgba(217, 194, 0, 0.75)', 'rgba(0, 179, 71, 0.75)', 'rgba(138, 0, 204, 0.75)']

  const strokeColors = isDark
    ? ['#00f0ff', '#ff0055', '#ffe600', '#00ff66', '#b000ff']
    : ['#ff0055', '#00a3ad', '#d9c200', '#00b347', '#8a00cc']

  const canvases = root.querySelectorAll<HTMLCanvasElement>('canvas[data-chart-config]')

  canvases.forEach((canvas) => {
    if (chartInstances.has(canvas)) return

    const raw = canvas.dataset.chartConfig
    if (!raw) return

    try {
      const config = JSON.parse(raw) as ChartConfiguration
      config.options = {
        ...config.options,
        responsive: true,
        maintainAspectRatio: false,
      }

      if (config.data?.datasets) {
        config.data.datasets.forEach((ds, index) => {
          const themeColor = strokeColors[index % strokeColors.length]
          const themeFill = fillColors[index % fillColors.length]
          const hasCustomColor = Boolean(ds.backgroundColor) || Boolean(ds.borderColor)

          if (!hasCustomColor) {
            ds.backgroundColor = config.type === 'pie' || config.type === 'doughnut' ? fillColors : themeFill
            ds.borderColor = themeColor
          } else {
            if (!ds.backgroundColor) ds.backgroundColor = ds.borderColor
            if (!ds.borderColor) ds.borderColor = ds.backgroundColor
          }
        })
      }

      const instance = new Chart(canvas, config)
      chartInstances.set(canvas, instance)
    } catch (err) {
      console.error('Failed to render chart:', err)
      const fallback = document.createElement('div')
      fallback.className = 'chart-error'
      fallback.textContent = '[SYS_CHART_ERROR] :: INVALID_CONFIG'
      canvas.replaceWith(fallback)
    }
  })
}

export function processFootnotes(markdown: string): string {
  const definitions = new Map<string, string>()

  const defRegex = /^\[\^([^\]]+)\]:[ \t]*(.+)$/gm
  const withoutDefs = markdown.replace(defRegex, (_match, rawLabel: string, text: string) => {
    definitions.set(rawLabel.trim(), text.trim())
    return ''
  })

  const order: string[] = []
  const scanRegex = /\[\^([^\]]+)\]/g
  let scanMatch: RegExpExecArray | null
  while ((scanMatch = scanRegex.exec(withoutDefs)) !== null) {
    const label = scanMatch[1]?.trim()
    if (label && definitions.has(label) && !order.includes(label)) {
      order.push(label)
    }
  }

  if (order.length === 0) {
    return withoutDefs.trim()
  }

  const replaceRegex = /\[\^([^\]]+)\]/g
  const withRefs = withoutDefs.replace(replaceRegex, (full, rawLabel: string) => {
    const label = rawLabel.trim()
    if (!definitions.has(label)) return full

    const number = order.indexOf(label) + 1
    return `<sup><a href="#fn-${label}" id="fnref-${label}">${number}</a></sup>`
  })

  const items = order
    .map((label) => {
      const text = definitions.get(label) ?? ''
      return `<li id="fn-${label}">${text} <a href="#fnref-${label}" class="footnote-backref">↩</a></li>`
    })
    .join('\n')

  const footnotesBlock = `\n\n<div class="footnotes">\n<ol>\n${items}\n</ol>\n</div>\n`

  return withRefs.trim() + footnotesBlock
}

export { marked }