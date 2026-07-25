#!/usr/bin/env node

const { spawn } = require('child_process');
const path = require('path');
const os = require('os');
const fs = require('fs');

const platform = os.platform();
const arch = os.arch();

let binaryName = '';

if (platform === 'win32') {
	binaryName = 'dbstudio-win-x64.exe';
} else if (platform === 'darwin') {
	binaryName = arch === 'arm64' ? 'dbstudio-darwin-arm64' : 'dbstudio-darwin-x64';
} else if (platform === 'linux') {
	binaryName = 'dbstudio-linux-x64';
} else {
	console.error(`❌ Unsupported operating system: ${platform}`);
	process.exit(1);
}

let binaryPath = path.join(__dirname, binaryName);

// Fallback to local dbstudio.exe or dbstudio if build target missing
if (!fs.existsSync(binaryPath)) {
	const defaultExe = platform === 'win32' ? 'dbstudio.exe' : 'dbstudio';
	const defaultPath = path.join(__dirname, defaultExe);
	if (fs.existsSync(defaultPath)) {
		binaryPath = defaultPath;
	} else {
		console.error(`❌ DBStudio binary not found at path: ${binaryPath}`);
		console.error(`   Please run 'npm run build' first.`);
		process.exit(1);
	}
}

// Spawn native binary with inherited arguments and stdio
const child = spawn(binaryPath, process.argv.slice(2), {
	stdio: 'inherit',
	env: process.env
});

child.on('error', (err) => {
	console.error(`❌ Failed to start DBStudio process:`, err);
	process.exit(1);
});

child.on('exit', (code) => {
	process.exit(code ?? 0);
});
