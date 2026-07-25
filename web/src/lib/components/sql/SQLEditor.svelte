<script lang="ts">
	import { onMount } from 'svelte';
	import { EditorView, basicSetup } from 'codemirror';
	import { sql } from '@codemirror/lang-sql';
	import { oneDark } from '@codemirror/theme-one-dark';
	import { executeRawQuery, fetchTables, type QueryResult } from '../../../data/services';
	import { exportToCSV, exportToJSON } from '../../utils/exportUtils';
	import { tablesList } from '../../stores/dbStore';
	import { Button } from '../shadcn/button';
	import { Play, Terminal, Download, AlertTriangle } from '@lucide/svelte';

	let editorContainer = $state<HTMLDivElement | null>(null);
	let editorView = $state<EditorView | null>(null);

	let queryInput = $state('SELECT * FROM information_schema.tables LIMIT 10;');
	let queryResult = $state<QueryResult | null>(null);
	let executing = $state(false);
	let errorMsg = $state<string | null>(null);

	// Danger Zone Modal State
	let showDangerModal = $state(false);
	let dangerWarning = $state('');

	onMount(() => {
		if (editorContainer) {
			editorView = new EditorView({
				doc: queryInput,
				extensions: [
					basicSetup,
					sql(),
					oneDark,
					EditorView.updateListener.of((update) => {
						if (update.docChanged) {
							queryInput = update.state.doc.toString();
						}
					})
				],
				parent: editorContainer
			});
		}

		return () => {
			editorView?.destroy();
		};
	});

	async function runQuery(force = false) {
		if (!queryInput.trim()) return;
		executing = true;
		errorMsg = null;

		try {
			const res = await executeRawQuery(queryInput, force);
			if (res.is428 && res.warning) {
				dangerWarning = res.warning;
				showDangerModal = true;
				return;
			}
			if (res.data) {
				queryResult = res.data;
				try {
					const { tables } = await fetchTables();
					tablesList.set(tables || []);
				} catch (e) {
					// Ignore background refresh errors
				}
			}
		} catch (err: any) {
			errorMsg = err.message || 'Failed to execute SQL query';
		} finally {
			executing = false;
		}
	}

	function confirmDangerousQuery() {
		showDangerModal = false;
		runQuery(true);
	}

	function handleExportCSV() {
		if (queryResult) {
			exportToCSV('query_result', queryResult.columns, queryResult.rows);
		}
	}

	function handleExportJSON() {
		if (queryResult) {
			exportToJSON('query_result', queryResult.rows);
		}
	}
</script>

<div class="flex-1 flex flex-col h-full bg-background text-foreground overflow-hidden">
	<!-- Editor Header & Toolbar -->
	<div class="px-6 py-3 border-b border-border bg-card/40 flex items-center justify-between">
		<div class="flex items-center gap-2">
			<Terminal class="size-4 text-primary" />
			<h2 class="text-xs font-bold text-foreground uppercase tracking-wider">SQL Workspace</h2>
		</div>

		<Button
			variant="default"
			size="sm"
			disabled={executing}
			onclick={() => runQuery(false)}
		>
			{#if executing}
				<div class="w-3.5 h-3.5 border-2 border-primary-foreground border-t-transparent rounded-full animate-spin mr-1.5"></div>
				<span>Running...</span>
			{:else}
				<Play class="size-3.5 mr-1.5 fill-current" />
				<span>Execute Query</span>
			{/if}
		</Button>
	</div>

	<!-- CodeMirror Editor Box -->
	<div class="h-56 border-b border-border relative bg-[#282c34] font-mono text-xs">
		<div bind:this={editorContainer} class="h-full w-full overflow-auto"></div>
	</div>

	<!-- Query Execution Results -->
	<div class="flex-1 flex flex-col overflow-hidden bg-background">
		<div class="px-6 py-2 border-b border-border bg-card/40 flex items-center justify-between text-xs text-muted-foreground">
			<span>Execution Results</span>
			{#if queryResult}
				<div class="flex items-center gap-3">
					<span class="text-success font-mono text-[11px] font-medium">
						Returned {queryResult.rows.length} rows ({queryResult.execution_ms} ms)
					</span>
					<div class="flex items-center rounded-md border border-border bg-card text-[11px]">
						<button
							type="button"
							onclick={handleExportCSV}
							class="px-2 py-0.5 text-muted-foreground hover:text-foreground hover:bg-secondary transition rounded-l-md border-r border-border font-mono cursor-pointer flex items-center gap-1"
						>
							<Download class="size-3" /> CSV
						</button>
						<button
							type="button"
							onclick={handleExportJSON}
							class="px-2 py-0.5 text-muted-foreground hover:text-foreground hover:bg-secondary transition rounded-r-md font-mono cursor-pointer flex items-center gap-1"
						>
							<Download class="size-3" /> JSON
						</button>
					</div>
				</div>
			{/if}
		</div>

		<div class="flex-1 overflow-auto p-6">
			{#if errorMsg}
				<div class="p-4 rounded-xl bg-destructive/10 border border-destructive/20 text-destructive text-xs font-mono">
					❌ {errorMsg}
				</div>
			{:else if !queryResult}
				<div class="h-40 flex flex-col items-center justify-center text-muted-foreground text-xs">
					Run a SQL query to see formatted tabular output.
				</div>
			{:else if queryResult.rows.length === 0}
				<div class="p-4 rounded-xl bg-card border border-border text-muted-foreground text-xs text-center">
					Query executed successfully. 0 rows returned.
				</div>
			{:else}
				<div class="border border-border rounded-xl overflow-hidden shadow-xs bg-card">
					<table class="w-full text-left text-xs border-collapse font-mono">
						<thead>
							<tr class="bg-secondary/60 text-muted-foreground font-semibold border-b border-border uppercase tracking-wider">
								{#each queryResult.columns as col}
									<th class="px-4 py-3 border-r border-border/60 text-primary">{col}</th>
								{/each}
							</tr>
						</thead>
						<tbody class="divide-y divide-border/60">
							{#each queryResult.rows as row}
								<tr class="hover:bg-secondary/40 transition-colors">
									{#each queryResult.columns as col}
										<td class="px-4 py-2.5 border-r border-border/60 max-w-xs truncate text-foreground">
											{row[col] ?? 'null'}
										</td>
									{/each}
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			{/if}
		</div>
	</div>
</div>

<!-- Danger Zone Safety Warning Modal (HTTP 428) -->
{#if showDangerModal}
	<div class="fixed inset-0 bg-black/60 backdrop-blur-xs flex items-center justify-center z-50 p-4">
		<div class="bg-card border border-warning/40 rounded-xl p-6 max-w-md w-full shadow-2xl">
			<div class="flex items-center gap-3 mb-3">
				<div class="w-10 h-10 rounded-full bg-warning/10 border border-warning/30 flex items-center justify-center text-warning font-bold">
					<AlertTriangle class="size-5" />
				</div>
				<div>
					<h3 class="text-base font-bold text-foreground">Danger Zone Warning</h3>
					<p class="text-xs text-warning font-medium">Potentially Destructive Query Detected</p>
				</div>
			</div>

			<p class="text-xs text-muted-foreground mb-4 leading-relaxed bg-background p-3 rounded-lg border border-border">
				{dangerWarning}
			</p>

			<div class="flex justify-end gap-3">
				<Button variant="outline" size="sm" onclick={() => (showDangerModal = false)}>
					Cancel
				</Button>
				<Button variant="destructive" size="sm" onclick={confirmDangerousQuery}>
					Yes, Execute Query
				</Button>
			</div>
		</div>
	</div>
{/if}
