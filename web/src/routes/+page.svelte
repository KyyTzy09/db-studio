<script lang="ts">
	import { onMount } from 'svelte';
	import { fetchConnectionStatus, fetchTables } from '$lib/api';
	import {
		connectionStatus,
		tablesList,
		selectedTable,
		activeTab,
		isLoading,
		errorMessage
	} from '$lib/stores/dbStore';

	import Sidebar from '$lib/components/Sidebar.svelte';
	import TableGrid from '$lib/components/TableGrid.svelte';
	import SchemaExplorer from '$lib/components/SchemaExplorer.svelte';
	import SQLEditor from '$lib/components/SQLEditor.svelte';
	import ConnectWizard from '$lib/components/ConnectWizard.svelte';

	onMount(() => {
		initAppData();
	});

	async function initAppData() {
		isLoading.set(true);
		errorMessage.set(null);

		try {
			const status = await fetchConnectionStatus();
			connectionStatus.set(status);

			if (status.connected) {
				const { tables } = await fetchTables();
				tablesList.set(tables);
				if (tables.length > 0 && !$selectedTable) {
					selectedTable.set(tables[0].name);
				}
			}
		} catch (err: any) {
			errorMessage.set(err.message || 'Failed to connect to backend server');
		} finally {
			isLoading.set(false);
		}
	}
</script>

<svelte:head>
	<title>DBStudio - One Command Database Studio</title>
</svelte:head>

<div class="flex h-screen w-screen bg-slate-950 text-slate-100 font-sans overflow-hidden antialiased">
	{#if $isLoading}
		<div class="flex-1 flex flex-col items-center justify-center gap-3">
			<div class="w-10 h-10 border-3 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
			<span class="text-xs text-slate-400 font-medium">Initializing DBStudio Web Studio...</span>
		</div>
	{:else if $errorMessage}
		<div class="flex-1 flex flex-col items-center justify-center p-6">
			<div class="max-w-md bg-red-500/10 border border-red-500/20 rounded-2xl p-6 text-center shadow-2xl">
				<div class="w-12 h-12 rounded-full bg-red-500/20 text-red-400 flex items-center justify-center text-2xl mx-auto mb-3">
					⚠️
				</div>
				<h2 class="text-base font-bold text-white mb-1">Connection Failure</h2>
				<p class="text-xs text-red-300 leading-relaxed mb-4">{$errorMessage}</p>
				<button
					onclick={initAppData}
					class="px-4 py-2 rounded-lg bg-red-600 hover:bg-red-500 text-xs font-bold text-white transition-colors cursor-pointer"
				>
					Retry Connection
				</button>
			</div>
		</div>
	{:else if $connectionStatus && !$connectionStatus.connected}
		<ConnectWizard />
	{:else}
		<!-- Main DBStudio Layout -->
		<Sidebar />

		<!-- Main Studio Container -->
		<main class="flex-1 flex flex-col h-screen overflow-hidden bg-slate-950">
			<!-- Header Navigation Tabs -->
			<header class="px-6 py-2.5 border-b border-slate-800 bg-slate-900/60 flex items-center justify-between select-none">
				<div class="flex items-center gap-2">
					{#if $selectedTable && $activeTab !== 'sql'}
						<span class="text-xs font-bold text-white font-mono bg-slate-800 px-2.5 py-1 rounded-md border border-slate-700">
							{$selectedTable}
						</span>
					{/if}

					<!-- Tabs -->
					<div class="flex bg-slate-950 p-1 rounded-lg border border-slate-800 text-xs font-medium">
						<button
							onclick={() => activeTab.set('data')}
							disabled={!$selectedTable}
							class="px-3 py-1 rounded-md transition-all cursor-pointer disabled:opacity-40 {$activeTab === 'data' ? 'bg-indigo-600 text-white font-semibold shadow-md' : 'text-slate-400 hover:text-slate-200'}"
						>
							Data Grid
						</button>
						<button
							onclick={() => activeTab.set('schema')}
							disabled={!$selectedTable}
							class="px-3 py-1 rounded-md transition-all cursor-pointer disabled:opacity-40 {$activeTab === 'schema' ? 'bg-indigo-600 text-white font-semibold shadow-md' : 'text-slate-400 hover:text-slate-200'}"
						>
							Schema Info
						</button>
						<button
							onclick={() => activeTab.set('sql')}
							class="px-3 py-1 rounded-md transition-all cursor-pointer {$activeTab === 'sql' ? 'bg-indigo-600 text-white font-semibold shadow-md' : 'text-slate-400 hover:text-slate-200'}"
						>
							SQL Workspace
						</button>
					</div>
				</div>

				<div class="text-xs text-slate-500 font-mono">
					DBStudio v0.1.0-mvp
				</div>
			</header>

			<!-- Dynamic Tab View Content -->
			<div class="flex-1 flex flex-col overflow-hidden">
				{#if $activeTab === 'data'}
					{#if $selectedTable}
						<TableGrid tableName={$selectedTable} />
					{:else}
						<div class="flex-1 flex items-center justify-center text-slate-500 text-xs">
							Select a table from the sidebar to view data.
						</div>
					{/if}
				{:else if $activeTab === 'schema'}
					{#if $selectedTable}
						<SchemaExplorer tableName={$selectedTable} />
					{:else}
						<div class="flex-1 flex items-center justify-center text-slate-500 text-xs">
							Select a table from the sidebar to view schema metadata.
						</div>
					{/if}
				{:else if $activeTab === 'sql'}
					<SQLEditor />
				{/if}
			</div>
		</main>
	{/if}
</div>
