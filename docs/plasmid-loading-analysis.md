# Plasmid Data Loading Analysis Report

Based on analysis of `internal/cli/stockcenter/plasmid.go` and related files.

## Input Files

The plasmid loader requires **4 input files**, configured via CLI flags:

| CLI Flag | Registry Key | File Purpose | Format |
|----------|-------------|--------------|---------|
| `--plasmid-input` (`-i`) | `PlasmidReader` | Main plasmid data | TSV |
| `--plasmid-annotator-input` (`-a`) | `PlasmidAnnotatorReader` | Annotator/timestamp mapping | CSV |
| `--plasmid-pub-input` (`-p`) | `PlasmidPubReader` | Publication links | TSV |
| `--plasmid-gene-input` (`-g`) | `PlasmidGeneReader` | Gene associations | TSV |

## Field Mapping

### Main Plasmid File (`plasmid-input`)
**Format**: Tab-separated values
**Structure** (from `internal/datasource/csv/stockcenter/plasmid.go:56-65`):

| Index | Field | Required | Validation | Notes |
|-------|-------|----------|------------|-------|
| 0 | `Id` | ✅ Yes | Must not be empty | Plasmid identifier |
| 1 | `Name` | ✅ Yes | Must not be empty | Auto-prefixed with 'p' if missing |
| 2 | `Summary` | ✅ Yes | Must not be empty | Plasmid description |

### Annotator File (`plasmid-annotator-input`)
**Format**: CSV
**Structure** (from `internal/datasource/csv/stockcenter/annotator.go:38-64`):

| Index | Field | Purpose | Notes |
|-------|-------|---------|-------|
| 0 | Stock ID | Key for lookup | Maps to plasmid ID |
| 1 | Annotator Code | User mapping | Mapped via `annMap` to email |
| 2 | Created Date | Creation timestamp | Format: `2006-01-02 15:04:05` |
| 3 | Updated Date | Update timestamp | Format: `2006-01-02 15:04:05` |

**Enriches**: `User`, `CreatedOn`, `UpdatedOn` fields

### Publication File (`plasmid-pub-input`)
**Format**: Tab-separated
**Structure** (from `internal/datasource/csv/stockcenter/publication.go:25-49`):

| Index | Field | Purpose | Notes |
|-------|-------|---------|-------|
| 0 | Stock ID | Key for lookup | Maps to plasmid ID |
| 1 | Publication ID | Publication reference | Skips entries starting with 'd' |

**Enriches**: `Publications` array (can have multiple entries per plasmid)

### Gene File (`plasmid-gene-input`)
**Format**: Tab-separated
**Structure** (from `internal/datasource/csv/stockcenter/gene.go:26-47`):

| Index | Field | Purpose | Notes |
|-------|-------|---------|-------|
| 0 | Stock ID | Key for lookup | Maps to plasmid ID |
| 1 | Gene ID | Gene identifier | - |

**Enriches**: `Genes` array (can have multiple entries per plasmid)

## Complete Field Summary

### Required Fields (Validation Failures)
From `internal/load/stockcenter/plasmid.go:186-221`:

1. **`Id`** - Required, validated in `validateRequiredFields`
2. **`Name`** - Required, validated in `validateRequiredFields`
3. **`User`** - Required, validated in `validateUserAssignment` (must have annotator data)

### Optional Fields (No Validation)

4. **`Summary`** - Parsed from main file index 2 (appears required by parser but could be empty)
5. **`CreatedOn`** - Enriched from annotator file (optional - missing if no annotator)
6. **`UpdatedOn`** - Enriched from annotator file (optional - missing if no annotator)
7. **`Publications`** - Enriched from publication file (optional - warning logged if empty)
8. **`Genes`** - Enriched from gene file (optional - silently omitted if empty)

## Validation Rules

From `internal/load/stockcenter/plasmid.go:185-238`:

1. **User Assignment** (`validateUserAssignment`):
   - `User` field must not be empty
   - This comes from annotator lookup, so annotator file is effectively required

2. **Required Fields** (`validateRequiredFields`):
   - `Id` must not be empty string
   - `Name` must not be empty string

3. **Publication Warning**:
   - If no publications found, warning logged: `"plasmid %s has no publication entry"` (plasmid.go:564)
   - Processing continues (not a hard error)

## Data Flow

1. **CLI Flags** (`plasmid.go:78-93`) → **Viper Config**
2. **PreRun Hook** (`plasmid.go:33-40`) → Initializes readers from files or S3
3. **Lookup Initialization** (`plasmid.go:80-156`):
   - Annotator lookup
   - Publication lookup
   - Gene lookup
4. **CSV Reader Creation** (`plasmid.go:158-171`) with all lookups injected
5. **Record Processing** (plasmid.go:389-443):
   - Parse main fields (Id, Name, Summary)
   - Enrich with annotator data
   - Enrich with publications
   - Enrich with genes
6. **API Call** (plasmid.go:240-316):
   - Check if plasmid exists
   - Create or update via gRPC

## Code References

- Main CLI: `internal/cli/stockcenter/plasmid.go`
- Load Logic: `internal/load/stockcenter/plasmid.go`
- Plasmid Struct: `internal/datasource/csv/stockcenter/plasmid.go:56-65`
- Annotator Lookup: `internal/datasource/csv/stockcenter/annotator.go`
- Publication Lookup: `internal/datasource/csv/stockcenter/publication.go`
- Gene Lookup: `internal/datasource/csv/stockcenter/gene.go`
- Registry Keys: `internal/registry/stockcenter/keys.go`

## Summary

**Minimum Required Files**: 2
- Main plasmid file (Id, Name, Summary)
- Annotator file (provides required User field)

**Optional But Recommended**: 2
- Publication file (warning if missing)
- Gene file (silently optional)

**Critical Fields**: Id, Name, User
**Optional Fields**: Summary (parsed but could be empty), Publications, Genes, timestamps
