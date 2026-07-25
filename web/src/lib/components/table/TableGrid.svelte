<script lang="ts">
	import { fetchTableData, fetchTableSchema, deleteTableRow } from '../../../data/services';
	import { exportToCSV, exportToJSON } from '../../utils/exportUtils';
	import InsertRowModal from '../modals/InsertRowModal.svelte';
	import EditRowModal from '../modals/EditRowModal.svelte';
	import ImportModal from '../modals/ImportModal.svelte';
	import { useEditRow } from '../../hooks/useEditRow.svelte';
	import type { QueryResult, TableSchema } from '$lib/api';
	import { Button } from '../shadcn/button';
	import { Input } from '../shadcn/input';
	import {
		AlertDialog,
		AlertDialogContent,
		AlertDialogHeader,
		AlertDialogTitle,
		AlertDialogDescription,
		AlertDialogFooter,
		AlertDialogCancel,
		AlertDialogAction
	} from '../shadcn/alert-dialog';
	import { Search, Plus, Upload, Download, RefreshCw, Pencil, Trash2, Database, AlertTriangle } from '@lucide/svelte';

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
	let showImportModal = $state(false);

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
			return schema.columns.map((c: any) => c.name);
		}
		return [];
	});

	let schemaColumns = $derived(schema?.columns || []);
	let rowsList = $derived(dataResult?.rows || []);

	// Initialize Edit Row Hook controller
	const editRowController = useEditRow(
		() => tableName,
		() => schemaColumns
	);

	let filteredRows = $derived(
		rowsList.filter((row: any) =>
			Object.values(row || {}).some((val) =>
				String(val ?? '').toLowerCase().includes(searchQuery.toLowerCase())
			)
		)
	);

	let totalPages = $derived(Math.ceil(filteredRows.length / pageSize) || 1);

	let paginatedRows = $derived(
		filteredRows.slice((currentPage - 1) * pageSize, currentPage * pageSize)
	);

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

