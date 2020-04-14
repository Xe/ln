# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.9.0]

### Changed
- Default logging format is now JSON, to get logfmt output use `LN_FORMATTER=text`

## [0.8.0]

### Added
- Add Fmt function for adding printf style messages

## [0.7.0]

### Added
- Write logs to syslog if `LN_OUT=<syslog>` is set in the environment

### Fixed
- Fix logfmt output for arrays and maps

