from os import system
from os import rename
from os import remove


system("go build shelf.go")

remove(u"C:\shelf\shelf.exe")

rename("shelf.exe", u"C:\shelf\shelf.exe")