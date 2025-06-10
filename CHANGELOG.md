# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
-

### Changed
-

### Deprecated
-

### Removed
-

### Fixed
-

### Security
-

## [0.1.0] - 2023-10-27

### Added
- Initial project setup with Fiber and Firebase Firestore.
- User registration and login endpoints (`/api/register`, `/api/login`).
- JWT authentication middleware.
- Basic user handlers (`/api/users`): Get All, Get by ID, Update.
- Basic task handlers (`/api/tasks`): Get All, Create, Get by ID, Update.
- Firestore configuration and client initialization.
- User and Task models.
- `.gitignore` file to exclude sensitive files and build artifacts.

### Fixed
- Corrected error handling for "document not found" in `GetTask` and `DeleteTask` using string comparison workaround.
- Implemented `DeleteUser` handler with "document not found" error handling.