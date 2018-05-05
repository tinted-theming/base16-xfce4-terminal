# base16-xfce4-terminal

This is a [``template``](https://github.com/chriskempson/base16-templates-source) for [base16](https://github.com/chriskempson/base16) which supports xfce4-terminal

## Installation

### Git
You can download this repo and copy all `.theme` files to the terminal folder.
```bash
git clone https://github.com/afg984/base16-xfce4-terminal.git ~/Downloads/base16-xfce4-terminal && \
cd ~/.local/share && \
mkdir -p xfce4/terminal && cd xfce4/terminal && \
cp -r ~/Downloads/base16-xfce4-terminal/colorschemes .
```
### Manual Build
You can use a builder like [this](https://github.com/chriskempson/base16-builder-php), instructions over there.

Or you can use `build.py`
```bash
git clone https://github.com/afg984/base16-xfce4-terminal.git ~/Downloads/base16-xfce4-terminal && \
cd ~/Downloads/base16-xfce4-terminal && \
python3 build.py
mkdir -p ~/.local/share/xfce4/terminal
cp -r ./colorschemes ~/.local/share/xfce4/terminal
```
