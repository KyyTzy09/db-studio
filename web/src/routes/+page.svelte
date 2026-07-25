<script lang="ts">
	import {
		connectionStatus,
		tablesList,
		selectedTable,
		activeTab,
		isLoading,
		errorMessage
	} from '$lib/stores/dbStore';

	import {
		Sidebar,
		TableGrid,
		SchemaExplorer,
		SQLEditor,
		ConnectWizard
	} from '$lib/components';

	import { Button } from '$lib/components/shadcn/button';
	import { Table, Database, Terminal, FileText, AlertOctagon } from '@lucide/svelte';

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

<div class="flex h-screen w-screen bg-background text-foreground font-sans overflow-hidden antialiased select-none">
	{#if $isLoading}
		<div class="flex-1 flex flex-col items-center justify-center gap-3">
			<div class="w-10 h-10 border-3 border-primary border-t-transparent rounded-full animate-spin"></div>
			<span class="text-xs text-muted-foreground font-medium">Initializing DBStudio Web Studio...</span>
		</div>
	{:else if $errorMessage}
		<div class="flex-1 flex flex-col items-center justify-center p-6">
			<div class="max-w-md bg-destructive/10 border border-destructive/20 rounded-2xl p-6 text-center shadow-xl">
				<div class="w-12 h-12 rounded-full bg-destructive/20 text-destructive flex items-center justify-center mx-auto mb-3">
					<AlertOctagon class="size-6" />
				</div>
				<h2 class="text-base font-bold text-foreground mb-1">Connection Failure</h2>
				<p class="text-xs text-destructive leading-relaxed mb-4">{$errorMessage}</p>
				<Button
					variant="destructive"
					size="sm"
					onclick={() => window.location.reload()}
				>
					Retry Connection
				</Button>
			</div>
		</div>
	{:else if $connectionStatus && !$connectionStatus.connected}
		<ConnectWizard />
	{:else}
		<!-- Main DBStudio Layout -->
		<Sidebar
			selectedTable={$selectedTable}
			onSelectTable={(name) => selectedTable.set(name)}
			activeTab={$activeTab === 'data' ? 'table' : $activeTab}
			onSelectTab={(tab) => activeTab.set(tab === 'table' ? 'data' : tab)}
		/>

		<!-- Main Studio Container -->
		<main class="flex-1 flex flex-col h-screen overflow-hidden bg-background">
			<!-- Header Navigation Tabs -->
			<header class="px-6 py-2.5 border-b border-border bg-card/60 flex items-center justify-between">
				<div class="flex items-center gap-3">
					{#if $selectedTable && $activeTab !== 'sql'}
						<span class="text-xs font-bold text-primary font-mono bg-primary/10 px-2.5 py-1 rounded-md border border-primary/20 flex items-center gap-1.5">
							<Table class="size-3.5" /> {$selectedTable}
						</span>
					{/if}

					<!-- Tabs -->
					<div class="flex bg-secondary p-1 rounded-lg border border-border text-xs font-medium">
						<button
							type="button"
							onclick={() => activeTab.set('data')}
							class="px-3 py-1 rounded-md transition-all cursor-pointer {$activeTab === 'data' ? 'bg-primary text-primary-foreground font-semibold shadow-xs' : 'text-muted-foreground hover:text-foreground'}"
						>
							Data Grid
						</button>
						<button
							type="button"
							onclick={() => activeTab.set('schema')}
							class="px-3 py-1 rounded-md transition-all cursor-pointer {$activeTab === 'schema' ? 'bg-primary text-primary-foreground font-semibold shadow-xs' : 'text-muted-foreground hover:text-foreground'}"
						>
							Schema Info
						</button>
						<button
							type="button"
							onclick={() => activeTab.set('sql')}
							class="px-3 py-1 rounded-md transition-all cursor-pointer {$activeTab === 'sql' ? 'bg-primary text-primary-foreground font-semibold shadow-xs' : 'text-muted-foreground hover:text-foreground'}"
						>
							SQL Workspace
						</button>
					</div>
				</div>

				<div class="text-xs text-muted-foreground font-mono">
					DBStudio v0.1.0
				</div>
			</header>

			<!-- Dynamic Tab View Content -->
			<div class="flex-1 flex flex-col overflow-hidden">
				{#if $activeTab === 'data'}
					{#if $selectedTable}
						<TableGrid tableName={$selectedTable} />
					{:else}
						<div class="flex-1 flex flex-col items-center justify-center text-muted-foreground text-xs gap-2">
							<Table class="size-8 text-muted-foreground/60" />
							<span>Pilih tabel dari sidebar sebelah kiri untuk melihat Data Grid.</span>
						</div>
					{/if}
				{:else if $activeTab === 'schema'}
					{#if $selectedTable}
						<SchemaExplorer tableName={$selectedTable} />
					{:else}
						<div class="flex-1 flex flex-col items-center justify-center text-muted-foreground text-xs gap-2">
							<FileText class="size-8 text-muted-foreground/60" />
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
