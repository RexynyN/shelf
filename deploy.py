import os
from os import system, rename, remove
from platform import system as sysname

def assert_file(path: str):
    return os.path.isfile(path)

def assert_path(pathl: str, create: bool = False) -> bool:
    exists = os.path.isdir(pathl)
    if not exists and create:
        os.mkdir(pathl)
    return exists

osname = sysname().lower()

system("go build shelf.go")

if osname == "windows":
    if assert_path(u"C:\shelf", create=True) and assert_file(u"C:\shelf\shelf.exe"):
        remove(u"C:\shelf\shelf.exe")
    rename("shelf.exe", u"C:\shelf\shelf.exe")
elif osname == "linux" or osname == "darwin":
    if assert_path(u"/bin/shelf", create=True) and assert_file(u"/bin/shelf/shelf"):
        remove(u"/bin/shelf/shelf")
    rename("shelf", u"/bin/shelf/shelf")
    system("source ~/.bashrc")

print("Executable successfully deployed!")