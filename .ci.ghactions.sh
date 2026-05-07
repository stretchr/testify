#!/usr/bin/env bash

# MIT License
#
# Copyright (c) 2026 Olivier Mengué and contributors.
# 
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
# AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
# LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
# OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
# SOFTWARE.

#
# Verify that hashes of GitHub actions match the declared tag in attached comment.
#

set -euo pipefail

declare -A seen
status=0

for w in .github/workflows/*.yml
do
	sed -n -e '/uses: / s!^ *-\{0,1\} uses: \([^@]*\)@\([0-9a-f][0-9a-f]*\) *# *\(v.*\)$!\1 \2 \3!p' "$w" | while read -r action hash tag
	do
		if (( ${seen["$action-$hash-$tag"]:-0} )); then
			printf "\e[1;32m%s: %s@%s == %s\e[m\n" "$w" "$action" "$tag" "$hash"
			continue
		fi
		seen["$action-$hash-$tag"]=1

		if eval "$( curl -s -H "Accept: application/vnd.github+json" \
			"https://api.github.com/repos/$action/commits/$tag" | jq -r '.sha == "'"$hash"'"' )"
		then
			printf "\e[1;32m%s: %s@%s == %s\e[m\n" "$w" "$action" "$tag" "$hash"
		else
			printf "\e[1;31m%s: %s@%s != %s\e[m\n" "$w" "$action" "$tag" "$hash"
			status=1
		fi
	done
done

exit $status
