package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

// TargetFile 削除対象ファイルの定義
type TargetFile struct {
	Path        string
	Description string
	IsDirectory bool
	DeleteAll   bool // ディレクトリ内全削除フラグ
}

// FM24Tool FM24実名化ツール
type FM24Tool struct {
	DBBasePath  string
	BackupDir   string
	TargetFiles []TargetFile
	Config      *Config
}

// NewFM24Tool ツールインスタンスを作成
func NewFM24Tool(config *Config) *FM24Tool {
	return &FM24Tool{
		Config: config,
		TargetFiles: []TargetFile{
			{Path: "lnc/all", Description: "lnc/all (全ファイル)", IsDirectory: true, DeleteAll: true},
			{Path: "lnc/greek", Description: "lnc/greek (全ファイル)", IsDirectory: true, DeleteAll: true},
			{Path: "edt/permanent/fake.edt", Description: "fake.edt", IsDirectory: false},
			{Path: "dbc/permanent/brazil_kits.dbc", Description: "brazil_kits.dbc", IsDirectory: false},
			{Path: "dbc/permanent/forbidden names.dbc", Description: "forbidden names.dbc", IsDirectory: false},
			{Path: "dbc/permanent/license.dbc", Description: "license.dbc", IsDirectory: false},
			{Path: "dbc/permanent/j league non player.dbc", Description: "j league non player.dbc", IsDirectory: false},
			{Path: "dbc/permanent/1_japan_removed_clubs.dbc", Description: "1_japan_removed_clubs.dbc", IsDirectory: false},
			{Path: "language/Licensing2.dbc", Description: "Licensing2.dbc", IsDirectory: false},
			{Path: "language/Licensing2_chn.dbc", Description: "Licensing2_chn.dbc", IsDirectory: false},
		},
	}
}

// DetectInstallation FM24のインストールパスを検出（設定ファイルベース）
func (t *FM24Tool) DetectInstallation(customPath string) error {
	osType := runtime.GOOS

	// カスタムパスが指定されている場合
	if customPath != "" {
		if _, err := os.Stat(customPath); err == nil {
			versionPath, err := t.detectVersionFolder(customPath)
			if err != nil {
				return fmt.Errorf("カスタムパスのバージョン検出エラー: %w", err)
			}
			t.DBBasePath = versionPath
			return nil
		}
		return fmt.Errorf("指定されたパスが存在しません: %s", customPath)
	}

	// 設定ファイルから現在のOSに対応するパスを検索
	for _, installPath := range t.Config.InstallPaths {
		// プラットフォームが一致する場合のみチェック
		if installPath.Platform != osType {
			continue
		}

		if _, err := os.Stat(installPath.Path); err == nil {
			versionPath, err := t.detectVersionFolder(installPath.Path)
			if err != nil {
				continue
			}
			t.DBBasePath = versionPath
			color.Cyan("検出: %s (%s)", installPath.Description, installPath.Name)
			return nil
		}
	}

	// 設定ファイルにない場合、自動スキャンを試行
	color.Yellow("設定ファイルに一致するパスが見つかりません。自動スキャンを開始します...")
	if foundPath, err := t.scanForInstallation(); err == nil {
		versionPath, err := t.detectVersionFolder(foundPath)
		if err == nil {
			t.DBBasePath = versionPath
			color.Green("✓ 自動検出: %s", foundPath)
			return nil
		}
	}

	return fmt.Errorf("FM24のインストールが見つかりません。設定ファイルを確認するか、--path オプションでパスを指定してください")
}

