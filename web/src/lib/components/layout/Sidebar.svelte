<script lang="ts">
	import { onMount } from 'svelte';
	import { connectionStatus, tablesList } from '../../stores/dbStore';
	import { fetchConnectionStatus, fetchTables, type TableInfo } from '../../../data/services';

	let { selectedTable, onSelectTable, activeTab, onSelectTab } = $props<{
		selectedTable: string | null;
		onSelectTable: (tableName: string) => void;
		activeTab: 'table' | 'schema' | 'sql';
		onSelectTab: (tab: 'table' | 'schema' | 'sql') => void;
	}>();

	let searchQuery = $state('');
	let loading = $state(true);

	onMount(async () => {
		try {
			const connStatus = await fetchConnectionStatus();
			connectionStatus.set(connStatus);

			if (connStatus.connected) {
				const { tables } = await fetchTables();
				tablesList.set(tables || []);
				if (tables && tables.length > 0 && !selectedTable) {
					onSelectTable(tables[0].name);
				}
			}
		} catch (err) {
			console.error('Failed to load sidebar data:', err);
		} finally {
			loading = false;
		}
	});

	let filteredTables = $derived(
		($tablesList || []).filter((t: TableInfo) =>
			t.name.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);
</script>

<aside class="w-64 bg-slate-900 border-r border-slate-800 flex flex-col h-full select-none">
	<!-- Brand Header -->
	<div class="h-14 px-4 border-b border-slate-800 flex items-center justify-between">
		<div class="flex items-center gap-2.5">
			<div class="w-8 h-8 rounded-lg bg-gradient-to-br from-indigo-500 to-purple-600 flex items-center justify-center text-white font-bold text-sm shadow-md shadow-indigo-500/20">
				⚡
			</div>
			<div>
				<h1 class="font-bold text-white text-sm tracking-wide">DBStudio</h1>
				<p class="text-[10px] text-slate-400 font-mono">v0.1.0-mvp</p>
			</div>
		</div>
		<span class="inline-flex items-center gap-1.5 px-2 py-0.5 rounded-full text-[10px] font-medium bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
			<span class="w-1.5 h-1.5 rounded-full bg-emerald-400 animate-pulse"></span>
			{$connectionStatus?.config?.driver || 'Connected'}
		</span>
	</div>

	<!-- Database Connection Badge -->
	<div class="px-4 py-3 bg-slate-950/60 border-b border-slate-800/80">
		<div class="text-[10px] font-semibold uppercase tracking-wider text-slate-400 mb-1">Active Target</div>
		<div class="flex items-center justify-between text-xs font-mono">
			<span class="text-indigo-300 font-medium truncate max-w-[140px]">
				{$connectionStatus?.config?.name || 'Local Database'}
			</span>
			<span class="text-[10px] text-slate-400">
				{$connectionStatus?.config?.host || 'localhost'}
			</span>
		</div>
	</div>

	<!-- Navigation Tabs -->
	<div class="px-3 pt-3">
		<div class="flex p-1 bg-slate-950 rounded-lg text-xs font-medium text-slate-400">
			<button
				onclick={() => onSelectTab('table')}
				class="flex-1 py-1.5 rounded-md transition-all text-center cursor-pointer {
					activeTab === 'table' ? 'bg-indigo-600 text-white shadow' : 'hover:text-slate-200'
				}"
			>
				Table
			</button>
			<button
				onclick={() => onSelectTab('schema')}
				class="flex-1 py-1.5 rounded-md transition-all text-center cursor-pointer {
					activeTab === 'schema' ? 'bg-indigo-600 text-white shadow' : 'hover:text-slate-200'
				}"
			>
				Schema
			</button>
			<button
				onclick={() => onSelectTab('sql')}
				class="flex-1 py-1.5 rounded-md transition-all text-center cursor-pointer {
					activeTab === 'sql' ? 'bg-indigo-600 text-white shadow' : 'hover:text-slate-200'
				}"
			>
				SQL
			</button>
		</div>
	</div>

	<!-- Tables Search Filter -->
	<div class="p-3">
		<div class="relative">
			<input
				type="text"
				bind:value={searchQuery}
				placeholder="Filter tables..."
				class="w-full bg-slate-950 border border-slate-800 text-xs rounded-lg pl-8 pr-3 py-1.5 focus:outline-none focus:border-indigo-500 text-slate-200 placeholder:text-slate-400"
			/>
			<svg class="w-3.5 h-3.5 text-slate-400 absolute left-2.5 top-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
			</svg>
		</div>
	</div>

	<!-- Tables List -->
	<div class="flex-1 overflow-auto custom-scrollbar px-2 space-y-0.5">
		<div class="px-2 py-1 text-[10px] font-semibold uppercase tracking-wider text-slate-400 flex justify-between">
			<span>Tables ({filteredTables.length})</span>
		</div>

		{#if loading}
			<div class="p-4 text-center text-xs text-slate-400 animate-pulse">Loading tables...</div>
		{:else if filteredTables.length === 0}
			<div class="p-4 text-center text-xs text-slate-400">No tables found</div>
		{:else}
			{#each filteredTables as t}
				<button
					onclick={() => {
						onSelectTable(t.name);
						if (activeTab === 'sql') onSelectTab('table');
					}}
					class="w-full px-2.5 py-2 rounded-lg text-left text-xs transition-colors flex items-center justify-between group cursor-pointer {
						selectedTable === t.name
							? 'bg-indigo-600/20 text-indigo-300 font-medium border border-indigo-500/30'
							: 'text-slate-300 hover:bg-slate-800/60 hover:text-white'
					}"
				>
					<div class="flex items-center gap-2 truncate">
						<svg class="w-3.5 h-3.5 text-slate-400 group-hover:text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M3 14h18m-9-4v8m-7 0h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
						</svg>
						<span class="truncate font-mono">{t.name}</span>
					</div>
					{#if t.type === 'VIEW'}
						<span class="text-[9px] px-1 py-0.2 rounded bg-purple-500/20 text-purple-300 font-mono">VIEW</span>
					{/if}
				</button>
			{/each}
		{/if}
	</div>
</aside>
