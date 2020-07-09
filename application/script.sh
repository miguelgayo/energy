nodejs application/addToWallet.js

e="0"
while [[ $e -lt $1 ]]
do 
nodejs application/applicationUser.js $e &
sleep 15
e=$[$e+1]
done

nodejs application/applicationAdmin.js $1