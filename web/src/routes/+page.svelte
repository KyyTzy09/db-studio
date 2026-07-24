<script lang="ts">
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

	let { data } = $props();

	$effect(() => {
		if (data) {
			connectionStatus.set(data.status);
			tablesList.set(data.tables || []);

			if (data.tables && data.tables.length > 0 && !$selectedTable) {
				selectedTable.set(data.tables[0].name);
			}
		}
	});
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
					onclick={() => window.location.reload()}
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
				<div class="flex items-center gap-3">
					{#if $selectedTable && $activeTab !== 'sql'}
						<span class="text-xs font-bold text-indigo-300 font-mono bg-indigo-500/10 px-2.5 py-1 rounded-md border border-indigo-500/20">
							📁 {$selectedTable}
						</span>
					{/if}

					<!-- Tabs (Always Enabled) -->
					<div class="flex bg-slate-950 p-1 rounded-lg border border-slate-800 text-xs font-medium">
						<button
							onclick={() => activeTab.set('data')}
							class="px-3 py-1 rounded-md transition-all cursor-pointer {$activeTab === 'data' ? 'bg-indigo-600 text-white font-semibold shadow-md' : 'text-slate-400 hover:text-slate-200'}"
						>
							Data Grid
						</button>
						<button
							onclick={() => activeTab.set('schema')}
							class="px-3 py-1 rounded-md transition-all cursor-pointer {$activeTab === 'schema' ? 'bg-indigo-600 text-white font-semibold shadow-md' : 'text-slate-400 hover:text-slate-200'}"
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
						<div class="flex-1 flex flex-col items-center justify-center text-slate-400 text-xs gap-2">
							<svg class="w-8 h-8 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
							</svg>
							<span>Pilih tabel dari sidebar sebelah kiri untuk melihat Data Grid.</span>
						</div>
					{/if}
				{:else if $activeTab === 'schema'}
					{#if $selectedTable}
						<SchemaExplorer tableName={$selectedTable} />
					{:else}
						<div class="flex-1 flex flex-col items-center justify-center text-slate-400 text-xs gap-2">
							<svg class="w-8 h-8 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
							</svg>
							<span>Pilih tabel dari sidebar sebelah kiri untuk melihat Schema Info.</span>
						</div>
					{/if}
				{:else if $activeTab === 'sql'}
					<SQLEditor />
				{/if}
			</div>
		</main>
	{/if}
</div>
