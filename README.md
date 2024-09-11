# jectmobi

## Execution

### Run


#### Local
For running the project we need to have go language installed in the host machine and available in the path.

```sh
  go mod tidy
  go run main.go
```
Here the command *go mod tidy* installs any libs necessary for the project to build. We are not using any third party library here. Only standard libs bundled with golang is used.

Finally *go run main.go* runs the code. The input json is embedded in the code inside the main function.

#### ReplIt
For replit import we need to import the github repo from the replit console and run the main.go file.
