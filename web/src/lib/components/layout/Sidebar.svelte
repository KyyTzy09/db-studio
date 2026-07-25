<script lang="ts">
	import { onMount } from 'svelte';
	import { connectionStatus, tablesList } from '../../stores/dbStore';
	import { fetchConnectionStatus, fetchTables } from '../../../data/services';
	import type { TableInfo } from '$lib/api';
	import { Input } from '../shadcn/input';
	import { Button } from '../shadcn/button';
	import { Database, Table, Search, Sun, Moon, RefreshCw } from '@lucide/svelte';

	let { selectedTable, onSelectTable, activeTab, onSelectTab } = $props<{
		selectedTable: string | null;
		onSelectTable: (tableName: string) => void;
		activeTab: 'table' | 'schema' | 'sql';
		onSelectTab: (tab: 'table' | 'schema' | 'sql') => void;
	}>();

	let searchQuery = $state('');
	let loading = $state(true);
	let isDark = $state(true);

	onMount(async () => {
		// Initialize theme state from html class or dark default
		if (typeof document !== 'undefined') {
			isDark = document.documentElement.classList.contains('dark') || true;
			if (isDark) {
				document.documentElement.classList.add('dark');
			}
		}

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

	function toggleTheme() {
		isDark = !isDark;
		if (isDark) {
			document.documentElement.classList.add('dark');
		} else {
			document.documentElement.classList.remove('dark');
		}
	}

	async function refreshTables() {
		loading = true;
		try {
			const { tables } = await fetchTables();
			tablesList.set(tables || []);
		} catch (err) {
			console.error('Refresh tables error:', err);
		} finally {
			loading = false;
		}
	}

	let filteredTables = $derived(
		($tablesList || []).filter((t: TableInfo) =>
			t.name.toLowerCase().includes(searchQuery.toLowerCase())
		)
	);
</script>

<aside class="w-64 bg-sidebar border-r border-sidebar-border flex flex-col h-full select-none transition-colors">
	<!-- Brand Header -->
	<div class="h-14 px-4 border-b border-sidebar-border flex items-center justify-between">
		<div class="flex items-center gap-2.5">
			<div class="w-8 h-8 rounded-lg bg-primary text-primary-foreground flex items-center justify-center font-bold text-sm shadow-sm">
				<Database class="size-4" />
			</div>
			<div>
				<h1 class="font-bold text-foreground text-sm tracking-wide">DBStudio</h1>
				<p class="text-[10px] text-muted-foreground font-mono">v0.1.0-sleek</p>
			</div>
		</div>

		<!-- Theme Toggle & Refresh Actions -->
		<div class="flex items-center gap-1">
			<Button variant="ghost" size="icon-xs" onclick={refreshTables} title="Refresh tables list">
				<RefreshCw class="size-3.5 text-muted-foreground hover:text-foreground" />
			</Button>
			<Button variant="ghost" size="icon-xs" onclick={toggleTheme} title="Toggle Dark/Light Mode">
				{#if isDark}
					<Sun class="size-3.5 text-warning" />
				{:else}
					<Moon class="size-3.5 text-primary" />
				{/if}
			</Button>
		</div>
	</div>

	<!-- Database Connection Status Badge -->
	<div class="px-4 py-3 bg-secondary/50 border-b border-sidebar-border">
		<div class="flex items-center justify-between text-[10px] font-semibold uppercase tracking-wider text-muted-foreground mb-1">
			<span>Active Target</span>
			<span class="inline-flex items-center gap-1 text-[10px] font-medium text-success">
				<span class="w-1.5 h-1.5 rounded-full bg-success animate-pulse"></span>
				{$connectionStatus?.config?.driver || 'Connected'}
			</span>
		</div>
		<div class="flex items-center justify-between text-xs font-mono">
			<span class="text-primary font-medium truncate max-w-[140px]">
				{$connectionStatus?.config?.name || 'Local Database'}
			</span>
			<span class="text-[10px] text-muted-foreground">
				{$connectionStatus?.config?.host || 'localhost'}
			</span>
		</div>
	</div>

	<!-- Navigation View Tabs -->
	<div class="px-3 pt-3">
		<div class="flex p-1 bg-secondary rounded-lg text-xs font-medium text-muted-foreground border border-border/50">
			<button
				type="button"
				onclick={() => onSelectTab('table')}
				class="flex-1 py-1 rounded-md transition-all text-center cursor-pointer {
					activeTab === 'table' ? 'bg-primary text-primary-foreground font-semibold shadow-xs' : 'hover:text-foreground'
				}"
			>
				Table
			</button>
			<button
				type="button"
				onclick={() => onSelectTab('schema')}
				class="flex-1 py-1 rounded-md transition-all text-center cursor-pointer {
					activeTab === 'schema' ? 'bg-primary text-primary-foreground font-semibold shadow-xs' : 'hover:text-foreground'
				}"
			>
				Schema
			</button>
			<button
				type="button"
				onclick={() => onSelectTab('sql')}
				class="flex-1 py-1 rounded-md transition-all text-center cursor-pointer {
					activeTab === 'sql' ? 'bg-primary text-primary-foreground font-semibold shadow-xs' : 'hover:text-foreground'
				}"
			>
				SQL
			</button>
		</div>
	</div>

	<!-- Tables Search Filter -->
	<div class="p-3">
		<div class="relative">
			<Search class="size-3.5 text-muted-foreground absolute left-2.5 top-2.5" />
			<Input
				type="text"
				bind:value={searchQuery}
				placeholder="Filter tables..."
				class="h-8 pl-8 text-xs bg-background"
			/>
		</div>
	</div>

	<!-- Tables List -->
	<div class="flex-1 overflow-auto px-2 space-y-0.5">
		<div class="px-2 py-1 text-[10px] font-semibold uppercase tracking-wider text-muted-foreground flex justify-between">
			<span>Tables ({filteredTables.length})</span>
		</div>

		{#if loading}
			<div class="p-4 text-center text-xs text-muted-foreground animate-pulse">Loading tables...</div>
		{:else if filteredTables.length === 0}
			<div class="p-4 text-center text-xs text-muted-foreground">No tables found</div>
		{:else}
			{#each filteredTables as t}
				<button
					type="button"
					onclick={() => {
						onSelectTable(t.name);
						if (activeTab === 'sql') onSelectTab('table');
					}}
					class="w-full px-2.5 py-1.5 rounded-lg text-left text-xs transition-colors flex items-center justify-between group cursor-pointer {
						selectedTable === t.name
							? 'bg-primary/10 text-primary font-semibold border border-primary/20'
							: 'text-foreground hover:bg-secondary/80'
					}"
				>
					<div class="flex items-center gap-2 truncate">
						<Table class="size-3.5 text-muted-foreground group-hover:text-primary transition-colors" />
						<span class="truncate font-mono">{t.name}</span>
					</div>
					{#if t.type === 'VIEW'}
						<span class="text-[9px] px-1 py-0.2 rounded bg-primary/20 text-primary font-mono font-medium">VIEW</span>
					{/if}
				</button>
			{/each}
		{/if}
	</div>
</aside>
