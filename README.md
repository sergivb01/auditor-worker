# Auditor Worker
[![pr_welcome][pr_welcome]][pr_welcome-url]

Auditor worker is the worker (also known as _slave_) for the auditor project.
It uses [GRPC](https://grpc.io/), although I am planning to use **[bokchoy](https://github.com/thoas/bokchoy)** ~~NSQ (:heart:), Kafka or RabbitMQ~~. 

The worker process consists of (**simplified**):
 * Workers get a _bunch_ of files.
 * The compilation process is executed using [CCache](https://ccache.dev/) to speed builds up.
 * A _bunch_ of tests or execution-examples are run in parallel (if compilation was successful).
 * Results from the compilation (warnings) and execution get returned. If any output mismatches were captured, they will be returned as errors.

To write/read from the program during executing we pipe the `stdin`, `stdout` and `stderr`.  

## 📖 License
**Auditor Worker** is licensed under the [MIT License](LICENSE).

<!-- PR Welcome -->
[pr_welcome]: https://img.shields.io/badge/PRs-welcome-brightgreen.svg
[pr_welcome-url]: https://github.com/sergivb01/auditor-worker/pulls
