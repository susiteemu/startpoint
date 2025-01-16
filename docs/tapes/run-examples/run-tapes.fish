#!/usr/local/bin/fish


for tape in *.tape
    set name (string replace '.tape' '' $tape)
    echo $tape
    sed -i "" "s/catppuccin-mocha/$name/g" ~/demo/.startpoint.yaml
    vhs $tape
    sed -i "" "s/$name/catppuccin-mocha/g" ~/demo/.startpoint.yaml
end
