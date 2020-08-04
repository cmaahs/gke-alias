#!/usr/bin/env bash

# Set sed command to gsed if running on mac
if [[ `uname` = "Darwin" ]]; then
    GSED=$(which gsed)
    if [[ "${GSED}" == "gsed not found" ]]; then
        echo "This script requires GNU-sed, please install with 'brew install gnu-sed'"
        exit 1
    fi
    sed_command="gsed"
    GDATE=$(which gdate)
    if [[ "${GDATE}" == "gdate not found" ]]; then
        echo "This script requires GNU-sed, please install with 'brew install gnu-sed'"
        exit 1
    fi
    date_command="gdate"

else
    sed_command="sed"
    date_command="date"
fi

declare semver="$1"

CURR_SEMVER=$(cat ../homebrew-admin-scripts/Formula/gke-alias.rb | grep version | awk '{print $2}' | jq -r)
CURR_REVISION=$(cat ../homebrew-admin-scripts/Formula/gke-alias.rb | grep revision | awk '{print $2}' | jq -r)
ADD_REV=1
if [[ "${CURR_SEMVER}" == "${semver}" ]]; then
  NEW_REVISION=$(expr ${CURR_REVISION} + ${ADD_REV})
  $sed_command -i 's|  revision.*|  revision '${NEW_REVISION}'|' ../homebrew-admin-scripts/Formula/gke-alias.rb
  semver+="_${NEW_REVISION}"
else
  $sed_command -i 's|  version.*|  version "'${semver}'"|' ../homebrew-admin-scripts/Formula/gke-alias.rb
  $sed_command -i 's|  revision.*|  revision 1|' ../homebrew-admin-scripts/Formula/gke-alias.rb
fi

BUILD_DATE=$($date_command --utc +%FT%T.%3NZ)
echo "Build Date: ${BUILD_DATE}"
echo "SemVer: ${semver}"
package_name="gke-alias"
# platforms=("darwin/386" "darwin/amd64" "freebsd/386" "freebsd/amd64" "linux/386" "linux/amd64" "netbsd/386" "netbsd/amd64" "openbsd/386" "openbsd/amd64" "solaris/amd64" "windows/386" "windows/amd64")
platforms=("darwin/amd64" "linux/amd64")
for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name="./bin/${GOOS}/${package_name}"
    destination_name="${GOOS}/${package_name}"
    echo "Building ${output_name}"

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name -ldflags "-X gke-alias/cmd.semVer=${semver} -X 'gke-alias/cmd.buildDate=${BUILD_DATE}' -X gke-alias/cmd.gitCommit=$(git rev-parse HEAD)"
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
    cp ${output_name} ../homebrew-admin-scripts/bin/${destination_name}
done
