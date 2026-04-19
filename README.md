# tiny-deck

tinygo-keebで作成したzero-kb02を使って
Xxx-deckみたいにMacをちょっとだけ便利にします。

マイコン：waveshare-rp2040-zero

## 実装済み
- [CPU・メモリ使用率をOLEDに表示する](https://github.com/rin2yh/tiny-deck/blob/1c83f4b4767cc1c7cc81922803662c129790b809/internal/display/metrics.go)
- [使用率に応じてLEDマトリクスにバー表示する](https://github.com/rin2yh/tiny-deck/blob/1c83f4b4767cc1c7cc81922803662c129790b809/internal/keyboard/metrics.go)
- [ロータリーエンコーダーで音量を調節する](https://github.com/rin2yh/tiny-deck/blob/1c83f4b4767cc1c7cc81922803662c129790b809/internal/keyboard/encoder.go)
- [joystickでマウスカーソルを動かす](https://github.com/rin2yh/tiny-deck/blob/1c83f4b4767cc1c7cc81922803662c129790b809/internal/joystick/joystick.go)
- [キー入力・長押しでレイヤーを切り替える](https://github.com/rin2yh/tiny-deck/blob/1c83f4b4767cc1c7cc81922803662c129790b809/internal/keyboard/layer.go)
- [macOSの通知を検知してLEDを光らせる](https://github.com/rin2yh/tiny-deck/blob/1c83f4b4767cc1c7cc81922803662c129790b809/internal/keyboard/animation.go)
- [USB接続で`host_daemon`を自動起動する](https://github.com/rin2yh/tiny-deck/blob/1c83f4b4767cc1c7cc81922803662c129790b809/scripts/install.sh)

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

- 10キーにする

