<script lang="ts">
	import type { useSchemaEditor } from '../../hooks/useSchemaEditor.svelte';
	import { tablesList, fetchColumnsForTable } from '../../stores/dbStore';
	import type { ColumnInfo } from '$lib/api';
	import { Button } from '../shadcn/button';
	import { Input } from '../shadcn/input';
	import { Label } from '../shadcn/label';
	import {
		Select,
		SelectTrigger,
		SelectContent,
		SelectItem
	} from '../shadcn/select';
	import {
		Dialog,
		DialogContent,
		DialogHeader,
		DialogTitle,
		DialogFooter,
		DialogDescription
	} from '../shadcn/dialog';
	import { Plus, Link2 } from '@lucide/svelte';

	let { controller } = $props<{
		controller: ReturnType<typeof useSchemaEditor>;
	}>();

	let targetColumnsMap = $state<Record<string, ColumnInfo[]>>({});

	async function handleFkTableSelect(selectedTbl: string) {
		controller.newColumn.fk_table = selectedTbl;
		if (selectedTbl) {
			if (!targetColumnsMap[selectedTbl]) {
				targetColumnsMap[selectedTbl] = await fetchColumnsForTable(selectedTbl);
			}
			const cols = targetColumnsMap[selectedTbl] || [];
			const pk = cols.find((c) => c.is_primary_key);
			controller.newColumn.fk_column = pk ? pk.name : cols[0]?.name || 'id';
		}
	}

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

<Dialog bind:open={controller.isAddColumnOpen}>
	<DialogContent class="max-w-md">
		<DialogHeader>
			<DialogTitle class="flex items-center gap-1.5 text-base font-bold">
				<Plus class="size-4 text-primary" /> Add Column to <span class="text-primary font-mono">{controller.targetTableForColumn}</span>
			</DialogTitle>
			<DialogDescription>
				Specify the new column parameters to alter table schema.
			</DialogDescription>
		</DialogHeader>

		<form
			id="add-column-form"
			onsubmit={(e) => {
				e.preventDefault();
				controller.submitAddColumn();
			}}
			class="space-y-4 text-xs py-2"
		>
			<div>
				<Label for="add-col-name" class="block text-xs font-semibold text-foreground mb-1">
					Column Name <span class="text-destructive">*</span>
				</Label>
				<Input
					id="add-col-name"
					type="text"
					bind:value={controller.newColumn.name}
					placeholder="e.g. user_id, bio, is_active"
					class="h-9 text-xs font-mono"
				/>
			</div>

			<div>
				<Label for="add-col-datatype" class="block text-xs font-semibold text-foreground mb-1">
					Data Type
				</Label>
				<Select type="single" bind:value={controller.newColumn.data_type}>
					<SelectTrigger id="add-col-datatype" class="w-full h-9 text-xs font-mono">
						{controller.newColumn.data_type || 'Select Data Type'}
					</SelectTrigger>
					<SelectContent>
						{#each commonDataTypes as dt}
							<SelectItem value={dt} label={dt} class="text-xs font-mono" />
						{/each}
					</SelectContent>
				</Select>
			</div>

			<div>
				<Label for="add-col-default" class="block text-xs font-semibold text-foreground mb-1">
					Default Value
				</Label>
				<Input
					id="add-col-default"
					type="text"
					bind:value={controller.newColumn.default_value}
					placeholder="Leave empty or specify e.g. 'active', 0"
					class="h-9 text-xs font-mono"
				/>
			</div>

			<div class="flex items-center gap-6 pt-1">
				<label class="flex items-center gap-2 cursor-pointer text-xs select-none">
					<input
						type="checkbox"
						bind:checked={controller.newColumn.is_nullable}
						class="h-4 w-4 rounded border-input bg-background text-primary focus-visible:ring-1 focus-visible:ring-ring cursor-pointer accent-primary"
					/>
					<span class="text-foreground font-medium">Allow NULL</span>
				</label>

				<label class="flex items-center gap-2 cursor-pointer text-xs select-none">
					<input
						type="checkbox"
						bind:checked={controller.newColumn.is_primary_key}
						class="h-4 w-4 rounded border-input bg-background text-primary focus-visible:ring-1 focus-visible:ring-ring cursor-pointer accent-primary"
					/>
					<span class="text-foreground font-medium">Primary Key</span>
				</label>

				<label class="flex items-center gap-2 cursor-pointer text-xs select-none">
					<input
						type="checkbox"
						checked={controller.newColumn.is_foreign_key}
						onchange={(e) => {
							controller.newColumn.is_foreign_key = (e.target as HTMLInputElement).checked;
							if (controller.newColumn.is_foreign_key && !controller.newColumn.fk_column) {
								controller.newColumn.fk_column = 'id';
							}
						}}
						class="h-4 w-4 rounded border-input bg-background text-primary focus-visible:ring-1 focus-visible:ring-ring cursor-pointer accent-primary"
					/>
					<span class="text-foreground font-medium flex items-center gap-1">
						<Link2 class="size-3 text-primary" /> Foreign Key
					</span>
				</label>
			</div>

			{#if controller.newColumn.is_foreign_key}
				<div class="p-3 bg-muted/40 border border-border rounded-lg space-y-3">
					<div class="flex items-center justify-between text-xs font-semibold text-foreground">
						<span>Foreign Key Reference</span>
					</div>

					<div class="grid grid-cols-2 gap-2">
						<div>
							<Label class="block text-[11px] text-muted-foreground mb-1">Target Table</Label>
							<Select
								type="single"
								bind:value={controller.newColumn.fk_table}
								onValueChange={handleFkTableSelect}
							>
								<SelectTrigger class="w-full h-8 text-xs font-mono">
									{controller.newColumn.fk_table || 'Select Table'}
								</SelectTrigger>
								<SelectContent>
									{#each $tablesList as tbl}
										<SelectItem value={tbl.name} label={tbl.name} class="text-xs font-mono" />
									{/each}
								</SelectContent>
							</Select>
						</div>

						<div>
							<Label class="block text-[11px] text-muted-foreground mb-1">Target Column</Label>
							<Select
								type="single"
								bind:value={controller.newColumn.fk_column}
							>
								<SelectTrigger class="w-full h-8 text-xs font-mono">
									{controller.newColumn.fk_column || 'Select Column'}
								</SelectTrigger>
								<SelectContent>
									{#if controller.newColumn.fk_table && targetColumnsMap[controller.newColumn.fk_table]}
										{#each targetColumnsMap[controller.newColumn.fk_table] as targetCol}
											<SelectItem
												value={targetCol.name}
												label={`${targetCol.name}${targetCol.is_primary_key ? ' (PK)' : ''}`}
												class="text-xs font-mono"
											/>
										{/each}
									{:else}
										<SelectItem value="id" label="id (PK)" class="text-xs font-mono" />
									{/if}
								</SelectContent>
							</Select>
						</div>
					</div>
				</div>
			{/if}
		</form>

		<DialogFooter class="pt-2 border-t border-border">
			<Button type="button" variant="outline" size="sm" onclick={controller.closeAddColumnModal}>
				Cancel
			</Button>
			<Button type="submit" form="add-column-form" size="sm" disabled={controller.isSubmitting}>
				{controller.isSubmitting ? 'Adding Column...' : 'Add Column'}
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
