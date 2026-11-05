# govox

A golang based voxel engine.

## Building, Running, Packaging

```bash
make run # This will build/rebuild on source changes and then run the project.
```

If you want to build the project without running it, you can use:

```bash
make build # Build executes go mod tidy every time, so watch out for unintended changes to the go.mod file.
```

You can create a binary for "all" platforms using:

```bash
make package # This will create a binary for linux, Windows (both amd64) and macos (arm64) in the ./dist directory.
```

Tests are runnable using:

```bash
make test # This will run all tests in the project.
```

## License

This project is licensed under the [GNU Affero General Public License v3.0](LICENSE). Any distributed version of this project must include a copy of the license and remain free software.
