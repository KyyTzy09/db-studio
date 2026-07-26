<script lang="ts">
	import type { TableSchema } from '$lib/api';
	import { Key, Table, Copy, CheckCheck } from '@lucide/svelte';
	import { fetchTableDDL } from '../../../data/services';

	let { schema, x, y, isSelected, onMouseDown, onSelectTable } = $props<{
		schema: TableSchema;
		x: number;
		y: number;
		isSelected: boolean;
		onMouseDown: (e: MouseEvent) => void;
		onSelectTable: (name: string) => void;
	}>();

	let isCopied = $state(false);

	async function handleCopyDDL(e: MouseEvent) {
		e.stopPropagation();
		try {
			const res = await fetchTableDDL(schema.table_name);
			if (res && res.ddl) {
				await navigator.clipboard.writeText(res.ddl);
				isCopied = true;
				setTimeout(() => (isCopied = false), 2000);
			}
		} catch (err) {
			console.error('Failed to copy table DDL:', err);
		}
	}
</script>

<div
	role="button"
	tabindex="0"
	onmousedown={onMouseDown}
	class="absolute z-10 w-64 rounded-xl border bg-card/95 backdrop-blur-md shadow-lg select-none transition-shadow cursor-grab active:cursor-grabbing {
		isSelected ? 'border-primary ring-2 ring-primary/30 shadow-xl' : 'border-border hover:border-border/80'
	}"
	style="transform: translate({x}px, {y}px);"
>
	<!-- Card Header -->
	<div class="flex items-center justify-between px-3.5 py-2.5 bg-muted/60 border-b border-border rounded-t-xl">
		<div class="flex items-center gap-2 overflow-hidden">
			<Table class="size-4 text-primary shrink-0" />
			<span class="font-bold text-xs text-foreground font-mono truncate" title={schema.table_name}>
				{schema.table_name}
			</span>
		</div>
		<div class="flex items-center gap-1">
			<button
				type="button"
				onclick={handleCopyDDL}
				class="text-[10px] font-semibold text-muted-foreground hover:text-foreground p-1 rounded hover:bg-muted transition-colors"
				title="Copy Table DDL"
			>
				{#if isCopied}
					<CheckCheck class="size-3 text-success" />
				{:else}
					<Copy class="size-3" />
				{/if}
			</button>

			<button
				type="button"
				onclick={(e) => {
					e.stopPropagation();
					onSelectTable(schema.table_name);
				}}
				class="text-[10px] font-semibold text-primary hover:underline px-1.5 py-0.5 rounded bg-primary/10 hover:bg-primary/20 transition-colors"
				title="View Table Data"
			>
				Data Grid
			</button>
		</div>
	</div>

	<!-- Columns List -->
	<div class="p-2 space-y-1 max-h-64 overflow-y-auto font-mono text-[11px]">
		{#each schema.columns as col}
			<div class="flex items-center justify-between px-2 py-1 rounded hover:bg-secondary/60 transition-colors group">
				<div class="flex items-center gap-1.5 overflow-hidden">
					{#if col.is_primary_key}
						<Key class="size-3 text-warning shrink-0" title="Primary Key" />
					{:else}
						<span class="size-3 block shrink-0"></span>
					{/if}

					<span class="text-foreground group-hover:text-primary transition-colors truncate {col.is_primary_key ? 'font-semibold' : ''}">
						{col.name}
					</span>
				</div>

				<span class="text-[10px] text-muted-foreground shrink-0 font-medium">
					{col.data_type}
				</span>
			</div>
		{/each}
	</div>
</div>
