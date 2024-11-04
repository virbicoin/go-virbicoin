# Go VirBiCoin (Gvbc)

本リポジトリはGo EthereumからフォークされたVirBiCoinプロジェクトのVirBiCoin Chainプロトコル公式Golang実装です。

[![API Reference](
https://camo.githubusercontent.com/915b7be44ada53c290eb157634330494ebe3e30a/68747470733a2f2f676f646f632e6f72672f6769746875622e636f6d2f676f6c616e672f6764646f3f7374617475732e737667
)](https://pkg.go.dev/github.com/virbicoin/go-virbicoin?tab=doc)
[![Go Report Card](https://goreportcard.com/badge/github.com/virbicoin/go-virbicoin)](https://goreportcard.com/report/github.com/virbicoin/go-virbicoin)
[![Travis](https://travis-ci.org/virbicoin/go-virbicoin.svg?branch=master)](https://travis-ci.org/virbicoin/go-virbicoin)
[![Discord](https://img.shields.io/badge/discord-join%20chat-blue.svg)](https://discord.gg/nthXNEv)

安定版リリースと最新のマスターブランチの自動ビルドが利用可能です。
バイナリアーカイブはhttps://github.com/virbicoin/go-virbicoin/releases で公開されています。

## ソースのビルド

前提条件と詳細なビルド手順については、wikiの[インストール手順](https://github.com/virbicoin/go-virbicoin/wiki/Installing-Gvbc)を参照してください。

`gvbc`のビルドには、Go（バージョン1.13以降）とCコンパイラの両方が必要です。これらはお好みのパッケージマネージャを使用してインストールできます。依存関係がインストールされたら、以下のコマンドを実行します。

```shell
make gvbc
```

または、ユーティリティのフルスイートをビルドするには：

```shell
make all
```

## 実行可能ファイル

go-virbicoinプロジェクトには、`cmd`ディレクトリにいくつかのラッパー/実行可能ファイルが含まれています。

|    コマンド    | 説明                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                                          |
| :-----------: | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
|  **`gvbc`**   | VirBiCoin のメインのクライアントCLIです。これは、VirBiCoin Chainネットワーク（メインネット、テストネット、プライベートネット）へのエントリーポイントであり、フルノード（デフォルト）、アーカイブノード（すべての履歴状態を保持）、またはライトノード（データをライブで取得）として動作できます。HTTP、WebSocket、および/またはIPCトランスポートの上に公開されたJSON RPCエンドポイントを介して、他のプロセスによってVirBiCoin Chainネットワークへのゲートウェイとして使用できます。コマンドラインオプションについては、`gvbc --help`および[CLI Wiki page](https://github.com/ethereum/go-ethereum/wiki/Command-Line-Options)を参照してください。          |
|   `abigen`    | VirBiCoin Chainコントラクト定義を使いやすく、コンパイル時に型安全なGoパッケージに変換するソースコードジェネレータです。これは、プレーンな[Ethereum contract ABIs](https://github.com/ethereum/wiki/wiki/Ethereum-Contract-ABI)で動作し、コントラクトバイトコードが利用可能な場合に拡張機能を提供します。ただし、Solidityソースファイルも受け入れるため、開発が非常にスムーズになります。詳細については、[Native DApps](https://github.com/ethereum/go-ethereum/wiki/Native-DApps:-Go-bindings-to-Ethereum-contracts) wikiページを参照してください。 |
|  `bootnode`   | ネットワークノード発見プロトコルにのみ参加し、上位レベルのアプリケーションプロトコルを実行しない、VirBiCoin Chainクライアント実装のストリップダウンバージョンです。プライベートネットワークでピアを見つけるのを支援する軽量ブートストラップノードとして使用できます。                                                                                                                                                                                                                                                                 |
|     `evm`     | EVM（Ethereum Virtual Machine）の開発者ユーティリティバージョンで、構成可能な環境と実行モード内でバイトコードスニペットを実行できます。その目的は、EVMオペコードの分離された詳細なデバッグを可能にすることです（例：`evm --code 60ff60ff --debug run`）。                                                                                                                                                                                                                                                                     |
| `gethrpctest` | [ethereum/rpc-test](https://github.com/ethereum/rpc-tests)テストスイートをサポートする開発者ユーティリティツールで、[Ethereum JSON RPC](https://github.com/ethereum/wiki/wiki/JSON-RPC)仕様へのベースライン準拠を検証します。詳細については、[test suite's readme](https://github.com/ethereum/rpc-tests/blob/master/README.md)を参照してください。                                                                                                                                                                                                     |
|   `rlpdump`   | バイナリRLP ([Recursive Length Prefix](https://github.com/ethereum/wiki/wiki/RLP))ダンプ（VirBiCoin Chainプロトコルによってネットワークおよびコンセンサスの両方で使用されるデータエンコーディング）をユーザーフレンドリーな階層表現に変換する開発者ユーティリティツールです（例：`rlpdump --hex CE0183FFFFFFC4C304050583616263`）。                                                                                                                                                                                                                                 |
|   `puppeth`   | 新しいVirBiCoin Chainネットワークの作成を支援するCLIウィザードです。                                                                                                                                                                                                                                                                                                                                                                                                                                                                                           |

## `gvbc`の実行

すべてのコマンドラインフラグを説明するのはここでは範囲外ですが（[CLI Wiki page](https://github.com/ethereum/go-ethereum/wiki/Command-Line-Options)を参照してください）、いくつかの一般的なパラメータの組み合わせを列挙して、`gvbc`インスタンスを迅速に実行する方法を説明します。

### メインVirBiCoin Chainネットワーク上のフルノード

最も一般的なシナリオは、単にVirBiCoin Chainネットワークと対話したい人々です：アカウントを作成する；資金を転送する；コントラクトをデプロイして対話する。この特定のユースケースでは、ユーザーは数年前の履歴データを気にしないため、ネットワークの現在の状態に迅速に同期できます。これを行うには：

```shell
$ gvbc console
```

このコマンドは次のことを行います：
 * `gvbc`を高速同期モード（デフォルト、`--syncmode`フラグで変更可能）で起動し、VirBiCoin Chainネットワークの全履歴を処理するのではなく、より多くのデータをダウンロードしてCPU負荷を軽減します。
 * `gvbc`の組み込みインタラクティブ[JavaScript console](https://github.com/ethereum/go-ethereum/wiki/JavaScript-Console)を起動し（末尾の`console`サブコマンドを介して）、すべての公式[`web3`メソッド](https://github.com/username/go-virbicoin/wiki/JavaScript-API)および`gvbc`独自の[management APIs](https://github.com/ethereum/go-ethereum/wiki/Management-APIs)を呼び出すことができます。このツールはオプションであり、省略した場合でも、`gvbc attach`を使用して既に実行中の`gvbc`インスタンスに接続できます。

### Görliテストネットワーク上のフルノード

開発者向けに、VirBiCoin Chainコントラクトの作成を試してみたい場合、システム全体に慣れるまで実際のお金を使わずに行いたいでしょう。つまり、メインネットワークに接続する代わりに、ノードを**テスト**ネットワークに参加させたいのです。これはメインネットワークと完全に同等ですが、プレイEtherのみを使用します。

```shell
$ gvbc --goerli console
```

`console`サブコマンドの意味は上記と同じで、テストネットでも同様に有用です。上記の説明をスキップした場合は、再度参照してください。

`--goerli`フラグを指定すると、`gvbc`インスタンスが少し再構成されます：

 * メインVirBiCoin Chainネットワークに接続する代わりに、クライアントはGörliテストネットワークに接続します。これには異なるP2Pブートノード、異なるネットワークID、およびジェネシス状態が使用されます。
 * デフォルトのデータディレクトリ（例：Linuxでは`~/.virbicoin`）の代わりに、`gvbc`は1段階深い`goerli`サブフォルダ（例：Linuxでは`~/.virbicoin/goerli`）に配置されます。注意：OSXおよびLinuxでは、テストネットノードに接続するにはカスタムエンドポイントを使用する必要があります。`gvbc attach`はデフォルトでプロダクションノードエンドポイントに接続しようとします。例：`gvbc attach <datadir>/goerli/gvbc.ipc`。Windowsユーザーは影響を受けません。

*注意：メインネットワークとテストネットワーク間でトランザクションが交差しないようにする内部保護措置がありますが、プレイマネーとリアルマネーのために常に別々のアカウントを使用するようにしてください。アカウントを手動で移動しない限り、`gvbc`はデフォルトで2つのネットワークを正しく分離し、アカウントを相互に利用できないようにします。*

### Rinkebyテストネットワーク上のフルノード

Go VirBiCoin Chainは、コミュニティメンバーによって運営されている古いプルーフ・オブ・オーソリティベースのテストネットワーク[*Rinkeby*](https://www.rinkeby.io)への接続もサポートしています。

```shell
$ gvbc --rinkeby console
```

### Ropstenテストネットワーク上のフルノード

GörliおよびRinkebyに加えて、Gvbcは古いRopstenテストネットもサポートしています。RopstenテストネットワークはEthashプルーフ・オブ・ワークコンセンサスアルゴリズムに基づいています。そのため、追加のオーバーヘッドがあり、ネットワークの低難易度/セキュリティのために再編成攻撃に対してより脆弱です。

```shell
$ gvbc --ropsten console
```

*注意：古いGvbc構成は、Ropstenデータベースを`testnet`サブディレクトリに保存します。*

### 設定

`gvbc`バイナリに多数のフラグを渡す代わりに、設定ファイルを渡すこともできます：

```shell
$ gvbc --config /path/to/your_config.toml
```

ファイルの見本を得るには、`dumpconfig`サブコマンドを使用して既存の設定をエクスポートできます：

```shell
$ gvbc --your-favourite-flags dumpconfig
```

*注意：これは`gvbc` v1.6.0以降でのみ動作します。*

#### Dockerクイックスタート

VirBiCoin Chainをマシンで迅速に稼働させる最も簡単な方法の1つは、Dockerを使用することです：

```shell
docker run -d --name VirBiCoin-Chain-node -v /Users/alice/VirBiCoin-Chain:/root \
           -p 8329:8329 -p 28329:28329 \
           virbicoin/go-virbicoin
```

これにより、上記のコマンドと同様に、DBメモリ許容量1GBで`gvbc`が高速同期モードで起動します。また、ブロックチェーンを保存するための永続ボリュームがホームディレクトリに作成され、デフォルトのポートがマッピングされます。スリムバージョンのイメージ用に`alpine`タグも利用可能です。

他のコンテナやホストからRPCにアクセスしたい場合は、`--http.addr 0.0.0.0`を忘れないでください。デフォルトでは、`gvbc`はローカルインターフェースにバインドされ、RPCエンドポイントは外部からアクセスできません。

### プログラムによる`gvbc`ノードとのインターフェース

開発者として、手動でコンソールを介してではなく、独自のプログラムを介して`gvbc`およびVirBiCoin Chainネットワークと対話したくなるでしょう。これを支援するために、`gvbc`にはJSON-RPCベースのAPI（[標準API](https://github.com/virbicoin/go-virbicoin/wiki/JSON-RPC)および`gvbc`固有の[management APIs](https://github.com/ethereum/go-ethereum/wiki/Management-APIs)）のサポートが組み込まれています。これらはHTTP、WebSockets、およびIPC（UNIXベースのプラットフォームではUNIXソケット、Windowsでは名前付きパイプ）を介して公開できます。

IPCインターフェースはデフォルトで有効になっており、`gvbc`がサポートするすべてのAPIを公開しますが、HTTPおよびWSインターフェースは手動で有効にする必要があり、セキュリティ上の理由からAPIのサブセットのみを公開します。これらは期待通りにオン/オフおよび構成できます。

HTTPベースのJSON-RPC APIオプション：

  * `--http` HTTP-RPCサーバーを有効にする
  * `--http.addr` HTTP-RPCサーバーのリスニングインターフェース（デフォルト：`localhost`）
  * `--http.port` HTTP-RPCサーバーのリスニングポート（デフォルト：`8329`）
  * `--http.api` HTTP-RPCインターフェースで提供されるAPI（デフォルト：`eth,net,web3`）
  * `--http.corsdomain` クロスオリジンリクエストを受け入れるドメインのカンマ区切りリスト（ブラウザで強制）
  * `--ws` WS-RPCサーバーを有効にする
  * `--ws.addr` WS-RPCサーバーのリスニングインターフェース（デフォルト：`localhost`）
  * `--ws.port` WS-RPCサーバーのリスニングポート（デフォルト：`8330`）
  * `--ws.api` WS-RPCインターフェースで提供されるAPI（デフォルト：`eth,net,web3`）
  * `--ws.origins` WebSocketsリクエストを受け入れるオリジン
  * `--ipcdisable` IPC-RPCサーバーを無効にする
  * `--ipcapi` IPC-RPCインターフェースで提供されるAPI（デフォルト：`admin,debug,eth,miner,net,personal,shh,txpool,web3`）
  * `--ipcpath` データディレクトリ内のIPCソケット/パイプのファイル名（明示的なパスはエスケープされます）

独自のプログラミング環境の機能（ライブラリ、ツールなど）を使用して、上記のフラグで構成された`gvbc`ノードにHTTP、WS、またはIPCを介して接続し、すべてのトランスポートで[JSON-RPC](https://www.jsonrpc.org/specification)を話す必要があります。複数のリクエストに同じ接続を再利用できます！

**注意：HTTP/WSベースのトランスポートを開く前に、そのセキュリティ上の影響を理解してください！インターネット上のハッカーは、公開されたAPIを持つVirBiCoin Chainノードを積極的にサブバートしようとしています！さらに、すべてのブラウザタブはローカルで実行されているWebサーバーにアクセスできるため、悪意のあるWebページがローカルで利用可能なAPIをサブバートしようとする可能性があります！**

### プライベートネットワークの運用

独自のプライベートネットワークを維持することは、公式ネットワークで当然とされる多くの設定を手動で行う必要があるため、より複雑です。

#### プライベートジェネシス状態の定義

まず、ネットワークのジェネシス状態を作成する必要があります。これは、すべてのノードが認識し、合意する必要がある小さなJSONファイル（例：`genesis.json`と呼びます）で構成されます：

```json
{
  "config": {
    "chainId": <arbitrary positive integer>,
    "homesteadBlock": 0,
    "eip150Block": 0,
    "eip155Block": 0,
    "eip158Block": 0,
    "byzantiumBlock": 0,
    "constantinopleBlock": 0,
    "petersburgBlock": 0,
    "istanbulBlock": 0
  },
  "alloc": {},
  "coinbase": "0x0000000000000000000000000000000000000000",
  "difficulty": "0x20000",
  "extraData": "",
  "gasLimit": "0x2fefd8",
  "nonce": "0x0000000000000042",
  "mixhash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "parentHash": "0x0000000000000000000000000000000000000000000000000000000000000000",
  "timestamp": "0x00"
}
```

上記のフィールドはほとんどの目的に適していますが、`nonce`をランダムな値に変更して、未知のリモートノードが接続できないようにすることをお勧めします。テストを容易にするためにいくつかのアカウントに事前に資金を提供したい場合は、アカウントを作成し、そのアドレスで`alloc`フィールドを埋めます。

```json
"alloc": {
  "0x0000000000000000000000000000000000000001": {
    "balance": "111111111"
  },
  "0x0000000000000000000000000000000000000002": {
    "balance": "222222222"
  }
}
```

上記のJSONファイルでジェネシス状態を定義したら、すべての`gvbc`ノードを起動する前にそれで初期化して、すべてのブロックチェーンパラメータが正しく設定されていることを確認する必要があります：

```shell
$ gvbc init path/to/genesis.json
```

#### ランデブーポイントの作成

希望するジェネシス状態に初期化されたすべてのノードを持って、他のノードがネットワーク内および/またはインターネット上でお互いを見つけるために使用できるブートストラップノードを起動する必要があります。クリーンな方法は、専用のブートノードを構成して実行することです：

```shell
$ bootnode --genkey=boot.key
$ bootnode --nodekey=boot.key
```

ブートノードがオンラインになると、他のノードが接続してピア情報を交換するために使用できる[`enode` URL](https://github.com/virbicoin/go-virbicoin/wiki/enode-url-format)が表示されます。表示されたIPアドレス情報（おそらく`[::]`）を外部からアクセス可能なIPに置き換えて、実際の`enode` URLを取得してください。

*注意：フル機能の`gvbc`ノードをブートノードとして使用することもできますが、これはあまり推奨されません。*

#### メンバーノードの起動

ブートノードが稼働し、外部からアクセス可能であることを確認したら（`telnet <ip> <port>`を試して実際にアクセス可能であることを確認できます）、ブートノードを指す`--bootnodes`フラグを使用して、すべての後続の`gvbc`ノードを起動します。プライベートネットワークのデータディレクトリを分離しておくことも望ましいでしょうので、カスタム`--datadir`フラグも指定してください。

```shell
$ gvbc --datadir=path/to/custom/data/folder --bootnodes=<bootnode-enode-url-from-above>
```

*注意：ネットワークがメインネットおよびテストネットから完全に切り離されるため、トランザクションを処理し、新しいブロックを作成するためにマイナーを構成する必要もあります。*

#### プライベートマイナーの実行

パブリックVirBiCoin Chainネットワークでのマイニングは、GPUを使用する必要があるため複雑な作業です。OpenCLまたはCUDA対応の`ethminer`インスタンスが必要です。このようなセットアップに関する情報は、[EtherMining subreddit](https://www.reddit.com/r/EtherMining/)および[ethminer](https://github.com/ethereum-mining/ethminer)リポジトリを参照してください。

ただし、プライベートネットワーク設定では、単一のCPUマイナーインスタンスで実用目的には十分です。これは、重いリソースを必要とせずに、正しい間隔で安定したブロックストリームを生成できます（複数のスレッドは必要なく、単一のスレッドで実行することを検討してください）。マイニング用に`gvbc`インスタンスを起動するには、通常のフラグに加えて以下を実行します：

```shell
$ gvbc <usual-flags> --mine --miner.threads=1 --etherbase=0x0000000000000000000000000000000000000000
```

これにより、単一のCPUスレッドでブロックとトランザクションのマイニングが開始され、すべての収益が`--etherbase`で指定されたアカウントにクレジットされます。デフォルトのガスリミットブロックの収束先（`--targetgaslimit`）およびトランザクションが受け入れられる価格（`--gasprice`）を変更することで、マイニングをさらに調整できます。

## 貢献

ソースコードの支援を検討していただきありがとうございます！インターネット上の誰からの貢献も歓迎し、最小の修正でも感謝します！

go-virbicoinに貢献したい場合は、フォークして修正し、コミットしてプルリクエストを送信し、メンテナがレビューしてメインコードベースにマージします。より複雑な変更を提出したい場合は、[gitterチャンネル](https://gitter.im/virbicoin/go-virbicoin)でコアデベロッパーに確認して、これらの変更がプロジェクトの一般的な哲学に沿っているかどうかを確認し、早期のフィードバックを得ることで、あなたの努力を軽減し、レビューとマージの手続きを迅速かつ簡単にすることができます。

貢献が以下のコーディングガイドラインに従っていることを確認してください：

 * コードは公式のGo[フォーマット](https://golang.org/doc/effective_go.html#formatting)ガイドラインに従う必要があります（つまり、[gofmt](https://golang.org/cmd/gofmt/)を使用します）。
 * コードは公式のGo[コメント](https://golang.org/doc/effective_go.html#commentary)ガイドラインに従って文書化される必要があります。
 * プルリクエストは`master`ブランチに基づいて開かれる必要があります。
 * コミットメッセージは変更するパッケージにプレフィックスを付ける必要があります。
   * 例："eth, rpc: make trace configs optional"

環境の設定、プロジェクト依存関係の管理、およびテスト手順の詳細については、[Developers' Guide](https://github.com/ethereum/go-ethereum/wiki/Developers'-Guide)を参照してください。

## ライセンス

go-virbicoinライブラリ（`cmd`ディレクトリ外のすべてのコード）は、[GNU Lesser General Public License v3.0](https://www.gnu.org/licenses/lgpl-3.0.en.html)の下でライセンスされています。リポジトリ内の`COPYING.LESSER`ファイルにも含まれています。

go-virbicoinバイナリ（`cmd`ディレクトリ内のすべてのコード）は、[GNU General Public License v3.0](https://www.gnu.org/licenses/gpl-3.0.en.html)の下でライセンスされています。リポジトリ内の`COPYING`ファイルにも含まれています。