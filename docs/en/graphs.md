# All Chart.js Types: Which One to Choose and When

Chart.js covers almost the entire set of basic data visualizations — from traditional bar charts to polar area charts and bubble maps. In this article, we will break down each type: what it shows, when it is appropriate (and when it is not), and how it looks in an actual post.

All examples below are functional `chart` blocks, identical to the ones you can use in any blog article.

---

## Bar — Bar Chart

The most versatile chart type. Ideal for comparing multiple categories against each other at a single point in time or dynamically across different time periods.

```chart
{
  "type": "bar",
  "data": {
    "labels": ["Jan", "Feb", "Mar", "Apr", "May"],
    "datasets": [
      { "label": "Sales", "data": [12, 19, 8, 15, 22], "backgroundColor": "rgba(54, 162, 235, 0.7)" },
      { "label": "Costs", "data": [8, 11, 6, 10, 14], "backgroundColor": "rgba(255, 99, 132, 0.7)" }
    ]
  }
}

```

When to use: comparing discrete categories — e.g., "Product A vs. Product B", "monthly revenue". When to avoid: if you have more than 10–15 categories, vertical bars become unreadable; a horizontal bar chart or a table is preferred instead.

**Horizontal variant** (`options.indexAxis: "y"`) is useful when category labels are long and do not fit comfortably beneath vertical bars:

```chart
{
  "type": "bar",
  "data": {
    "labels": ["North American Market", "European Market", "Asia-Pacific Region"],
    "datasets": [{ "label": "Revenue, $M", "data": [420, 310, 275], "backgroundColor": "rgba(75, 192, 192, 0.7)" }]
  },
  "options": { "indexAxis": "y" }
}

```

**Stacked bar** — when not only the total value matters, but also its composition:

```chart
{
  "type": "bar",
  "data": {
    "labels": ["Q1", "Q2", "Q3", "Q4"],
    "datasets": [
      { "label": "Product A", "data": [10, 15, 12, 18], "backgroundColor": "#36a2eb" },
      { "label": "Product B", "data": [8, 9, 14, 10], "backgroundColor": "#ff6384" }
    ]
  },
  "options": {
    "scales": { "x": { "stacked": true }, "y": { "stacked": true } }
  }
}

```

---

## Line — Line Chart

Displays value changes over time. This is the default choice for time series data, as human eyes perceive trends far better along continuous lines than across discrete bars.

```chart
{
  "type": "line",
  "data": {
    "labels": ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"],
    "datasets": [{
      "label": "Visitors",
      "data": [120, 190, 150, 220, 180, 90, 60],
      "borderColor": "rgb(75, 192, 192)",
      "backgroundColor": "rgba(75, 192, 192, 0.15)",
      "fill": true,
      "tension": 0.3
    }]
  }
}

```

`tension` smooths the line between data points (0 means straight line segments, values closer to 1 produce smooth curves). `fill: true` shades the area beneath the line — helpful when both momentum and cumulative "weight" are important, though this can become visually cluttered with multiple datasets (in which case `fill` should be disabled).

Multiple lines on a single chart allow comparing trends across different entities over the same timeframe:

```chart
{
  "type": "line",
  "data": {
    "labels": ["2020", "2021", "2022", "2023", "2024"],
    "datasets": [
      { "label": "Product A", "data": [30, 45, 42, 58, 71], "borderColor": "#36a2eb", "fill": false },
      { "label": "Product B", "data": [20, 25, 38, 40, 39], "borderColor": "#ff6384", "fill": false }
    ]
  }
}

```

---

## Pie / Doughnut — Pie and Doughnut Charts

These charts represent the proportional breakdown of categories within a whole. They function best with a small number of segments (up to 5–6) — beyond that, it becomes difficult to visually distinguish adjacent slice sizes.

```chart
{
  "type": "pie",
  "data": {
    "labels": ["Direct", "Search", "Social", "Referral"],
    "datasets": [{
      "data": [40, 30, 20, 10],
      "backgroundColor": ["#4285F4", "#FF7139", "#00C7B1", "#9E9E9E"]
    }]
  }
}

```

Doughnut follows the same logic with a hollow center. While often chosen purely for aesthetics, the open center also serves as a convenient spot for displaying total metrics or summary text (rendered via layout overlays, rather than directly in Chart.js):

```chart
{
  "type": "doughnut",
  "data": {
    "labels": ["Chrome", "Firefox", "Safari", "Other"],
    "datasets": [{
      "data": [55, 15, 20, 10],
      "backgroundColor": ["#4285F4", "#FF7139", "#00C7B1", "#9E9E9E"]
    }]
  }
}

```

A common mistake is attempting to illustrate time-series trends using pie or doughnut charts. This is outside their purpose: as soon as data represents a trend over time rather than a single point-in-time breakdown, a line or bar chart should be used instead.

---

## Radar — Radar (Spider) Chart

