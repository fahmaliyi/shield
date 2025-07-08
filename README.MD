# Shield

Minimal, secure session management package for Go.

Features:

- 🔐 Single-session per user (no concurrent logins)
- 💤 Sliding expiration (sessions expire after inactivity)
- 💾 Pluggable storage (default: in-memory)
- 🌐 Framework-agnostic, with optional HTTP middleware
- 🧪 Fully tested, production-ready

## Install

```bash
go get github.com/fahmcode/shield
```
