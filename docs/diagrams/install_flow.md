# Install Flow

```mermaid
flowchart TD
    A[EnsureInstalled] --> B{Binary found?\nGetPath}
    B -->|no| INST[Install]
    B -->|yes| C{tinygo version\n== required?}
    C -->|yes| RET[Return path]
    C -->|no — system install\nnot under installDir| WARN[Log: system version differs\ncannot manage — return path]
    WARN --> RET
    C -->|no — local install\nunder installDir| REINST[Log: version mismatch\nReinstall]
    REINST --> INST

    INST --> GOOS{runtime.GOOS?}
    GOOS -->|linux / darwin| TAR[Download .tar.gz\nfrom GitHub releases]
    GOOS -->|windows| SCOOP{scoop in PATH?}

    SCOOP -->|no| INSCOOP[Install Scoop\npowershell one-liner]
    INSCOOP --> SCOOPINSTALL[scoop install tinygo]
    SCOOP -->|yes| SCOOPINSTALL

    TAR --> EXTRACT[extractTarGz → installDir]
    EXTRACT --> CHECK{Binary exists\nat binPath?}
    CHECK -->|no| ERR[Cleanup partial\nreturn error]
    CHECK -->|yes| VER{tinygo version ok?}
    VER -->|fail| ERR2[Cleanup install\nreturn error]
    VER -->|ok| RET

    SCOOPINSTALL --> VER2{tinygo version ok?}
    VER2 -->|fail| ERR3[return error]
    VER2 -->|ok| RET
```

## Version matching rule

`DefaultVersion = "0.40.1"` is the required version.

| Situation | Action |
|-----------|--------|
| Installed version == `DefaultVersion` | Use as-is |
| Installed version ≠ `DefaultVersion` — **system install** (found in PATH, not under `installDir`) | Log warning, use as-is — cannot manage system packages |
| Installed version ≠ `DefaultVersion` — **local install** (under `installDir`) | Reinstall the correct version |
| Not installed | Install |