Compares an entity across multiple parameters simultaneously, where each parameter forms its own axis radiating from the center. A classic choice for comparing character stats in gaming, but equally useful for evaluating products, candidates, or configuration options against multiple criteria.

```chart
{
  "type": "radar",
  "data": {
    "labels": ["Speed", "Strength", "Agility", "Endurance", "Skill"],
    "datasets": [
      { "label": "Player A", "data": [65, 59, 90, 81, 56], "borderColor": "#ff6384", "backgroundColor": "rgba(255, 99, 132, 0.2)" },
      { "label": "Player B", "data": [28, 48, 40, 19, 96], "borderColor": "#36a2eb", "backgroundColor": "rgba(54, 162, 235, 0.2)" }
    ]
  }
}

```

Works best for 2–3 datasets and 4–8 axes. Comparing more than two entities causes overlapping shapes, degrading readability.

---

## PolarArea — Polar Area Chart

A hybrid between pie and radar charts: sectors have equal angular width, but varying radial length (unlike pie charts, where sector radius is uniform and angles vary). It displays both proportion and magnitude simultaneously.

```chart
{
  "type": "polarArea",
  "data": {
    "labels": ["Red", "Green", "Yellow", "Grey", "Blue"],
    "datasets": [{
      "data": [11, 16, 7, 3, 14],
      "backgroundColor": ["#ff6384", "#4bc0c0", "#ffce56", "#e7e9ed", "#36a2eb"]
    }]
  }
}

```

Used less frequently than other types — visually striking, though not always immediately intuitive for first-time readers. Highly effective in articles where a non-standard visual shape provides design emphasis rather than serving as the primary source for precise numerical analysis.

---

## Scatter — Scatter Plot

Shows relationships between two numerical variables: each point represents an individual observation with `x` and `y` coordinates. This is the primary tool for detecting correlations and outliers within data sets.

```chart
{
  "type": "scatter",
  "data": {
    "datasets": [{
      "label": "Observations",
      "data": [
        { "x": -10, "y": 0 }, { "x": 0, "y": 10 }, { "x": 10, "y": 5 },
        { "x": 20, "y": 15 }, { "x": 5, "y": -5 }, { "x": 15, "y": 8 }
      ],
      "backgroundColor": "#e91e63"
    }]
  },
  "options": {
    "scales": { "x": { "type": "linear", "position": "bottom" } }
  }
}

```

Important: The `x` axis must be explicitly configured as `"type": "linear"`. By default, Chart.js treats the x-axis as categorical, which causes scatter points to plot incorrectly. This step is unique to scatter and bubble charts and is unnecessary for line or bar charts.

---

## Bubble — Bubble Chart

Extends the scatter plot by adding a third dimension: point radius (`r`). Useful when presenting three variables at once — for example: x = price, y = rating, bubble radius = sales volume.

```chart
{
  "type": "bubble",
  "data": {
    "datasets": [{
      "label": "Products",
      "data": [
        { "x": 20, "y": 30, "r": 15 },
        { "x": 40, "y": 10, "r": 25 },
        { "x": 15, "y": 50, "r": 8 },
        { "x": 35, "y": 35, "r": 18 }
      ],
      "backgroundColor": "rgba(255, 159, 64, 0.6)"
    }]
  }
}

```

The primary risk with bubble charts is overcomplication. Reading three data dimensions simultaneously is significantly more demanding than reading two, so bubble charts should be clearly labeled and limited in point density.

---

## Summary Table: Which Type to Choose

| Type | Shows | Max. Categories/Datasets |
| --- | --- | --- |
| Bar | Category comparison | Up to 10–15 categories |
| Line | Trends over time | 2–4 lines without loss of clarity |
| Pie / Doughnut | Proportions of a whole | Up to 5–6 segments |
| Radar | Multi-criteria comparison | 2–3 datasets, 4–8 axes |
| PolarArea | Proportion + magnitude | Up to 6–8 sectors |
| Scatter | Relationship between 2 variables | No strict limit (requires sufficient data density) |
| Bubble | Relationship between 3 variables | 10–20 points (beyond this, readability degrades) |

---

## How to Format on the Page

Any chart above can be inserted with text alignment modifiers so that it wraps alongside prose instead of expanding across the full width of the post:

```chart:right
{
  "type": "line",
  "data": {
    "labels": ["Jan", "Feb", "Mar", "Apr"],
    "datasets": [{ "label": "Trend", "data": [10, 25, 18, 30], "borderColor": "#0b7a30", "fill": false }]
  }
}

```

This chart wraps text to its left. Notice how this entire paragraph sits alongside it without requiring manual wrapper elements, made possible simply by appending the `:right` modifier inside the code fence declaration.

A full guide to all alignment modifiers (`:left`, `:right`, `:float`, `:full`) and general `chart` block syntax is available in the Syntax Reference for Authors.