
```

          __| |  |_  )
          _|  │  │  / 
        _|   ____│___|


```

Fuzzy search the current directory for files. `fuz` is a simple utility
written in golang, which lists the files in the current directory.
You can search through the filenames and open the file that you are 
looking for with ease.

<img src="screenshot.png" width="1153">

# Installation

TBA

# Usage

Fuz has 2 modes,

 - Normal Mode 
 - Search Mode

Use `Esc` to toggle the mode. `fuz` starts up in search mode. In
search mode, you can type in you search string, the file list
will get filtered accordingly.

In normal mode, `j` and `k` keys can be used to move up and down
the file list. On pressing `Enter`, the file will be opened up
in your `$EDITOR`. Use `q` to quit.

To summarize,

```
j       Move downwards
k       Move upwards
Enter   Open file
q       Quit
Esc     Toggle mode
```
 
The current mode, will be displayed next to the search space,

```
[S]     Search Mode
[N]     Normal Mode
```


