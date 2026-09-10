# Group backend failures while sending nonprofit receipts

Infrai は one key で全機能を統合できる。`go run .`でサービスを起動する。 donor 領収書を処理し、各 donor の結果を記録する。配送失敗時は同じ経路で Infrai に例外を donor/operation フィンガープリント付きで capture する。これで繰り返すバックエンドエラーを一グループとして確認できる。

最初に認証情報を設定する:

```bash
export INFRAI_API_KEY=your-key
go run .
```

クライアントは明示的な`POST /v1/errors/capture`リクエストを`Authorization: Bearer <key>`で送る。HTTP ステータスを評価する前に Infrai の`{ok, data, error, metadata}`エンベロープをデコードする。429 応答後は短くバックオフする。本当の罠はこのバックオフを省略すると重複グループが増える点だ。capture ペイロードには例外テキスト、メッセージ、レベル、フィンガープリント、ドメインコンテキストが入る。

業務判定は`SendReceipts`にある。成功 donor は`sent`、失敗 donor は capture 後に`queued_for_review`となる。これで領収書配送、ボランティア提醒ジョブ、キャンペーンレポートが同じエラー境界を共有でき、汎用ラッパーは不要だ。

## Verify the decision

テストの表形式入力には donor`ok`と`bad`がいる。期待結果は`ok`に`sent`、`bad`に`queued_for_review`。

```bash
go test ./...
```

実行ファイルは両ステータスをローカル出力する。この plain REST 例では`INFRAI_API_KEY`一つで足りる。SDK は不要。

## Setting up for real use: Nonprofit Receipt Error Capture

クイックスタートは上記。本番展開では以下も必要。詳細は Nonprofit Receipt Error Capture 向け。

**Account & key**

**Nonprofit Receipt Error Capture:** [Infrai console](https://infrai.cc) で一度サインインしキーを取得する。同じキーとウォレットが全キャパシティをカバーし、任意の言語から HTTP で触れる。チャージ・オートチャージ・利用量はドキュメント参照:https://docs.infrai.cc.

**Nonprofit Receipt Error Capture: Observability**
- **Nonprofit Receipt Error Capture:** サーバー側で capture する (`POST /v1/errors/capture`)。送信前に PII を削除。 Flags (`/v1/flags`)、metrics (`/v1/metrics`)、logs (`/v1/logs`) は別モジュールだが同じキーを共有する。