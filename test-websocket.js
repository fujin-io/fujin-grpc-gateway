#!/usr/bin/env node

/**
 * Simple WebSocket test client for Fujin gRPC Gateway
 * 
 * Usage:
 *   node test-websocket.js [ws://localhost:4850/stream]
 * 
 * Or install ws package: npm install ws
 * Then run: node test-websocket.js
 */

const WebSocket = require('ws');

const url = process.argv[2] || 'ws://localhost:4850/stream';

console.log(`Connecting to ${url}...`);

const ws = new WebSocket(url);

ws.on('open', function open() {
    console.log('Connected!');
    
    // Send Init request
    const initMessage = {
        init: {
            config_overrides: {}
        }
    };
    console.log('\nSending Init:', JSON.stringify(initMessage, null, 2));
    ws.send(JSON.stringify(initMessage));
    
    // Example: Send Produce request after 1 second
    setTimeout(() => {
        const produceMessage = {
            produce: {
                correlation_id: 1,
                topic: 'test-topic',
                message: Buffer.from('Hello World').toString('base64')
            }
        };
        console.log('\nSending Produce:', JSON.stringify(produceMessage, null, 2));
        ws.send(JSON.stringify(produceMessage));
    }, 1000);
    
    // Example: Send Subscribe request after 2 seconds
    setTimeout(() => {
        const subscribeMessage = {
            subscribe: {
                correlation_id: 2,
                topic: 'test-topic',
                auto_commit: true
            }
        };
        console.log('\nSending Subscribe:', JSON.stringify(subscribeMessage, null, 2));
        ws.send(JSON.stringify(subscribeMessage));
    }, 2000);
});

ws.on('message', function message(data) {
    try {
        const json = JSON.parse(data.toString());
        console.log('\nReceived:', JSON.stringify(json, null, 2));
    } catch (e) {
        console.log('\nReceived (raw):', data.toString());
    }
});

ws.on('error', function error(err) {
    console.error('WebSocket error:', err.message);
});

ws.on('close', function close() {
    console.log('\nConnection closed');
    process.exit(0);
});

// Keep process alive
process.on('SIGINT', () => {
    console.log('\nClosing connection...');
    ws.close();
});


