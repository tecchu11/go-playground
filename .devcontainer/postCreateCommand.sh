#!/bin/sh

# setup github-cli completion
gh completion -s zsh > /usr/local/share/zsh/site-functions/_gh

# setup asdf
git clone https://github.com/asdf-vm/asdf.git ~/.asdf --branch v0.14.0
cat << 'EOF' >> ~/.zshrc
. "$HOME/.asdf/asdf.sh"
fpath=(${ASDF_DIR}/completions $fpath)
autoload -Uz compinit && compinit
EOF
