## k8s completion zsh

Generate the autocompletion script for zsh

### Synopsis

Generate the autocompletion script for the zsh shell.

If shell completion is not already enabled in your environment you will need
to enable it.  You can execute the following once:

	echo "autoload -U compinit; compinit" >> ~/.zshrc

To load completions in your current shell session:

	source <(k8s completion zsh); compdef _k8s k8s

To load completions for every new session, execute once:

#### Linux:

	k8s completion zsh > "${fpath[1]}/_k8s"

#### macOS:

	k8s completion zsh > $(brew --prefix)/share/zsh/site-functions/_k8s

You will need to start a new shell for this setup to take effect.


```
k8s completion zsh [flags]
```

### Options

```
  -h, --help              help for zsh
      --no-descriptions   disable completion descriptions
```

### Options inherited from parent commands

```
      --access-key string             access key for S3 server
      --kubeconfig string             path to the kubernetes client(kubeconfig) file[REQUIRED]
      --log-file string               file for log output other than standard output, written to a temp folder by default
      --log-file-bucket string        S3 bucket for log file (default "dictybase")
      --log-file-bucket-path string   S3 path inside the bucket for storing log file (default "import/log")
      --log-format string             format of the logging out, either of json or text (default "json")
      --log-level string              log level for the application (default "error")
      --namespace string              kubernetes namespace (default "dictybase")
      --s3-bucket string              S3 bucket for input files (default "dictybase")
      --s3-bucket-path string         path inside S3 bucket for input files[REQUIRED]
      --s3-server string              S3 server endpoint (default "minio")
      --s3-server-port string         S3 server port
      --secret-key string             secret key for S3 server
```

### SEE ALSO

* [k8s completion](k8s_completion.md)	 - Generate the autocompletion script for the specified shell

###### Auto generated by spf13/cobra on 15-Jul-2022