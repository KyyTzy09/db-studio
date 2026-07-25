<script lang="ts">
	import { exportToCSV, exportToJSON } from '../../utils/exportUtils';
	import InsertRowModal from '../modals/InsertRowModal.svelte';
	import EditRowModal from '../modals/EditRowModal.svelte';
	import ImportModal from '../modals/ImportModal.svelte';
	import { useEditRow } from '../../hooks/useEditRow.svelte';
	import { useInsertRow } from '../../hooks/useInsertRow.svelte';
	import { useImportData } from '../../hooks/useImportData.svelte';
	import { useTableData, type FilterOperator } from '../../hooks/useTableData.svelte';
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
	import {
		Search,
		Plus,
		Upload,
		Download,
		RefreshCw,
		Pencil,
		Trash2,
		Database,
		AlertTriangle,
		Filter,
		FilterX,
		ArrowUp,
		ArrowDown,
		ArrowUpDown,
		X
	} from '@lucide/svelte';

	let { tableName } = $props<{ tableName: string }>();

	// Table Data Controller Hook
	const tableData = useTableData(() => tableName);

	$effect(() => {
		if (tableName) {
			tableData.loadData();
		}
	});

	let columns = $derived(() => {
		if (tableData.dataResult && tableData.dataResult.columns && tableData.dataResult.columns.length > 0) {
			return tableData.dataResult.columns;
		}
		if (tableData.schema && tableData.schema.columns && tableData.schema.columns.length > 0) {
			return tableData.schema.columns.map((c: any) => c.name);
		}
		return [];
	});

	let schemaColumns = $derived(tableData.schema?.columns || []);

	// Controllers for Modals
	const editRowController = useEditRow(
		() => tableName,
		() => schemaColumns
	);

	const insertRowController = useInsertRow(
		() => tableName,
		() => schemaColumns
	);

	const importDataController = useImportData(
		() => tableName
	);

	// Filter Form State
	let selectedCol = $state('');
	let selectedOp = $state<FilterOperator>('contains');
	let filterValue = $state('');

	$effect(() => {
		const cols = columns();
		if (cols.length > 0 && !selectedCol) {
			selectedCol = cols[0];
		}
	});

	function handleAddFilter() {
		if (!selectedCol) return;
		if (selectedOp !== 'is_null' && !filterValue.trim()) return;

		tableData.addFilterRule({
			column: selectedCol,
			operator: selectedOp,
			value: filterValue
		});

		filterValue = '';
	}

	function handleExportCSV() {
		exportToCSV(tableName, columns(), tableData.filteredRows);
	}

	function handleExportJSON() {
		exportToJSON(tableName, tableData.filteredRows);
	}
</script>