// detectVersionFolder データベースバージョンフォルダを検出（例: 2400, 2410など）
func (t *FM24Tool) detectVersionFolder(basePath string) (string, error) {
	entries, err := os.ReadDir(basePath)
	if err != nil {
		return "", err
	}

	var versions []int
	for _, entry := range entries {
		if entry.IsDir() {
			if version, err := strconv.Atoi(entry.Name()); err == nil {
				versions = append(versions, version)
			}
		}
	}

	if len(versions) == 0 {
		return "", fmt.Errorf("バージョンフォルダが見つかりません")
	}

	// 最新バージョンを選択
	sort.Ints(versions)
	latestVersion := versions[len(versions)-1]

	return filepath.Join(basePath, strconv.Itoa(latestVersion)), nil
}

// scanForInstallation システムをスキャンしてFM24のインストールを自動検出
func (t *FM24Tool) scanForInstallation() (string, error) {
	osType := runtime.GOOS
	home, _ := os.UserHomeDir()

	// スキャン対象パスのリスト
	var scanPaths []string

	if osType == "windows" {
		// Windows: 一般的なインストール場所をスキャン
		scanPaths = []string{
			`C:\Program Files (x86)\Steam\steamapps\common\Football Manager 2024\data\database\db`,
			`C:\Program Files\Steam\steamapps\common\Football Manager 2024\data\database\db`,
			`C:\Program Files\Epic Games\Football Manager 2024\data\database\db`,
			`C:\Program Files\Epic Games\FootballManager2024\data\database\db`,
			`C:\XboxGames\Football Manager 2024\Content\data\database\db`,
			// Steamライブラリフォルダ（複数のドライブをチェック）
			`D:\SteamLibrary\steamapps\common\Football Manager 2024\data\database\db`,
			`E:\SteamLibrary\steamapps\common\Football Manager 2024\data\database\db`,
			`F:\SteamLibrary\steamapps\common\Football Manager 2024\data\database\db`,
		}
	} else if osType == "darwin" {
		// macOS: 一般的なインストール場所をスキャン
		scanPaths = []string{
			filepath.Join(home, "Library/Application Support/Steam/steamapps/common/Football Manager 2024/data/database/db"),
			filepath.Join(home, "Library/Application Support/Sports Interactive/Football Manager 2024/data/database/db"),
			"/Users/Shared/Epic Games/FootballManager2024/data/database/db",
			filepath.Join(home, "Library/Application Support/Steam/steamapps/common/Football Manager 2024/database/data/db"),
		}
	}

	// 各パスをチェック
	for _, scanPath := range scanPaths {
		if _, err := os.Stat(scanPath); err == nil {
			// バージョンフォルダが存在するか確認
			if _, err := t.detectVersionFolder(scanPath); err == nil {
				return scanPath, nil
			}
		}
	}

	// Steamライブラリフォルダを動的に検索（macOS/Windows共通）
	if steamPath := t.findSteamLibraryPath(); steamPath != "" {
		fmPath := filepath.Join(steamPath, "steamapps/common/Football Manager 2024/data/database/db")
		if _, err := os.Stat(fmPath); err == nil {
			if _, err := t.detectVersionFolder(fmPath); err == nil {
				return fmPath, nil
			}
		}
		// macOSの場合の別パス
		if osType == "darwin" {
			fmPath = filepath.Join(steamPath, "steamapps/common/Football Manager 2024/database/data/db")
			if _, err := os.Stat(fmPath); err == nil {
				if _, err := t.detectVersionFolder(fmPath); err == nil {
					return fmPath, nil
				}
			}
		}
	}

	return "", fmt.Errorf("自動スキャンでインストールが見つかりませんでした")
}

