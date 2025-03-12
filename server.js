// Import required modules
const express = require('express');
const { Gateway, Wallets } = require('fabric-network');
const path = require('path');
const fs = require('fs');

const app = express();
const port = 5000;

// Middleware for JSON request handling
app.use(express.json());

// Load the connection profile
const ccpPath = path.resolve(__dirname, 'connection-org1.json');
const ccp = JSON.parse(fs.readFileSync(ccpPath, 'utf8'));

// Initialize wallet and gateway (assuming wallet setup is done)
async function connectToNetwork() {
    const walletPath = path.join(__dirname, 'wallet');
    const wallet = await Wallets.newFileSystemWallet(walletPath);

    const gateway = new Gateway();
    await gateway.connect(ccp, {
        wallet,
        identity: 'admin', // Make sure 'admin' identity is already enrolled
        discovery: { enabled: true, asLocalhost: true }
    });

    return gateway;
}

// ** ADD THIS NEW GET API ROUTE**
app.get('/api/query', async (req, res) => {
    try {
        const gateway = await connectToNetwork();
        const network = await gateway.getNetwork('mychannel'); // Change to your channel name
        const contract = network.getContract('farmer-chain.go'); // Change to your chaincode name

        // Example: Query all assets
        const result = await contract.evaluateTransaction('queryAllProducts'); // Change to your function
        await gateway.disconnect();

        res.json({ success: true, data: JSON.parse(result.toString()) });
    } catch (error) {
        console.error(`Error: ${error}`);
        res.status(500).json({ success: false, message: error.message });
    }
});
