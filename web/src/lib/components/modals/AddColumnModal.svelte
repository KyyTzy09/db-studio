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
	import { Plus } from '@lucide/svelte';

	let { controller } = $props<{
		controller: ReturnType<typeof useSchemaEditor>;
	}>();

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
					placeholder="e.g. bio, is_active, total_price"
					class="h-9 text-xs font-mono"
				/>
			</div>

			<div>
				<Label for="add-col-datatype" class="block text-xs font-semibold text-foreground mb-1">
					Data Type
				</Label>
				<select
					id="add-col-datatype"
					bind:value={controller.newColumn.data_type}
					class="h-9 w-full rounded-md border border-border bg-background px-3 text-xs font-mono focus:outline-hidden focus:ring-1 focus:ring-primary"
				>
					{#each commonDataTypes as dt}
						<option value={dt}>{dt}</option>
					{/each}
				</select>
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
				<label class="flex items-center gap-2 cursor-pointer text-xs">
					<input
						type="checkbox"
						bind:checked={controller.newColumn.is_nullable}
						class="rounded border-border accent-primary size-4"
					/>
					<span>Allow NULL</span>
				</label>

				<label class="flex items-center gap-2 cursor-pointer text-xs">
					<input
						type="checkbox"
						bind:checked={controller.newColumn.is_primary_key}
						class="rounded border-border accent-primary size-4"
					/>
					<span>Primary Key</span>
				</label>
			</div>
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