<div class="flex-1 flex flex-col h-full bg-background text-foreground overflow-hidden">
	<!-- Toolbar -->
	<div class="px-6 py-3 border-b border-border bg-card/40 flex flex-wrap items-center justify-between gap-4">
		<div class="flex items-center gap-3">
			<!-- Global Search Input -->
			<div class="relative w-56">
				<Search class="size-3.5 text-muted-foreground absolute left-2.5 top-2.5" />
				<Input
					type="text"
					bind:value={tableData.searchQuery}
					placeholder="Search rows..."
					class="h-8 pl-8 text-xs bg-background"
				/>
			</div>

			<!-- Filter Toggle Button with Badge -->
			<Button
				variant={tableData.showFilterPanel || tableData.filterRules.length > 0 ? 'default' : 'outline'}
				size="sm"
				onclick={() => (tableData.showFilterPanel = !tableData.showFilterPanel)}
				class="relative"
			>
				<Filter class="size-3.5 mr-1" />
				Filter
				{#if tableData.filterRules.length > 0}
					<span class="ml-1.5 rounded-full bg-primary-foreground text-primary font-bold px-1.5 py-0.2 text-[10px]">
						{tableData.filterRules.length}
					</span>
				{/if}
			</Button>

			<span class="text-xs text-muted-foreground font-medium">
				{tableData.filteredRows.length} rows found
			</span>
		</div>

		<!-- Action Buttons -->
		<div class="flex items-center gap-2">
			<Button
				variant="default"
				size="sm"
				onclick={() => insertRowController.openModal()}
			>
				<Plus class="size-3.5 mr-1" />
				Add Row
			</Button>

			<Button
				variant="outline"
				size="sm"
				onclick={() => importDataController.openModal()}
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
				onclick={() => tableData.loadData()}
			>
				<RefreshCw class="size-3.5 mr-1" />
				Refresh
			</Button>
		</div>
	</div>

	<!-- Collapsible Filter Builder Panel -->
	{#if tableData.showFilterPanel}
		<div class="px-6 py-3 border-b border-border bg-secondary/30 space-y-3 transition-all animate-in fade-in-50">
			<div class="flex flex-wrap items-center gap-2 text-xs">
				<span class="font-semibold text-foreground mr-1 flex items-center gap-1">
					<Filter class="size-3.5 text-primary" /> Filter Rules:
				</span>

				<!-- Column Selector -->
				<select
					bind:value={selectedCol}
					class="h-8 rounded-md border border-border bg-background px-2.5 text-xs font-mono text-foreground focus:outline-none focus:ring-1 focus:ring-primary"
				>
					{#each columns() as col}
						<option value={col}>{col}</option>
					{/each}
				</select>

				<!-- Operator Selector -->
				<select
					bind:value={selectedOp}
					class="h-8 rounded-md border border-border bg-background px-2.5 text-xs text-foreground focus:outline-none focus:ring-1 focus:ring-primary"
				>
					<option value="contains">contains</option>
					<option value="equals">equals (=)</option>
					<option value="gt">greater than (&gt;)</option>
					<option value="lt">less than (&lt;)</option>
					<option value="starts">starts with</option>
					<option value="is_null">is NULL / Empty</option>
				</select>

				<!-- Filter Value Input -->
				{#if selectedOp !== 'is_null'}
					<Input
						type="text"
						bind:value={filterValue}
						placeholder="Enter value..."
						class="h-8 w-44 text-xs font-mono bg-background"
						onkeydown={(e) => {
							if (e.key === 'Enter') handleAddFilter();
						}}
					/>
				{/if}

				<!-- Add Condition Button -->
				<Button type="button" variant="default" size="xs" onclick={handleAddFilter}>
					<Plus class="size-3 mr-1" /> Add Condition
				</Button>

				{#if tableData.filterRules.length > 0 || tableData.searchQuery}
					<Button type="button" variant="ghost" size="xs" onclick={() => tableData.clearAllFilters()} class="text-destructive hover:bg-destructive/10">
						<FilterX class="size-3 mr-1" /> Clear All
					</Button>
				{/if}
			</div>

			<!-- Active Filter Rule Chips -->
			{#if tableData.filterRules.length > 0}
				<div class="flex flex-wrap items-center gap-1.5 pt-1">
					{#each tableData.filterRules as rule}
						<div class="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-md bg-card border border-border text-xs font-mono shadow-2xs">
							<span class="text-primary font-bold">{rule.column}</span>
							<span class="text-muted-foreground text-[11px]">{rule.operator}</span>
							{#if rule.operator !== 'is_null'}
								<span class="text-foreground font-semibold">"{rule.value}"</span>
							{/if}
							<button
								type="button"
								onclick={() => tableData.removeFilterRule(rule.id)}
								class="text-muted-foreground hover:text-destructive transition ml-1 cursor-pointer"
								title="Remove filter rule"
							>
								<X class="size-3" />
							</button>
						</div>
					{/each}
				</div>
			{/if}
		</div>
	{/if}

	<!-- Main Data Grid Container -->
	<div class="flex-1 overflow-auto p-6">
		{#if tableData.loading}
			<div class="h-64 flex flex-col items-center justify-center text-muted-foreground gap-2">
				<div class="w-6 h-6 border-2 border-primary border-t-transparent rounded-full animate-spin"></div>
				<span class="text-xs">Loading table data...</span>
			</div>
		{:else if tableData.errorMsg}
			<div class="p-4 rounded-xl bg-destructive/10 border border-destructive/20 text-destructive text-xs">
				❌ Error: {tableData.errorMsg}
			</div>
		{:else}
			<div class="border border-border rounded-xl overflow-hidden shadow-xs bg-card">
				<table class="w-full text-left text-xs border-collapse">
					<thead>
						<tr class="bg-secondary/60 text-muted-foreground font-semibold border-b border-border uppercase tracking-wider">
							<th class="px-4 py-3 border-r border-border/60 w-12 text-center">#</th>
							{#each columns() as col}
								<th class="px-4 py-3 border-r border-border/60 font-mono text-[11px]">
									<button
										type="button"
										onclick={() => tableData.toggleSort(col)}
										class="flex items-center justify-between w-full hover:text-primary transition-colors cursor-pointer group"
									>
										<span class={tableData.sortState.column === col ? 'text-primary font-bold' : ''}>
											{col}
										</span>
										{#if tableData.sortState.column === col}
											{#if tableData.sortState.direction === 'asc'}
												<ArrowUp class="size-3.5 text-primary" />
											{:else if tableData.sortState.direction === 'desc'}
												<ArrowDown class="size-3.5 text-primary" />
											{/if}
										{:else}
											<ArrowUpDown class="size-3 text-muted-foreground/40 group-hover:text-primary transition-colors" />
										{/if}
									</button>
								</th>
							{/each}
							<th class="px-4 py-3 w-20 text-center">Actions</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-border/60">
						{#if tableData.paginatedRows.length === 0}
							<tr>
								<td colspan={columns().length + 2} class="px-6 py-12 text-center text-muted-foreground">
									<div class="flex flex-col items-center justify-center gap-2">
										<Database class="size-8 text-muted-foreground/60" />
										<span class="font-medium text-foreground">No rows found matching current filters.</span>
										<span class="text-[11px] text-muted-foreground">Try clearing filters or search keywords.</span>
									</div>
								</td>
							</tr>
						{:else}
							{#each tableData.paginatedRows as row, i}
								<tr class="hover:bg-secondary/40 transition-colors">
									<td class="px-4 py-2.5 border-r border-border/60 text-muted-foreground text-center font-mono">
										{(tableData.currentPage - 1) * tableData.pageSize + i + 1}
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
											onclick={() => tableData.openDeleteModal(row)}
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
	{#if columns().length > 0 && tableData.filteredRows.length > 0}
		<div class="px-6 py-3 border-t border-border bg-card/40 flex items-center justify-between text-xs text-muted-foreground">
			<div>
				Showing Page <span class="font-bold text-foreground">{tableData.currentPage}</span> of <span class="font-bold text-foreground">{tableData.totalPages}</span>
			</div>
			<div class="flex items-center gap-2">
				<Button
					variant="outline"
					size="xs"
					disabled={tableData.currentPage === 1}
					onclick={() => tableData.currentPage--}
				>
					Previous
				</Button>
				<Button
					variant="outline"
					size="xs"
					disabled={tableData.currentPage >= tableData.totalPages}
					onclick={() => tableData.currentPage++}
				>
					Next
				</Button>
			</div>
		</div>
	{/if}
</div>

<!-- Modals with Shadcn Dialog & Controller Props -->
<InsertRowModal
	{tableName}
	columns={schemaColumns}
	controller={insertRowController}
	onSuccess={() => tableData.loadData()}
/>

<EditRowModal
	{tableName}
	columns={schemaColumns}
	controller={editRowController}
	onSuccess={() => tableData.loadData()}
/>

<ImportModal
	{tableName}
	controller={importDataController}
	onSuccess={() => tableData.loadData()}
/>

<!-- Delete Confirmation Shadcn AlertDialog -->
<AlertDialog bind:open={tableData.showDeleteModal}>
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
			<AlertDialogCancel onclick={() => tableData.closeDeleteModal()}>
				Cancel
			</AlertDialogCancel>
			<AlertDialogAction onclick={() => tableData.confirmDelete()} disabled={tableData.deleting} class="bg-destructive hover:bg-destructive/90 text-destructive-foreground">
				{tableData.deleting ? 'Deleting...' : 'Delete Row'}
			</AlertDialogAction>
		</AlertDialogFooter>
	</AlertDialogContent>
</AlertDialog>
