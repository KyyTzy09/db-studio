<script lang="ts">
	import { fetchTableData, fetchTableSchema, deleteTableRow, type QueryResult, type TableSchema } from '$lib/api';
	import { exportToCSV, exportToJSON } from '$lib/utils/exportUtils';
	import InsertRowModal from './InsertRowModal.svelte';
	import EditRowModal from './EditRowModal.svelte';
	import ImportModal from './ImportModal.svelte';

	let { tableName } = $props<{ tableName: string }>();

	let dataResult = $state<QueryResult | null>(null);
	let schema = $state<TableSchema | null>(null);
	let loading = $state(true);
	let errorMsg = $state<string | null>(null);
	let searchQuery = $state('');

	// Pagination
	let currentPage = $state(1);
	let pageSize = $state(15);

	// Modals State
	let showInsertModal = $state(false);
	let showEditModal = $state(false);
	let showImportModal = $state(false);
	let selectedRowToEdit = $state<Record<string, any> | null>(null);

	// Delete Confirmation Modal
	let showDeleteModal = $state(false);
	let selectedRowToDelete = $state<Record<string, any> | null>(null);
	let deleting = $state(false);

	$effect(() => {
		if (tableName) {
			loadData(tableName);
		}
	});

	async function loadData(name: string) {
		loading = true;
		errorMsg = null;
		currentPage = 1;
		try {
			const [dataRes, schemaRes] = await Promise.all([
				fetchTableData(name).catch(() => null),
				fetchTableSchema(name).catch(() => null)
			]);
			dataResult = dataRes;
			schema = schemaRes;
		} catch (err: any) {
			errorMsg = err.message || 'Failed to load table data';
		} finally {
			loading = false;
		}
	}

	let columns = $derived(() => {
		if (dataResult && dataResult.columns && dataResult.columns.length > 0) {
			return dataResult.columns;
		}
		if (schema && schema.columns && schema.columns.length > 0) {
			return schema.columns.map((c) => c.name);
		}
		return [];
	});

	let schemaColumns = $derived(schema?.columns || []);

	let rowsList = $derived(dataResult?.rows || []);

	let filteredRows = $derived(
		rowsList.filter((row) =>
			Object.values(row || {}).some((val) =>
				String(val ?? '').toLowerCase().includes(searchQuery.toLowerCase())
			)
		)
	);

	let totalPages = $derived(Math.ceil(filteredRows.length / pageSize) || 1);

	let paginatedRows = $derived(
		filteredRows.slice((currentPage - 1) * pageSize, currentPage * pageSize)
	);

	function openEditModal(row: Record<string, any>) {
		selectedRowToEdit = row;
		showEditModal = true;
	}

	function openDeleteModal(row: Record<string, any>) {
		selectedRowToDelete = row;
		showDeleteModal = true;
	}

	async function confirmDelete() {
		if (!selectedRowToDelete || !dataResult) return;
		deleting = true;
		try {
			await deleteTableRow(tableName, selectedRowToDelete);
			showDeleteModal = false;
			await loadData(tableName);
		} catch (err: any) {
			alert(err.message || 'Failed to delete row');
		} finally {
			deleting = false;
		}
	}

	function handleExportCSV() {
		exportToCSV(tableName, columns(), filteredRows);
	}

	function handleExportJSON() {
		exportToJSON(tableName, filteredRows);
	}
</script>

