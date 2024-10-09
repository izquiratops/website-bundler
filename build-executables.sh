version=$1

platforms=(
    "darwin/amd64"
    "darwin/arm64"
    "linux/amd64"
    "linux/arm"
    "linux/arm64"
    "windows/amd64"
)

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}

    output_name="dobunezumi-${version}-${GOOS}-${GOARCH}"
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi

    echo "Building release/$output_name..."
    env GOOS=$GOOS GOARCH=$GOARCH go build \
      -ldflags "-X github.com/izquiratops/dobunezumi/commands.Version=$version" \
      -o release/$output_name

    if [ $? -ne 0 ]; then
        echo 'An error has occurred during the build process.'
        exit 1
    fi

    zip_name="dobunezumi-${version}-${GOOS}-${GOARCH}"
    # Change to the release directory and hide the output
    pushd release > /dev/null
    if [ $GOOS = "windows" ]; then
        zip -r $zip_name.zip $output_name
        rm $output_name
    else
        chmod a+x $output_name
        gzip $output_name
    fi

    # Change back to the previous directory and hide the output
    popd > /dev/null
done
