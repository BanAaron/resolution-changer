# Resolution Changer

A simple system tray app to quickly switch between display resolutions and refresh rates.

## How it Works

Right click to open the menu and select the desired resolution.

![A screenshot of the resolution changer UI. Shows possible resolutions and refresh rates.](assets/demo.png
"Resolution Changer UI")

You can click the tray icon to switch between the current and previously selected resolution.

## Install
1. clone the repo
    ```shell
    git clone git@github.com:BanAaron/resolution-changer.git
    ```
2. change directory
    ```shell
    cd resolution-changer
    ```
3. go build
    ```shell
    go build -ldflags -H=windowsgui
    ```
4. run the `.exe`
   ```shell
   .\resolution-changer.exe
   ```
