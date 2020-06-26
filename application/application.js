/*
 * SPDX-License-Identifier: Apache-2.0
 */

'use strict';

const { Gateway, Wallets } = require('fabric-network');
const fs = require('fs');
const path = require('path');

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
        const gateway = new Gateway();
        await gateway.connect(ccp, { wallet, identity: 'admin', discovery: { enabled: true, asLocalhost: false } });

        // Get the network (channel) our contract is deployed to.
        const network = await gateway.getNetwork('mychannel');

        // Get the contract from the network.
        const contract = network.getContract('market');

        await contract.submitTransaction('createOffer','OFFER1','-100','16','USER1',);
        console.log('Battery read');
        var result = await contract.evaluateTransaction('queryAllOffers');
        await console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        await contract.submitTransaction('matchOffers');
        console.log('Battery read');
        result = await contract.evaluateTransaction('queryAllMatches');
        await console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        await contract.submitTransaction('createUser','USER2','BERTOLDO','false');
        console.log('Battery read');
        result = await contract.evaluateTransaction('queryAllUsers');
        await console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        // Disconnect from the gateway.
        setTimeout(function() {gateway.disconnect()},1000);
         

    } catch (error) {
        console.error(`Failed to submit transaction: ${error}`);
        process.exit(1);
    }
}

main();
