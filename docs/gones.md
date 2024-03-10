## gones

NES emulator written in Go

```
gones ROM [flags]
```

### Options

```
  -a, --audio             Enabled audio output (default true)
  -c, --config string     Config file (default is $HOME/.config/gones/config.yaml)
      --debug             Start with step debugging enabled
  -f, --fullscreen        Start in fullscreen
  -h, --help              help for gones
      --palette string    Optional palette (.pal) file to use
      --pause-unfocused   Pauses when the window loses focus. Optional, but audio will be glitchy when the game is running in the background. (default true)
      --resume            Automatically resume where you left off (default true)
      --scale float       Default UI scale (default 3)
      --trace             Enable trace logging
```

