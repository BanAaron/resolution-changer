class ResolutionChanger:
    def __init__(self):
        import PIL.Image
        import pywintypes
        import screeninfo
        import win32api
        import win32con
        import pystray

        self.py_win_types = pywintypes
        self.screen_info = screeninfo
        self.win_32_api = win32api
        self.win_32_con = win32con
        self.dev_mode = self.py_win_types.DEVMODEType()

        self.icon_image = PIL.Image.open("icon.png")
        self.menu = (
            pystray.MenuItem(
                "Toggle Resolution", self.toggle_resolution, default=True, visible=False
            ),
            pystray.MenuItem("3840x1080", lambda: self.set_resolution(3840, 1080)),
            pystray.MenuItem("1920x1080", lambda: self.set_resolution(1920, 1080)),
            pystray.MenuItem("1280x720", lambda: self.set_resolution(1280, 720)),
            pystray.MenuItem("Quit", self.on_quit),
        )
        self.icon = pystray.Icon("Name", self.icon_image, "ResChanger", self.menu)
        self.icon.run()

    def get_resolution(self) -> tuple[int, int]:
        """
        Get the current resolution of the display
        :return: Tuple[int, int]
        """
        monitor = self.screen_info.get_monitors()[0]
        return monitor.width, monitor.height

    def set_resolution(self, width: int, height: int):
        """
        Sets the display to resolution specified within the parameters
        :param width: resolution width in pixels
        :param height: resolution height in pixels
        """
        self.dev_mode.PelsWidth = width
        self.dev_mode.PelsHeight = height

        self.dev_mode.Fields = (
            self.win_32_con.DM_PELSWIDTH | self.win_32_con.DM_PELSHEIGHT
        )
        self.win_32_api.ChangeDisplaySettings(self.dev_mode, 0)

    def toggle_resolution(self):
        """
        Toggles the display resolution between two preset resolutions
        """
        resolution = self.get_resolution()
        if resolution == (3840, 1080):
            self.dev_mode.PelsWidth = 1920
        elif resolution == (1920, 1080):
            self.dev_mode.PelsWidth = 3840
        else:
            self.dev_mode.PelsWidth = 3840

        self.dev_mode.Fields = (
            self.win_32_con.DM_PELSWIDTH | self.win_32_con.DM_PELSHEIGHT
        )
        self.win_32_api.ChangeDisplaySettings(self.dev_mode, 0)

    def on_quit(self):
        """
        Quit the program
        """
        self.icon.visible = False
        self.icon.stop()


if __name__ == "__main__":
    res_changer = ResolutionChanger
    res_changer()
