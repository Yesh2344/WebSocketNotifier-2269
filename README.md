# WebSocketNotifier

[![Go Report Card](https://goreportcard.com/badge/github.com/example/websocket-notifier)](https://goreportcard.com/report/github.com/example/websocket-notifier)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

WebSocketNotifier is a production‑ready Go service that provides real‑time notifications over WebSocket. It demonstrates clean architecture, robust error handling, structured logging, and comprehensive testing.

## Features

- **Scalable client manager** – thread‑safe registration, deregistration, and broadcasting.
- **Configurable** via environment variables (`.env`).
- **Structured JSON logging** with Logrus.
- **Graceful shutdown** handling.
- **Full test coverage** for core server logic.

## Prerequisites

- Go 1.20 or newer
- `git` (to clone the repository)

## Installation