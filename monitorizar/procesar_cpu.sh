cat cpu.txt | grep "$1:" | cut -d ":" -f 2 | cut -d "%"  -f 1 > "$1"_cpu.txt
cat cpu.txt | grep "$1:" | cut -d " " -f 1 > "$1"_cpu_TIME.txt
