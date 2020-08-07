function run() {
    watchexec -r go run main.go
}

function deploy() {
    gcloud app deploy
}

$1
