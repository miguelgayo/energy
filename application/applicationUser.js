/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');
const gateway = new Gateway();
var T = 5;
const user = process.argv[2]

async function timestamp(){
    var date = new Date()
    var hour = date.getUTCHours();
    var min = date.getMinutes()
    var sec = date.getSeconds()
    var ms = date.getMilliseconds()
    process.stdout.write(hour+':'+min+':'+sec+':'+ms+' -->')
}

async function main() {
    try {
        // load the network configuration
        const ccpPath = path.resolve(__dirname,'..' ,'organizations', 'peerOrganizations', 'org1.example.com', 'connection-org1.json');
        let ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

        // Create a new file system based wallet for managing identities.
        const walletPath = path.join(__dirname, 'wallet');
        const wallet = await Wallets.newFileSystemWallet(walletPath);
        console.log(`Wallet path: ${walletPath}`);

        // Check to see if we've already enrolled the user.
        const identity = await wallet.get('admin');
        if (!identity) {
            console.log('An identity for the user "admin" does not exist in the wallet');
            console.log('Run the addToWallet.js application before retrying');
            return;
        }

        // Create a new gateway for connecting to our peer node.
        //const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'admin', discovery: { enabled: true, asLocalhost: false } });
        timestamp()
        console.log('Started application of user'+user)
        var wait_1=60;
        var wait_start = setInterval(function(){
            var date = new Date()
            var min = date.getMinutes()
            var wait_min = min%T;
            var wait = T - wait_min;
            var sec = date.getSeconds()
            wait_1 = 60 - sec + 60*(wait-1);
            if (wait_1 == 0 ||wait_1==T*60){
                process.stdout.clearLine();
                insert_offer();
                setInterval(insert_offer,T*1000*60)
                clearInterval(wait_start)
            }
        },1000)
        
        // var result = await contract.evaluateTransaction('queryAllOffers');
        // await console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        // await contract.submitTransaction('matchOffers');
        // console.log('Battery read');
        // result = await contract.evaluateTransaction('queryAllMatches');
        // await console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        // await contract.submitTransaction('createUser','USER2','BERTOLDO','false');
        // console.log('Battery read');
        // result = await contract.evaluateTransaction('queryAllUsers');
        // await console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        // Disconnect from the gateway.
        
         

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

async function insert_offer(){
    try{
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('market');
        var quantity = Math.floor(Math.random() * 200 - 100) // esto simplemente simula el cálculo de una oferta
        await contract.submitTransaction('createOffer', 'OFFER'+user.toString(), quantity.toString(), '16', 'USER'+user.toString());
        timestamp()
        console.log('Offer inserted ' + quantity);
    } catch (error){
        console.error(`Failed to insert offer: ${error}`);
        process.exit(1);
    }
}



main();
