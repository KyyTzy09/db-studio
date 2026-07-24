<script lang="ts">
	let driver = $state('postgres');
	let host = $state('localhost');
	let port = $state(5432);
	let user = $state('postgres');
	let password = $state('');
	let database = $state('');
	let filePath = $state('./local.db');

	let saving = $state(false);

	function updatePortDefault() {
		if (driver === 'postgres') port = 5432;
		if (driver === 'mysql') port = 3306;
	}
</script>

<div class="flex-1 flex flex-col items-center justify-center bg-slate-950 p-6">
	<div class="max-w-md w-full bg-slate-900 border border-slate-800 rounded-2xl p-6 shadow-2xl">
		<div class="text-center mb-6">
			<div class="w-12 h-12 rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white text-2xl font-bold mx-auto mb-3 shadow-lg shadow-indigo-500/20">
				🛠️
			</div>
			<h2 class="text-lg font-bold text-white tracking-tight">Connect to Database</h2>
			<p class="text-xs text-slate-400 mt-1">Manual fallback setup for DBStudio</p>
		</div>

		<div class="space-y-4 text-xs">
			<div>
				<label for="driver-select" class="block text-slate-300 font-medium mb-1">Database Type</label>
				<select
					id="driver-select"
					bind:value={driver}
					onchange={updatePortDefault}
					class="w-full bg-slate-950 border border-slate-800 text-slate-200 rounded-lg px-3 py-2 focus:outline-none focus:border-indigo-500"
				>
					<option value="postgres">PostgreSQL</option>
					<option value="mysql">MySQL / MariaDB</option>
					<option value="sqlite">SQLite</option>
				</select>
			</div>

			{#if driver === 'sqlite'}
				<div>
					<label for="filepath-input" class="block text-slate-300 font-medium mb-1">SQLite File Path</label>
					<input
						id="filepath-input"
						type="text"
						bind:value={filePath}
						placeholder="./local.db"
						class="w-full bg-slate-950 border border-slate-800 text-slate-200 rounded-lg px-3 py-2 focus:outline-none focus:border-indigo-500"
					/>
				</div>
			{:else}
				<div class="grid grid-cols-3 gap-3">
					<div class="col-span-2">
						<label for="host-input" class="block text-slate-300 font-medium mb-1">Host</label>
						<input
							id="host-input"
							type="text"
							bind:value={host}
							class="w-full bg-slate-950 border border-slate-800 text-slate-200 rounded-lg px-3 py-2 focus:outline-none focus:border-indigo-500"
						/>
					</div>
					<div>
						<label for="port-input" class="block text-slate-300 font-medium mb-1">Port</label>
						<input
							id="port-input"
							type="number"
							bind:value={port}
							class="w-full bg-slate-950 border border-slate-800 text-slate-200 rounded-lg px-3 py-2 focus:outline-none focus:border-indigo-500"
						/>
					</div>
				</div>

				<div class="grid grid-cols-2 gap-3">
					<div>
						<label for="user-input" class="block text-slate-300 font-medium mb-1">Username</label>
						<input
							id="user-input"
							type="text"
							bind:value={user}
							class="w-full bg-slate-950 border border-slate-800 text-slate-200 rounded-lg px-3 py-2 focus:outline-none focus:border-indigo-500"
						/>
					</div>
					<div>
						<label for="pass-input" class="block text-slate-300 font-medium mb-1">Password</label>
						<input
							id="pass-input"
							type="password"
							bind:value={password}
							class="w-full bg-slate-950 border border-slate-800 text-slate-200 rounded-lg px-3 py-2 focus:outline-none focus:border-indigo-500"
						/>
					</div>
				</div>

				<div>
					<label for="dbname-input" class="block text-slate-300 font-medium mb-1">Database Name</label>
					<input
						id="dbname-input"
						type="text"
						bind:value={database}
						placeholder="my_database"
						class="w-full bg-slate-950 border border-slate-800 text-slate-200 rounded-lg px-3 py-2 focus:outline-none focus:border-indigo-500"
					/>
				</div>
			{/if}

			<button
				disabled={saving}
				class="w-full mt-2 py-2.5 rounded-lg bg-indigo-600 hover:bg-indigo-500 text-xs font-bold text-white transition-all cursor-pointer shadow-lg shadow-indigo-500/20"
			>
				{saving ? 'Connecting...' : 'Save & Connect'}
			</button>
		</div>
	</div>
</div>
