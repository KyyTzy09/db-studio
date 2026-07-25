<script lang="ts">
	import { Input } from '../shadcn/input';
	import { Button } from '../shadcn/button';
	import { Label } from '../shadcn/label';
	import { Database, Wrench } from '@lucide/svelte';

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

<div class="flex-1 flex flex-col items-center justify-center bg-background p-6">
	<div class="max-w-md w-full bg-card border border-border rounded-2xl p-6 shadow-xl">
		<div class="text-center mb-6">
			<div class="w-12 h-12 rounded-xl bg-primary/10 border border-primary/20 flex items-center justify-center text-primary mx-auto mb-3 shadow-xs">
				<Wrench class="size-6" />
			</div>
			<h2 class="text-base font-bold text-foreground tracking-tight">Connect to Database</h2>
			<p class="text-xs text-muted-foreground mt-1">Manual fallback setup for DBStudio</p>
		</div>

		<div class="space-y-4 text-xs">
			<div>
				<Label for="driver-select" class="block font-medium mb-1.5 text-foreground">Database Type</Label>
				<select
					id="driver-select"
					bind:value={driver}
					onchange={updatePortDefault}
					class="w-full bg-background border border-border text-foreground rounded-lg px-3 py-2 text-xs focus:outline-none focus:ring-2 focus:ring-ring"
				>
					<option value="postgres">PostgreSQL</option>
					<option value="mysql">MySQL / MariaDB</option>
					<option value="sqlite">SQLite</option>
				</select>
			</div>

			{#if driver === 'sqlite'}
				<div>
					<Label for="filepath-input" class="block font-medium mb-1.5 text-foreground">SQLite File Path</Label>
					<Input
						id="filepath-input"
						type="text"
						bind:value={filePath}
						placeholder="./local.db"
						class="h-9 text-xs"
					/>
				</div>
			{:else}
				<div class="grid grid-cols-3 gap-3">
					<div class="col-span-2">
						<Label for="host-input" class="block font-medium mb-1.5 text-foreground">Host</Label>
						<Input
							id="host-input"
							type="text"
							bind:value={host}
							class="h-9 text-xs"
						/>
					</div>
					<div>
						<Label for="port-input" class="block font-medium mb-1.5 text-foreground">Port</Label>
						<Input
							id="port-input"
							type="number"
							bind:value={port}
							class="h-9 text-xs"
						/>
					</div>
				</div>

				<div class="grid grid-cols-2 gap-3">
					<div>
						<Label for="user-input" class="block font-medium mb-1.5 text-foreground">Username</Label>
						<Input
							id="user-input"
							type="text"
							bind:value={user}
							class="h-9 text-xs"
						/>
					</div>
					<div>
						<Label for="pass-input" class="block font-medium mb-1.5 text-foreground">Password</Label>
						<Input
							id="pass-input"
							type="password"
							bind:value={password}
							class="h-9 text-xs"
						/>
					</div>
				</div>

				<div>
					<Label for="dbname-input" class="block font-medium mb-1.5 text-foreground">Database Name</Label>
					<Input
						id="dbname-input"
						type="text"
						bind:value={database}
						placeholder="my_database"
						class="h-9 text-xs"
					/>
				</div>
			{/if}

			<Button
				disabled={saving}
				class="w-full mt-2"
			>
				<Database class="size-4 mr-1.5" />
				{saving ? 'Connecting...' : 'Save & Connect'}
			</Button>
		</div>
	</div>
</div>
