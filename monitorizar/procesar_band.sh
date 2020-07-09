cat band.txt | grep "$1:" | cut -d ":" -f 2 | cut -d "/"  -f 1 > "$1"_band_input.txt
cat band.txt | grep "$1:" | cut -d " " -f 1 > "$1"_band_input_TIME.txt
cat band.txt | grep "$1:" | cut -d ":" -f 2 | cut -d "/"  -f 2 > "$1"_band_output.txt
cat band.txt | grep "$1:" | cut -d " " -f 1 > "$1"_band_output_TIME.txt
