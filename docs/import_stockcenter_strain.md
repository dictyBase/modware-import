## import stockcenter strain

load stockcenter strain data

### Synopsis

load stockcenter strain data

```
import stockcenter strain [flags]
```

### Options

```
  -h, --help                            help for strain
  -i, --input string                    csv file with strain data
      --stock-grpc-host string          grpc host address for stock service (default "stock-api")
      --stock-grpc-port string          grpc port for stock service
  -a, --strain-annotator-input string   csv file that provides mapping among strain identifier, annotator and annotation timestamp
  -g, --strain-gene-input string        csv file that maps strains to gene identifiers
  -p, --strain-pub-input string         csv file that maps strains to publication identifiers
```

### Options inherited from parent commands

```
      --access-key string             access key for S3 server
      --input-source string           source of the file, could be one of bucket or folder (default "bucket")
      --log-file string               file for log output other than standard output, written to a temp folder by default
      --log-file-bucket-path string   S3 bucket path where log file will be stored (default "dictybase/import/log")
      --log-format string             format of the logging out, either of json or text (default "json")
      --log-level string              log level for the application (default "error")
      --s3-bucket-path string         S3 bucket path where all stockcenter data will be kept (default "dictybase/import/stockcenter")
      --s3-server string              S3 server endpoint (default "minio")
      --s3-server-port string         S3 server port
      --secret-key string             secret key for S3 server
```

### SEE ALSO

* [import stockcenter](import_stockcenter.md)	 - subcommand for stockcenter data loading

###### Auto generated by spf13/cobra on 19-Jun-2019