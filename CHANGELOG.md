# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Project Items To Complete

- [ ] primary: Consider changing pagination from origin/limit to [keyset](https://use-the-index-luke.com/no-offset) or [seek](https://www.moesif.com/blog/technical/api-design/REST-API-Design-Filtering-Sorting-and-Pagination/)
- [ ] Primary: Refactor SelectError() handling which adds an Echo dependency to platform level code.
- [ ] cli: Add pagination logic to functions which return lists (which requires a loop around those calls). 

## [0.0.1] - 2019-11-11
### Added
- Used [GORSK](https://github.com/ribice/gorsk) as starting template with modified project structure.
