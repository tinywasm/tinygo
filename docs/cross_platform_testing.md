# Cross-platform Testing with Docker — Feasibility Analysis

## Linux
**Viable.** Docker nativo. Imagen base `golang:1.25` es suficiente.

```dockerfile
FROM golang:1.25
COPY . /app
WORKDIR /app
RUN go test ./...
```

Multiple distros testables: `ubuntu`, `debian`, `alpine`, `fedora`.

## Windows
**Parcialmente viable.** Docker tiene Windows containers pero:
- Requieren host Windows con Hyper-V (`docker run --isolation=hyperv`)
- No funcionan en hosts Linux/macOS
- Imagenes grandes (~5-15GB): `mcr.microsoft.com/windows/servercore`
- **Alternativa real**: GitHub Actions con `runs-on: windows-latest` — funciona bien, gratis para open source

## macOS
**No viable con Docker.** Apple no licencia macOS para virtualización fuera de hardware Apple.
- No existen imágenes Docker oficiales de macOS
- Existen proyectos como `sickcodes/Docker-OSX` pero violan EULA de Apple y son inestables
- **Alternativa real**: GitHub Actions con `runs-on: macos-latest` — funciona bien, gratis para open source

## Recomendación

| Plataforma | Docker | GitHub Actions | Costo |
|------------|--------|---------------|-------|
| Linux (multi-distro) | Si | Si | Gratis |
| Windows | Solo en host Windows | Si | Gratis (open source) |
| macOS | No | Si | Gratis (open source) |

**GitHub Actions es la única solución que cubre los 3 OS de forma real.** Docker solo sirve para Linux multi-distro.

Ejemplo de matrix en GitHub Actions:

```yaml
strategy:
  matrix:
    os: [ubuntu-latest, macos-latest, windows-latest]
runs-on: ${{ matrix.os }}
steps:
  - uses: actions/checkout@v4
  - uses: actions/setup-go@v5
    with:
      go-version: '1.25'
  - run: go test ./...
```
