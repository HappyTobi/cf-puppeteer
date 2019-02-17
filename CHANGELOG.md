# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [UNRELEASED]
### Added
- will look for new features bringing in new versions
- try to implement pull requests from original repository

## [0.0.11] - 2019-02-17
### Changed
- Rename plugin because of [CF-Plugin-Repo issue](https://github.com/cloudfoundry/cli-plugin-repo/pull/282#issuecomment-463328661)

## [0.0.10] - 2019-02-14
### Added
- Add new argument -t to specify a push timeout
- Add new manifest parser to get some informations out of the manifest instead of passing them through the cmd
- Add new feature to push application without appName if manifest was provided
- Fixtures for test

### Changed
- Changed version to 0.0.10
- Changed ParseArgs method

### Removed
- Unused code

## [0.0.9] - 2019-01-17
### Added
- Changelog.md file to privide a better overview about the releases and changed stuff.

### Changed
- Switched dependency management to [govendor](https://github.com/kardianos/govendor)
- Notice.md copyright
- Packagenames was renamed from contraband to happytobi
- build.yml to build the autopilot binary
- Update cf-cli dependency to newest version
- Changed version to 0.0.9

### Removed
- Remove go dep dependency management.


## [0.0.1 - 0.0.8] - 2015-12-31
- Original version from Contraband. [contraband](https://github.com/contraband/autopilot)
