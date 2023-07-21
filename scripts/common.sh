#!/usr/bin/env bash


# This will canonicalize the path
SCRIPT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")"/.. && pwd -P)

source "${SCRIPT_ROOT}/scripts/lib/init.sh"
source "${SCRIPT_ROOT}/scripts/install/environment.sh"

# Here we map the output directories across both the local and remote _output
# directories:
#
# *_OUTPUT_ROOT    - the base of all output in that environment.
# *_OUTPUT_SUBPATH - location where golang stuff is built/cached.  Also
#                    persisted across docker runs with a volume mount.
# *_OUTPUT_BINPATH - location where final binaries are placed.  If the remote
#                    is really remote, this is the stuff that has to be copied
#                    back.
# OUT_DIR can come in from the Makefile, so honor it.
readonly LOCAL_OUTPUT_ROOT="${SCRIPT_ROOT}/${OUT_DIR:-_output}"
readonly LOCAL_OUTPUT_SUBPATH="${LOCAL_OUTPUT_ROOT}/platforms"
readonly LOCAL_OUTPUT_BINPATH="${LOCAL_OUTPUT_SUBPATH}"
readonly LOCAL_OUTPUT_GOPATH="${LOCAL_OUTPUT_SUBPATH}/go"
readonly LOCAL_OUTPUT_IMAGE_STAGING="${LOCAL_OUTPUT_ROOT}/images"

