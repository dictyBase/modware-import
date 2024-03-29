## import ontology load

load obojson formatted ontologies to arangodb

```
import ontology load [flags]
```

### Options

```
      --cv-collection string     arangodb collection for storing ontology information (default "cv")
  -f, --folder string            input folder with obojson format files
      --group string             file belong to this ontology group will be uploaded. Only works for S3 storage
  -h, --help                     help for load
      --obograph string          arangodb named graph for managing ontology graph (default "obograph")
      --rel-collection string    arangodb collection for storing cvterm relationships (default "cvterm_relationship")
      --term-collection string   arangodb collection for storing ontoloy terms (default "cvterm")
```

### Options inherited from parent commands

```
      --access-key string             access key for S3 server
      --arangodb-database string      arangodb database name
  -H, --arangodb-host string          arangodb database host (default "arangodb")
  -p, --arangodb-pass string          arangodb database password
      --arangodb-port int             arangodb database port (default 8529)
  -u, --arangodb-user string          arangodb database user
      --input-source string           source of the file, could be one of bucket or folder (default "bucket")
      --is-secure                     flag for secured or unsecured arangodb endpoint
      --log-file string               file for log output other than standard output, written to a temp folder by default
      --log-file-bucket string        S3 bucket for log file (default "dictybase")
      --log-file-bucket-path string   S3 path inside the bucket for storing log file (default "import/log")
      --log-format string             format of the logging out, either of json or text (default "json")
      --log-level string              log level for the application (default "error")
      --s3-bucket string              S3 bucket for input files (default "dictybase")
      --s3-bucket-path string         path inside S3 bucket for input files[REQUIRED]
      --s3-server string              S3 server endpoint (default "minio")
      --s3-server-port string         S3 server port
      --secret-key string             secret key for S3 server
```

### SEE ALSO

* [import ontology](import_ontology.md)	 - subcommand for ontology management

###### Auto generated by spf13/cobra on 15-Jul-2022
