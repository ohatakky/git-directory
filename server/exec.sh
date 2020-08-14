function run() {
    watchexec -r go run main.go
}

function pre() {
    go mod vendor &&
    cp -r ./ $GOPATH/src/github.com/ohatakky/git-directory/server
    rm -f $GOPATH/src/github.com/ohatakky/git-directory/server/go.sum
    rm -f $GOPATH/src/github.com/ohatakky/git-directory/server/go.mod
}

function deploy() {
    gcloud app deploy
}

$1
