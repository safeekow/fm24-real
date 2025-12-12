package main

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config 設定ファイル構造
type Config struct {
	InstallPaths []InstallPath `yaml:"install_paths"`
	Backup       BackupConfig  `yaml:"backup"`
}

// InstallPath FM24インストールパス設定
type InstallPath struct {
	Name        string `yaml:"name"`
	Path        string `yaml:"path"`
	Platform    string `yaml:"platform"`
	Description string `yaml:"description,omitempty"`
}

// BackupConfig バックアップ設定
type BackupConfig struct {
	Enabled   bool   `yaml:"enabled"`
	Directory string `yaml:"directory,omitempty"`
}

// DefaultConfig デフォルト設定を生成
func DefaultConfig() *Config {
	home, _ := os.UserHomeDir()

	return &Config{
		InstallPaths: []InstallPath{
			// Windows Steam
			{
				Name:        "windows-steam",
				Path:        `C:\Program Files (x86)\Steam\steamapps\common\Football Manager 2024\data\database\db`,
				Platform:    "windows",
				Description: "Windows Steam版",
			},
			// Windows Epic Games
			{
				Name:        "windows-epic",
				Path:        `C:\Program Files\Epic Games\Football Manager 2024\data\database\db`,
				Platform:    "windows",
				Description: "Windows Epic Games版",
			},
			// macOS Steam
			{
				Name:        "macos-steam",
				Path:        filepath.Join(home, "Library/Application Support/Steam/steamapps/common/Football Manager 2024/data/database/db"),
				Platform:    "darwin",
				Description: "macOS Steam版",
			},
			// macOS App Store
			{
				Name:        "macos-appstore",
				Path:        filepath.Join(home, "Library/Application Support/Sports Interactive/Football Manager 2024/data/database/db"),
				Platform:    "darwin",
				Description: "macOS App Store版",
			},
		},
		Backup: BackupConfig{
			Enabled:   true,
			Directory: filepath.Join(home, "FM24_Backup"),
		},
	}
}

// LoadConfig 設定ファイルを読み込み
func LoadConfig(configPath string) (*Config, error) {
	// 設定ファイルが存在しない場合はデフォルト設定を使用
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("設定ファイル読み込みエラー: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("設定ファイル解析エラー: %w", err)
	}

	// バックアップディレクトリのデフォルト設定
	if config.Backup.Directory == "" {
		home, _ := os.UserHomeDir()
		config.Backup.Directory = filepath.Join(home, "FM24_Backup")
	}

	return &config, nil
}

// SaveConfig 設定ファイルを保存
func SaveConfig(configPath string, config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("設定ファイル生成エラー: %w", err)
	}

	// ディレクトリが存在しない場合は作成
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("ディレクトリ作成エラー: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("設定ファイル保存エラー: %w", err)
	}

	return nil
}

// GetDefaultConfigPath デフォルト設定ファイルパスを取得
func GetDefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".config", "fm24-real", "config.yaml")
}

// GenerateDefaultConfig デフォルト設定ファイルを生成
func GenerateDefaultConfig() error {
	configPath := GetDefaultConfigPath()

	// 既に存在する場合は上書き確認
	if _, err := os.Stat(configPath); err == nil {
		fmt.Printf("設定ファイルが既に存在します: %s\n", configPath)
		fmt.Print("上書きしますか? (y/n): ")
		var response string
		fmt.Scanln(&response)
		if response != "y" && response != "Y" {
			fmt.Println("❌ キャンセルしました")
			return nil
		}
	}

	config := DefaultConfig()
	if err := SaveConfig(configPath, config); err != nil {
		return err
	}

	fmt.Printf("✅ デフォルト設定ファイルを生成しました: %s\n", configPath)
	return nil
}

// FindInstallPathFromConfig 設定ファイルからインストールパスを検出
func FindInstallPathFromConfig(config *Config, customPath string) (string, error) {
	// カスタムパスが指定されている場合
	if customPath != "" {
		if _, err := os.Stat(customPath); err == nil {
			return customPath, nil
		}
		return "", fmt.Errorf("指定されたパスが存在しません: %s", customPath)
	}

	// 設定ファイルから検索
	for _, installPath := range config.InstallPaths {
		if _, err := os.Stat(installPath.Path); err == nil {
			return installPath.Path, nil
		}
	}

	return "", fmt.Errorf("FM24のインストールが見つかりません")
}
