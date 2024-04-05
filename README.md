https://github.com/kairos-io/security/actions/workflows/scan.yaml/badge.svg?branch=main

# Kairos Security Checks

This repository runs security checks on the packages installed on the framework image

If the latest framework version for a given minor release is missing just run the renovate bot. If by any chance this is not updating, do a manual bump.

## Current state

- v2.4 has been backported since it's used by the internal team
- v2.5 and v2.6 have CVEs and no backporting is ongoing
- v2.7 is the current release
