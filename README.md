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
prompt-security monitor
```

Or build from source:

```bash
git clone https://github.com/happytaoer/prompt-security.git
cd prompt-security
go mod tidy
go build -o prompt-security
./prompt-security monitor
```

---

## 🔥 Features

- **Real-time clipboard monitoring**
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

---

## 📝 Configuration

When you run the application for the first time, a configuration file will be automatically created at `~/.prompt-security/config.json`. You can manually edit this file; changes take effect after restarting the program.

**About the configuration file:**
- Path: `~/.prompt-security/config.json`
- Automatically created on first run; if corrupted or deleted, it will be restored to default
- Main configuration logic is in `internal/config/config.go`
- Main fields:

```json
{
  "detect_emails": true,
  "detect_phones": true,
  "detect_credit_cards": true,
  "detect_ssns": true,
  "detect_ipv4": true,
  "string_match_patterns": [
    {
      "name": "company_name",
      "pattern": "Acme Corporation",
      "enabled": true,
      "replacement": "[COMPANY NAME]" 
    },
    {
      "name": "internal_project",
      "pattern": "Project Phoenix",
      "enabled": true,
      "replacement": "[PROJECT NAME]"
    }
  ],
  "email_replacement": "security@example.com",
  "phone_replacement": "+1-555-123-4567",
  "credit_card_replacement": "XXXX-XXXX-XXXX-XXXX",
  "ssn_replacement": "XXX-XX-XXXX",
  "ipv4_replacement": "0.0.0.0",
  "monitoring_interval_ms": 500,
  "notify_on_filter": true
}
```

---

## 🧩 Pattern Types

### 1. Regular Expression Detection
- Email: `[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`
- Phone: `(\+\d{1,3}[\s-]?)?\(?\d{3}\)?[\s.-]?\d{3}[\s.-]?\d{4}`
- Credit Card: `\b(?:\d{4}[- ]?){3}\d{4}\b`
- SSN: `\b\d{3}-\d{2}-\d{4}\b`
- IPv4: `\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`

### 2. Exact String Match
Configure custom sensitive words, project names, company names, etc. in `string_match_patterns`:

```json
{
  "name": "pattern_name",
  "pattern": "exact text",
  "enabled": true,
  "replacement": "replacement"
}
```

---

## ⚙️ How It Works

1. Periodically checks clipboard content
2. Filters sensitive information using regular expressions and custom rules
3. Replaces matches with safe placeholders
4. Automatically writes the filtered content back to the clipboard

This way, any sensitive content you copy will be safely replaced before pasting into LLMs or web pages.

---

## 🖥️ CLI Usage

- Start monitoring:
  ```bash
  ./prompt-security monitor
  ```
- View current configuration:
  ```bash
  ./prompt-security config
  ```

---



## 🔒 Security & Privacy Statement

- All clipboard content is processed locally; no network connection, no uploads
- Open source and fully auditable—use with confidence



---

## License

MIT
