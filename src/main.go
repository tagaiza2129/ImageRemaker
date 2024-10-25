package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/tidwall/gjson"
)

func main() {
	AppDir, err := filepath.Abs("../")
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		return
	}
	fmt.Println("ApplicationDirectory:", AppDir)
	ConfigDir := filepath.Join(AppDir, "config.json")
	if _, err := os.ReadFile(ConfigDir); err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	Config, err := os.ReadFile(ConfigDir)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	fmt.Println("ConfigFile:", ConfigDir)
	resp, err := http.Get("http://" + gjson.Get(string(Config), "IPADDR.2").String() + ":" + gjson.Get(string(Config), "port").String() + "/OSList")
	if err != nil {
		fmt.Println("Error making GET request:", err)
		return
	}
	defer resp.Body.Close()
	OSINFO, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	myApp := app.New()
	myWindow := myApp.NewWindow("ImageRemaker")
	myWindow.Resize(fyne.NewSize(400, 500))
	myWindow.SetFixedSize(true)
	AppIconPath := filepath.Join(AppDir, "assets", "favicon.png")
	imageData, err := os.ReadFile(AppIconPath)
	if err != nil {
		fmt.Println("Error reading image file:", err)
		return
	}
	myWindow.SetIcon(fyne.NewStaticResource("favicon.png", imageData))
	//そもそも下の構文が違うっぽい
	var osOptionsList []string
	getKeys(string(OSINFO), "", &osOptionsList)
	fmt.Println("OSOptions:", osOptionsList)
	var DistributionSelect *widget.Select
	//以下の変数はLinuxの場合のみ有効
	var DistributionOptions []string
	var versionOptions []string
	var versionSelect *widget.Select
	var MediaSelect *widget.Select
	osSelect := widget.NewSelect(osOptionsList, func(selected string) {
		fmt.Println("Selected OS:", selected)
	})
	osSelect.PlaceHolder = "Select OS"
	DistributionSelect = widget.NewSelect(DistributionOptions, func(selected string) {
		fmt.Println("Selected Distribution:", selected)
	})
	DistributionSelect.PlaceHolder = "Distribution"
	MediaSelect = widget.NewSelect([]string{"ISO", "USB"}, func(selected string) {
		fmt.Println("Selected Media:", selected)
	})
	MediaSelect.PlaceHolder = "インストールメディアを選択してください"

	DistributionSelect.OnChanged = func(selected string) {
		if selected == "Ubuntu" {
			versionSelect.ClearSelected()
			versionOptions = []string{}
			versions := gjson.Get(string(OSINFO), "Linux.Ubuntu.server")
			fmt.Println("Versions:", versions)
			serverVersions := []string{}
			desktopVersions := []string{}
			versions.ForEach(func(key, value gjson.Result) bool {
				serverVersions = append(serverVersions, value.String())
				return true
			})
			versions = gjson.Get(string(OSINFO), "Linux.Ubuntu.desktop")
			versions.ForEach(func(key, value gjson.Result) bool {
				desktopVersions = append(desktopVersions, value.String())
				return true
			})
			// Interleave server and desktop versions
			maxLen := len(serverVersions)
			if len(desktopVersions) > maxLen {
				maxLen = len(desktopVersions)
			}
			for i := 0; i < maxLen; i++ {
				if i < len(serverVersions) {
					versionOptions = append(versionOptions, serverVersions[i])
				}
				if i < len(desktopVersions) {
					versionOptions = append(versionOptions, desktopVersions[i])
				}
			}
			versionSelect.Options = versionOptions
			versionSelect.Refresh()
		} else if selected == "Kubuntu" {
			versionOptions = []string{}
			versions := gjson.Get(string(OSINFO), "Linux.Kubuntu")
			versions.ForEach(func(key, value gjson.Result) bool {
				versionOptions = append(versionOptions, value.String())
				return true
			})
			versionSelect.ClearSelected()
			versionSelect.Options = versionOptions
			versionSelect.Refresh()
		}
		versionSelect.Options = versionOptions
		versionSelect.Refresh()
	}

	versionSelect = widget.NewSelect(versionOptions, func(selected string) {
		fmt.Println("Selected Version:", selected)
	})
	versionSelect.PlaceHolder = "version"

	osSelect.OnChanged = func(selected string) {
		fmt.Println("Selected OS:", selected)
		if selected == "Windows" {
			DistributionOptions = []string{}
			versionOptions = []string{}
			versions := gjson.Get(string(OSINFO), "Windows")
			versions.ForEach(func(key, value gjson.Result) bool {
				versionOptions = append(versionOptions, value.String())
				return true
			})
			fmt.Println("VersionOptions:", versionOptions)
			DistributionSelect.ClearSelected()
			versionSelect.ClearSelected()
			DistributionSelect.Options = DistributionOptions
			DistributionSelect.Refresh()
			versionSelect.Options = versionOptions
			versionSelect.Refresh()
		} else {
			DistributionOptions = []string{}
			getKeys(string(OSINFO), "Linux", &DistributionOptions)
			versionOptions = []string{}
			versionSelect.Options = versionOptions
			versionSelect.Refresh()
			DistributionSelect.ClearSelected()
			versionSelect.ClearSelected()
			DistributionSelect.Options = DistributionOptions
			DistributionSelect.Refresh()
		}
	}
	myWindow.SetContent(container.NewVBox(
		container.NewCenter(widget.NewLabel("インストールするOSを選択してください")),
		container.NewCenter(container.NewHBox(
			osSelect,
			DistributionSelect,
			versionSelect,
		)),
		container.NewCenter(widget.NewLabel("インストールするメディアを選択してください")),
		container.NewCenter(container.NewHBox(MediaSelect)),
		widget.NewButton("Click Me", func() {
			println("Button clicked!")
		}),
	))

	myWindow.ShowAndRun()
}

// keyを取得するだけだったらデフォルトの関数を使ったほうが簡単
func getKeys(jsonString string, key string, keys *[]string) {
	var JsonToMap map[string]interface{}
	if err := json.Unmarshal([]byte(jsonString), &JsonToMap); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}
	if key != "" {
		if PathToKey, ok := JsonToMap[key].(map[string]interface{}); ok {
			for distro := range PathToKey {
				*keys = append(*keys, distro)
			}
		} else {
			fmt.Println("指定されたキーが存在しませんでした")
		}
	} else {
		for key := range JsonToMap {
			*keys = append(*keys, key)
		}
	}
}
