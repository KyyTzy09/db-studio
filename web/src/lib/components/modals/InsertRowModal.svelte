<script lang="ts">
	import type { ColumnInfo } from '../../../data/models';
	import { insertTableRow } from '../../../data/services';
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

	let {
		tableName,
		columns = [],
		isOpen = $bindable(false),
		onClose,
		onSuccess
	} = $props<{
		tableName: string;
		columns?: ColumnInfo[];
		isOpen: boolean;
		onClose?: () => void;
		onSuccess?: () => void;
	}>();

	let formData = $state<Record<string, any>>({});
	let isSubmitting = $state(false);
	let errorMessage = $state('');

	function handleOpenChange(open: boolean) {
		if (open) {
			errorMessage = '';
			const initial: Record<string, any> = {};
			columns.forEach((col: any) => {
				initial[col.name] = '';
			});
			formData = initial;
		} else {
			handleClose();
		}
	}

	function handleClose() {
		isOpen = false;
		if (onClose) onClose();
	}

	async function handleSubmit() {
		isSubmitting = true;
		errorMessage = '';
		try {
			const payload: Record<string, any> = {};
			for (const col of columns) {
				const val = formData[col.name];
				if (val !== '' && val !== null && val !== undefined) {
					payload[col.name] = val;
				}
			}

			await insertTableRow(tableName, payload);
			if (onSuccess) onSuccess();
			handleClose();
		} catch (err: any) {
			errorMessage = err.message || 'Failed to insert row';
		} finally {
			isSubmitting = false;
		}
	}
</script>

<Dialog bind:open={isOpen} onOpenChange={handleOpenChange}>
	<DialogContent class="max-w-lg">
		<DialogHeader>
			<DialogTitle class="flex items-center gap-1.5 text-base font-bold">
				<Plus class="size-4 text-success" /> Insert Row to <span class="text-primary font-mono">{tableName}</span>
			</DialogTitle>
			<DialogDescription>
				Fill in column values to insert a new record into database.
			</DialogDescription>
		</DialogHeader>

		{#if errorMessage}
			<div class="rounded-lg bg-destructive/10 border border-destructive/20 p-3 text-xs text-destructive">
				⚠️ {errorMessage}
			</div>
		{/if}

		<form
			id="insert-row-form"
			onsubmit={(e) => {
				e.preventDefault();
				handleSubmit();
			}}
			class="max-h-[60vh] overflow-y-auto pr-1 space-y-4 text-xs py-2"
		>
			{#each columns as col}
				<div>
					<Label
						for={`insert-${col.name}`}
						class="block text-xs font-semibold text-foreground mb-1"
					>
						{col.name}
						<span class="text-[10px] text-muted-foreground font-mono font-normal">({col.data_type})</span>
						{#if col.is_primary_key}
							<span class="ml-1 rounded bg-warning/20 px-1 py-0.5 text-[9px] text-warning font-bold">PK (Auto-ID)</span>
						{/if}
						{#if !col.is_nullable && !col.is_primary_key && !col.default_value}
							<span class="text-destructive">*</span>
						{/if}
					</Label>

					<Input
						id={`insert-${col.name}`}
						type="text"
						bind:value={formData[col.name]}
						placeholder={col.is_primary_key
							? 'Auto-generated / optional'
							: col.default_value
								? `Default: ${col.default_value}`
								: col.is_nullable
									? 'Leave empty for NULL'
									: 'Enter value...'}
						class="h-9 text-xs font-mono"
					/>
				</div>
			{/each}
		</form>

		<DialogFooter class="pt-2 border-t border-border">
			<Button type="button" variant="outline" size="sm" onclick={handleClose}>
				Cancel
			</Button>
			<Button type="submit" form="insert-row-form" size="sm" disabled={isSubmitting}>
				{isSubmitting ? 'Inserting...' : 'Insert Row'}
			</Button>
		</DialogFooter>
	</DialogContent>
</Dialog>
