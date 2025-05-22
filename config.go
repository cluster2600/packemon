package packemon

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// Config represents the configuration for Packemon
// ConfigはPackemonの設定を表します
type Config struct {
	// General settings
	// 一般設定
	DefaultInterface string `json:"defaultInterface"` // Default network interface / デフォルトのネットワークインターフェース
	
	// Packet templates
	// パケットテンプレート
	Templates map[string]PacketTemplate `json:"templates"` // Named packet templates / 名前付きパケットテンプレート
	
	// UI settings
	// UI設定
	UI UIConfig `json:"ui"` // UI configuration / UI設定
	
	// Keyboard shortcuts
	// キーボードショートカット
	KeyboardShortcuts KeyboardShortcutConfig `json:"keyboardShortcuts"` // Keyboard shortcut configuration / キーボードショートカット設定
}

// PacketTemplate represents a template for a packet
// PacketTemplateはパケットのテンプレートを表します
type PacketTemplate struct {
	Description string                 `json:"description"` // Description of the template / テンプレートの説明
	Layers      map[string]interface{} `json:"layers"`      // Layer configurations / レイヤー設定
}

// UIConfig represents the UI configuration
// UIConfigはUI設定を表します
type UIConfig struct {
	Theme            string `json:"theme"`            // UI theme / UIテーマ
	ShowStatistics   bool   `json:"showStatistics"`   // Whether to show statistics / 統計情報を表示するかどうか
	MaxPacketHistory int    `json:"maxPacketHistory"` // Maximum number of packets to keep in history / 履歴に保持するパケットの最大数
}

// KeyboardShortcutConfig represents the keyboard shortcut configuration
// KeyboardShortcutConfigはキーボードショートカット設定を表します
type KeyboardShortcutConfig struct {
	SendPacket    string `json:"sendPacket"`    // Shortcut for sending a packet / パケット送信のショートカット
	ClearHistory  string `json:"clearHistory"`  // Shortcut for clearing history / 履歴クリアのショートカット
	SwitchToLayer map[string]string `json:"switchToLayer"` // Shortcuts for switching to layers / レイヤー切り替えのショートカット
	SaveTemplate  string `json:"saveTemplate"`  // Shortcut for saving a template / テンプレート保存のショートカット
	LoadTemplate  string `json:"loadTemplate"`  // Shortcut for loading a template / テンプレート読み込みのショートカット
}

// DefaultConfig returns the default configuration
// デフォルト設定を返します
func DefaultConfig() *Config {
	return &Config{
		DefaultInterface: "eth0",
		Templates: make(map[string]PacketTemplate),
		UI: UIConfig{
			Theme:            "dark",
			ShowStatistics:   true,
			MaxPacketHistory: 1000,
		},
		KeyboardShortcuts: KeyboardShortcutConfig{
			SendPacket:   "Ctrl+S",
			ClearHistory: "Ctrl+L",
			SwitchToLayer: map[string]string{
				"Ethernet": "Alt+1",
				"IPv4":     "Alt+2",
				"IPv6":     "Alt+3",
				"TCP":      "Alt+4",
				"UDP":      "Alt+5",
				"ICMP":     "Alt+6",
				"ICMPv6":   "Alt+7",
				"DNS":      "Alt+8",
				"HTTP":     "Alt+9",
			},
			SaveTemplate: "Ctrl+T",
			LoadTemplate: "Ctrl+O",
		},
	}
}

// GetConfigDir returns the directory where configuration files are stored
// 設定ファイルが保存されるディレクトリを返します
func GetConfigDir() (string, error) {
	// Get the user's home directory
	// ユーザーのホームディレクトリを取得
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	
	// Create the .packemon directory if it doesn't exist
	// .packemonディレクトリが存在しない場合は作成
	configDir := filepath.Join(homeDir, ".packemon")
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.Mkdir(configDir, 0755); err != nil {
			return "", fmt.Errorf("failed to create config directory: %v", err)
		}
	}
	
	return configDir, nil
}

// LoadConfig loads the configuration from the default location
// デフォルトの場所から設定を読み込みます
func LoadConfig() (*Config, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return nil, err
	}
	
	configPath := filepath.Join(configDir, "config.json")
	
	// If the config file doesn't exist, create a default one
	// 設定ファイルが存在しない場合は、デフォルトの設定ファイルを作成
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		config := DefaultConfig()
		if err := config.Save(); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
		return config, nil
	}
	
	// Read the config file
	// 設定ファイルを読み込む
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}
	
	// Parse the config file
	// 設定ファイルを解析
	config := &Config{}
	if err := json.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}
	
	return config, nil
}

// Save saves the configuration to the default location
// デフォルトの場所に設定を保存します
func (c *Config) Save() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}
	
	configPath := filepath.Join(configDir, "config.json")
	
	// Marshal the config to JSON
	// 設定をJSONにマーシャル
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}
	
	// Write the config file
	// 設定ファイルを書き込む
	if err := ioutil.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}
	
	return nil
}

// SaveTemplate saves a packet template
// パケットテンプレートを保存します
func (c *Config) SaveTemplate(name string, template PacketTemplate) error {
	if c.Templates == nil {
		c.Templates = make(map[string]PacketTemplate)
	}
	
	c.Templates[name] = template
	return c.Save()
}

// LoadTemplate loads a packet template by name
// 名前でパケットテンプレートを読み込みます
func (c *Config) LoadTemplate(name string) (PacketTemplate, error) {
	template, ok := c.Templates[name]
	if !ok {
		return PacketTemplate{}, fmt.Errorf("template not found: %s", name)
	}
	
	return template, nil
}

// ListTemplates returns a list of all template names
// すべてのテンプレート名のリストを返します
func (c *Config) ListTemplates() []string {
	names := make([]string, 0, len(c.Templates))
	for name := range c.Templates {
		names = append(names, name)
	}
	return names
}

// DeleteTemplate deletes a template by name
// 名前でテンプレートを削除します
func (c *Config) DeleteTemplate(name string) error {
	if _, ok := c.Templates[name]; !ok {
		return fmt.Errorf("template not found: %s", name)
	}
	
	delete(c.Templates, name)
	return c.Save()
}

// GetKeyboardShortcuts returns the keyboard shortcuts
// キーボードショートカットを返します
func (c *Config) GetKeyboardShortcuts() KeyboardShortcutConfig {
	return c.KeyboardShortcuts
}

// GetShortcutHelp returns a formatted string with keyboard shortcut help
// キーボードショートカットのヘルプを含むフォーマットされた文字列を返します
func (c *Config) GetShortcutHelp() string {
	help := "Keyboard Shortcuts:\n"
	help += fmt.Sprintf("  %s: Send packet\n", c.KeyboardShortcuts.SendPacket)
	help += fmt.Sprintf("  %s: Clear history\n", c.KeyboardShortcuts.ClearHistory)
	help += fmt.Sprintf("  %s: Save template\n", c.KeyboardShortcuts.SaveTemplate)
	help += fmt.Sprintf("  %s: Load template\n", c.KeyboardShortcuts.LoadTemplate)
	help += "\nLayer Shortcuts:\n"
	
	for layer, shortcut := range c.KeyboardShortcuts.SwitchToLayer {
		help += fmt.Sprintf("  %s: Switch to %s\n", shortcut, layer)
	}
	
	return help
}
