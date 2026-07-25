<script lang="ts">
	import type { TableSchema } from '$lib/api';
	import type { useSchemaEditor } from '../../hooks/useSchemaEditor.svelte';
	import { Button } from '../shadcn/button';
	import { Plus, Trash2, Key, Layers } from '@lucide/svelte';

	let { schema, tableName, schemaEditorController } = $props<{
		schema: TableSchema | null;
		tableName: string;
		schemaEditorController: ReturnType<typeof useSchemaEditor>;
	}>();
</script>

<div class="p-4 space-y-4">
	<div class="flex items-center justify-between">
		<div>
			<h3 class="text-sm font-bold flex items-center gap-2 text-foreground">
				<Layers class="size-4 text-primary" /> Schema Structure: <span class="font-mono text-primary">{tableName}</span>
			</h3>
			<p class="text-xs text-muted-foreground">
				Inspect column metadata, data types, constraints, and alter schema.
			</p>
		</div>

		<Button
			size="sm"
			class="gap-1.5 text-xs font-semibold"
			onclick={() => schemaEditorController.openAddColumnModal(tableName)}
		>
			<Plus class="size-3.5" /> Add Column
		</Button>
	</div>

	{#if !schema || !schema.columns || schema.columns.length === 0}
		<div class="rounded-lg border border-border bg-card p-8 text-center text-muted-foreground text-xs">
			No column metadata available for table '{tableName}'.
		</div>
	{:else}
		<div class="border border-border rounded-lg overflow-hidden bg-card shadow-xs">
			<table class="w-full text-left text-xs border-collapse">
				<thead class="bg-muted text-muted-foreground text-[11px] uppercase font-semibold border-b border-border">
					<tr>
						<th class="px-4 py-2.5">Column Name</th>
						<th class="px-4 py-2.5">Data Type</th>
						<th class="px-4 py-2.5 text-center">Primary Key</th>
						<th class="px-4 py-2.5 text-center">Nullable</th>
						<th class="px-4 py-2.5">Default Value</th>
						<th class="px-4 py-2.5 text-center">Action</th>
					</tr>
				</thead>
				<tbody class="divide-y divide-border font-mono">
					{#each schema.columns as col}
						<tr class="hover:bg-muted/50 transition-colors">
							<td class="px-4 py-2.5 font-bold text-foreground flex items-center gap-1.5">
								{#if col.is_primary_key}
									<Key class="size-3.5 text-amber-500 shrink-0" />
								{/if}
								{col.name}
							</td>
							<td class="px-4 py-2.5 text-primary">
								{col.data_type}
							</td>
							<td class="px-4 py-2.5 text-center">
								{#if col.is_primary_key}
									<span class="inline-block rounded bg-amber-500/20 px-1.5 py-0.5 text-[10px] text-amber-500 font-bold">YES</span>
								{:else}
									<span class="text-muted-foreground text-[11px]">-</span>
								{/if}
							</td>
							<td class="px-4 py-2.5 text-center">
								{#if col.is_nullable}
									<span class="inline-block rounded bg-blue-500/10 px-1.5 py-0.5 text-[10px] text-blue-400">NULL</span>
								{:else}
									<span class="inline-block rounded bg-destructive/10 px-1.5 py-0.5 text-[10px] text-destructive">NOT NULL</span>
								{/if}
							</td>
							<td class="px-4 py-2.5 text-muted-foreground">
								{col.default_value ? col.default_value : '-'}
							</td>
							<td class="px-4 py-2.5 text-center">
								<button
									type="button"
									class="p-1.5 rounded-md hover:bg-destructive/10 text-muted-foreground hover:text-destructive transition-colors"
									title="Drop Column"
									onclick={() => schemaEditorController.submitDropColumn(tableName, col.name)}
								>
									<Trash2 class="size-3.5" />
								</button>
							</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</div>
	{/if}
</div>
