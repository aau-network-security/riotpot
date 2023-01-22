# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

[^docker-image] : Once we decide how, who and where to publish it, the honeypot will be available as a Docker image. For now, the image can be built from source using the `Dockerfile` included in the **`build > docker`** folder (there is also a `docker-compose` file ready to use).

## [v0.1.2] 2023-01-22

### Added

- Embeds the UI in the binary
  
### Changed

- All the plugins now use a `localhost:<port>` address to listen for incoming connections (i.e., even if exposed, they should not provide any service).
- Optimised the Docker image. Now the UI folder only includes the required files, instead of the whole directory contents.
  
    > **Note:** The Docker image has not been published yet. This will come soon [^docker-image]

- Changed flows for the pipeline
  - Now the UI uses [Vite](https://vitejs.dev/) to build the project
  - Using [Goreleaser](https://goreleaser.com/) - again - to make releases
