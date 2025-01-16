Set Theme catppuccin-mocha
Set WindowBar Colorful
Set Margin 10
Set BorderRadius 10
Set Shell fish
Set FontFamily "Hack"
Set FontSize 16
Set MarginFill "#D3DBD8"
Set Padding 4
Set Width 800
Set Height 600
Output startpoint-add-request.gif
Sleep 0s
Hide
Type "startpoint -w ~/test/"
Enter
Show
Sleep 3s
Type "a"
Sleep 3s
Type "y"
Sleep 1s
Type@100ms "A Request"
Sleep 2s
Enter
Type ":4"
Enter
Type "i"
Sleep 2.5s
Right
Sleep 500ms
Type@100ms " http://localhost:8000/api/echo"
Sleep 2.5s
Escape
Type ":x"
Sleep 500ms
Enter
Sleep 3s
Type ":q"
Sleep 3s
Type "q"
