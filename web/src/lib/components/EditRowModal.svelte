<script lang="ts">
	import type { ColumnInfo } from '../api';
	import { updateTableRow } from '../api';

	let { tableName, columns = [], rowData = null, isOpen = false, onClose, onSuccess } = $props<{
		tableName: string;
		columns?: ColumnInfo[];
		rowData?: Record<string, any> | null;
		isOpen?: boolean;
		onClose: () => void;
		onSuccess: () => void;
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

			columns.forEach((col) => {
				if (col.is_primary_key && rowData[col.name] !== undefined) {
					pkData[col.name] = rowData[col.name];
				}
			});

			if (Object.keys(pkData).length === 0) {
				pkData = { ...rowData };
			}
		}
	});

	async function handleSubmit() {
		isSubmitting = true;
		errorMessage = '';
		try {
			// Clean update payload: omit primary key fields and empty strings
			const cleanData: Record<string, any> = {};
			for (const col of columns) {
				if (col.is_primary_key) continue; // Don't update PK in SET clause
				const val = formData[col.name];
				if (val !== '' && val !== null && val !== undefined) {
					cleanData[col.name] = val;
				}
			}

			await updateTableRow(tableName, pkData, cleanData);
			onSuccess();
			onClose();
		} catch (err: any) {
			errorMessage = err.message || 'Failed to update row';
		} finally {
			isSubmitting = false;
		}
	}
</script>

{#if isOpen && rowData}
	<!-- Modal Backdrop (Click outside to close) -->
	<div
		onclick={(e) => { if (e.target === e.currentTarget) onClose(); }}
		class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4 cursor-pointer"
		role="button"
		tabindex="-1"
		onkeydown={(e) => { if (e.key === 'Escape') onClose(); }}
	>
		<div class="w-full max-w-lg rounded-xl border border-slate-700 bg-slate-900 p-6 shadow-2xl cursor-default">
			<div class="flex items-center justify-between border-b border-slate-800 pb-4">
				<h3 class="text-lg font-bold text-slate-100">
					✏️ Edit Row in <span class="text-sky-400 font-mono">{tableName}</span>
				</h3>
				<button type="button" onclick={onClose} class="text-slate-400 hover:text-slate-200 text-xl font-bold">✕</button>
			</div>

			{#if errorMessage}
				<div class="mt-4 rounded-lg bg-rose-950/70 border border-rose-800/80 p-3 text-xs text-rose-300">
					⚠️ {errorMessage}
				</div>
			{/if}

			<form onsubmit={(e) => { e.preventDefault(); handleSubmit(); }} class="mt-4 max-h-[60vh] overflow-y-auto pr-1 space-y-4">
				{#each columns as col}
					<div>
						<label for={`edit-${col.name}`} class="block text-xs font-semibold text-slate-300 mb-1">
							{col.name}
							<span class="text-[10px] text-slate-500 font-mono font-normal">({col.data_type})</span>
							{#if col.is_primary_key}
								<span class="ml-1 rounded bg-amber-500/20 px-1 py-0.5 text-[9px] text-amber-300">PK (Locked)</span>
							{/if}
						</label>

						<input
							id={`edit-${col.name}`}
							type="text"
							bind:value={formData[col.name]}
							disabled={col.is_primary_key}
							class="w-full rounded-lg border border-slate-700 bg-slate-950 px-3 py-2 text-sm text-slate-100 font-mono focus:border-sky-500 focus:outline-none disabled:opacity-50 disabled:bg-slate-900/50"
						/>
					</div>
				{/each}

				<div class="pt-4 flex justify-end gap-3 border-t border-slate-800">
					<button
						type="button"
						onclick={onClose}
						class="rounded-lg px-4 py-2 text-xs font-medium text-slate-400 hover:bg-slate-800 hover:text-slate-200 transition"
					>
						Cancel
					</button>
					<button
						type="submit"
						disabled={isSubmitting}
						class="rounded-lg bg-sky-600 px-4 py-2 text-xs font-semibold text-white shadow-lg hover:bg-sky-500 disabled:opacity-50 transition"
					>
						{isSubmitting ? 'Updating...' : 'Save Changes'}
					</button>
				</div>
			</form>
		</div>
	</div>
{/if}
