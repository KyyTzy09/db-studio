<script lang="ts">
	import { tablesList, selectedTable, activeTab, connectionStatus } from '$lib/stores/dbStore';

	let searchQuery = $state('');

	let filteredTables = $derived(
		($tablesList || []).filter((t) => t.name.toLowerCase().includes(searchQuery.toLowerCase()))
	);

	function selectTable(name: string) {
		selectedTable.set(name);
		activeTab.set('data');
	}

	function openSQLEditor() {
		activeTab.set('sql');
	}
</script>

<aside class="w-64 bg-slate-900 border-r border-slate-800 flex flex-col h-screen text-slate-300 select-none">
	<!-- Branding Header -->
	<div class="p-4 border-b border-slate-800 flex items-center justify-between">
		<div class="flex items-center gap-2.5">
			<div class="w-8 h-8 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-lg shadow-lg shadow-indigo-500/20">
				⚡
			</div>
			<div>
				<h1 class="font-bold text-white tracking-wide text-base leading-tight">DBStudio</h1>
				<p class="text-xs text-slate-400">Zero-Config Local UI</p>
			</div>
		</div>
	</div>

	<!-- Connection Status Pill -->
	{#if $connectionStatus}
		<div class="px-4 py-3 border-b border-slate-800/60 bg-slate-950/40">
			<div class="flex items-center justify-between text-xs">
				<span class="text-slate-400 font-medium truncate max-w-[130px]" title={$connectionStatus.config?.name || 'Database'}>
					{$connectionStatus.config?.name || 'Database'}
				</span>
				{#if $connectionStatus.connected}
					<span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
						<span class="w-1.5 h-1.5 rounded-full bg-emerald-400 mr-1 animate-pulse"></span>
						{$connectionStatus.config?.driver.toUpperCase()}
					</span>
				{:else}
					<span class="inline-flex items-center px-2 py-0.5 rounded-full text-[10px] font-semibold bg-amber-500/10 text-amber-400 border border-amber-500/20">
						Disconnected
					</span>
				{/if}
			</div>
		</div>
	{/if}

	<!-- SQL Editor Button -->
	<div class="p-3">
		<button
			onclick={openSQLEditor}
			class="w-full py-2 px-3 rounded-lg bg-indigo-600/20 hover:bg-indigo-600/30 text-indigo-300 font-medium text-xs flex items-center justify-center gap-2 border border-indigo-500/30 transition-all cursor-pointer shadow-sm"
		>
			<svg class="w-4 h-4 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
			</svg>
			SQL Workspace
		</button>
	</div>

	<!-- Search Input -->
	<div class="px-3 pb-2">
		<div class="relative">
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Filter tables..."
				class="w-full bg-slate-950 text-slate-200 text-xs rounded-md pl-8 pr-3 py-1.5 border border-slate-800 focus:outline-none focus:border-indigo-500 transition-all placeholder:text-slate-500"
			/>
			<svg class="w-3.5 h-3.5 text-slate-500 absolute left-2.5 top-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
			</svg>
		</div>
	</div>

	<!-- Tables List -->
	<div class="flex-1 overflow-y-auto px-2 py-1 space-y-0.5 custom-scrollbar">
		<div class="px-2 py-1.5 text-[11px] font-semibold tracking-wider text-slate-400 uppercase flex items-center justify-between">
			<span>Tables ({filteredTables.length})</span>
		</div>

		{#if filteredTables.length === 0}
			<div class="px-3 py-4 text-center text-slate-400 text-xs">
				No tables found
			</div>
		{:else}
			{#each filteredTables as table}
				<button
					onclick={() => selectTable(table.name)}
					class="w-full text-left px-2.5 py-1.5 rounded-md text-xs font-medium flex items-center justify-between transition-colors cursor-pointer {table.name === $selectedTable && $activeTab !== 'sql' ? 'bg-indigo-600/20 text-indigo-300 font-semibold border-l-2 border-indigo-500' : 'hover:bg-slate-800/60 text-slate-400 hover:text-slate-200'}"
				>
					<div class="flex items-center gap-2 truncate">
						<svg class="w-3.5 h-3.5 text-slate-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M3 14h18m-9-4v8m-7 0h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
						</svg>
						<span class="truncate">{table.name}</span>
					</div>
					{#if table.type === 'VIEW'}
						<span class="text-[9px] uppercase px-1 py-0.5 bg-purple-500/10 text-purple-400 rounded border border-purple-500/20">view</span>
					{/if}
				</button>
			{/each}
		{/if}
	</div>
</aside>
