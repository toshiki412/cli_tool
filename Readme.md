- goのインストールをする

- モジュール作成
$ go mod init url
    $ go mod init github.com/toshiki412/cli_tool
go.modが作成される

cli_tool.goを作成
package main
import "fmt"
func main() {
	fmt.Println("Hello, World!")
}

$ go run .
Hello worldが出力

$ go mod tidy
依存関係を見に行って更新してくれる

- cobraのインストール
go install github.com/spf13/cobra-cli@latest

デフォルトでホームディレクトリの下にgoというフォルダができてその中に入る

パスを通していない場合は
$ ~/go/bin/cobra-cli init

通している場合は
$ cobra-cli init