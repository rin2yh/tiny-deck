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
- [Gameレイヤーで選んだゲームをOLEDで遊ぶ](internal/arcade/launcher.go)

## 展望
- 10キーにする

## Gameレイヤー

レイヤーキー（`KeyC3R2`）の長押しは `NumPad → Temp → Game → NumPad` と巡回する。
`Game` に入るとゲーム選択画面になり、OLEDはゲームが占有する。

| 場面 | 操作 |
| --- | --- |
| 選択画面 | エンコーダー回転でカーソル移動、押し込みで開始 |
| プレイ中 | レイヤーキー以外のどのキーでもジャンプ、エンコーダー押し込みで終了して選択画面へ |
| どちらでも | レイヤーキーの長押しで `NumPad` に戻る |

`Game` レイヤーの間は音量調節・マウス操作・LEDの使用率表示を止める。
ゲーム本体は [rin2yh/dinosaur-game](https://github.com/rin2yh/dinosaur-game) の
`game` パッケージをそのまま使っていて、`internal/arcade` は選択画面と
ssd1306 への描画だけを持つ。追加するゲームが 128x64 の 1bit 画面に
`SetPixel` で描き、`Update(pressed bool)` で 1 フレーム進む形なら
`internal/arcade/arcade.go` の `entries` に足すだけで選択画面に並ぶ。

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

### 再インストール (多重起動してしまった場合の掃除付き)

plist の IOKit マッチング設定ミス等で複数プロセスが同時起動してしまった場合、
`install.sh` を走らせる前に稼働中のデーモンを明示的に止めてから入れ直す。

```sh
launchctl bootout "gui/${UID}/com.github.rin2yh.tiny-deck" 2>/dev/null || true
pkill -f "${HOME}/.local/bin/tiny-deck-host" || true
mise run install
```

背景は [docs/2026-04-20-usb-attach-kernel-panic.md](docs/2026-04-20-usb-attach-kernel-panic.md) を参照。
