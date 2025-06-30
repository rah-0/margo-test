![MarGO logo](https://github.com/rah-0/margo-test/blob/master/margo.png "MariaDB's Sea Lion with Golang's Gopher")

[![Go Report Card](https://goreportcard.com/badge/github.com/rah-0/margo?v=1)](https://goreportcard.com/report/github.com/rah-0/margo)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

<a href="https://www.buymeacoffee.com/rah.0" target="_blank"><img src="https://cdn.buymeacoffee.com/buttons/v2/arial-orange.png" alt="Buy Me A Coffee" height="50"></a>

# MarGO Test & Benchmarks

This repository contains tests and benchmarks for [MarGO](https://github.com/rah-0/margo), a simple, reflection-free ORM that maps MariaDB table schemas to Go structs.

The `dbs` directory contains the output of MarGO's code generation, with structures and database interaction methods automatically generated from the database schema.

## About

MarGO (MariaDB + GO) has the following features:

- Reflection-free design for improved performance
- Direct mapping between MariaDB schemas and Go structs
- Minimal overhead compared to raw SQL
- Simple, intuitive API

## Benchmarks

This repository includes benchmarks comparing MarGO with other popular Go ORMs:

- Raw SQL (baseline)
- MarGO
- Bun
- Ent
- GORM

Check out [BENCHMARKS.md](./BENCHMARKS.md) for performance comparisons.

## Key Results

MarGO performs well in benchmarks, with performance metrics close to raw SQL:

Unlike other ORMs that can introduce significant performance penalties, MarGO maintains speed while providing some convenience.

[![Buy Me A Coffee](https://cdn.buymeacoffee.com/buttons/default-orange.png)](https://www.buymeacoffee.com/rah.0)
