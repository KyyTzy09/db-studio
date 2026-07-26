<script lang="ts">
	import type { TableSchema } from '$lib/api';
	import { fetchTableSchema, fetchTableDDL } from '../../../data/services';
	import { Button } from '../shadcn/button';
	import { RefreshCw, FileText, Key, Check, X, Copy, CheckCheck } from '@lucide/svelte';

	let { tableName } = $props<{ tableName: string }>();

	let schema = $state<TableSchema | null>(null);
	let loading = $state(true);
	let errorMsg = $state<string | null>(null);
	let isCopiedDDL = $state(false);

	$effect(() => {
		if (tableName) {
			loadSchema(tableName);
		}
	});

	async function loadSchema(name: string) {
		loading = true;
		errorMsg = null;
		try {
			schema = await fetchTableSchema(name);
		} catch (err: any) {
			errorMsg = err.message || 'Failed to load schema';
		} finally {
			loading = false;
		}
	}

	async function handleCopyDDL() {
		try {
			const data = await fetchTableDDL(tableName);
			if (data && data.ddl) {
				await navigator.clipboard.writeText(data.ddl);
				isCopiedDDL = true;
				setTimeout(() => (isCopiedDDL = false), 2000);
			}
		} catch (err) {
			console.error('Failed to copy DDL:', err);
		}
	}
</script>

<div class="flex-1 flex flex-col h-full bg-background text-foreground overflow-hidden">
	<div class="px-6 py-3 border-b border-border bg-card/40 flex items-center justify-between">
		<h2
			class="text-xs font-bold text-muted-foreground uppercase tracking-wider flex items-center gap-2"
		>
			<FileText class="size-4 text-primary" /> Schema for
			<span class="text-primary font-mono">{tableName}</span>
		</h2>
		<div class="flex items-center gap-2">
			<Button variant="outline" size="xs" onclick={handleCopyDDL}>
				{#if isCopiedDDL}
					<CheckCheck class="size-3 mr-1 text-success" /> Copied DDL
				{:else}
					<Copy class="size-3 mr-1" /> Copy DDL
				{/if}
			</Button>

			<Button variant="outline" size="xs" onclick={() => loadSchema(tableName)}>
				<RefreshCw class="size-3 mr-1" /> Refresh
			</Button>
		</div>
	</div>

	<div class="flex-1 overflow-auto p-6">
		{#if loading}
			<div class="h-64 flex flex-col items-center justify-center text-muted-foreground gap-2">
				<div
					class="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin"
				></div>
				<span class="text-xs">Loading schema metadata...</span>
			</div>
		{:else if errorMsg}
			<div
				class="p-4 rounded-xl bg-destructive/10 border border-destructive/20 text-destructive text-xs font-mono"
			>
				❌ {errorMsg}
			</div>
		{:else if schema}
			<div class="border border-border rounded-xl overflow-hidden shadow-xs bg-card">
				<table class="w-full text-left text-xs border-collapse font-mono">
					<thead>
						<tr
							class="bg-secondary/60 text-muted-foreground font-semibold border-b border-border uppercase tracking-wider"
						>
							<th class="px-4 py-3 border-r border-border/60">Column Name</th>
							<th class="px-4 py-3 border-r border-border/60">Data Type</th>
							<th class="px-4 py-3 border-r border-border/60 text-center">Nullable</th>
							<th class="px-4 py-3 border-r border-border/60 text-center">Primary Key</th>
							<th class="px-4 py-3">Default Value</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-border/60">
						{#each schema.columns as col}
							<tr class="hover:bg-secondary/40 transition-colors">
								<td class="px-4 py-2.5 border-r border-border/60 font-semibold text-foreground">
									{col.name}
								</td>
								<td class="px-4 py-2.5 border-r border-border/60 text-primary font-medium">
									{col.data_type}
								</td>
								<td class="px-4 py-2.5 border-r border-border/60 text-center">
									{#if col.is_nullable}
										<span
											class="px-2 py-0.5 rounded bg-success/10 text-success text-[10px] font-semibold inline-flex items-center gap-1"
										>
											<Check class="size-3" /> YES
										</span>
									{:else}
										<span
											class="px-2 py-0.5 rounded bg-destructive/10 text-destructive text-[10px] font-semibold inline-flex items-center gap-1"
										>
											<X class="size-3" /> NO
										</span>
									{/if}
								</td>
								<td class="px-4 py-2.5 border-r border-border/60 text-center">
									{#if col.is_primary_key}
										<span
											class="px-2 py-0.5 rounded bg-warning/20 text-warning text-[10px] font-bold inline-flex items-center gap-1"
										>
											<Key class="size-3" /> PK
										</span>
									{:else}
										<span class="text-muted-foreground/60">-</span>
									{/if}
								</td>
								<td class="px-4 py-2.5 text-muted-foreground">
									{col.default_value || '-'}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>