// findSteamLibraryPath Steamライブラリパスを検索
func (t *FM24Tool) findSteamLibraryPath() string {
	osType := runtime.GOOS
	home, _ := os.UserHomeDir()

	if osType == "windows" {
		// Windows: Steamの設定ファイルからライブラリパスを取得
		steamConfig := filepath.Join(home, "AppData/Local/Steam/steamapps/libraryfolders.vdf")
		if data, err := os.ReadFile(steamConfig); err == nil {
			content := string(data)
			// "path" キーを探してパスを抽出（簡易パース）
			lines := strings.Split(content, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.Contains(line, "\"path\"") {
					// "path"		"C:\\..." の形式からパスを抽出
					parts := strings.Split(line, "\"")
					for _, part := range parts {
						if strings.Contains(part, ":\\") || strings.Contains(part, ":/") {
							// パスらしい文字列を抽出
							path := strings.Trim(part, " \t\"")
							path = strings.ReplaceAll(path, "\\\\", "\\")
							if _, err := os.Stat(path); err == nil {
								return path
							}
						}
					}
				}
			}
		}
		// デフォルトのSteamパス
		defaultPaths := []string{
			`C:\Program Files (x86)\Steam`,
			`C:\Program Files\Steam`,
		}
		for _, p := range defaultPaths {
			if _, err := os.Stat(p); err == nil {
				return p
			}
		}
	} else if osType == "darwin" {
		// macOS: デフォルトのSteamパス
		steamPath := filepath.Join(home, "Library/Application Support/Steam")
		if _, err := os.Stat(steamPath); err == nil {
			return steamPath
		}
	}

	return ""
}

// CheckStatus 実名化対応されているかチェック
func (t *FM24Tool) CheckStatus(customPath string) error {
	color.Cyan("==========================================================")
	color.Cyan("FM24 実名化状態チェック")
	color.Cyan("==========================================================\n")

	// インストールパス検出
	if err := t.DetectInstallation(customPath); err != nil {
		return err
	}

	color.Green("✓ FM24データベース検出: %s\n", t.DBBasePath)

	// 対象ファイルの存在チェック
	existCount := 0
	deletedCount := 0

	fmt.Println("\n📋 ライセンスファイル状態:")
	fmt.Println()

	for _, target := range t.TargetFiles {
		fullPath := filepath.Join(t.DBBasePath, target.Path)
		exists := false

		if target.IsDirectory {
			if stat, err := os.Stat(fullPath); err == nil && stat.IsDir() {
				// ディレクトリ内のファイル数をチェック
				entries, _ := os.ReadDir(fullPath)
				if len(entries) > 0 {
					exists = true
					color.Yellow("  ⊘ %s (%d個のファイル存在)", target.Description, len(entries))
				} else {
					color.Green("  ✓ %s (空)", target.Description)
				}
			} else {
				color.Green("  ✓ %s (ディレクトリなし)", target.Description)
			}
		} else {
			if _, err := os.Stat(fullPath); err == nil {
				exists = true
				color.Yellow("  ⊘ %s (存在)", target.Description)
			} else {
				color.Green("  ✓ %s (削除済み)", target.Description)
			}
		}

		if exists {
			existCount++
		} else {
			deletedCount++
		}
	}

	// 日本関連ファイルチェック
	japanFiles, _ := t.findJapanFiles()
	if len(japanFiles) > 0 {
		color.Yellow("  ⊘ 日本関連ファイル (%d個存在)", len(japanFiles))
		existCount += len(japanFiles)
	} else {
		color.Green("  ✓ 日本関連ファイル (削除済み)")
	}

	// 結果サマリー
	fmt.Println()
	color.Cyan("==========================================================")
	fmt.Printf("ライセンスファイル: %d個存在 / %d個削除済み\n", existCount, deletedCount)

	if existCount > 0 {
		color.Yellow("\n⚠️  実名化は未適用です")
		color.White("実名化を適用するには: fm24-real --apply")
	} else {
		color.Green("\n✅ 実名化が適用されています")
	}
	color.Cyan("==========================================================")

	return nil
}

