# Release Draft — v0.1.0 (First Release)

## Suggested release title
`v0.1.0 — Initial Public Release`

## Suggested tag
`v0.1.0`

## Release notes (copy/paste)
Welcome to the first release of **Memstore Seeder** 🎉

This release publishes prebuilt binaries for Linux, macOS, and Windows so you can run the service without building from source.

### What’s included
- Initial public release of the Memstore Seeder service.
- Cross-platform binaries for 64-bit systems.
- Basic project tooling for build, proto generation, and certificate generation.

### Downloads
Attach the following artifacts to this release:

- `seeder-linux-amd64`
- `seeder-macos-amd64`
- `seeder-windows-amd64.exe`

### Checksums (recommended)
For supply-chain safety, publish SHA256 checksums for all binaries in a `checksums.txt` file.

Example format:

```txt
<sha256>  seeder-linux-amd64
<sha256>  seeder-macos-amd64
<sha256>  seeder-windows-amd64.exe
```

### Quick start
After downloading the binary for your OS, run:

```bash
./seeder --help
```

> On Windows (PowerShell):

```powershell
.\seeder-windows-amd64.exe --help
```

### Notes
- This is the **first release**, so some flags/configuration and behavior may evolve.
- Please open issues for bugs, platform-specific behavior, or feature requests.

---

## Optional “GitHub Release” form fields

### Release name
`v0.1.0 — Initial Public Release`

### Description
Use the “Release notes (copy/paste)” section above.

### Assets to upload
- Linux: `seeder-linux-amd64`
- macOS: `seeder-macos-amd64`
- Windows: `seeder-windows-amd64.exe`
- Optional: `checksums.txt`
