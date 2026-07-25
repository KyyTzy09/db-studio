<script lang="ts">
	import type { ColumnInfo } from '../../../data/models';
	import { updateTableRow } from '../../../data/services';
	import { Button } from '../shadcn/button';
	import { Input } from '../shadcn/input';
	import { Label } from '../shadcn/label';
	import { Pencil, X } from '@lucide/svelte';

	let {
		tableName,
		columns = [],
		rowData = null,
		isOpen = $bindable(false),
		onClose,
		onSuccess
	} = $props<{
		tableName: string;
		columns?: ColumnInfo[];
		rowData?: Record<string, any> | null;
		isOpen?: boolean;
		onClose?: () => void;
		onSuccess?: () => void;
	}>();

	let formData = $state<Record<string, any>>({});
	let pkData = $state<Record<string, any>>({});
	let isSubmitting = $state(false);
	let errorMessage = $state('');

	$effect(() => {
		if (isOpen && rowData) {
			errorMessage = '';
			formData = { ...rowData };
			pkData = {};

			columns.forEach((col: any) => {
				if (col.is_primary_key && rowData[col.name] !== undefined) {
					pkData[col.name] = rowData[col.name];
				}
			});

			if (Object.keys(pkData).length === 0) {
				pkData = { ...rowData };
			}
		}
	});

	function handleClose() {
		isOpen = false;
		if (onClose) onClose();
	}

	async function handleSubmit() {
		isSubmitting = true;
		errorMessage = '';
		try {
			const cleanData: Record<string, any> = {};
			for (const col of columns) {
				if (col.is_primary_key) continue;
				const val = formData[col.name];
				if (val !== '' && val !== null && val !== undefined) {
					cleanData[col.name] = val;
				}
			}

			await updateTableRow(tableName, pkData, cleanData);
			if (onSuccess) onSuccess();
			handleClose();
		} catch (err: any) {
			errorMessage = err.message || 'Failed to update row';
		} finally {
			isSubmitting = false;
		}
	}
</script>

{#if isOpen && rowData}
	<div
		onclick={(e) => {
			if (e.target === e.currentTarget) handleClose();
		}}
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4 cursor-pointer"
		role="dialog"
		aria-modal="true"
		tabindex="-1"
		onkeydown={(e) => {
			if (e.key === 'Escape') handleClose();
		}}
	>
		<div
			class="w-full max-w-lg rounded-xl border border-border bg-card p-6 shadow-2xl cursor-default"
			onclick={(e) => e.stopPropagation()}
			role="document"
			tabindex="-1"
		>
			<div class="flex items-center justify-between border-b border-border pb-4">
				<h3 class="text-base font-bold text-foreground flex items-center gap-1.5">
					<Pencil class="size-4 text-primary" /> Edit Row in <span class="text-primary font-mono">{tableName}</span>
				</h3>
				<Button type="button" variant="ghost" size="icon-xs" onclick={handleClose}>
					<X class="size-4" />
				</Button>
			</div>

			{#if errorMessage}
				<div class="mt-4 rounded-lg bg-destructive/10 border border-destructive/20 p-3 text-xs text-destructive">
					⚠️ {errorMessage}
				</div>
			{/if}

			<form
				onsubmit={(e) => {
					e.preventDefault();
					handleSubmit();
				}}
				class="mt-4 max-h-[60vh] overflow-y-auto pr-1 space-y-4 text-xs"
			>
				{#each columns as col}
					<div>
						<Label for={`edit-${col.name}`} class="block text-xs font-semibold text-foreground mb-1">
							{col.name}
							<span class="text-[10px] text-muted-foreground font-mono font-normal">({col.data_type})</span>
							{#if col.is_primary_key}
								<span class="ml-1 rounded bg-warning/20 px-1 py-0.5 text-[9px] text-warning font-bold">PK (Locked)</span>
							{/if}
						</Label>

						<Input
							id={`edit-${col.name}`}
							type="text"
							bind:value={formData[col.name]}
							disabled={col.is_primary_key}
							class="h-9 text-xs font-mono disabled:opacity-50"
						/>
					</div>
				{/each}

				<div class="pt-4 flex justify-end gap-3 border-t border-border">
					<Button type="button" variant="outline" size="sm" onclick={handleClose}>
						Cancel
					</Button>
					<Button type="submit" size="sm" disabled={isSubmitting}>
						{isSubmitting ? 'Updating...' : 'Save Changes'}
					</Button>
				</div>
			</form>
		</div>
	</div>
{/if}
