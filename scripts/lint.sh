PKGS=("cmd" "internal")
ROOT=$(pwd)
CODE=0

for pkg in "${PKGS[@]}" ; do
  cd $ROOT
  echo -e "\n>>> Running linting on $pkg" ;
  folders=$(go list -f '{{.Dir}}' $ROOT/$pkg/... | grep -v -e "mock")
  for folder in $folders; do
    echo $folder
    GOGC=20 golangci-lint run "$folder/..."
    if [[ "$?" != "0" ]]; then
      CODE=1
    fi
  done
done

if [[ "$CODE" != 0 ]]; then
  exit 1
fi
