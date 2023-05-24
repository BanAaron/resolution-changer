import pywintypes
import screeninfo
import win32api
import win32con


class ResolutionManager:
    def __init__(self):
        self.win_types = pywintypes
        self.win_32_api = win32api
        self.win_32_con = win32con
        self.dev_mode = self.win_types.DEVMODEType()

        width, height = self.get_current_resolution()
        self.current_resolution: dict[str, int] = dict(width=width, height=height)
        self.previous_resolution: dict[str, int] = dict(width=width, height=height)

    @staticmethod
    def get_current_resolution() -> tuple[int, int]:
        """
        gets the current resolution of the primary monitor.
        @return: monitor width and height in pixel's
        """
        monitor = [
            monitor
            for monitor in screeninfo.get_monitors()
            if monitor.is_primary is True
        ][0]
        return monitor.width, monitor.height

    def set_previous_resolution(self) -> object:
        """
        stores the previous resolution.
        @return: self
        """
        width, height = self.get_current_resolution()
        self.previous_resolution["width"] = width
        self.previous_resolution["height"] = height
        return self

    def change_resolution(self, width: int, height: int) -> object:
        """
        changes the primary monitor resolution.
        @param width: width in pixel's
        @param height: height in pixel's
        @return: self
        """
        self.set_previous_resolution()
        self.dev_mode.PelsWidth = width
        self.dev_mode.PelsHeight = height
        self.dev_mode.Fields = (
            self.win_32_con.DM_PELSWIDTH | self.win_32_con.DM_PELSHEIGHT
        )
        self.win_32_api.ChangeDisplaySettings(self.dev_mode, 0)
        return self

    def toggle_resolution(self) -> object:
        """
        toggles the resolution between the current resolution and the previously selected
        @return: self
        """
        width = self.previous_resolution["width"]
        height = self.previous_resolution["height"]
        self.change_resolution(width, height)
        return self


if __name__ == "__main__":
    rm = ResolutionManager()
