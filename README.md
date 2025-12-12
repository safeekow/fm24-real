# FM24 実名化ツール

Football Manager 2024のライセンス制限ファイルを削除して、実名表示を有効化するツールです。

## 機能

- ✅ **実名化チェック** - 現在の実名化状態を確認
- ⚡ **実名化適用** - ライセンスファイルを削除して実名化を実施
- 🔄 **実名化更新** - ゲームアップデート後の再適用
- 💾 **自動バックアップ** - 削除前に全ファイルを自動バックアップ
- 🔍 **自動インストール検出** - 設定ファイルにない場合も自動スキャンで検出
- 🖥️ **クロスプラットフォーム** - Windows/macOS対応

## インストール

### 方法1: ビルド済みバイナリ

リリースページから環境に合わせたバイナリをダウンロードしてください。

### 方法2: ソースからビルド

```bash
git clone https://github.com/safeekow/fm24-real.git
cd fm24-real
go mod download
go build -o fm24-real
```

## 初回セットアップ

**自動スキャン機能により、設定ファイルなしでも動作します！**

一般的な場所（Steam/Epic Gamesのデフォルトパス）にインストールしている場合は、そのまま使用できます。

設定ファイルを生成したい場合（カスタムパスの追加、バックアップ設定の変更など）：

```bash
fm24-real --init
```

これにより `~/.config/fm24-real/config.yaml` に設定ファイルが作成されます。

## 使用方法

### 基本コマンド

```bash
# 設定ファイル生成（初回のみ）
fm24-real --init
fm24-real -i

# 実名化状態をチェック
fm24-real --check
fm24-real -c

# 実名化を適用
fm24-real --apply
fm24-real -a

# 実名化を更新（ゲームアップデート後）
fm24-real --update
fm24-real -u

# バージョン表示
fm24-real --version
fm24-real -v

# ヘルプ表示
fm24-real --help
fm24-real -h
```

### 高度な使用方法

```bash
# カスタムパスを指定
fm24-real --check --path /custom/path/to/db/2400

# カスタム設定ファイルを使用
fm24-real --config /path/to/custom-config.yaml --check

# カスタムパスで実名化適用
fm24-real --apply -p /custom/path/to/db/2400
```

### 使用例

#### 1. 初回実名化

```bash
# まず状態を確認
$ fm24-real --check
==========================================================
FM24 実名化状態チェック
==========================================================

✓ FM24データベース検出: /Users/username/Library/Application Support/Steam/...

📋 ライセンスファイル状態:

  ⊘ lnc/all (全ファイル) (15個のファイル存在)
  ⊘ lnc/greek (全ファイル) (8個のファイル存在)
  ⊘ fake.edt (存在)
  ...

⚠️  実名化は未適用です
実名化を適用するには: fm24-real --apply

# 実名化を適用
$ fm24-real --apply
```

#### 2. ゲームアップデート後

```bash
# アップデート後、ライセンスファイルが復活した場合
$ fm24-real --update
==========================================================
FM24 実名化更新（再適用）
==========================================================

ゲームアップデート後にライセンスファイルが復活した場合に使用します

# 状態チェック → 確認 → 再適用
```

## 設定ファイル

設定ファイルは YAML 形式で、FM24のインストールパスやバックアップ設定を管理します。

### デフォルト設定ファイル

- 場所: `~/.config/fm24-real/config.yaml`
- サンプル: [config.example.yaml](config.example.yaml)

### 設定ファイルの構造

```yaml
# FM24インストールパス設定
install_paths:
  - name: windows-steam
    path: C:\Program Files (x86)\Steam\steamapps\common\Football Manager 2024\data\database\db
    platform: windows
    description: Windows Steam版

  - name: macos-steam
    path: ~/Library/Application Support/Steam/steamapps/common/Football Manager 2024/data/database/db
    platform: darwin
    description: macOS Steam版

# バックアップ設定
backup:
  enabled: true
  directory: ~/FM24_Backup
```

### カスタムインストールパスの追加

設定ファイルに独自のインストールパスを追加できます：

```yaml
install_paths:
  - name: my-custom-install
    path: /path/to/your/fm24/data/database/db
    platform: darwin  # または windows
    description: カスタムインストール
```

