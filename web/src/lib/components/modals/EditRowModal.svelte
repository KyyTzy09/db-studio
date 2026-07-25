<script lang="ts">
	import type { ColumnInfo } from '../../../data/models';
	import type { useEditRow } from '../../hooks/useEditRow.svelte';
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
	import { Pencil } from '@lucide/svelte';

	let {
		tableName,
		columns = [],
		controller,
		onSuccess
	} = $props<{
		tableName: string;
		columns?: ColumnInfo[];
		controller: ReturnType<typeof useEditRow>;
		onSuccess?: () => void;
	}>();
</script>

<Dialog bind:open={controller.isOpen}>
	<DialogContent class="max-w-lg">
		<DialogHeader>
			<DialogTitle class="flex items-center gap-1.5 text-base font-bold">
				<Pencil class="size-4 text-primary" /> Edit Row in <span class="text-primary font-mono">{tableName}</span>
			</DialogTitle>
			<DialogDescription>
				Update column values for this record. Primary Keys are locked.
			</DialogDescription>
		</DialogHeader>

		{#if controller.errorMessage}
			<div class="rounded-lg bg-destructive/10 border border-destructive/20 p-3 text-xs text-destructive">
				⚠️ {controller.errorMessage}
			</div>
		{/if}

		<form
			id="edit-row-form"
			onsubmit={(e) => {
				e.preventDefault();
				controller.submitForm(onSuccess);
			}}
			class="max-h-[60vh] overflow-y-auto pr-1 space-y-4 text-xs py-2"
		>
			{#each columns as col}
				<div>
					<Label for={`edit-${col.name}`} class="block text-xs font-semibold text-foreground mb-1">
						{col.name}
						<span class="text-[10px] text-muted-foreground font-mono font-normal">({col.data_type})</span>
						{#if col.is_primary_key}
							<span class="ml-1 rounded bg-warning/20 px-1 py-0.5 text-[9px] text-warning font-bold">PK (Locked)</span>
						{/if}
						{#if col.is_foreign_key}
							<span class="ml-1 rounded bg-primary/20 px-1 py-0.5 text-[9px] text-primary font-bold">FK</span>
						{/if}
					</Label>

					<Input
						id={`edit-${col.name}`}
						type="text"
						bind:value={controller.formData[col.name]}
						disabled={col.is_primary_key}
						class="h-9 text-xs font-mono disabled:opacity-50"
					/>
				</div>
			{/each}
		</form>

		<DialogFooter class="pt-2 border-t border-border">
			<Button type="button" variant="outline" size="sm" onclick={controller.closeModal}>
				Cancel
			</Button>
			<Button type="submit" form="edit-row-form" size="sm" disabled={controller.isSubmitting}>
				{controller.isSubmitting ? 'Updating...' : 'Save Changes'}
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
