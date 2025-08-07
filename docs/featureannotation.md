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
