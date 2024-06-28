# Pattern-Signature generation

`apicactl` is equipped with log Pattern-Signature (PS) generation and post-PS statistical analysis. You can automatically derive common text patterns in all the logs ingested by the `apicactl` client using the `-g` flag. You can run Pattern-Signature generation using the `-g` flag in the following commands:

- [`apicactl logs`](logs/apicactl_logs.md)
- [`apicactl logs search`](/logs/apicactl_logs_search)
- [`apicactl tail`](/tail/apicactl_tail)

The executable binary [`psmod`](https://github.com/logiqai/logiqctl/releases/tag/2.1.2) processes PS generation. In order to run the `psmod` binary, ensure that you copy the binary that's suitable for your platform's architecture or operating system, rename it to `psmod`, and place it in the same location as `apicactl`. For example, if you're if you're running `apicactl` on a Darwin AMD64-based machine, copy the `apicactl_darwin_amd64` and `psmod_darwin_amd64` binaries to a folder and rename `psmod_darwin_amd64` to `psmod` before running `apicactl`. 

Generated pattern signatures are transferred to and available in the `ps_stat.out` file.