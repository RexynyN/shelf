from os import system
from os import rename
from os import remove

os = "linux" 

system("go build shelf.go")

if os == "win":
    remove(u"C:\shelf\shelf.exe")
    rename("shelf.exe", u"C:\shelf\shelf.exe")
elif os == "linux":
    remove(u"~/bin/shelf")
    rename("shelf", u"~/bin/shelf")
    system("source ~/.bashrc")
