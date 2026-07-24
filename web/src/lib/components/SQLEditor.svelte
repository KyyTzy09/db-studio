<script lang="ts">
	import { onMount } from 'svelte';
	import { EditorView, basicSetup } from 'codemirror';
	import { sql } from '@codemirror/lang-sql';
	import { oneDark } from '@codemirror/theme-one-dark';
	import { executeRawQuery, fetchTables, type QueryResult } from '$lib/api';
	import { tablesList } from '$lib/stores/dbStore';

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
				// Auto-refresh tables list if a DDL query (CREATE, DROP, ALTER) was executed
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
</script>

<div class="flex-1 flex flex-col h-full bg-slate-950 text-slate-200 overflow-hidden">
	<!-- Editor Header & Toolbar -->
	<div class="px-6 py-3 border-b border-slate-800 bg-slate-900/50 flex items-center justify-between">
		<div class="flex items-center gap-2">
			<svg class="w-4 h-4 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
				<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l3 3-3 3m5 0h3M5 20h14a2 2 0 002-2V6a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
			</svg>
			<h2 class="text-xs font-bold text-white uppercase tracking-wider">SQL Workspace</h2>
		</div>

		<button
			onclick={() => runQuery(false)}
			disabled={executing}
			class="px-4 py-1.5 rounded-lg bg-indigo-600 hover:bg-indigo-500 text-xs font-bold text-white shadow-lg shadow-indigo-500/20 flex items-center gap-2 transition-all cursor-pointer disabled:opacity-50"
		>
			{#if executing}
				<div class="w-3.5 h-3.5 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
				<span>Running...</span>
			{:else}
				<svg class="w-3.5 h-3.5 fill-current" viewBox="0 0 24 24">
					<path d="M8 5v14l11-7z" />
				</svg>
				<span>Execute Query</span>
			{/if}
		</button>
	</div>

	<!-- CodeMirror Editor Box -->
	<div class="h-56 border-b border-slate-800 relative bg-[#282c34] font-mono text-xs">
		<div bind:this={editorContainer} class="h-full w-full overflow-auto"></div>
	</div>

	<!-- Query Execution Results -->
	<div class="flex-1 flex flex-col overflow-hidden bg-slate-950">
		<div class="px-6 py-2 border-b border-slate-800 bg-slate-900/30 flex items-center justify-between text-xs text-slate-400">
			<span>Execution Results</span>
			{#if queryResult}
				<span class="text-emerald-400 font-mono text-[11px]">
					Returned {queryResult.rows.length} rows ({queryResult.execution_ms} ms)
				</span>
			{/if}
		</div>

		<div class="flex-1 overflow-auto custom-scrollbar p-6">
			{#if errorMsg}
				<div class="p-4 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-xs font-mono">
					❌ {errorMsg}
				</div>
			{:else if !queryResult}
				<div class="h-40 flex flex-col items-center justify-center text-slate-400 text-xs">
					Run a SQL query to see formatted tabular output.
				</div>
			{:else if queryResult.rows.length === 0}
				<div class="p-4 rounded-xl bg-slate-900 border border-slate-800 text-slate-400 text-xs text-center">
					Query executed successfully. 0 rows returned.
				</div>
			{:else}
				<div class="border border-slate-800 rounded-xl overflow-hidden shadow-2xl bg-slate-900/60">
					<table class="w-full text-left text-xs border-collapse font-mono">
						<thead>
							<tr class="bg-slate-900/90 text-slate-400 font-semibold border-b border-slate-800 uppercase tracking-wider">
								{#each queryResult.columns as col}
									<th class="px-4 py-3 border-r border-slate-800/80 text-indigo-300">{col}</th>
								{/each}
							</tr>
						</thead>
						<tbody class="divide-y divide-slate-800/60">
							{#each queryResult.rows as row}
								<tr class="hover:bg-slate-800/40 transition-colors">
									{#each queryResult.columns as col}
										<td class="px-4 py-2.5 border-r border-slate-800/60 max-w-xs truncate text-slate-300">
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
	<div class="fixed inset-0 bg-slate-950/80 backdrop-blur-xs flex items-center justify-center z-50 p-4">
		<div class="bg-slate-900 border border-amber-500/40 rounded-xl p-6 max-w-md w-full shadow-2xl shadow-amber-500/10">
			<div class="flex items-center gap-3 mb-3">
				<div class="w-10 h-10 rounded-full bg-amber-500/10 border border-amber-500/30 flex items-center justify-center text-amber-400 text-xl font-bold">
					⚠️
				</div>
				<div>
					<h3 class="text-base font-bold text-white">Danger Zone Warning</h3>
					<p class="text-xs text-amber-400 font-medium">Potentially Destructive Query Detected</p>
				</div>
			</div>

			<p class="text-xs text-slate-300 mb-4 leading-relaxed bg-slate-950 p-3 rounded-lg border border-slate-800">
				{dangerWarning}
			</p>

			<div class="flex justify-end gap-3">
				<button
					onclick={() => (showDangerModal = false)}
					class="px-4 py-2 rounded-lg bg-slate-800 hover:bg-slate-700 text-xs font-medium text-slate-300 transition-colors"
				>
					Cancel
				</button>
				<button
					onclick={confirmDangerousQuery}
					class="px-4 py-2 rounded-lg bg-amber-600 hover:bg-amber-500 text-xs font-bold text-white transition-colors shadow-lg shadow-amber-500/20"
				>
					Yes, Execute Query
				</button>
			</div>
		</div>
	</div>
{/if}
