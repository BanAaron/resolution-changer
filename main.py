import darkdetect
import PIL.Image
import pystray
from resolution_manager import ResolutionManager


class ResolutionChanger:
    def __init__(self):
        self.res_manager = ResolutionManager()
        self.menu = (
            pystray.MenuItem(
                "Toggle Resolution",
                self.res_manager.toggle_resolution,
                default=True,
                visible=False,
            ),
            pystray.MenuItem(
                "3840x1080", lambda: self.res_manager.change_resolution(3840, 1080)
            ),
            pystray.MenuItem(
                "2560x1080", lambda: self.res_manager.change_resolution(2560, 1080)
            ),
            pystray.MenuItem(
                "1920x1080", lambda: self.res_manager.change_resolution(1920, 1080)
            ),
            pystray.MenuItem(
                "Quit", lambda: self.quit()
            ),
        )

        if darkdetect.isDark():
            self.icon_image = PIL.Image.open(r"img\icon_white.png")
        else:
            self.icon_image = PIL.Image.open(r"img\icon_black.png")

        self.icon = pystray.Icon(
            "Resolution Changer", self.icon_image, "ResChanger", self.menu
        )
        self.icon.run()

    def quit(self):
        """
        exits the application.
        """
        self.icon.stop()


if __name__ == "__main__":
    res_changer = ResolutionChanger
    res_changer()
