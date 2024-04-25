refact:
	go mod edit -module github.com/noolingo/api-gw
	-- rename all imported module
	find . -type f -name '*.go' \
  	-exec sed -i -e 's,github.com/MelnikovNA/noolingo-api-gw,github.com/noolingo/api-gw,g' {} \;