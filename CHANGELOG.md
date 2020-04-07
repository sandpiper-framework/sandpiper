# Changelog

All notable changes to this project are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Want to Contribute?

Here are some areas where we could use some help.

- [ ] testing: Swap out [dockertest](github.com/fortytw2/dockertest) (used for db mocking) with something that also works on Windows.
- [ ] primary: Consider changing pagination from origin/limit to [keyset](https://use-the-index-luke.com/no-offset) or [seek](https://www.moesif.com/blog/technical/api-design/REST-API-Design-Filtering-Sorting-and-Pagination/) pattern.
- [ ] primary: Refactor SelectError() handling which adds an Echo dependency to platform-level code.
- [ ] cli: Add pagination logic to functions which return lists (requires a loop around those calls). This is really a bug (albeit intentional).

## [0.1.x] public - 2020-??-??
### Added
- Initial public release

## [0.1.0] beta - 2020-??-??
### Added
- Initial beta release

## [0.0.1] pre-alpha - 2019-11-11
### Added
- Project started using [GORSK](https://github.com/ribice/gorsk) as basic template with modified project structure.