// Apply 実名化対応を実施
func (t *FM24Tool) Apply(customPath string) error {
	color.Cyan("==========================================================")
	color.Cyan("FM24 実名化適用")
	color.Cyan("==========================================================\n")

	// インストールパス検出
	if err := t.DetectInstallation(customPath); err != nil {
		return err
	}

	color.Green("✓ FM24データベース検出: %s\n", t.DBBasePath)

	// バックアップディレクトリ作成
	if err := t.createBackupDir(); err != nil {
		return err
	}

	color.Cyan("📦 バックアップディレクトリ: %s\n", t.BackupDir)

	// 確認
	color.Yellow("\n⚠️  警告: ライセンスファイルを削除します")
	fmt.Println("バックアップは自動的に作成されますが、自己責任で実行してください")
	fmt.Print("\n続行しますか? (y/n): ")

	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		color.Red("❌ 処理をキャンセルしました")
		return nil
	}

	// 実名化処理実行
	totalFiles, deletedCount, err := t.executeRealNameProcess()
	if err != nil {
		return err
	}

	// レポート生成
	t.generateReport(totalFiles, deletedCount)

	return nil
}

// Update 実名化対応を更新（再適用）
func (t *FM24Tool) Update(customPath string) error {
	color.Cyan("==========================================================")
	color.Cyan("FM24 実名化更新（再適用）")
	color.Cyan("==========================================================\n")

	color.Yellow("ゲームアップデート後にライセンスファイルが復活した場合に使用します\n")

	// 状態チェック
	if err := t.CheckStatus(customPath); err != nil {
		return err
	}

	fmt.Println()
	fmt.Print("実名化を再適用しますか? (y/n): ")

	var response string
	fmt.Scanln(&response)
	if response != "y" && response != "Y" {
		color.Red("❌ 処理をキャンセルしました")
		return nil
	}

	// Apply処理を実行（確認なしで実行）
	if err := t.DetectInstallation(customPath); err != nil {
		return err
	}

	if err := t.createBackupDir(); err != nil {
		return err
	}

	color.Cyan("\n📦 バックアップディレクトリ: %s\n", t.BackupDir)

	totalFiles, deletedCount, err := t.executeRealNameProcess()
	if err != nil {
		return err
	}

	t.generateReport(totalFiles, deletedCount)

	return nil
}

// createBackupDir バックアップディレクトリを作成
func (t *FM24Tool) createBackupDir() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	timestamp := time.Now().Format("20060102_150405")
	t.BackupDir = filepath.Join(home, "FM24_Backup", timestamp)

	return os.MkdirAll(t.BackupDir, 0755)
}

// backupFile ファイルをバックアップ
func (t *FM24Tool) backupFile(srcPath string) error {
	relPath, err := filepath.Rel(t.DBBasePath, srcPath)
	if err != nil {
		return err
	}

	dstPath := filepath.Join(t.BackupDir, relPath)

	// ディレクトリ作成
	dstDir := filepath.Dir(dstPath)
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return err
	}

	// ファイル情報取得
	srcInfo, err := os.Stat(srcPath)
	if err != nil {
		return err
	}

	if srcInfo.IsDir() {
		return t.backupDirectory(srcPath, dstPath)
	}

	// ファイルコピー
	srcFile, err := os.ReadFile(srcPath)
	if err != nil {
		return err
	}

	return os.WriteFile(dstPath, srcFile, srcInfo.Mode())
}

// backupDirectory ディレクトリを再帰的にバックアップ
func (t *FM24Tool) backupDirectory(srcDir, dstDir string) error {
	entries, err := os.ReadDir(srcDir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		srcPath := filepath.Join(srcDir, entry.Name())
		dstPath := filepath.Join(dstDir, entry.Name())

		if entry.IsDir() {
			os.MkdirAll(dstPath, 0755)
			if err := t.backupDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			srcFile, err := os.ReadFile(srcPath)
			if err != nil {
				continue
			}
			info, _ := entry.Info()
			os.WriteFile(dstPath, srcFile, info.Mode())
		}
	}

	return nil
}

