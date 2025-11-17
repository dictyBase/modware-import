This file is a merged representation of a subset of the codebase, containing files not matching ignore patterns, combined into a single document by Repomix.
The content has been processed where comments have been removed, empty lines have been removed, content has been compressed (code blocks are separated by ⋮---- delimiter).

# File Summary

## Purpose
This file contains a packed representation of a subset of the repository's contents that is considered the most important context.
It is designed to be easily consumable by AI systems for analysis, code review,
or other automated processes.

## File Format
The content is organized as follows:
1. This summary section
2. Repository information
3. Directory structure
4. Repository files (if enabled)
5. Multiple file entries, each consisting of:
  a. A header with the file path (## File: path/to/file)
  b. The full contents of the file in a code block

## Usage Guidelines
- This file should be treated as read-only. Any changes should be made to the
  original repository files, not this packed version.
- When processing this file, use the file path to distinguish
  between different files in the repository.
- Be aware that this file may contain sensitive information. Handle it with
  the same level of security as you would the original repository.

## Notes
- Some files may have been excluded based on .gitignore rules and Repomix's configuration
- Binary files are not included in this packed representation. Please refer to the Repository Structure section for a complete list of file paths, including binary files
- Files matching these patterns are excluded: **/**organism**
- Files matching patterns in .gitignore are excluded
- Files matching default ignore patterns are excluded
- Code comments have been removed from supported file types
- Empty lines have been removed from all files
- Content has been compressed - code blocks are separated by ⋮---- delimiter
- Files are sorted by Git change count (files with more changes are at the bottom)

# Directory Structure
```
stock.proto
```

# Files

## File: stock.proto
```protobuf
syntax = "proto3";

package dictybase.stock;

import "dictybase/api/upload/file.proto";
import "github.com/mwitkow/go-proto-validators/validator.proto";
import "google/protobuf/empty.proto";
import "google/protobuf/timestamp.proto";

option cc_enable_arenas = true;
option go_package = "github.com/dictyBase/go-genproto/dictybaseapis/stock;stock";
option java_multiple_files = true;
option java_outer_classname = "StockProto";
option java_package = "org.dictybase.stock";
option objc_class_prefix = "DICTYAPI";

// The stock service specification
service StockService {
  // Retrieves strain by ID
  rpc GetStrain(StockId) returns (Strain) {}
  // Retrieves stock by ID
  rpc GetPlasmid(StockId) returns (Plasmid) {}
  // Create a new strain
  rpc CreateStrain(NewStrain) returns (Strain) {}
  // Create a new plasmid
  rpc CreatePlasmid(NewPlasmid) returns (Plasmid) {}
  // Update an existing strain
  rpc UpdateStrain(StrainUpdate) returns (Strain) {}
  // Update an existing plasmid
  rpc UpdatePlasmid(PlasmidUpdate) returns (Plasmid) {}
  // Remove an existing stock
  rpc RemoveStock(StockId) returns (google.protobuf.Empty) {}
  // List strains using pagination, ten entries are retrieved by default
  rpc ListStrains(StockParameters) returns (StrainCollection) {}
  // List strains using strain id without any pagination
  rpc ListStrainsByIds(StockIdList) returns (StrainList) {}
  // List plasmids using pagination, ten entries are retrieved by default
  rpc ListPlasmids(StockParameters) returns (PlasmidCollection) {}
  // Load existing strain
  rpc LoadStrain(ExistingStrain) returns (Strain) {}
  // Load existing plasmid
  rpc LoadPlasmid(ExistingPlasmid) returns (Plasmid) {}
  // Upload obojson formatted file through client side streaming
  rpc OboJSONFileUpload(stream dictybase.api.upload.FileUploadRequest) returns (dictybase.api.upload.FileUploadResponse) {}
}

message StockId {
  // Unique identifier for the stock
  string id = 1 [(validator.field) = {string_not_empty: true}];
}

// Definition for list of unique stock identifier
message StockIdList {
  repeated string id = 1 [(validator.field) = {regex: "^DB(P|S)[0-9]{5,}$"}];
}

// Definition of an individual strain
message Strain {
  message Data {
    // Resource name
    string type = 1;
    // Unique identifier for the strain
    string id = 2;
    StrainAttributes attributes = 3;
  }
  Data data = 1 [(validator.field) = {msg_exists: true}];
}

// Definition of an individual plasmid
message Plasmid {
  message Data {
    // Resource name
    string type = 1;
    // Unique identifier for the plasmid
    string id = 2;
    PlasmidAttributes attributes = 3;
  }
  Data data = 1 [(validator.field) = {msg_exists: true}];
}

// Definition of various strain attributes
message StrainAttributes {
  // Timestamp for creation
  google.protobuf.Timestamp created_at = 1 [(validator.field) = {msg_exists: true}];
  // Timestamp for update
  google.protobuf.Timestamp updated_at = 2 [(validator.field) = {msg_exists: true}];
  // User who created stock entry
  string created_by = 3 [(validator.field) = {string_not_empty: true}];
  // User who updated stock entry
  string updated_by = 4 [(validator.field) = {string_not_empty: true}];
  // Summary of the stock
  string summary = 5;
  // Editable version of the stock summary (Slate JSON format)
  string editable_summary = 6;
  // Depositor of the stock
  string depositor = 7;
  // List of associated genes
  repeated string genes = 8;
  // List of database cross references
  repeated string dbxrefs = 9;
  // List of related publications
  repeated string publications = 10;
  // Descriptor for the strain, a quick overview of its key genetic modifications
  string label = 11 [(validator.field) = {string_not_empty: true}];
  // Species of the strain
  string species = 12 [(validator.field) = {string_not_empty: true}];
  // Related plasmid for the strain
  string plasmid = 13;
  // Parent of the strain
  string parent = 14;
  // List of names for the strain
  repeated string names = 15;
  // dictybase specific strain property that will
  // map to dicty_strain_property ontology
  string dicty_strain_property = 16;
}

// Definition of various stock attributes
message PlasmidAttributes {
  // Timestamp for creation
  google.protobuf.Timestamp created_at = 1 [(validator.field) = {msg_exists: true}];
  // Timestamp for update
  google.protobuf.Timestamp updated_at = 2 [(validator.field) = {msg_exists: true}];
  // User who created stock entry
  string created_by = 3 [(validator.field) = {string_not_empty: true}];
  // User who updated stock entry
  string updated_by = 4 [(validator.field) = {string_not_empty: true}];
  // Summary of the stock
  string summary = 5;
  // Editable version of the stock summary (Slate JSON format)
  string editable_summary = 6;
  // Depositor of the stock
  string depositor = 7;
  // List of associated genes
  repeated string genes = 8;
  // List of database cross references
  repeated string dbxrefs = 9;
  // List of related publications
  repeated string publications = 10;
  // Image map for the plasmid
  string image_map = 11;
  // Sequence for the plasmid
  string sequence = 12;
  // Unambiguous name for the plasmid
  string name = 13 [(validator.field) = {string_not_empty: true}];
  // dictybase specific plasmid property that will
  // map to dicty_plasmid_keyword ontology
  string dicty_plasmid_property = 14;
}

// Definition for creating a new strain
message NewStrain {
  message Data {
    // Resource name
    string type = 1;
    NewStrainAttributes attributes = 2;
  }
  Data data = 1 [(validator.field) = {msg_exists: true}];
}

// Definition for creating a new plasmid
message NewPlasmid {
  message Data {
    // Resource name
    string type = 1;
    NewPlasmidAttributes attributes = 2;
  }
  Data data = 1 [(validator.field) = {msg_exists: true}];
}

// Defines attributes for creating a new strain
message NewStrainAttributes {
  // User who created stock entry
  string created_by = 1 [(validator.field) = {string_not_empty: true}];
  // User who updated stock entry
  string updated_by = 2 [(validator.field) = {string_not_empty: true}];
  // Summary of the stock
  string summary = 3;
  // Editable version of the stock summary (Slate JSON format)
  string editable_summary = 4;
  // List of associated genes
  repeated string genes = 5;
  // List of database cross references
  repeated string dbxrefs = 6;
  // Depositor of the stock
  string depositor = 7 [(validator.field) = {string_not_empty: true}];
  // List of related publications
  repeated string publications = 8;
  // Descriptor for the strain, a quick overview of its key genetic modifications
  string label = 9 [(validator.field) = {string_not_empty: true}];
  // Species of the strain
  string species = 10 [(validator.field) = {string_not_empty: true}];
  // Related plasmid for the strain
  string plasmid = 11;
  // Parent of the strain
  string parent = 12;
  // List of names for the strain
  repeated string names = 13;
  // dictybase specific strain property that will
  // map to dicty_strain_property ontology
  string dicty_strain_property = 14;
}

// Defines attributes for creating a new plasmid
message NewPlasmidAttributes {
  // User who created stock entry
  string created_by = 1 [(validator.field) = {string_not_empty: true}];
  // User who updated stock entry
  string updated_by = 2 [(validator.field) = {string_not_empty: true}];
  // Summary of the stock
  string summary = 3;
  // Editable version of the stock summary (Slate JSON format)
  string editable_summary = 4;
  // List of associated genes
  repeated string genes = 5;
  // List of database cross references
  repeated string dbxrefs = 6;
  // Depositor of the stock
  string depositor = 7 [(validator.field) = {string_not_empty: true}];
  // List of related publications
  repeated string publications = 8;
  // Image map for the plasmid
  string image_map = 9;
  // Sequence for the plasmid
  string sequence = 10;
  // Unambiguous name for the plasmid
  string name = 11 [(validator.field) = {string_not_empty: true}];
  // dictybase specific plasmid property that will
  // map to dicty_plasmid_keyword ontology
  string dicty_plasmid_property = 12;
}

// Definition for loading an existing strain
message ExistingStrain {
  message Data {
    // Resource name
    string type = 1;
    // Existing strain ID
    string id = 2;
    ExistingStrainAttributes attributes = 3;
  }
  Data data = 1 [(validator.field) = {msg_exists: true}];
}

// Definition for loading an existing plasmid
message ExistingPlasmid {
  message Data {
    // Resource name
    string type = 1;
    // Existing plasmid ID
    string id = 2;
    ExistingPlasmidAttributes attributes = 3;
  }
  Data data = 1 [(validator.field) = {msg_exists: true}];
}

// Defines attributes for loading an existing strain
message ExistingStrainAttributes {
  // Timestamp for creation
  google.protobuf.Timestamp created_at = 1 [(validator.field) = {msg_exists: true}];
  // Timestamp for update
  google.protobuf.Timestamp updated_at = 2 [(validator.field) = {msg_exists: true}];
  // User who created stock entry
  string created_by = 3 [(validator.field) = {string_not_empty: true}];
  // User who updated stock entry
  string updated_by = 4 [(validator.field) = {string_not_empty: true}];
  // Summary of the stock
  string summary = 5;
  // Editable version of the stock summary (Slate JSON format)
  string editable_summary = 6;
  // List of associated genes
  repeated string genes = 7;
  // List of database cross references
  repeated string dbxrefs = 8;
  // Depositor of the stock
  string depositor = 9;
  // List of related publications
  repeated string publications = 10;
  // Descriptor for the strain, a quick overview of its key genetic modifications
  string label = 11 [(validator.field) = {string_not_empty: true}];
  // Species of the strain
  string species = 12 [(validator.field) = {string_not_empty: true}];
  // Related plasmid for the strain
  string plasmid = 13;
  // Parent of the strain
  string parent = 14;
  // List of names for the strain
  repeated string names = 15;
  // dictybase specific strain property that will
  // map to dicty_strain_property ontology
  string dicty_strain_property = 16;
}

// Defines attributes for loading an existing plasmid
message ExistingPlasmidAttributes {
  // Timestamp for creation
  google.protobuf.Timestamp created_at = 1 [(validator.field) = {msg_exists: true}];
  // Timestamp for update
  google.protobuf.Timestamp updated_at = 2 [(validator.field) = {msg_exists: true}];
  // User who created stock entry
  string created_by = 3 [(validator.field) = {string_not_empty: true}];
  // User who updated stock entry
  string updated_by = 4;
  // Summary of the stock
  string summary = 5;
  // Editable version of the stock summary (Slate JSON format)
  string editable_summary = 6;
  // List of associated genes
  repeated string genes = 7;
  // List of database cross references
  repeated string dbxrefs = 8;
  // Depositor of the stock
  string depositor = 9;
  // List of related publications
  repeated string publications = 10;
  // Image map for the plasmid
  string image_map = 11;
  // Sequence for the plasmid
  string sequence = 12;
  // Unambiguous name for the plasmid
  string name = 13 [(validator.field) = {string_not_empty: true}];
  // dictybase specific plasmid property that will
  // map to dicty_plasmid_keyword ontology
  string dicty_plasmid_property = 14;
}

// Definition for creating a new strain
message StrainUpdate {
  message Data {
    // Resource name
    string type = 1;
    // Unique ID for strain
    string id = 2 [(validator.field) = {string_not_empty: true}];
    StrainUpdateAttributes attributes = 3;
  }
  Data data = 1 [(validator.field) = {msg_exists: true}];
}

// Definition for creating a new plasmid
message PlasmidUpdate {
  message Data {
    // Resource name
    string type = 1;
    // Unique ID for plasmid
    string id = 2 [(validator.field) = {string_not_empty: true}];
    PlasmidUpdateAttributes attributes = 3;
  }
  Data data = 1 [(validator.field) = {msg_exists: true}];
}

// Defines attributes for updating a strain
message StrainUpdateAttributes {
  // User who updated stock entry
  string updated_by = 1 [(validator.field) = {string_not_empty: true}];
  // Summary of the stock
  string summary = 2;
  // Editable version of the stock summary (Slate JSON format)
  string editable_summary = 3;
  // Depositor of the stock
  string depositor = 4;
  // List of associated genes
  repeated string genes = 5;
  // List of database cross references
  repeated string dbxrefs = 6;
  // List of related publications
  repeated string publications = 7;
  // Descriptor for the strain, a quick overview of its key genetic modifications
  string label = 8;
  // Species of the strain
  string species = 9;
  // Related plasmid for the strain
  string plasmid = 10;
  // Parent of the strain
  string parent = 11;
  // List of names for the strain
  repeated string names = 12;
  // dictybase specific strain property that will
  // map to dicty_strain_property ontology
  string dicty_strain_property = 13;
}

// Defines attributes for updating a plasmid
message PlasmidUpdateAttributes {
  // User who updated stock entry
  string updated_by = 1 [(validator.field) = {string_not_empty: true}];
  // Summary of the stock
  string summary = 2;
  // Editable version of the stock summary (Slate JSON format)
  string editable_summary = 3;
  // Depositor of the stock
  string depositor = 4;
  // List of associated genes
  repeated string genes = 5;
  // List of database cross references
  repeated string dbxrefs = 6;
  // List of related publications
  repeated string publications = 7;
  // Image map for the plasmid
  string image_map = 8;
  // Sequence for the plasmid
  string sequence = 9;
  // Unambiguous name for the plasmid
  string name = 10;
  // dictybase specific plasmid property that will
  // map to dicty_plasmid_keyword ontology
  string dicty_plasmid_property = 11;
}

// List of strains without any metadata for pagination
message StrainList {
  message Data {
    // Resource name
    string type = 1;
    // Unique identifier for the stock
    string id = 2;
    StrainAttributes attributes = 3;
  }
  repeated Data data = 1;
}

// List of strains
message StrainCollection {
  message Data {
    // Resource name
    string type = 1;
    // Unique identifier for the stock
    string id = 2;
    StrainAttributes attributes = 3;
  }
  repeated Data data = 1;
  Meta meta = 2 [(validator.field) = {msg_exists: true}];
}

// List of plasmids
message PlasmidCollection {
  message Data {
    // Resource name
    string type = 1;
    // Unique identifier for the stock
    string id = 2;
    PlasmidAttributes attributes = 3;
  }
  repeated Data data = 1;
  Meta meta = 2 [(validator.field) = {msg_exists: true}];
}

// StockParameters defines fields for manipulating output of Stock collection
message StockParameters {
  // A unique pointer to the next set of result in the list (default is 0)
  int64 cursor = 1;
  // Maximum number of records that can be fetch per request (default is 10)
  int64 limit = 2;
  // The `filter` field restricts the data return by the collection. To use
  // it, supply one or multiple allowed fields to filter followed
  // by a filter expression. It uses the following syntax...
  //        field_name operator expression
  //
  // The following fields of `StockAttributes` definition are allowed to
  // be used for filtering
  //   * depositor          - Depositor of the stock (string)
  //   * parent             - Parental strain (string) (currently not implemented)
  //   * plasmid            - Related plasmid for the strain (string)
  //   * species            - The species of the strain (string)
  //   * summary            - Summary of the stock (string)
  //   * name               - Name used for strain (string), searches in the "names" attribute
  //   * descriptor         - Descriptor for the strain (string), searches in the "label" attribute
  //   * plasmid_name       - Name used for plasmid (string)
  //   * created_at         - Date the stock was created (number), can be in the
  //                          following formats:
  //                          YYYY-MM-DD, YYYY-MM, YYYY
  //   * updated_at         - Date the stock was updated (number), can be in the
  //                          following formats:
  //                          YYYY-MM-DD, YYYY-MM, YYYY
  //
  // field_name - Any one of the allowed field_name of the `StockAttributes` definition.
  // operator - Defines the type of filter match to use. It could be any of
  // the following four and all of them should be URL-encoded for http request.
  //
  //        Operators for strings
  //              =~   Contains substring
  //              !~   Not contains substring
  //              ===  Equals
  //              !=   Not equals
  //
  //        Operators for number
  //              ==  Equals
  //              >   Greater than
  //              <   Less than
  //              <=  Less than equal to
  //              >=  Greater than equal to
  //
  //        Operators for dates
  //              $==  Equals
  //              $>   Greater than
  //              $<   Less than
  //              $<=  Less than equal to
  //              $>=  Greater than equal to
  //
  //        Operators for items in arrays
  //              @=~   Contains substring
  //              @!~   Not contains substring (not implemented yet)
  //              @==   Equals
  //              @!=   Not equals
  //
  // expression - The value that will be included or excluded from the
  // result. URL-reserved characters must be URL-encoded for http request.
  //
  //           filter: "created_at$>=2018-12-01"
  //           filter: "depositor===Costanza"
  //
  // Filter can be combined using OR or AND boolean logic.
  //   * The OR is represented using a comma(,).
  //   * The AND is represented using a semi-colon(;).
  //   * AND and OR operators can be combined and AND takes precedence over OR.
  //
  //           filter: "depositor===Benes;created_at$>=2018-12-01"
  //
  string filter = 3;
  // The sort field allow to sort the data return by the collection based on
  // fields of `StockAttributes. To use it, supply a comma separated one
  // or more allowed field from the definition of `StockAttributes`.
}

// Metadata definition for traversing the collection
message Meta {
  // A unique pointer to the next set of result in the collection. Set the
  // cursor value parameter to the value of next_cursor to retrieve the next
  // set of collection using the same method
  int64 next_cursor = 1;
  // Maximum number of records that can be fetch per request
  int64 limit = 2;
  // Total number of records in the collection.
  int64 total = 3;
}
```
