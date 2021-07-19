# Pattern-Signature generation

`logiqctl` is equipped with log Pattern-Signature (PS) generation and post-PS statistical analysis. You can automatically derive common text patterns in all the logs ingested by the `logiqctl` client using the `-g` flag. You can run Pattern-Signature generation using the `-g` flag in the following commands:

- [`logiqctl logs`](logs/logiqctl_logs.md)
- [`logiqctl logs search`](/logs/logiqctl_logs_search)
- [`logiqctl tail`](/tail/logiqctl_tail)

The executable binary [`psmod`](https://github.com/logiqai/logiqctl/releases/tag/2.1.2) processes PS generation. In order to run the `psmod` binary, ensure that you copy the binary that's suitable for your platform's architecture or operating system, rename it to `psmod`, and place it in the same location as `logiqctl`. For example, if you're if you're running `logiqctl` on a Darwin AMD64-based machine, copy the `logiqctl_darwin_amd64` and `psmod_darwin_amd64` binaries to a folder and rename `psmod_darwin_amd64` to `psmod` before running `logiqctl`. 

Generated pattern signatures are transferred to and available in the `ps_stat.out` file.