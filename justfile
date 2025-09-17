branch := `git rev-parse --abbrev-ref HEAD`
feature := replace_regex(branch, ".*/", "")
llm := "gpt-oss:latest"

[doc("Install in user bin folder")]
install:
    make install

[working-directory: './c64']
_make_c64: install
    make clean all

[working-directory: './c64']
test_koala: _make_c64
    open koala.prg

[working-directory: './c64']
test_hires: _make_c64
    open hires.prg

test_koala_png: install
    ./scripts/koala.sh ./scripts/paintface.src.png
    open ./scripts/paintface.src.png.koala.prg

test_hires_png: install
    ./scripts/hires.sh ./scripts/paintface.src.png
    open ./scripts/paintface.src.png.hires.prg

coverage:
    go test ./... -coverprofile=cover.out
    go tool cover -html=cover.out

[group("AI")]
make_commit_message:
    git diff --cached | ollama run {{llm}} --hidethinking 'Generate a short commit message that summarizes the most important changes. Make sure it starts with {{feature}} followed by a space, not a colon.' | pbcopy

