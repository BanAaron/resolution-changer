from PIL import Image
from darkdetect import isDark
from pystray import Icon, Menu, MenuItem

from resolution_manager import ResolutionManager


class ResolutionChanger:
    def __init__(self):
        self.res_manager = ResolutionManager()
        self.menu = (
            MenuItem(
                "Toggle Resolution",
                self.res_manager.toggle_resolution,
                default=True,
                visible=False,
            ),
            MenuItem(
                "3840x1080", lambda: self.res_manager.change_resolution(3840, 1080)
            ),
            MenuItem(
                "2560x1080", lambda: self.res_manager.change_resolution(2560, 1080)
            ),
            MenuItem(
                "1920x1080", lambda: self.res_manager.change_resolution(1920, 1080)
            ),
            MenuItem(
                "Refresh Rate",
                Menu(
                    MenuItem("144", lambda: self.res_manager.change_refresh_rate(144)),
                    MenuItem("120", lambda: self.res_manager.change_refresh_rate(120)),
                    MenuItem("100", lambda: self.res_manager.change_refresh_rate(100)),
                    MenuItem("60", lambda: self.res_manager.change_refresh_rate(60)),
                ),
            ),
            MenuItem("Quit", lambda: self.quit()),
        )

        if isDark():
            self.icon_image = Image.open(r"img\icon_white.png")
        else:
            self.icon_image = Image.open(r"img\icon_black.png")

        self.icon = Icon("Resolution Changer", self.icon_image, "ResChanger", self.menu)
        self.icon.run()

    def quit(self):
        """
        exits the application.
        """
        self.icon.stop()


if __name__ == "__main__":
    res_changer = ResolutionChanger
    res_changer()
