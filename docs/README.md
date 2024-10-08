# Bernstein
[![ci](https://github.com/guillego/bernstein/actions/workflows/ci.yml/badge.svg)](https://github.com/guillego/bernstein/actions/workflows/ci.yml)

This project is a minimal container orchestrator built with Go. It is designed to manage container deployment, execution, and monitoring across a small cluster of nodes.

## Features

- Node registration and management
- Container deployment
- Basic scheduling (round-robin)
- Container lifecycle management
- Monitoring and logging
- Fault tolerance

## Getting Started

### Prerequisites

- Go 1.18+
- Docker

### Installation

1. Clone the repository
2. Run `make run`

### Configuration

- Environment variables are stored in the `.env` file.

