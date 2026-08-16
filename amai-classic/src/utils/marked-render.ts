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

function applyChartDefaults(theme: 'light' | 'dark') {
  Chart.defaults.color = theme === 'dark' ? '#c9ccd1' : '#57606a'
  Chart.defaults.borderColor = theme === 'dark' ? '#3a3f45' : '#d0d7de'
}

export function initChartTheme(): void {
  const { theme } = useTheme()

  watch(
    theme,
    (val) => {
      applyChartDefaults(val)
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

const markExtension = {
  name: 'mark',
  level: 'inline' as const,
  start(src: string) {
    return src.match(/==/)?.index
  },
  tokenizer(this: { lexer: { inlineTokens: (s: string) => Token[] } }, src: string) {
    const rule = /^==([^=]+)==/
    const match = rule.exec(src)
    const text = match?.[1]
    if (match && text) {
      return {
        type: 'mark',
        raw: match[0],
        text,
        tokens: this.lexer.inlineTokens(text),
      }
    }
    return undefined
  },
  renderer(this: { parser: { parseInline: (tokens: Token[]) => string } }, token: {
    tokens: Token[]
  }) {
    return `<mark>${this.parser.parseInline(token.tokens)}</mark>`
  },
}

const markRedExtension = {
  name: 'markRed',
  level: 'inline' as const,
  start(src: string) {
    return src.match(/--/)?.index
  },
  tokenizer(this: { lexer: { inlineTokens: (s: string) => Token[] } }, src: string) {
    const rule = /^--([^-]+)--/
    const match = rule.exec(src)
    const text = match?.[1]
    if (match && text) {
      return {
        type: 'markRed',
        raw: match[0],
        text,
        tokens: this.lexer.inlineTokens(text),
      }
    }
    return undefined
  },
  renderer(this: { parser: { parseInline: (tokens: Token[]) => string } }, token: {
    tokens: Token[]
  }) {
    return `<mark class="mark-red">${this.parser.parseInline(token.tokens)}</mark>`
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
    markRedExtension,
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
  const canvases = root.querySelectorAll<HTMLCanvasElement>('canvas[data-chart-config]')

  canvases.forEach((canvas) => {
    if (chartInstances.has(canvas)) return

    const raw = canvas.dataset.chartConfig
    if (!raw) return

    try {
      const config = JSON.parse(raw) as ChartConfiguration
      const instance = new Chart(canvas, config)
      chartInstances.set(canvas, instance)
    } catch (err) {
      console.error('Failed to render chart:', err)
      const fallback = document.createElement('div')
      fallback.className = 'chart-error'
      fallback.textContent = 'Invalid chart config'
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