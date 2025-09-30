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
