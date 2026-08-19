#!/bin/sh
set -eu
exec gh pr-quality-gate audit --repo "${1:-.}"
