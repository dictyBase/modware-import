# modware-import
[![License](https://img.shields.io/badge/License-BSD%202--Clause-blue.svg)](LICENSE)   
![Build](https://github.com/dictyBase/modware-import/workflows/Build/badge.svg)
![Last commit](https://badgen.net/github/last-commit/dictyBase/modware-import/develop)   
[![Funding](https://badgen.net/badge/Funding/Rex%20L%20Chisholm,dictyBase,DCR/yellow?list=|)](https://reporter.nih.gov/project-details/10024726)

Cli application for importing dictybase data.

## Components

### Data Importer
Cli application for importing dictybase data.

### Mock gRPC Server for Feature Annotation Service
A command-line mock gRPC server that implements the `FeatureAnnotationService` for integration testing of gRPC clients.

#### Features

- ✅ **All 8 gRPC Methods Implemented**:
  - `CreateFeatureAnnotation` - Create new feature annotations
  - `GetFeatureAnnotation` - Retrieve annotations by ID
  - `UpdateFeatureAnnotation` - Update existing annotations
  - `DeleteFeatureAnnotation` - Delete annotations (soft delete or purge)
  - `AddTag` - Add tags to annotations
  - `UpdateTag` - Update existing tags
  - `RemoveTag` - Remove tags from annotations
  - `ListFeatureAnnotationsByPubmedId` - Query by PubMed ID
  - `ListFeatureAnnotationsByDOI` - Query by DOI

- ✅ **Thread-Safe In-Memory Storage** with indexes for efficient lookups
- ✅ **Realistic Mock Data** - Pre-loaded with 5 sample feature annotations
- ✅ **Comprehensive Validation** - Email format, DOI patterns, required fields
- ✅ **CLI Interface** with configurable port, logging, and TLS
- ✅ **gRPC Reflection** enabled for debugging with tools like `grpcurl`
- ✅ **Graceful Shutdown** with signal handling

#### Quick Start

##### Build and Run

```bash
# Build the server
cd cmd/mock-grpc-server
go build -o mock-grpc-server .

# Run with default settings (port 9000)
./mock-grpc-server

# Run with custom settings
./mock-grpc-server --port 9001 --log-level debug
```

##### CLI Options

```bash
NAME:
   mock-grpc-server - Mock gRPC server for feature annotation service integration testing

USAGE:
   mock-grpc-server [global options]

GLOBAL OPTIONS:
   --port value, -p value       Server port (default: 9000) [$GRPC_PORT]
   --log-level value, -l value  Log level (debug, info, warn, error) (default: "info") [$LOG_LEVEL]
   --tls, -t                    Enable TLS (default: false) [$TLS_ENABLED]
   --help, -h                   show help
```

#### Testing with grpcurl

The server has gRPC reflection enabled, so you can use `grpcurl` to explore and test the API:

```bash
# List available services
grpcurl -plaintext localhost:9000 list

# List methods for FeatureAnnotationService
grpcurl -plaintext localhost:9000 list dictybase.feature_annotation.FeatureAnnotationService

# Get a feature annotation (using pre-loaded mock data)
grpcurl -plaintext -d '{"id": "DDB_G0267398"}' localhost:9000 dictybase.feature_annotation.FeatureAnnotationService/GetFeatureAnnotation

# Create a new feature annotation
grpcurl -plaintext -d '{
  "type": "gene",
  "id": "TEST_001",
  "attributes": {
    "name": "testGene",
    "publications": ["10.1000/test.2023.001"],
    "pubmed": ["12345678"]
  },
  "created_by": "test@dictybase.org"
}' localhost:9000 dictybase.feature_annotation.FeatureAnnotationService/CreateFeatureAnnotation

# List annotations by PubMed ID
grpcurl -plaintext -d '{"id": "12345678"}' localhost:9000 dictybase.feature_annotation.FeatureAnnotationService/ListFeatureAnnotationsByPubmedId
```

#### Pre-loaded Mock Data

The server starts with 5 realistic feature annotations:

1. **actA** (DDB_G0267398) - Actin gene with cytoskeleton function
2. **myoB** (DDB_G0275199) - Myosin II heavy chain B with motor activity
3. **pakA** (DDB_G0282525) - P21-activated kinase A with kinase activity
4. **rasG** (DDB_G0283471) - Ras protein G with GTPase activity
5. **discoidin1** (DDB_G0291234) - Discoidin with carbohydrate binding

Each annotation includes:
- Realistic gene names and synonyms
- Valid DOI and PubMed ID references
- Tag properties with functions and cellular locations
- Database cross-references (DbLinks)
- Proper timestamps and user information

#### Architecture

```
Mock gRPC Server
├── CLI Interface (urfave/cli v2)
├── FeatureAnnotationService (8 methods)
├── In-Memory Storage Layer
│   ├── Thread-Safe Maps
│   └── Indexes (ID, Name, PubmedID, DOI)
└── Validation & Mock Data Generation
```

## Documentation 
* [Importer](docs/import.md)
* [k8s](docs/k8s.md)


# Misc Badges
![Issues](https://badgen.net/github/issues/dictyBase/modware-import)
![Open Issues](https://badgen.net/github/open-issues/dictyBase/modware-import)
![Closed Issues](https://badgen.net/github/closed-issues/dictyBase/modware-import)
![Total PRS](https://badgen.net/github/prs/dictyBase/modware-import)
![Open PRS](https://badgen.net/github/open-prs/dictyBase/modware-import)
![Closed PRS](https://badgen.net/github/closed-prs/dictyBase/modware-import)
![Commits](https://badgen.net/github/commits/dictyBase/modware-import/develop)
![Branches](https://badgen.net/github/branches/dictyBase/modware-import)
![Tags](https://badgen.net/github/tags/dictyBase/modware-import)   
![GitHub repo size](https://img.shields.io/github/repo-size/dictyBase/modware-import?style=plastic)
![GitHub code size in bytes](https://img.shields.io/github/languages/code-size/dictyBase/modware-import?style=plastic)
[![Lines of Code](https://badgen.net/codeclimate/loc/dictyBase/modware-import)](https://codeclimate.com/github/dictyBase/modware-import/code)   

