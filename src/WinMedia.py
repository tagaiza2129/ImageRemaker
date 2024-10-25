import os
def OSList():
    #USBを取得する
    RETURNUSBINFO = os.popen("powershell -Command \"Get-CimInstance Win32_Volume | Where-Object {$_.DriveType -eq 2}\"").read().split("\n")
    USBINFO = {item.split(":")[0].strip(): item.split(":")[1].strip() for item in RETURNUSBINFO if ":" in item}
    return USBINFO
if __name__ == '__main__':
    print(OSList())