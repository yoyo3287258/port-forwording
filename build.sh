#!/usr/bin/env bash
appName="port-forwarding"

ldflags="\
-w -s \
"
BuildDev() {
  rm -rf .git/
  mkdir -p "dist"
  muslflags="--extldflags '-static -fpic' $ldflags"
  BASE="https://musl.nn.ci/"
  FILES=(x86_64-linux-musl-cross aarch64-linux-musl-cross)
  for i in "${FILES[@]}"; do
    url="${BASE}${i}.tgz"
    curl -L -o "${i}.tgz" "${url}"
    sudo tar xf "${i}.tgz" --strip-components 1 -C /usr/local
  done
  OS_ARCHES=(linux-musl-amd64 linux-musl-arm64)
  CGO_ARGS=(x86_64-linux-musl-gcc aarch64-linux-musl-gcc)
  for i in "${!OS_ARCHES[@]}"; do
    os_arch=${OS_ARCHES[$i]}
    cgo_cc=${CGO_ARGS[$i]}
    echo building for ${os_arch}
    export GOOS=${os_arch%%-*}
    export GOARCH=${os_arch##*-}
    export CC=${cgo_cc}
    export CGO_ENABLED=1
    go build -o ./dist/$appName-$os_arch -ldflags="$muslflags" -tags=jsoniter .
  done
  xgo -targets=windows/amd64,darwin/amd64,darwin/arm64 -out "$appName" -ldflags="$ldflags" -tags=jsoniter .
  mv port-forwarding-* dist
  cd dist
  cp ./port-forwarding-windows-amd64.exe ./port-forwarding-windows-amd64-upx.exe
  upx -9 ./port-forwarding-windows-amd64-upx.exe
  find . -type f -print0 | xargs -0 md5sum >md5.txt
  cat md5.txt
}

if [ "$1" = "dev" ]; then
  BuildDev
elif [ "$1" = "release" ]; then
  if [ "$2" = "linux_musl_arm" ]; then
    BuildReleaseLinuxMuslArm
    MakeRelease "md5-linux-musl-arm.txt"
  elif [ "$2" = "linux_musl" ]; then
    BuildReleaseLinuxMusl
    MakeRelease "md5-linux-musl.txt"
  else
    BuildRelease
    MakeRelease "md5.txt"
  fi
else
  echo -e "Parameter error"
fi