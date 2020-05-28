# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## Upcoming features [1.x.x]

- log deployment time
- refactor push - add rewind to all push things
- multi manifest pushed in parallel
- multi manifest works correct.
- check space quota before deployment (if user pass the option)
- set timeout how long the deployment will wait for more / free space

## [1.2.2] - 2020-04-30

### Fixed
- push application with routes that contains a path element
- vars file multiple placeholder replacement in one line

## [1.2.1] - 2020-05-24

### Added
- release action (github actions)

### Fixed
- find matching routes
- pss a vars file combined with 'legacy-push'

### Changed
- add more tests for filed issues
- add information about 'Specifying Routes' (contribution)

### Changed
- add more trace logging
- make tests more stable

## [1.2.0] - 2020-02-24

### Added
- --vars-file argument was available again

### Changed
- move all v2 operations into the v2 package
- delete cloud controller version check
- all routes will be switched lazy now (after starting up the new application)  
- dropped CF_PUPPETEER_TRACE and replace them with CF_TRACE (read the README "Known issues section" for using it)
- clean up code
- move temp manifest file generator to the manifest file

### Fixed
- error while passing environment variables on legacy push (--env)

## [1.1.3] - 2019-11-26 

### Changed
- add more logging
- add error message for manifest path declarative (when the path contains a wildcard)

### Fixed
- venerable application will be deleted correct
- typos, update error messages 

## [1.1.2] - 2019-11-21

### Changed

- refactor environment parsing
- clean code
- documentation update
- back port "--no-route" and "--route-only" option to "--legacy-push" (v2).
- docker support was available again
- vendor-option dropped
- options was sorted right now (go 1.12)

### Fixed

- environment parsing, now it's more stable
- error handling while uploading artifact

## [1.1.1] - 2019-06-30

### Changed

- vendor-option to venerable-action - deprecation message will be written if old argument was used
- no-route option set venerable-action to none as default

## [1.1.0] - 2019-06-29

### Added

- no-routes option added - ignore route switching, should be combined with vendor-option
- route-only add routes only - took routes from manifest and add them to the application (without vendor extension)
- no-start option for new deployed application

### Changed

- vendor-option argument now supports the options: "stop,delete,none"
- env argument changed to POSIX style. (--env)

### Fixed

- upload file with v3 api temp file path generation was wrong
- delete application and rename vendor app when upload fails

## [1.0.0.rc.0 - rc.2 - Final] - 2019-06-13

### Added

- new argument to show crash log before old application will be deleted
- push application with v3 api

### Changed

- cleanup code (refactor the complete plugin)
- colorized logging - use default cf terminal logging with color
- remove some parameters (like var, vars-file etc.)
- update ParseArgs to struct instead of multi return values
- code splitting / complete refactoring

### Fixed

- publishing applications
- fix some linter errors
- implement new manifest parser to understand new manifest format (https://docs.cloudfoundry.org/devguide/deploy-apps/manifest-attributes.html#buildpack)

### Features

- add argument to stop old service only instead of deletion
- choose between v2 and v3 push (v2 = legacy option)

## [0.0.14] - 2019-03-25

### Added

- health check settings (v3): fallback on app manifest for health-check-type and health-check-http-endpoint settings
- print more information on Cloud Foundry API calls (if env CF_PUPPETEER_TRACE = true)

### Fixed

- health check settings (v3): support for applications with empty command
- go linting issues

### Changed

- improve documentation

## [0.0.13] - 2019-03-08

### Added

- add argument to stop old service instance instead of deletion
- add argument to set health check invocation timeout (v3)

### Fixed

- dynamic env issue when passing key value pairs
- argument handling

### Changed

- cleanup manifest parser
- add better documentation

## [0.0.12] - 2019-02-21

### Changed

- Change buildscript to generate sha1 sums for each file

### Fixed

- Fix timeout issue when passing no -t argument but a specified timout at manifest (default if nothing was provided)

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