## 対応プラットフォーム

### Windows
- Steam版: `C:\Program Files (x86)\Steam\steamapps\common\Football Manager 2024\`
- Epic Games版: `C:\Program Files\Epic Games\Football Manager 2024\`

### macOS
- Steam版: `~/Library/Application Support/Steam/steamapps/common/Football Manager 2024/`
- App Store版: `~/Library/Application Support/Sports Interactive/Football Manager 2024/`

## 削除対象ファイル

| ファイル/ディレクトリ | 説明 |
|---------------------|------|
| `lnc/all/*` | 全ライセンスファイル |
| `lnc/greek/*` | ギリシャライセンスファイル |
| `edt/permanent/fake.edt` | 偽名定義ファイル |
| `dbc/permanent/brazil_kits.dbc` | ブラジルキット制限 |
| `dbc/permanent/forbidden names.dbc` | 禁止名前リスト |
| `dbc/permanent/license.dbc` | ライセンスデータ |
| `dbc/permanent/j league non player.dbc` | Jリーグ非選手ライセンス |
| `dbc/permanent/japan*` | 日本関連ライセンス（japan.dbc, japan_loans.dbc, japan_fake.dbc等） |
| `dbc/permanent/1_japan_removed_clubs.dbc` | 日本削除クラブリスト（24.1.1追加） |
| `language/Licensing2.dbc` | ライセンス言語ファイル |
| `language/Licensing2_chn.dbc` | 中国語ライセンス言語ファイル |

## バックアップ

削除されたファイルは自動的に以下の場所にバックアップされます：

```
~/FM24_Backup/YYYYMMDD_HHMMSS/
```

例: `~/FM24_Backup/20240118_143022/`

## 注意事項

⚠️ **重要な注意点**

1. **ゲームアップデート時**: Steam/Epic Gamesでゲームがアップデートされると、ライセンスファイルが復活する場合があります。その場合は `fm24-real --update` を再実行してください。

2. **バックアップ**: 削除前に自動的にバックアップが作成されますが、自己責任で使用してください。

3. **ゲーム再起動**: 実名化適用後は、必ずゲームを再起動してください。

4. **セーブデータ**: 既存のセーブデータには影響しません。

## トラブルシューティング

### インストールが検出されない

1. **自動スキャンを試行**
   設定ファイルに一致するパスがない場合、ツールは自動的に一般的なインストール場所をスキャンします：
   - Steamライブラリフォルダ
   - Epic Gamesインストールフォルダ
   - 一般的なデフォルトパス

2. **設定ファイルを確認**
   ```bash
   cat ~/.config/fm24-real/config.yaml
   ```

3. **カスタムパスを指定**
   ```bash
   fm24-real --check --path /path/to/your/fm24/data/database/db
   ```

4. **設定ファイルを編集**
   `~/.config/fm24-real/config.yaml` を開いて、正しいインストールパスを追加：
   ```yaml
   install_paths:
     - name: my-install
       path: /実際のインストールパス/data/database/db
       platform: darwin  # または windows
   ```

### 設定ファイルエラー

設定ファイルが見つからない、または壊れている場合：

```bash
# 設定ファイルを再生成
fm24-real --init
```

### 実名化が反映されない

1. ゲームを完全に再起動
2. `fm24-real --check` で状態を確認
3. 必要に応じて `fm24-real --update` で再適用

## 技術仕様

- **言語**: Go 1.21+
- **依存パッケージ**:
  - `github.com/spf13/pflag` - GNU形式コマンドラインフラグ
  - `github.com/fatih/color` - カラー出力
  - `gopkg.in/yaml.v3` - YAML設定ファイル

## ライセンス

MIT License

## 免責事項

このツールは個人的な使用を目的としています。使用は自己責任でお願いします。
Football Managerおよび関連する商標は、Sports Interactive Limitedの商標です。

## 参考

- [FM2024 偽名化解除方法など](https://dosukoi.bulog.jp/2023/11/12/fm2024-%E5%81%BD%E5%90%8D%E5%8C%96%E8%A7%A3%E9%99%A4%E6%96%B9%E6%B3%95%E3%81%AA%E3%81%A9/)
