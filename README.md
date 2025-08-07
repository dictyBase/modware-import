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

### Feature Annotation CLI
A command-line application for managing feature annotations.

#### `load-feature-annotation`
This subcommand loads feature annotations from an ArangoDB instance into the feature annotation service via gRPC.

**Usage:**
```bash
featureannotation load-feature-annotation [command options]
```

**Options:**
| Flag | Description | Environment Variable | Default | Required |
|---|---|---|---|---|
| `--arangodb-user` | ArangoDB user name | `ARANGODB_USER` | | Yes |
| `--arangodb-pass` | ArangoDB password | `ARANGODB_PASS` | | Yes |
| `--arangodb-database` | ArangoDB database name | `ARANGODB_DATABASE` | | Yes |
| `--arangodb-host` | ArangoDB host | `ARANGODB_SERVICE_HOST` | `arangodb` | No |
| `--arangodb-port` | ArangoDB port | `ARANGODB_SERVICE_PORT` | `8529` | No |
| `--is-secure` | Use TLS for ArangoDB connection | `ARANGODB_IS_SECURE` | `false` | No |
| `--feature-annotation-grpc-host` | Feature annotation gRPC host | `ANNO_FEAT_API_SERVICE_HOST` | `anno-feat-api` | No |
| `--feature-annotation-grpc-port` | Feature annotation gRPC port | `ANNO_FEAT_API_SERVICE_PORT` | `9250` | No |

---

#### `load-csv-to-arangodb`
This subcommand updates an ArangoDB collection from a CSV file.

**Usage:**
```bash
featureannotation load-csv-to-arangodb [command options]
```

**Options:**
| Flag | Description | Environment Variable | Default | Required |
|---|---|---|---|---|
| `--csv-file` | Path to CSV file to load | | | Yes |
| `--collection` | ArangoDB collection name | | `featureprop` | No |
| `--delimiter` | CSV delimiter character | | `,` | No |
| `--batch-size` | Documents to update per batch | | `40` | No |
| `--workers` | Concurrent workers for batching | | `4` | No |
| `--arangodb-user` | ArangoDB user name | `ARANGODB_USER` | | Yes |
| `--arangodb-pass` | ArangoDB password | `ARANGODB_PASS` | | Yes |
| `--arangodb-database` | ArangoDB database name | `ARANGODB_DATABASE` | | Yes |
| `--arangodb-host` | ArangoDB host | `ARANGODB_SERVICE_HOST` | `arangodb` | No |
| `--arangodb-port` | ArangoDB port | `ARANGODB_SERVICE_PORT` | `8529` | No |
| `--is-secure` | Use TLS for ArangoDB connection | `ARANGODB_IS_SECURE` | `false` | No |

---

#### `gene-updater`
This subcommand updates gene annotations by stripping HTML from properties and using a gRPC API.

**Usage:**
```bash
featureannotation gene-updater [command options]
```

**Options:**
| Flag | Description | Environment Variable | Default | Required |
|---|---|---|---|---|
| `--aql-query` | AQL query to fetch gene data | `AQL_QUERY` | (See source) | No |
| `--processing-workers`| HTML processing workers | `PROCESSING_WORKERS` | `4` | No |
| `--grpc-workers` | gRPC update workers | `GRPC_WORKERS` | `8` | No |
| `--arangodb-user` | ArangoDB user name | `ARANGODB_USER` | | Yes |
| `--arangodb-pass` | ArangoDB password | `ARANGODB_PASS` | | Yes |
| `--arangodb-database` | ArangoDB database name | `ARANGODB_DATABASE` | | Yes |
| `--arangodb-host` | ArangoDB host | `ARANGODB_SERVICE_HOST` | `arangodb` | No |
| `--arangodb-port` | ArangoDB port | `ARANGODB_SERVICE_PORT` | `8529` | No |
| `--is-secure` | Use TLS for ArangoDB connection | `ARANGODB_IS_SECURE` | `false` | No |
| `--feature-annotation-grpc-host` | Feature annotation gRPC host | `ANNO_FEAT_API_SERVICE_HOST` | `anno-feat-api` | No |
| `--feature-annotation-grpc-port` | Feature annotation gRPC port | `ANNO_FEAT_API_SERVICE_PORT` | `9250` | No |

---

#### `gene-product-updater`
This subcommand updates gene products from a legacy database to the feature annotation service.

**Usage:**
```bash
featureannotation gene-product-updater [command options]
```

**Options:**
| Flag | Description | Environment Variable | Default | Required |
|---|---|---|---|---|
| `--legacy-database` | Legacy database name | `LEGACY_DATABASE` | `cgm_ddb` | No |
| `--legacy-workers` | Legacy DB query workers | `LEGACY_WORKERS` | `4` | No |
| `--grpc-workers` | gRPC update workers | `GRPC_WORKERS` | `8` | No |
| `--arangodb-user` | ArangoDB user name | `ARANGODB_USER` | | Yes |
| `--arangodb-pass` | ArangoDB password | `ARANGODB_PASS` | | Yes |
| `--arangodb-database` | ArangoDB database name | `ARANGODB_DATABASE` | | Yes |
| `--arangodb-host` | ArangoDB host | `ARANGODB_SERVICE_HOST` | `arangodb` | No |
| `--arangodb-port` | ArangoDB port | `ARANGODB_SERVICE_PORT` | `8529` | No |
| `--is-secure` | Use TLS for ArangoDB connection | `ARANGODB_IS_SECURE` | `false` | No |
| `--feature-annotation-grpc-host` | Feature annotation gRPC host | `ANNO_FEAT_API_SERVICE_HOST` | `anno-feat-api` | No |
| `--feature-annotation-grpc-port` | Feature annotation gRPC port | `ANNO_FEAT_API_SERVICE_PORT` | `9250` | No |

## Documentation 
* [Importer](docs/import.md)
* [k8s](docs/k8s.md)
* [Feature Annotation CLI](docs/featureannotation.md)


# Misc Badges
![Open Issues](https://badgen.net/github/open-issues/dictyBase/modware-import)
![Open PRS](https://badgen.net/github/open-prs/dictyBase/modware-import)

