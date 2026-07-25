<script lang="ts">
	import type { useSchemaEditor } from '../../hooks/useSchemaEditor.svelte';
	import { Button } from '../shadcn/button';
	import { Input } from '../shadcn/input';
	import { Label } from '../shadcn/label';
	import {
		Dialog,
		DialogContent,
		DialogHeader,
		DialogTitle,
		DialogFooter,
		DialogDescription
	} from '../shadcn/dialog';
	import { Plus, Trash2, Code2, Table } from '@lucide/svelte';

	let { controller } = $props<{
		controller: ReturnType<typeof useSchemaEditor>;
	}>();

	let activeTab = $state<'columns' | 'sql'>('columns');

	const commonDataTypes = [
		'INTEGER',
		'BIGINT',
		'VARCHAR(255)',
		'TEXT',
		'BOOLEAN',
		'TIMESTAMP',
		'DATE',
		'JSON',
		'REAL'
	];
</script>

<Dialog bind:open={controller.isCreateTableOpen}>
	<DialogContent class="max-w-3xl">
		<DialogHeader>
			<DialogTitle class="flex items-center justify-between text-base font-bold">
				<div class="flex items-center gap-2">
					<Table class="size-4 text-primary" /> Create New Table
				</div>
				<div class="flex items-center gap-1 rounded-lg bg-muted p-1 text-xs">
					<button
						type="button"
						class={`px-3 py-1 rounded-md font-medium transition-colors ${activeTab === 'columns' ? 'bg-background text-foreground shadow-xs' : 'text-muted-foreground hover:text-foreground'}`}
						onclick={() => (activeTab = 'columns')}
					>
						Columns Builder
					</button>
					<button
						type="button"
						class={`px-3 py-1 rounded-md font-medium transition-colors flex items-center gap-1 ${activeTab === 'sql' ? 'bg-background text-foreground shadow-xs' : 'text-muted-foreground hover:text-foreground'}`}
						onclick={() => (activeTab = 'sql')}
					>
						<Code2 class="size-3.5" /> DDL Preview
					</button>
				</div>
			</DialogTitle>
			<DialogDescription>
				Design your database table structure and specify column parameters visually.
			</DialogDescription>
		</DialogHeader>

		<div class="space-y-4 py-2">
			<div>
				<Label for="create-table-name" class="block text-xs font-semibold text-foreground mb-1">
					Table Name <span class="text-destructive">*</span>
				</Label>
				<Input
					id="create-table-name"
					type="text"
					bind:value={controller.tableName}
					placeholder="e.g. users, orders, products"
					class="h-9 text-xs font-mono"
				/>
			</div>

			{#if activeTab === 'columns'}
				<div class="border border-border rounded-lg overflow-hidden">
					<div class="max-h-[50vh] overflow-y-auto">
						<table class="w-full text-left text-xs border-collapse">
							<thead class="bg-muted text-muted-foreground text-[11px] uppercase sticky top-0 border-b border-border font-medium">
								<tr>
									<th class="px-3 py-2">Column Name</th>
									<th class="px-3 py-2">Data Type</th>
									<th class="px-2 py-2 text-center">PK</th>
									<th class="px-2 py-2 text-center">Null</th>
									<th class="px-2 py-2 text-center">Auto Inc</th>
									<th class="px-3 py-2">Default</th>
									<th class="px-2 py-2 text-center">Action</th>
								</tr>
							</thead>
							<tbody class="divide-y divide-border">
								{#each controller.columns as col, idx}
									<tr class="hover:bg-muted/50 transition-colors">
										<td class="px-3 py-2">
											<Input
												type="text"
												bind:value={col.name}
												placeholder="col_name"
												class="h-8 text-xs font-mono"
											/>
										</td>
										<td class="px-3 py-2">
											<select
												bind:value={col.data_type}
												class="h-8 w-full rounded-md border border-border bg-background px-2 text-xs font-mono focus:outline-hidden focus:ring-1 focus:ring-primary"
											>
												{#each commonDataTypes as dt}
													<option value={dt}>{dt}</option>
												{/each}
											</select>
										</td>
										<td class="px-2 py-2 text-center">
											<input
												type="checkbox"
												bind:checked={col.is_primary_key}
												class="rounded border-border accent-primary size-4"
											/>
										</td>
										<td class="px-2 py-2 text-center">
											<input
												type="checkbox"
												bind:checked={col.is_nullable}
												class="rounded border-border accent-primary size-4"
											/>
										</td>
										<td class="px-2 py-2 text-center">
											<input
												type="checkbox"
												bind:checked={col.auto_increment}
												class="rounded border-border accent-primary size-4"
											/>
										</td>
										<td class="px-3 py-2">
											<Input
												type="text"
												bind:value={col.default_value}
												placeholder="NULL / val"
												class="h-8 text-xs font-mono"
											/>
										</td>
										<td class="px-2 py-2 text-center">
											<button
												type="button"
												class="p-1.5 rounded-md hover:bg-destructive/10 text-muted-foreground hover:text-destructive transition-colors disabled:opacity-30"
												disabled={controller.columns.length <= 1}
												onclick={() => controller.removeColumnRow(idx)}
											>
												<Trash2 class="size-3.5" />
											</button>
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>

					<div class="p-2 bg-muted/40 border-t border-border flex justify-between items-center">
						<Button
							type="button"
							variant="outline"
							size="sm"
							class="text-xs gap-1"
							onclick={controller.addColumnRow}
						>
							<Plus class="size-3.5" /> Add Column
						</Button>
						<span class="text-[11px] text-muted-foreground font-mono">
							{controller.columns.length} Column(s) Defined
						</span>
					</div>
				</div>
			{:else}
				<div class="rounded-lg border border-border bg-slate-950 p-4 font-mono text-xs text-emerald-400 overflow-x-auto min-h-[220px]">
					<pre>{controller.generatedSql}</pre>
				</div>
			{/if}
		</div>

		<DialogFooter class="pt-2 border-t border-border">
			<Button type="button" variant="outline" size="sm" onclick={controller.closeCreateTableModal}>
				Cancel
			</Button>
			<Button type="button" size="sm" disabled={controller.isSubmitting} onclick={controller.submitCreateTable}>
				{controller.isSubmitting ? 'Creating Table...' : 'Create Table'}
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
