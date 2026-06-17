.PHONY: fmt lint test verify

fmt:
ifdef ZED_FORMAT_FILE
	cat > "$(ZED_FORMAT_FILE)"
	dune fmt >/dev/null
	cat "$(ZED_FORMAT_FILE)"
else
	dune fmt
endif

lint:
	dune build @all

test:
	dune runtest

verify: fmt lint test
	git diff --exit-code
