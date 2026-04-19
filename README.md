# tiny-deck

tinygo-keebで作成したzero-kb02を使って
Xxx-deckみたいにMacをちょっとだけ便利にします。

マイコン：waveshare-rp2040-zero

## 実装済み
- CPU使用率を表示する
- メモリ使用率を表示する

## セットアップ

### ファームウェア書き込み

```sh
mise run flash
```

### ホスト常駐デーモンの自動起動 (macOS)

デバイスを USB 接続したタイミングで `host_daemon` を起動し、切断で終了する LaunchAgent を登録する。

```sh
mise run install
```

- バイナリ: `~/.local/bin/tiny-deck-host`
- plist: `~/Library/LaunchAgents/com.github.rin2yh.tiny-deck.plist`
- ログ: `~/Library/Logs/tiny-deck.log`

通知監視のため macOS の Full Disk Access を `tiny-deck-host` に許可する必要がある
(システム設定 → プライバシーとセキュリティ → フルディスクアクセス)。

アンインストール:

```sh
mise run uninstall
```

## 展望

- 使用率のやばさ具合でキーボードを光らせる
- ロータリーエンコーダーで音量を調節する
- joystickでカーソルを動かす
- streamdeckみたいに使えるようにする？
- 10キーにする？

