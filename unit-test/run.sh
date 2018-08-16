#!/bin/bash -e
#
# Copyright IBM Corp. All Rights Reserved.
# Copyright DSiSc Group. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#   
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

# TODO: regexes for packages to exclude from unit test
excluded_packages=(
)

# TODO: regexes for packages that must be run serially
serial_packages=(
log
)

# packages which need to be tested with build tag pluginsenabled
plugin_packages=(
)

# obtain packages changed since some git refspec
packages_diff() {
    git -C "${GOPATH}/src/github.com/DSiSc/apigateway" diff --no-commit-id --name-only -r "${1:-HEAD}" |
        grep '.go$' | grep -Ev '^vendor/|^build/' | \
        sed 's%/[^/]*$%/%' | sort -u | \
        awk '{print "github.com/DSiSc/apigateway/"$1"..."}'
}

# "go list" packages and filter out excluded packages
list_and_filter() {
    go list $@ 2>/dev/null | grep -Ev $(local IFS='|' ; echo "${excluded_packages[*]}") || true
}

# remove packages that must be tested serially
parallel_test_packages() {
    echo "$@" | grep -Ev $(local IFS='|' ; echo "${serial_packages[*]}") || true
}

# get packages that must be tested serially
serial_test_packages() {
    echo "$@" | grep -E $(local IFS='|' ; echo "${serial_packages[*]}") || true
}

# "go test" the provided packages. Packages that are not prsent in the serial package list
# will be tested in parallel
run_tests() {
    echo ${GO_TAGS}
    flags="-cover"
    if [ -n "${VERBOSE}" ]; then
      flags="-v -cover"
    fi

    local parallel=$(parallel_test_packages "$@")
    if [ -n "${parallel}" ]; then
        time go test ${flags} -tags "$GO_TAGS" -ldflags "$GO_LDFLAGS" ${parallel[@]} -short -timeout=20m
    fi

    local serial=$(serial_test_packages "$@")
    if [ -n "${serial}" ]; then
        time go test ${flags} -tags "$GO_TAGS" -ldflags "$GO_LDFLAGS" ${serial[@]} -short -p 1 -timeout=20m
    fi
}

# "go test" the provided packages and generate code coverage reports. Until go 1.10 is released,
# profile reports can only be generated one package at a time.
run_tests_with_coverage() {
    # Initialize profile.cov
    for pkg in $@; do
        :> profile_tmp.cov
        go test -cover -coverprofile=profile_tmp.cov -tags "$GO_TAGS" -ldflags "$GO_LDFLAGS" $pkg -timeout=20m
        tail -n +2 profile_tmp.cov >> profile.cov || echo "Unable to append coverage for $pkg"
    done

    # convert to cobertura format
    gocov convert profile.cov | gocov-xml > report.xml
}

main() {
    # default behavior is to run all tests
    local package_spec=${TEST_PKGS:-github.com/DSiSc/apigateway/...}

    # when running a "verify" job, only test packages that have changed
    if [ "${JOB_TYPE}" = "VERIFY" ]; then
        # first check for uncommitted changes
        package_spec=$(packages_diff HEAD)
        if [ -z "${package_spec}" ]; then
            # next check for changes in the latest commit - typically this will
            # be for CI only, but could also handle a committed change before
            # pushing to Gerrit
            package_spec=$(packages_diff HEAD^)
        fi
    fi

    # expand the package spec into an array of packages
    local -a packages=$(list_and_filter ${package_spec})

    if [ -z "${packages}" ]; then
        echo "Nothing to test!!!"
    elif [ "${JOB_TYPE}" = "PROFILE" ]; then
        echo "mode: set" > profile.cov
        run_tests_with_coverage "${packages[@]}"
        GO_TAGS="${GO_TAGS} pluginsenabled" run_tests_with_coverage "${plugin_packages[@]}"
    else
        run_tests "${packages[@]}"
        GO_TAGS="${GO_TAGS} pluginsenabled" run_tests "${plugin_packages[@]}"
    fi
}

main
