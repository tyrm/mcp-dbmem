# MCP-DBMEM

A poorly written Go-based memory management system for the Metoro Control Protocol (MCP), providing a persistent knowledge graph implementation.

## Overview

MCP-DBMEM is a memory management service that implements a knowledge graph using PostgreSQL as the backend storage. It provides a set of tools for managing entities, relations, and observations within the knowledge graph.

## Features

- Create and manage entities in the knowledge graph
- Create and manage relations between entities
- Add and manage observations for entities
- Search and query the knowledge graph
- Delete entities, relations, and observations
- Full graph reading capabilities

## Tools

The service provides the following tools:

- `create_entities`: Create multiple new entities in the knowledge graph
- `create_relations`: Create multiple new relations between entities (in active voice)
- `add_observations`: Add new observations to existing entities
- `delete_entities`: Delete multiple entities and their associated relations
- `delete_observations`: Delete specific observations from entities
- `delete_relations`: Delete multiple relations from the graph
- `read_graph`: Read the entire knowledge graph
- `search_nodes`: Search for nodes based on a query
- `open_nodes`: Open specific nodes by their names

## Installation

### Claude

#### Docker

```json
{
  "mcpServers": {
    "memory": {
      "command": "docker",
      "env": {
        "DB_ADDRESS": "localhost",
        "DB_USER": "mcp-dbmem",
        "DB_PASSWORD": "mcp-dbmem",
        "DB_DATABASE": "mcp-dbmem"
      },
      "args": [
        "run",
        "-i",
        "--rm",
        "tyrm/mcp-dbmem",
        "direct"
      ]
    }
  }
}
```

#### Binary

```json
{
  "mcpServers": {
    "memory": {
      "command": "/path/to/mcp-dbmem",
      "env": {
        "DB_ADDRESS": "localhost",
        "DB_USER": "mcp-dbmem",
        "DB_PASSWORD": "mcp-dbmem",
        "DB_DATABASE": "mcp-dbmem"
      },
      "args": [
        "direct"
      ]
    }
  }
}
```


## Development

### Prerequisites

- Go 1.x
- PostgreSQL
- Make (for build commands)

### Building

```bash
make build
```

### Running

```bash
./bin/mcp-dbmem
```

## Project Structure

