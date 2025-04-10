# Next.js Routing Helper CLI 🛠️

A simple but powerful CLI tool to speed up your workflow when working with Next.js projects.  
Generate pages, components, and routing structures automatically with a single command.

## 🚀 Features

- 📄 **Page Generator**: Instantly scaffold new pages with or without the `'use client'` directive.
- 🧩 **Component Style Options**: Choose between `function` or `const` component styles.
- 🌿 **App / Pages Routers Support**: Both routers in Next.js are supported.
- ⚙️ **Configurable**: Adjust defaults via a config file to match your project’s standards.
- 🧼 **Visualize Structure**: Visualize your project's directory structure.

## 📦 Installation

```zsh
go install github.com/bllakcn/nextjs-routing-helper-cli@latest
```

## 📄 Usage

1. Initialize the config file

```zsh
nextjs-routing-helper init
```

2. Add a page

```zsh
nextjs-routing-helper add <dir/pageName> [flags]
```

## 🛤️ Roadmap

- [ ] Add tests and better error handling
- [ ] Custom templating support
- [ ] Add pages interactively
- [ ] Generate API routes
- [ ] Git hook integration for consistency checks

## 🤝 Contributing

Contributions are welcome! Feel free to open issues, request features, or submit PRs.

## 📄 License

MIT
