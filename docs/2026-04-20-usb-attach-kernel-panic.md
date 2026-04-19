# USB 接続/切断時のカーネルパニック事件簿 (2026-04-20)

## TL;DR

- **症状**: tiny-deck を USB に挿す / 抜くたびに Mac が再起動するレベルで落ちる。
- **真因**: macOS 26.3.1 の `AppleSerialShim` / `IOSerialFamily` / `AppleUSBXHCI`
  の USB-CDC 系カーネルドライバにバグがある。**OS 側の問題なので我々では直せない**。
- **回避策**:
  1. LaunchAgent の IOKit マッチング条件を絞り、**1 回の attach で 1 プロセスだけ**
     起動するようにした。
  2. それでも同時起動が発生しないよう、デーモン側で **flock による単一インスタンス
     保証**を入れた。
  3. シリアル切断時の `log.Fatalf` を **clean exit** に変更。
  4. TCC 拒否されたまま回り続けていた通知 DB ポーラを、**初回失敗で静かに停止**。
- 結果、複数プロセスがシリアルを奪い合って CDC スタックを過剰にドライブする状況
  が消え、panic が誘発されなくなった。

## 発覚までの時系列

1. USB 挿抜の度に Mac がハードクラッシュ（画面ブラックアウト → Apple ロゴ →
   再起動）。
2. `/Library/Logs/DiagnosticReports/Retired/` に以下の panic レポートを発見:
   - `panic-full-2026-04-20-011620.0002.panic`
   - `panic-full-2026-04-20-011751.0002.panic`
3. panic string を読むと `Kernel data abort` と `wdog,reset_in_1` の watchdog
   reset。スタックトレースに以下の Apple 製カーネル拡張が並んでいた:
   - `AppleSerialShim`
   - `IOSerialFamily`
   - `AppleUSBXHCI`
4. `~/Library/Logs/tiny-deck.log` を見ると:
   - 1 秒間隔で **同じ attach ログが 2 行ずつ出ている** → 同時に 2 プロセス起動
     している疑い。
   - `notification DB query error: unable to open database file (14)` が 2Hz で
     延々出続けてログが 608KB まで肥大化。

## 原因の切り分け

### なぜ 2 プロセス起動していたか

`scripts/com.github.rin2yh.tiny-deck.plist` の IOKit マッチング条件:

```xml
<key>com.apple.iokit.matching</key>
<dict>
    <key>com.apple.device-attach</key>
    <dict>
        <key>idVendor</key><integer>11914</integer>
        <key>idProduct</key><integer>3</integer>
    </dict>
</dict>
```

`IOProviderClass` を指定していないため、**同じ VID/PID に該当する IOKit
オブジェクト全てにマッチ**する。USB-CDC デバイスは 1 回 attach すると:

- `IOUSBHostDevice`（デバイス本体）
- `IOUSBHostInterface`（CDC の各インタフェース）
- `IOSerialBSDClient`（/dev/cu.* を生やす側）

といった複数の IOKit オブジェクトを生成するため、それぞれに対して launchd が
「条件に合致した」と判断し、**1 attach イベントで複数回 launch** していた。

### なぜそれがカーネルパニックを引き起こすか（推測）

2 つの tiny-deck-host プロセスが同時に `/dev/cu.usbmodem*` を open しようと
して奪い合いになり、AppleSerialShim / IOSerialFamily の排他制御経路で
race を踏む。正常系なら `EBUSY` などで上位にエラーが返るはずだが、macOS
26.3.1 の当該パスにバグがあり、カーネル側で参照外しをしてパニックに至る。

**同じ挙動を Apple Silicon 搭載の他機種 (M2, M3) でも再現**したので、機体固有
ではなく OS のリグレッションと判断。

### なぜ切断でも落ちるか

片方のプロセスが write 中にデバイスが消えると、`sp.Write` が失敗 →
`log.Fatalf` で exit code 1。もう一方のプロセスも直後に同様に落ちる。
この過程で上位からの close が CDC ドライバのまだ reference を持っている
path と競合し、detach 側の panic を誘発する。

### 通知 DB のログスパム

`~/Library/Group Containers/group.com.apple.usernoted/db2/db` は TCC の
Full Disk Access 配下で、未許可だと sqlite が `SQLITE_CANTOPEN (14)` を返す。
既存コードは失敗時 `continue` していたため 2Hz で永遠にリトライし、ログを
肥大化させていた。**panic の直接原因ではない**が、プロセスが 2 つあると
合わせて 4 Hz になり I/O 負荷も高かったので合わせて直した。

## 対策の内訳

### 1. plist: マッチング条件を 1 オブジェクトに絞る

`scripts/com.github.rin2yh.tiny-deck.plist`

- `IOProviderClass` = `IOUSBHostDevice` を追加
  （デバイス本体 1 つだけにマッチ）
- ストリームキーを任意の `com.github.rin2yh.tiny-deck.attach` へ改名
  （`com.apple.device-attach` は "attach イベント" のように誤読されるが、実は
  ただのユーザー定義ラベル）
- `ThrottleInterval` = `30` で短時間の多重起動を抑止
- `ProcessType` = `Background` を追加（スケジューラヒント）

### 2. デーモン: flock による単一インスタンス保証

`cmd/host_daemon/main.go`

- `$TMPDIR/tiny-deck-host.lock` に対して `syscall.Flock(LOCK_EX|LOCK_NB)` を取得。
- 取れなかったら別インスタンスが走っているので `os.Exit(0)` で静かに抜ける。
- これで plist 設定の変更が効かない環境でも **ポート奪い合いが物理的に不可能**
  になる。深いところでの保険。

### 3. 切断時を clean exit に

`sp.Write` 失敗時に `log.Fatalf`（exit 1）で落ちていたのを、`log.Printf` +
`return` に変更。launchd のログが綺麗になり、将来 `KeepAlive` 等を導入した際の
挙動も明確になる。

### 4. 通知 DB ポーラ: 初回失敗で無効化

`watchNotifications` の初回 `SELECT` で失敗したら warning 1 行だけ出して
`return`。以後のリトライは無し。Full Disk Access 運用に切り替えれば復活する。

## 学び

- **launchd の IOKit マッチングで `IOProviderClass` は実質必須**。VID/PID だけ
  だと attach 1 回に対して複数オブジェクトにヒットし得る。
- 外部要因で 2 重起動が発生し得る系では、アプリ側でも **flock 等で単一
  インスタンス保証を入れる**のが保険として強い。plist に全てを任せない。
- `log.Fatalf` は "想定外のバグで停止" を示すシグナル。**接続切断など想定内の
  終了条件**には使わない。
- "ログが出続ける"は負荷だけでなく、**本当の 1 発ログを埋もれさせる** という
  デバッグ上のコストも大きい。2Hz で同じエラーが出ている状況は原則バグ。
- macOS アップデート直後に挙動が変わった場合、**カーネル側のリグレッションも
  疑う**こと。panic ログのスタックにベンダー製品名が並んでいたら自分のコードは
  シロの可能性が高い。

## 参考リンク

- `/Library/Logs/DiagnosticReports/Retired/panic-full-2026-04-20-011620.0002.panic`
- `/Library/Logs/DiagnosticReports/Retired/panic-full-2026-04-20-011751.0002.panic`
- `~/Library/Logs/tiny-deck.log`（事件当時、608KB まで肥大化）
