# Shizuka

Shizuka is a minimalist static site generator.

Shizuka is designed to be very general; whether you’re setting up a personal blog, a project showcase, or just a space to collect your thoughts, Shizuka is designed to make the process of building and maintaining your site simple.

## Features

- **Markdown/HTML first**: Write your content in markdown, design page templates in HTML, and Shizuka takes care of the rest.
- **Minimalist Design**: No unnecessary bloat – just the essentials for building a fast, lightweight site.
- **Out of Your Way**: The main focus of Shizuka's design was to stay out of the way when you are designing a site.
- **Customizable**: Shizuka includes a simple default template to get you started, but supports easy and powerful customization for your own unique design.
- **CLI Tooling**: Manage your site with ease using Shizuka’s commands:
   - `init`: Set up a new project with a single command.
   - `build`: Compile your markdown into a fully functional static site.
   - `dev`: Start a local development server with live reloading.
   - `serve`: Preview a built site before deploying.

---

## Installation

To install Shizuka, simply use `go install`:

```bash
go install github.com/e74000/shizuka@latest
```

Make sure `$GOPATH/bin` is in your `PATH` so you can run the `shizuka` command globally.

---

## Getting Started

### 1. Initialize a Project

Run the following command to scaffold a new project in the current directory:

```bash
shizuka init
```

This creates a project directory with the following structure:

```
.
├── shizuka_conf.json
└── site
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

- **`content/`**: Your markdown files go here.
- **`static/`**: Place CSS, images, and other static assets here.
- **`templates/`**: Define how your content is rendered into HTML.

### 2. Start Developing

Start a live development server with:

```bash
shizuka dev
```

This will watch for changes, rebuild your site automatically, and serve it locally. By default, it runs on port `8080` (you can change this with the `--port` flag).

### 3. Build for Production

When you’re ready to deploy your site, use:

```bash
shizuka build
```

This compiles your site into the `dist/` folder (or as specified in `shizuka_conf.json`), ready to be uploaded to your hosting provider.

---

## Contributing

Contributions are welcome! There's many ways we could improve this, so please feel free to contribute.

---

## License

Shizuka is open-source and available under the [MIT License](LICENSE).
