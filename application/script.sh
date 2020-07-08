e="0"
while [[ $i -lt $1 ]]
do 
nodejs application/application.js &
sleep 5
echo "he llegado aqui"
e=$[$e+1]
done