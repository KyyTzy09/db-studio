<script lang="ts">
	import type { ColumnInfo } from '../../../data/models';
	import { insertTableRow } from '../../../data/services';
	import { Button } from '../shadcn/button';
	import { Input } from '../shadcn/input';
	import { Label } from '../shadcn/label';
	import { Plus, X } from '@lucide/svelte';

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

	$effect(() => {
		if (isOpen) {
			formData = {};
			errorMessage = '';
			columns.forEach((col: any) => {
				formData[col.name] = '';
			});
		}
	});

	function handleClose() {
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

{#if isOpen}
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
		<!-- Card Konten Modal -->
		<div
			role="document"
			tabindex="-1"
			onclick={(e) => e.stopPropagation()}
			onkeydown={(e) => e.stopPropagation()}
			class="w-full max-w-lg rounded-xl border border-border bg-card p-6 shadow-2xl cursor-default"
		>
			<!-- Header Modal -->
			<div class="flex items-center justify-between border-b border-border pb-4">
				<h3 class="text-base font-bold text-foreground flex items-center gap-1.5">
					<Plus class="size-4 text-success" /> Insert Row to <span class="text-primary font-mono">{tableName}</span>
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

				<div class="pt-4 flex justify-end gap-3 border-t border-border">
					<Button type="button" variant="outline" size="sm" onclick={handleClose}>
						Cancel
					</Button>
					<Button type="submit" size="sm" disabled={isSubmitting}>
						{isSubmitting ? 'Inserting...' : 'Insert Row'}
					</Button>
				</div>
			</form>
		</div>
	</div>
{/if}
