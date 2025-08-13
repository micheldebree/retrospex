[doc("Install in user bin folder")]
install:
    make install

test_koala: install
    ./scripts/koala.sh ./scripts/paintface.src.png
    open ./scripts/paintface.src.png.prg


coverage:
    go test ./... -coverprofile=cover.out
    go tool cover -html=cover.out
