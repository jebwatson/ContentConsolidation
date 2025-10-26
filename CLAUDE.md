# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based content consolidation application that uses MongoDB to store and manage website content. The application connects to MongoDB and provides basic CRUD operations for website entries.

## Development Commands

### Running the Application
```bash
go run main.go
```

### Building the Application
```bash
go build -o content-consolidation main.go
```

### Installing Dependencies
```bash
go mod download
go mod tidy
```

## Environment Setup

The application requires a MongoDB connection string set via the `MONGODB_URI` environment variable:
```bash
export MONGODB_URI="mongodb://localhost:27017"
go run main.go
```

## Architecture

### Database Schema
- **Database**: `content_consolidation_db`
- **Collection**: `websites`
- **Document Structure**: `Website` struct with `Uri` (string) and `Description` (string) fields

### Code Structure
The codebase is a single-file application (`main.go`) with the following key components:

- **MongoDB Connection Management**: `connectToMongo()` and `closeMongoDB()` handle database lifecycle
- **CRUD Operations**:
  - `CreateContentEntry()` - Inserts website records (currently unused in main flow)
  - `ReadContentEntry()` - Queries and displays all website records from the collection
- **Global State**: `mongoClient` maintains the MongoDB connection throughout the application lifecycle

### Known Issues
- `connectToMongo()` at line 76 incorrectly passes `mongodb_env` (constant string) instead of the `uri` parameter to `ApplyURI()`
- This will cause connection failures as it tries to connect to literal string "MONGODB_URI" rather than the actual connection URI
