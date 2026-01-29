# OTel Collector running on Wasm+JS

Setup steps:

1. Create a read-only GH PAT token, and paste it in a file called `gh_pat.key` at the repo root.
2. Put your username in a file called `gh_org`. This chooses with GitHub org/user to grab metrics from.
3. Run `make build` to build the Collector (Requires Go).
4. Run `make serve` (requires Python).
5. Navigate your browser to <http://localhost:8000>, open the console in your browser's dev tools, and click the "Run" button.
6. Wait a moment for metrics to be printed to the console.