<div class="flex-1 flex flex-col h-full bg-slate-950 text-slate-200 overflow-hidden">
	<!-- Toolbar -->
	<div class="px-6 py-3 border-b border-slate-800 bg-slate-900/50 flex flex-wrap items-center justify-between gap-4">
		<div class="flex items-center gap-3">
			<div class="relative w-56">
				<input
					type="text"
					bind:value={searchQuery}
					placeholder="Search rows..."
					class="w-full bg-slate-900 border border-slate-700/60 text-xs rounded-lg pl-8 pr-3 py-1.5 focus:outline-none focus:border-indigo-500 text-slate-200 placeholder:text-slate-500"
				/>
				<svg class="w-3.5 h-3.5 text-slate-500 absolute left-2.5 top-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
				</svg>
			</div>
			<span class="text-xs text-slate-400 font-medium">
				{filteredRows.length} rows found
			</span>
		</div>

		<!-- Actions (Insert, Import, Export, Refresh) -->
		<div class="flex items-center gap-2">
			<button
				onclick={() => (showInsertModal = true)}
				class="px-3 py-1.5 rounded-lg bg-emerald-600 hover:bg-emerald-500 text-xs font-semibold text-white shadow-lg transition flex items-center gap-1 cursor-pointer"
			>
				➕ Add Row
			</button>

			<button
				onclick={() => (showImportModal = true)}
				class="px-3 py-1.5 rounded-lg bg-purple-600 hover:bg-purple-500 text-xs font-semibold text-white shadow-lg transition flex items-center gap-1 cursor-pointer"
			>
				📥 Import Data
			</button>

			<!-- Export Dropdown / Buttons -->
			<div class="flex items-center rounded-lg border border-slate-700 bg-slate-900 text-xs">
				<button
					onclick={handleExportCSV}
					class="px-2.5 py-1.5 text-slate-300 hover:text-white hover:bg-slate-800 transition rounded-l-lg border-r border-slate-700 font-mono cursor-pointer"
				>
					CSV ⬇
				</button>
				<button
					onclick={handleExportJSON}
					class="px-2.5 py-1.5 text-slate-300 hover:text-white hover:bg-slate-800 transition rounded-r-lg font-mono cursor-pointer"
				>
					JSON ⬇
				</button>
			</div>

			<button
				onclick={() => loadData(tableName)}
				class="px-3 py-1.5 rounded-lg bg-slate-800 hover:bg-slate-700 text-xs font-medium text-slate-300 flex items-center gap-1.5 transition-colors cursor-pointer"
			>
				<svg class="w-3.5 h-3.5 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
				</svg>
				Refresh
			</button>
		</div>
	</div>

	<!-- Main Table Data Grid Container -->
	<div class="flex-1 overflow-auto custom-scrollbar p-6">
		{#if loading}
			<div class="h-64 flex flex-col items-center justify-center text-slate-400 gap-2">
				<div class="w-6 h-6 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
				<span class="text-xs">Loading table data...</span>
			</div>
		{:else if errorMsg}
			<div class="p-4 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-xs">
				❌ Error: {errorMsg}
			</div>
		{:else}
			<div class="border border-slate-800 rounded-xl overflow-hidden shadow-2xl bg-slate-900/60">
				<table class="w-full text-left text-xs border-collapse">
					<thead>
						<tr class="bg-slate-900/90 text-slate-400 font-semibold border-b border-slate-800 uppercase tracking-wider">
							<th class="px-4 py-3 border-r border-slate-800/80 w-12 text-center">#</th>
							{#each columns() as col}
								<th class="px-4 py-3 border-r border-slate-800/80 font-mono text-[11px] text-indigo-300">
									{col}
								</th>
							{/each}
							<th class="px-4 py-3 w-20 text-center">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-800/60">
						{#if paginatedRows.length === 0}
							<tr>
								<td colspan={columns().length + 2} class="px-6 py-12 text-center text-slate-500">
									<div class="flex flex-col items-center justify-center gap-2">
										<svg class="w-8 h-8 text-slate-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
											<path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
										</svg>
										<span class="font-medium text-slate-400">Table "{tableName}" has no rows (0 rows).</span>
										<span class="text-[11px] text-slate-600">Structure is displayed above. Use "+ Add Row" or "Import Data" to populate rows.</span>
									</div>
								</td>
							</tr>
						{:else}
							{#each paginatedRows as row, i}
								<tr class="hover:bg-slate-800/40 transition-colors group">
									<td class="px-4 py-2.5 border-r border-slate-800/60 text-slate-500 text-center font-mono">
										{(currentPage - 1) * pageSize + i + 1}
									</td>
									{#each columns() as col}
										<td class="px-4 py-2.5 border-r border-slate-800/60 max-w-xs truncate font-mono text-slate-300">
											{#if row[col] === null}
												<span class="text-slate-600 italic">null</span>
											{:else}
												{String(row[col] ?? '')}
											{/if}
										</td>
									{/each}
									<td class="px-4 py-2.5 text-center flex justify-center gap-1">
										<button
											onclick={() => openEditModal(row)}
											class="p-1 text-slate-400 hover:text-sky-400 hover:bg-sky-500/10 rounded transition-all cursor-pointer"
											title="Edit row"
										>
											✏️
										</button>
										<button
											onclick={() => openDeleteModal(row)}
											class="p-1 text-slate-400 hover:text-red-400 hover:bg-red-500/10 rounded transition-all cursor-pointer"
											title="Delete row"
										>
											🗑️
										</button>
									</td>
								</tr>
							{/each}
						{/if}
					</tbody>
				</table>
			</div>
		{/if}
	</div>

	<!-- Pagination Controls -->
	{#if columns().length > 0 && filteredRows.length > 0}
		<div class="px-6 py-3 border-t border-slate-800 bg-slate-900/50 flex items-center justify-between text-xs text-slate-400">
			<div>
				Showing Page <span class="font-bold text-white">{currentPage}</span> of <span class="font-bold text-white">{totalPages}</span>
			</div>
			<div class="flex items-center gap-2">
				<button
					disabled={currentPage === 1}
					onclick={() => currentPage--}
					class="px-3 py-1 rounded bg-slate-800 hover:bg-slate-700 disabled:opacity-40 text-slate-200 transition-colors"
				>
					Previous
				</button>
				<button
					disabled={currentPage >= totalPages}
					onclick={() => currentPage++}
					class="px-3 py-1 rounded bg-slate-800 hover:bg-slate-700 disabled:opacity-40 text-slate-200 transition-colors"
				>
					Next
				</button>
			</div>
		</div>
	{/if}
</div>

<!-- Modals -->
<InsertRowModal
	{tableName}
	columns={schemaColumns}
	bind:isOpen={showInsertModal}
	onClose={() => (showInsertModal = false)}
	onSuccess={() => loadData(tableName)}
/>

<EditRowModal
	{tableName}
	columns={schemaColumns}
	rowData={selectedRowToEdit}
	bind:isOpen={showEditModal}
	onClose={() => (showEditModal = false)}
	onSuccess={() => loadData(tableName)}
/>

<ImportModal
	{tableName}
	bind:isOpen={showImportModal}
	onClose={() => (showImportModal = false)}
	onSuccess={() => loadData(tableName)}
/>

<!-- Delete Confirmation Modal -->
{#if showDeleteModal}
	<div class="fixed inset-0 bg-slate-950/80 backdrop-blur-xs flex items-center justify-center z-50 p-4 cursor-pointer" onclick={(e) => { if (e.target === e.currentTarget) showDeleteModal = false; }} role="dialog" tabindex="-1">
		<div class="bg-slate-900 border border-slate-800 rounded-xl p-6 max-w-md w-full shadow-2xl cursor-default" onclick={(e) => e.stopPropagation()} role="document" tabindex="-1">
			<h3 class="text-base font-bold text-white mb-2 flex items-center gap-2">
				⚠️ Confirm Delete Row
			</h3>
			<p class="text-xs text-slate-400 mb-4 leading-relaxed">
				Are you sure you want to delete this row from table <span class="font-bold text-white">"{tableName}"</span>?
			</p>
			<div class="flex justify-end gap-3">
				<button
					onclick={() => (showDeleteModal = false)}
					class="px-4 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-xs font-medium text-slate-300 transition-colors cursor-pointer"
				>
					Cancel
				</button>
				<button
					onclick={confirmDelete}
					disabled={deleting}
					class="px-4 py-2 rounded-lg bg-red-600 hover:bg-red-500 text-xs font-semibold text-white transition-colors cursor-pointer"
				>
					{deleting ? 'Deleting...' : 'Delete Row'}
				</button>
			</div>
		</div>
	</div>
{/if}
