---
title: "How To Use Shizuka"
date: "1970-02-02"
template: "post.tmpl"
---

# Welcome to Shizuka!

Shizuka is a small, quiet, and thoughtful tool for building your own website.

Here’s a quick overview of how Shizuka works, using a demo project to guide you through its features and structure.

---

## Project Structure

A typical Shizuka site project looks like this:

```
site
├── content
│   ├── index.md
│   └── posts
│       ├── 1.md
│       ├── 2.md
│       └── 3.md
├── static
│   └── styles.css
└── templates
    ├── index.tmpl
    └── post.tmpl
```

### Key Directories

- **`content/`**  
  This folder contains all your markdown pages.
    - Files named `index.md` are rendered as `index.html` directly in their folder.
    - Other markdown files (e.g., `1.md`) are placed in a new folder (e.g., `posts/1/index.html`).

- **`templates/`**  
  Templates define how your markdown content is transformed into HTML. Each markdown file specifies its template using the `template` field in its frontmatter.

- **`static/`**  
  All static files (CSS, JS, images, etc.) are copied to the destination folder, preserving the folder structure.

---

## Configuration

The configuration for Shizuka lives in a tiny `shizuka_conf.json` file at the project root:

```json
{
  "src": "site",
  "dst": "dist",
  "port": "8080"
}
```

- **`src`**: The source directory containing your site files.
- **`dst`**: The destination directory where the site will be built.
- **`port`**: The port used by the development server (`shizuka dev`).

---

## Markdown with Frontmatter

Shizuka processes markdown files with frontmatter to create beautiful, dynamic HTML pages. Here's an example of what a markdown file might look like:

```markdown
---
title: "My First Blog Post"
description: "An introduction to my blog and what I'll be writing about."
author: "Shizuka Author"
date: "2024-12-07"
tags:
  - blogging
  - static site generators
  - shizuka

meta_title: "My First Blog Post - Shizuka Blog"
meta_description: "A brief introduction to my blog, powered by Shizuka."
meta_keywords: "blogging, static site generator, Shizuka"

data:
  custom_field_1: "Value for a custom field"
  custom_field_2: 42
  custom_field_3:
    - item1
    - item2

template: "post.tmpl"
---

This is the body of my first blog post. Here, I'll share my thoughts, ideas, and updates about using Shizuka to build a simple and elegant website.
```

### Key Frontmatter Fields

- **`title`**: The title of the page or post.
- **`description`**: A short summary of the content.
- **`author`**: The author's name.
- **`date`**: The publication date in `YYYY-MM-DD` format.
- **`tags`**: A list of tags associated with the post.
- **`meta_*`**: Metadata for SEO.
- **`data`**: Custom fields that can be used in templates.
- **`template`**: The template file to render this page (e.g., `post.tmpl`).

`template` is the only required field to include in the frontmatter. All other fields are optional (assuming that the template you are using can cope with that).

