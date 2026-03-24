## statefaker

---

A vibecoded utlitiy to generate fake state data that is not at all reflective of the real aws provider but highly realistic.

Used to generate stressful state payloads to test with HCP Terraform and friends

#### Usage

`make build`

`statefaker -outputs 2000 -resources 60000 > huggggggge.tfstate`

Some resources will contain multiple instances using a string index key. Some resources will be in modules. There are many other options! Use `statefaker -help` for more configuration.

#### Development

Requires terraform to run tests. Use `make test`

`make fmt`

## Uploading Multiple State Versions

To simulate a busy workspace with many state versions, use the included Ruby script `scripts/push_multiple_states.rb`.

This script will:
1. Generate a new random state file using `statefaker`.
2. Patch the `serial` and `lineage` to match your target workspace.
3. Push the state file to TFC.
4. Repeat for N iterations.

### Usage

1. **Build statefaker** (if you haven't already):
   ```bash
   go build -o statefaker ./cmd/statefaker
   ```

2. **Run the script**:
   Point the script to the directory containing your Terraform configuration (where `terraform init` was run).

   ```bash
   # Push 5 new state versions to the workspace configured in ../my-infra
   ruby scripts/push_multiple_states.rb --cwd ../my-infra --iterations 5
   ```

   **Options:**
   - `--cwd DIR`: Directory containing Terraform config (default: current dir).
   - `--iterations N`: Number of state versions to push (default: 5).
   - `--resources N`: Number of resources to generate per state (default: 1000).
   - `--outputs N`: Number of outputs to generate per state (default: 0).
   - `--bin PATH`: Path to statefaker binary (default: searches PATH/cwd).
   - `--verbose`: Run verbosely.

   **Note:** The script automatically fetches the current `serial` and `lineage` from the remote workspace before pushing. Ensure you are authenticated with Terraform Cloud/Enterprise.
