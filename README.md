# changelog

A tool help generate CHANGELOG.md from git.

## Installation

```
go install github.com/xpunch/changelog
```

## Usage

### under target git repository

```
changelog
```

### print logs

```
changelog --verbose
```

### fetch latest repository

```
changelog --fetch
```

### set target repository folder and target CHANGELOG.md path

```
changelog --source ~/gitrepo --output ~/gitrepo/CHANGELOG.md --fetch --verbose
```
