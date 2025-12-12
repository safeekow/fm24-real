package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/spf13/pflag"
)

var (
	checkFlag   bool
	applyFlag   bool
	updateFlag  bool
	initFlag    bool
	configPath  string
	customPath  string
	showVersion bool
	version     = "1.0.0"
)

func init() {
	pflag.BoolVarP(&checkFlag, "check", "c", false, "実名化対応されているかチェック")
	pflag.BoolVarP(&applyFlag, "apply", "a", false, "実名化対応を実施")
	pflag.BoolVarP(&updateFlag, "update", "u", false, "実名化対応を更新（再適用）")
	pflag.BoolVarP(&initFlag, "init", "i", false, "デフォルト設定ファイルを生成")
	pflag.StringVar(&configPath, "config", "", "設定ファイルパス（デフォルト: ~/.config/fm24-real/config.yaml）")
	pflag.StringVarP(&customPath, "path", "p", "", "FM24データベースのカスタムパス")
	pflag.BoolVarP(&showVersion, "version", "v", false, "バージョン情報を表示")

	pflag.Usage = printUsage
}

func main() {
	pflag.Parse()

	// バージョン表示
	if showVersion {
		fmt.Printf("fm24-real version %s\n", version)
		os.Exit(0)
	}

	// init コマンド（設定ファイル生成）
	if initFlag {
		if err := GenerateDefaultConfig(); err != nil {
			color.Red("❌ エラー: %v", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// ヘルプまたは引数なし
	if pflag.NFlag() == 0 {
		printUsage()
		os.Exit(0)
	}

	// 設定ファイル読み込み
	if configPath == "" {
		configPath = GetDefaultConfigPath()
	}

	config, err := LoadConfig(configPath)
	if err != nil {
		color.Red("❌ 設定ファイル読み込みエラー: %v", err)
		color.Yellow("💡 'fm24-real --init' でデフォルト設定ファイルを生成できます")
		os.Exit(1)
	}

	tool := NewFM24Tool(config)

	// コマンド実行（優先順位: check > apply > update）
	if checkFlag {
		if err := tool.CheckStatus(customPath); err != nil {
			color.Red("❌ エラー: %v", err)
			os.Exit(1)
		}
	} else if applyFlag {
		if err := tool.Apply(customPath); err != nil {
			color.Red("❌ エラー: %v", err)
			os.Exit(1)
		}
	} else if updateFlag {
		if err := tool.Update(customPath); err != nil {
			color.Red("❌ エラー: %v", err)
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Printf("Football Manager 2024 実名化ツール v%s\n\n", version)
	fmt.Println("使用方法:")
	fmt.Println("  fm24-real [オプション]")
	fmt.Println()
	fmt.Println("オプション:")
	pflag.PrintDefaults()
	fmt.Println()
	fmt.Println("例:")
	fmt.Println("  fm24-real --init                      # 設定ファイルを生成")
	fmt.Println("  fm24-real --check                     # 現在の状態を確認")
	fmt.Println("  fm24-real -c                          # 現在の状態を確認（短縮形）")
	fmt.Println("  fm24-real --apply                     # 実名化を適用")
	fmt.Println("  fm24-real --update                    # アップデート後に再適用")
	fmt.Println("  fm24-real --apply --path /path/to/db  # カスタムパスで実名化")
	fmt.Println("  fm24-real --config custom.yaml -c     # カスタム設定ファイル使用")
	fmt.Println()
	fmt.Println("設定ファイル:")
	fmt.Printf("  デフォルト: %s\n", GetDefaultConfigPath())
	fmt.Println()
}
