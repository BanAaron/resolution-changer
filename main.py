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

        self.previous_resolution = ()

        self.icon_image = PIL.Image.open("icon.png")
        self.menu = (
            pystray.MenuItem(
                "Toggle Resolution", self.toggle_resolution, default=True, visible=False
            ),
            pystray.MenuItem("3840x1080", lambda: self.change_resolution((3840, 1080))),
            pystray.MenuItem("2560x1080", lambda: self.change_resolution((2560, 1080))),
            pystray.MenuItem("1920x1080", lambda: self.change_resolution((1920, 1080))),
            pystray.MenuItem("1680x1050", lambda: self.change_resolution((1680, 1050))),
            pystray.MenuItem("1280x720", lambda: self.change_resolution((1280, 720))),
            pystray.MenuItem("Quit", self.on_quit),
        )
        self.icon = pystray.Icon(
            "Resolution Changer", self.icon_image, "ResChanger", self.menu
        )
        self.icon.run()

    def change_resolution(self, resolution: tuple[int, int]):
        """
        Change the resolution of the display
        :param resolution: Width and height in pixels
        :return: None
        """
        self.set_previous_resolution()
        self.dev_mode.PelsWidth = resolution[0]
        self.dev_mode.PelsHeight = resolution[1]
        self.dev_mode.Fields = (
            self.win_32_con.DM_PELSWIDTH | self.win_32_con.DM_PELSHEIGHT
        )
        self.win_32_api.ChangeDisplaySettings(self.dev_mode, 0)

    def toggle_resolution(self):
        """
        Toggle between the current resolution and the previous selected
        :return: None
        """
        if len(self.get_previous_resolution()) == 2:
            self.change_resolution(self.get_previous_resolution())

    def on_quit(self):
        """
        Quits the programs main loop
        :return: None
        """
        self.icon.visible = False
        self.icon.stop()

    def get_current_resolution(self) -> tuple:
        """
        Gets the current display resolution
        :return: Tuple of two ints. Width and height in pixels
        """
        screen_info = self.screen_info.get_monitors()[0]
        return screen_info.width, screen_info.height

    def set_current_resolution(self):
        """
        Unused
        :return: None
        """
        pass

    def get_previous_resolution(self) -> tuple:
        """
        Gets the previously set resolution
        :return: Tuple of two ints. Width and height in pixels
        """
        return self.previous_resolution

    def set_previous_resolution(self):
        """
        Sets the previous resolution
        :return: None
        """
        screen_info = self.screen_info.get_monitors()[0]
        self.previous_resolution = (screen_info.width, screen_info.height)


if __name__ == "__main__":
    res_changer = ResolutionChanger
    res_changer()
