# Snipet AI Runtime Architecture

## Overview

The Snipet runtime is responsible for orchestrating the complete execution of a user request. Instead of a fixed chatbot flow, every message is executed through a planning and execution pipeline.

The architecture is divided into two main concepts:

- **Knowledge** → Static or semi-static information used by the AI.
- **State** → Dynamic information that represents the current conversation.

---

# Core Concepts

## Client

Represents an external application integrated with Snipet.

Examples:

- Website
- Mobile App
- Discord Bot
- WhatsApp
- Internal API

A client owns multiple sessions.

---

## Session

Represents a conversation.

A session contains:

- Current conversation state
- Associated bot
- Participants (users)
- Messages

The session itself does **not** contain business logic.

---

## Bot

Defines the AI behavior.

A bot contains:

- Prompt / Persona
- Available models
- Tools
- Knowledge sources
- Planner configuration

Bots do not own conversation history.

---

# Knowledge

Knowledge represents information that can be queried by the AI but is not modified during a conversation.

Examples:

- RAG
- Vector databases
- SQL databases
- Knowledge Graphs
- KV stores
- Documentation
- FAQ
- Company data

A bot may have multiple knowledge sources.

```text
Bot
 ├── Company Docs
 ├── FAQ
 ├── PostgreSQL
 ├── Milvus
 └── Graph Database
```

The planner decides which knowledge sources should be used for each request.

---

# State

State represents everything that changes during a conversation.

Examples:

- Conversation history
- Conversation summary
- Extracted facts
- Variables
- Current goals
- Temporary context
- User preferences

Each session owns exactly one state.

---

# Runtime

The runtime executes the complete request lifecycle.

```text
User
    │
    ▼
Session
    │
    ▼
Runtime
```

The runtime is composed of several components.

---

# Planner

The planner analyzes the request and creates an execution plan.

Responsibilities:

- Select model
- Select knowledge sources
- Select tools
- Decide whether RAG is needed
- Decide how much conversation state is required
- Configure generation parameters
- Optimize cost and latency

The planner never executes actions.

It only creates a plan.

Example:

```json
{
  "model": "gpt-5-mini",
  "knowledge": [
    "docs",
    "faq"
  ],
  "tools": [
    "search_orders"
  ],
  "state": {
    "summary": true,
    "history": true
  }
}
```

---

# Context Builder

Builds the final prompt using the execution plan.

Possible inputs:

- System prompt
- Conversation summary
- Conversation history
- Retrieved knowledge
- Tool outputs
- User message

The output is the complete context sent to the model.

---

# Executor

Executes the execution plan produced by the planner.

Responsibilities:

- Retrieve knowledge
- Execute tools
- Call models
- Collect intermediate results

The executor contains no planning logic.

---

# State Manager

Responsible for updating the session state after the response is generated.

Responsibilities:

- Store messages
- Update summaries
- Extract facts
- Update variables
- Persist long-term conversation information

The State Manager decides **what should or should not be remembered**.

Examples:

Store:

> "My name is John."

Ignore:

> "Thanks!"

---

# Complete Flow

```text
User
    │
    ▼
Session
    │
    ▼
Runtime
    │
    ▼
Planner
    │
    ▼
Context Builder
    │
    ▼
Executor
    │
    ├── Knowledge
    ├── Tools
    └── Model
    │
    ▼
Response
    │
    ▼
State Manager
    │
    ▼
Persist State
    │
    ▼
Return Response
```

---

# Future Vision

The runtime should be implemented as a pipeline instead of hardcoded logic.

This allows every step to be replaced or customized in the future, enabling:

- RAG pipelines
- AI Agents
- Multi-agent systems
- Tool orchestration
- MCP integrations
- Custom workflows
- User-defined execution flows

The default chatbot becomes simply the default runtime pipeline.