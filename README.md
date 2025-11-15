# Prompt Security

![Go Version](https://img.shields.io/badge/Go-1.21%2B-blue)
![License](https://img.shields.io/badge/License-MIT-green)
![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen)

> **Protect your clipboard before you paste!**
> 
> Prompt Security is a blazing-fast, zero-config Go tool that automatically detects and filters sensitive data from your clipboard—before it can leak to large language models (LLMs) or online tools.

---

## 🚀 Quick Start

```bash
go install github.com/happytaoer/prompt-security@latest
prompt-security
```

Or build from source:

```bash
git clone https://github.com/happytaoer/prompt-security.git
cd prompt-security
go mod tidy
go build -o prompt-security
./prompt-security
```

---

## 🔥 Features

- **Real-time clipboard monitoring**
- **🎨 Web GUI** for configuration and log monitoring
- **Automatic filtering** of:
  - Email addresses
  - Phone numbers
  - Credit card numbers
  - Social Security Numbers (SSN)
  - IPv4 addresses
  - Custom string patterns (exact match)
- **Configurable rules and replacements**
- **Easy CLI, zero config required to start**
- **Safe placeholder replacements**
- **Cross-platform** (Windows, macOS, Linux)

---

## 🛡️ Why Prompt Security?

- Prevents accidental data leaks to ChatGPT, Copilot, Gemini, etc.
- Enforces company security policies
- Protects privacy and sensitive business information
- Lets you use AI tools with peace of mind

---

## 💡 Typical Use Cases

- Developers using LLMs for code review or documentation
- Security teams enforcing clipboard data policies
- AI power users handling confidential information
- Anyone worried about copy-paste data leaks


## 🔒 Security & Privacy Statement

- All clipboard content is processed locally; no network connection, no uploads
- Open source and fully auditable—use with confidence



---

## License

MIT
