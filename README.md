# Next.js Routing Helper CLI 🛠️

A simple CLI tool to speed up your workflow when working with Next.js projects.  
Generate pages, components, and routing structures automatically with a single command to speed up the scaffolding.

## 🚀 Features

- 📄 **Page Generator**: Instantly scaffold new pages with or without the `'use client'` directive.
- 🧩 **Component Style Options**: Choose between `function` or `const` component styles.
- 🌿 **App / Pages Routers Support**: Both routers in Next.js are supported.
- ⚙️ **Configurable**: Adjust defaults via a config file to match your project’s standards.
- 🧼 **Visualize Structure**: Visualize your project's directory structure.

## 📦 Installation

```zsh
$ go install github.com/bllakcn/nextjs-routing-helper-cli@latest
```

## 📄 Usage

1. Initialize the config file in your root directory of your Nextjs project.

```zsh
$ nextjs-routing-helper init
```

This will create a `.nextjs_routing_helper.json`, where the cli will hold the necessary preferences.

2. Add a Page

```zsh
$ nextjs-routing-helper add [route/subroute] [flags]
```

This command creates a new page under the specified directory.

- In **App Router** projects, it generates a `page.tsx` file under `app/route/subroute/`.
- In **Pages Router** projects, it generates a `index.tsx` file under `pages/route/subroute/`.

You can optionally pass flags like `--use-client` to include the `'use client';` directive in App Router components:

```zsh
$ nextjs-routing-helper add dashboard/home --use-client
```

## 🛤️ Roadmap

- [ ] Add support for dynamic routes
- [ ] Add pages interactively
- [ ] Custom templating support
- [ ] Generate API routes
- [ ] Git hook integration for consistency checks

## 🤝 Contributing

Contributions are welcome! Feel free to open issues, request features, or submit PRs.

## 📄 License

MIT
