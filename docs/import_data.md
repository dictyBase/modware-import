## import data

subcommand for managing data files

### Options

```
  -h, --help   help for data
```

### Options inherited from parent commands

```
      --access-key string             access key for S3 server
      --input-source string           source of the file, could be one of bucket or folder (default "bucket")
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

* [import](import.md)	 - cli for importing dictybase data for migration
* [import data refresh](import_data_refresh.md)	 - updates data files in S3 storage

###### Auto generated by spf13/cobra on 15-Jul-2022
