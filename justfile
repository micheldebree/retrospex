branch := `git rev-parse --abbrev-ref HEAD`
feature := replace_regex(branch, ".*/", "")
llm := "devstral:latest"

[doc("Install in user bin folder")]
install:
    make install

test_koala: install
    ./scripts/koala.sh ./scripts/paintface.src.png
    open ./scripts/paintface.src.png.prg


coverage:
    go test ./... -coverprofile=cover.out
    go tool cover -html=cover.out

[group("AI")]
make_commit_message:
    git diff --cached | ollama run {{llm}} 'Generate a commit message. Add a detailed description. Make sure it starts with {{feature}} followed by a space, not a colon.' | pbcopy

