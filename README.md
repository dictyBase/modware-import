# modware-import

[![License](https://img.shields.io/badge/License-BSD%202--Clause-blue.svg)](LICENSE)   
![Build](https://github.com/dictyBase/modware-import/workflows/Build/badge.svg)
![Last commit](https://badgen.net/github/last-commit/dictyBase/modware-import/develop)   
[![Funding](https://badgen.net/badge/Funding/Rex%20L%20Chisholm,dictyBase,DCR/yellow?list=|)](https://reporter.nih.gov/project-details/10024726)

CLI application suite for importing and managing dictyBase data, including feature annotations, gene data, and development tools.

## Table of Contents

- [Overview](#overview)
- [Components](#components)
  - [Feature Annotation CLI](#feature-annotation-cli)
  - [Import CLI](#import-cli)
  - [Kubernetes Deployment CLI](#kubernetes-deployment-cli)
  - [Mock gRPC Server](#mock-grpc-server)
    - [Features](#features)
    - [Quick Start](#quick-start)
    - [Testing](#testing)
    - [Architecture](#architecture)
- [Project Status](#project-status)

## Overview

This repository contains a comprehensive suite of CLI tools for managing dictyBase data operations:

- **Feature Annotation CLI**: Modern gRPC-based tools for loading, updating, and managing feature annotations from various data sources
- **Import CLI**: Legacy data migration tools with S3 integration for importing dictyBase data during migration
- **Kubernetes Deployment CLI**: Tools for deploying and running import commands in Kubernetes clusters
- **Mock gRPC Server**: Development and testing tool that provides a complete implementation of the FeatureAnnotationService

## Components

### Feature Annotation CLI

A production CLI application for importing and managing dictyBase feature annotation data. The CLI provides six main commands for different data operations:

- **`load-feature-annotation`** - Load feature annotations from ArangoDB to gRPC service
- **`load-csv-to-arangodb`** - Update ArangoDB collection from CSV file with batch processing
- **`gene-updater`** - Update gene annotations by stripping HTML and using gRPC
- **`gene-product-updater`** - Update gene products from legacy database
- **`load-gene-product-from-csv`** - Load gene products from CSV files
- **`load-synonyms`** - Load synonyms from ArangoDB to gRPC service

For detailed usage instructions, command options, and examples, see the [Feature Annotation CLI Reference](docs/featureannotation.md).

### Import CLI

A legacy command-line application for importing dictyBase data during migration. Supports importing various types of data including:

- **ArangoDB Management** - Database schema and data management operations
- **Data File Processing** - Handling and transformation of data files  
- **Ontology Management** - Import and management of ontological data
- **Stock Center Data** - Loading of stock center information
- **UniProt ID Mapping** - Protein identifier mapping operations

The tool supports S3 integration for file storage and can process CSV formatted data from various sources.

For complete command reference and usage examples, see the [Importer Documentation](docs/import.md).

### Kubernetes Deployment CLI

A command-line tool for deploying and running import commands in Kubernetes clusters. Provides:

- **Cluster Deployment** - Deploy import operations as Kubernetes jobs
- **S3 Integration** - Configure S3 access for input files and log storage
- **Namespace Management** - Target specific Kubernetes namespaces
- **Log Management** - Centralized logging with S3 backup

For deployment guides and configuration details, see the [Kubernetes Deployment](docs/k8s.md).

### Mock gRPC Server

A development and testing tool that implements a complete `FeatureAnnotationService` for integration testing of gRPC clients.

#### Features

- **Complete gRPC Implementation**: All 8 FeatureAnnotationService methods (Create, Get, Update, Delete, AddTag, UpdateTag, RemoveTag, List operations)
- **Thread-Safe Storage**: In-memory storage with indexes for efficient lookups
- **Mock Data**: Pre-loaded with 5 realistic feature annotations (actA, myoB, pakA, rasG, discoidin1)
- **gRPC Reflection**: Enabled for easy testing with `grpcurl`
- **Configurable**: CLI options for port, logging, and TLS

#### Quick Start

```bash
# Build and run
cd cmd/mock-grpc-server
go build -o mock-grpc-server .
./mock-grpc-server --port 9000 --log-level info

# Test with grpcurl
grpcurl -plaintext localhost:9000 list
grpcurl -plaintext -d '{"id": "DDB_G0267398"}' \
  localhost:9000 dictybase.feature_annotation.FeatureAnnotationService/GetFeatureAnnotation
```

#### Testing

The server supports comprehensive testing scenarios:

- **Service Discovery**: `grpcurl -plaintext localhost:9000 list`
- **Data Retrieval**: Pre-loaded annotations for immediate testing
- **CRUD Operations**: Full create, read, update, delete functionality
- **Validation Testing**: Email formats, DOI patterns, required fields

#### Architecture

Thread-safe in-memory storage with indexed lookups, comprehensive validation, and graceful shutdown handling.


## Project Status

![Open Issues](https://badgen.net/github/open-issues/dictyBase/modware-import)
![Open PRS](https://badgen.net/github/open-prs/dictyBase/modware-import)