<div class="flex-1 flex flex-col h-full bg-background text-foreground overflow-hidden">
	<!-- Toolbar -->
	<div class="px-6 py-3 border-b border-border bg-card/40 flex flex-wrap items-center justify-between gap-4">
		<div class="flex items-center gap-3">
			<div class="relative w-56">
				<Search class="size-3.5 text-muted-foreground absolute left-2.5 top-2.5" />
				<Input
					type="text"
					bind:value={searchQuery}
					placeholder="Search rows..."
					class="h-8 pl-8 text-xs bg-background"
				/>
			</div>
			<span class="text-xs text-muted-foreground font-medium">
				{filteredRows.length} rows found
			</span>
		</div>

		<!-- Action Buttons -->
		<div class="flex items-center gap-2">
			<Button
				variant="default"
				size="sm"
				onclick={() => (showInsertModal = true)}
			>
				<Plus class="size-3.5 mr-1" />
				Add Row
			</Button>

			<Button
				variant="outline"
				size="sm"
				onclick={() => (showImportModal = true)}
			>
				<Upload class="size-3.5 mr-1" />
				Import
			</Button>

			<!-- Export Buttons -->
			<div class="flex items-center rounded-md border border-border bg-card text-xs">
				<button
					type="button"
					onclick={handleExportCSV}
					class="px-2.5 py-1 text-muted-foreground hover:text-foreground hover:bg-secondary transition rounded-l-md border-r border-border font-mono cursor-pointer flex items-center gap-1"
				>
					<Download class="size-3" /> CSV
				</button>
				<button
					type="button"
					onclick={handleExportJSON}
					class="px-2.5 py-1 text-muted-foreground hover:text-foreground hover:bg-secondary transition rounded-r-md font-mono cursor-pointer flex items-center gap-1"
				>
					<Download class="size-3" /> JSON
				</button>
			</div>

			<Button
				variant="ghost"
				size="sm"
				onclick={() => loadData(tableName)}
			>
				<RefreshCw class="size-3.5 mr-1" />
				Refresh
			</Button>
		</div>
	</div>

	<!-- Main Data Grid Container -->
	<div class="flex-1 overflow-auto p-6">
		{#if loading}
			<div class="h-64 flex flex-col items-center justify-center text-muted-foreground gap-2">
				<div class="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
				<span class="text-xs">Loading table data...</span>
			</div>
		{:else if errorMsg}
			<div class="p-4 rounded-xl bg-destructive/10 border border-destructive/20 text-destructive text-xs">
				❌ Error: {errorMsg}
			</div>
		{:else}
			<div class="border border-border rounded-xl overflow-hidden shadow-xs bg-card">
				<table class="w-full text-left text-xs border-collapse">
					<thead>
						<tr class="bg-secondary/60 text-muted-foreground font-semibold border-b border-border uppercase tracking-wider">
							<th class="px-4 py-3 border-r border-border/60 w-12 text-center">#</th>
							{#each columns() as col}
								<th class="px-4 py-3 border-r border-border/60 font-mono text-[11px] text-primary">
									{col}
								</th>
							{/each}
							<th class="px-4 py-3 w-20 text-center">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-border/60">
						{#if paginatedRows.length === 0}
							<tr>
								<td colspan={columns().length + 2} class="px-6 py-12 text-center text-muted-foreground">
									<div class="flex flex-col items-center justify-center gap-2">
										<Database class="size-8 text-muted-foreground/60" />
										<span class="font-medium text-foreground">Table "{tableName}" has 0 rows.</span>
										<span class="text-[11px] text-muted-foreground">Use "+ Add Row" or "Import" to populate records.</span>
									</div>
								</td>
							</tr>
						{:else}
							{#each paginatedRows as row, i}
								<tr class="hover:bg-secondary/40 transition-colors">
									<td class="px-4 py-2.5 border-r border-border/60 text-muted-foreground text-center font-mono">
										{(currentPage - 1) * pageSize + i + 1}
									</td>
									{#each columns() as col}
										<td class="px-4 py-2.5 border-r border-border/60 max-w-xs truncate font-mono text-foreground">
											{#if row[col] === null}
												<span class="text-muted-foreground/60 italic">null</span>
											{:else}
												{String(row[col] ?? '')}
											{/if}
										</td>
									{/each}
									<td class="px-4 py-2.5 text-center flex justify-center gap-1">
										<button
											type="button"
											onclick={() => editRowController.openModal(row)}
											class="p-1 text-muted-foreground hover:text-primary hover:bg-primary/10 rounded transition-all cursor-pointer"
											title="Edit row"
										>
											<Pencil class="size-3.5" />
										</button>
										<button
											type="button"
											onclick={() => openDeleteModal(row)}
											class="p-1 text-muted-foreground hover:text-destructive hover:bg-destructive/10 rounded transition-all cursor-pointer"
											title="Delete row"
										>
											<Trash2 class="size-3.5" />
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
		<div class="px-6 py-3 border-t border-border bg-card/40 flex items-center justify-between text-xs text-muted-foreground">
			<div>
				Showing Page <span class="font-bold text-foreground">{currentPage}</span> of <span class="font-bold text-foreground">{totalPages}</span>
			</div>
			<div class="flex items-center gap-2">
				<Button
					variant="outline"
					size="xs"
					disabled={currentPage === 1}
					onclick={() => currentPage--}
				>
					Previous
				</Button>
				<Button
					variant="outline"
					size="xs"
					disabled={currentPage >= totalPages}
					onclick={() => currentPage++}
				>
					Next
				</Button>
			</div>
		</div>
	{/if}
</div>

<!-- Modals with Shadcn Dialog & Two-Way State Binding -->
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
	controller={editRowController}
	onSuccess={() => loadData(tableName)}
/>

<ImportModal
	{tableName}
	bind:isOpen={showImportModal}
	onClose={() => (showImportModal = false)}
	onSuccess={() => loadData(tableName)}
/>

<!-- Delete Confirmation Shadcn AlertDialog -->
<AlertDialog bind:open={showDeleteModal}>
	<AlertDialogContent class="max-w-md">
		<AlertDialogHeader>
			<AlertDialogTitle class="flex items-center gap-2 text-foreground">
				<AlertTriangle class="size-5 text-warning" /> Confirm Delete Row
			</AlertDialogTitle>
			<AlertDialogDescription>
				Are you sure you want to delete this row from table <span class="font-bold text-foreground font-mono">"{tableName}"</span>?
			</AlertDialogDescription>
		</AlertDialogHeader>
		<AlertDialogFooter>
			<AlertDialogCancel onclick={() => (showDeleteModal = false)}>
				Cancel
			</AlertDialogCancel>
			<AlertDialogAction onclick={confirmDelete} disabled={deleting} class="bg-destructive hover:bg-destructive/90 text-destructive-foreground">
				{deleting ? 'Deleting...' : 'Delete Row'}
			</AlertDialogAction>
		</AlertDialogFooter>
	</AlertDialogContent>
</AlertDialog>