// findJapanFiles 日本関連ファイルを検索
func (t *FM24Tool) findJapanFiles() ([]string, error) {
	japanDir := filepath.Join(t.DBBasePath, "dbc/permanent")
	if _, err := os.Stat(japanDir); os.IsNotExist(err) {
		return nil, nil
	}

	entries, err := os.ReadDir(japanDir)
	if err != nil {
		return nil, err
	}

	var japanFiles []string
	for _, entry := range entries {
		if !entry.IsDir() && len(entry.Name()) >= 5 && entry.Name()[:5] == "japan" {
			japanFiles = append(japanFiles, filepath.Join(japanDir, entry.Name()))
		}
	}

	return japanFiles, nil
}

// deleteDirectoryContents ディレクトリ内の全ファイルを削除
func (t *FM24Tool) deleteDirectoryContents(dirPath string) (int, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, err
	}

	deletedCount := 0
	for _, entry := range entries {
		fullPath := filepath.Join(dirPath, entry.Name())

		// バックアップ
		t.backupFile(fullPath)

		// 削除
		if err := os.RemoveAll(fullPath); err != nil {
			color.Yellow("  ⚠️  削除失敗: %s - %v", entry.Name(), err)
			continue
		}

		deletedCount++
	}

	return deletedCount, nil
}

// executeRealNameProcess 実名化処理を実行
func (t *FM24Tool) executeRealNameProcess() (int, int, error) {
	color.Cyan("\n🔄 実名化処理を開始します...\n")

	totalFiles := 0
	deletedCount := 0

	// ターゲットファイル処理
	for _, target := range t.TargetFiles {
		fullPath := filepath.Join(t.DBBasePath, target.Path)

		if target.IsDirectory && target.DeleteAll {
			// ディレクトリ内全削除
			if stat, err := os.Stat(fullPath); err == nil && stat.IsDir() {
				count, _ := t.deleteDirectoryContents(fullPath)
				color.Green("  ✓ %s: %d個のファイルを削除", target.Description, count)
				deletedCount += count
			} else {
				color.White("  ⊘ %s: ディレクトリが見つかりません", target.Description)
			}
			totalFiles++
		} else {
			// 個別ファイル削除
			if _, err := os.Stat(fullPath); err == nil {
				// バックアップ
				t.backupFile(fullPath)

				// 削除
				if err := os.RemoveAll(fullPath); err != nil {
					color.Yellow("  ⚠️  削除失敗: %s - %v", target.Description, err)
				} else {
					color.Green("  ✓ %s: 削除完了", target.Description)
					deletedCount++
				}
			} else {
				color.White("  ⊘ %s: ファイルが見つかりません", target.Description)
			}
			totalFiles++
		}
	}

	// 日本関連ファイル削除
	japanFiles, _ := t.findJapanFiles()
	for _, jpFile := range japanFiles {
		// バックアップ
		t.backupFile(jpFile)

		// 削除
		if err := os.Remove(jpFile); err != nil {
			color.Yellow("  ⚠️  削除失敗: %s - %v", filepath.Base(jpFile), err)
		} else {
			color.Green("  ✓ %s: 削除完了", filepath.Base(jpFile))
			deletedCount++
		}
		totalFiles++
	}

	return totalFiles, deletedCount, nil
}

// generateReport 処理結果レポートを生成
func (t *FM24Tool) generateReport(totalFiles, deletedCount int) {
	fmt.Println()
	color.Cyan("==========================================================")
	color.Cyan("📊 実名化処理レポート")
	color.Cyan("==========================================================")
	fmt.Printf("対象ファイル数: %d\n", totalFiles)
	color.Green("削除成功: %d", deletedCount)
	color.Yellow("削除失敗: %d", totalFiles-deletedCount)
	fmt.Printf("バックアップ場所: %s\n", t.BackupDir)
	color.Cyan("==========================================================")

	color.Green("\n✅ 実名化処理が完了しました")
	color.Yellow("⚠️  ゲームを再起動して変更を反映してください")
	color.Yellow("⚠️  アップデート後はファイルが復活する可能性があります")
	color.White("    その場合は 'fm24-real --update' を実行してください")
}
