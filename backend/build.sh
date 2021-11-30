#!/usr/bin/env bash

package=$1
if [[ -z "$package" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi
package_split=(${package//\// })
package_name=${package_split[-1]}

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/386" "linux/amd64" "linux/arm" "linux/arm64" "darwin/arm64" "freebsd/386" "freebsd/amd64" "freebsd/arm" "openbsd/386" "openbsd/amd64" "openbsd/arm" "netbsd/386" "netbsd/amd64" "netbsd/arm" "solaris/amd64" "dragonfly/amd64" "linux/ppc64" "linux/ppc64le" "linux/mips" "linux/mipsle" "linux/mips64" "linux/mips64le")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    output_name=$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    env GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name $package
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done
