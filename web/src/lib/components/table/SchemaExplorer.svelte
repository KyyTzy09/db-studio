<script lang="ts">
	import { fetchTableSchema, type TableSchema } from '../../../data/services';

	let { tableName } = $props<{ tableName: string }>();

	let schema = $state<TableSchema | null>(null);
	let loading = $state(true);
	let errorMsg = $state<string | null>(null);

	$effect(() => {
		if (tableName) {
			loadSchema(tableName);
		}
	});

	async function loadSchema(name: string) {
		loading = true;
		errorMsg = null;
		try {
			schema = await fetchTableSchema(name);
		} catch (err: any) {
			errorMsg = err.message || 'Failed to load schema';
		} finally {
			loading = false;
		}
	}
</script>

<div class="flex-1 flex flex-col h-full bg-slate-950 text-slate-200 overflow-hidden">
	<div class="px-6 py-3 border-b border-slate-800 bg-slate-900/50 flex items-center justify-between">
		<h2 class="text-xs font-bold text-slate-400 uppercase tracking-wider">
			Schema for <span class="text-indigo-400 font-mono">{tableName}</span>
		</h2>
		<button
			onclick={() => loadSchema(tableName)}
			class="px-3 py-1 rounded bg-slate-800 hover:bg-slate-700 text-xs text-slate-300 transition"
		>
			Refresh
		</button>
	</div>

	<div class="flex-1 overflow-auto custom-scrollbar p-6">
		{#if loading}
			<div class="h-64 flex flex-col items-center justify-center text-slate-400 gap-2">
				<div class="w-6 h-6 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
				<span class="text-xs">Loading schema metadata...</span>
			</div>
		{:else if errorMsg}
			<div class="p-4 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-xs font-mono">
				❌ {errorMsg}
			</div>
		{:else if schema}
			<div class="border border-slate-800 rounded-xl overflow-hidden shadow-2xl bg-slate-900/60">
				<table class="w-full text-left text-xs border-collapse font-mono">
					<thead>
						<tr class="bg-slate-900/90 text-slate-400 font-semibold border-b border-slate-800 uppercase tracking-wider">
							<th class="px-4 py-3 border-r border-slate-800/80">Column Name</th>
							<th class="px-4 py-3 border-r border-slate-800/80">Data Type</th>
							<th class="px-4 py-3 border-r border-slate-800/80 text-center">Nullable</th>
							<th class="px-4 py-3 border-r border-slate-800/80 text-center">Primary Key</th>
							<th class="px-4 py-3">Default Value</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-800/60">
						{#each schema.columns as col}
							<tr class="hover:bg-slate-800/40 transition-colors">
								<td class="px-4 py-2.5 border-r border-slate-800/60 font-semibold text-slate-200">
									{col.name}
								</td>
								<td class="px-4 py-2.5 border-r border-slate-800/60 text-purple-400">
									{col.data_type}
								</td>
								<td class="px-4 py-2.5 border-r border-slate-800/60 text-center">
									{#if col.is_nullable}
										<span class="px-2 py-0.5 rounded bg-emerald-500/10 text-emerald-400 text-[10px]">YES</span>
									{:else}
										<span class="px-2 py-0.5 rounded bg-rose-500/10 text-rose-400 text-[10px]">NO</span>
									{/if}
								</td>
								<td class="px-4 py-2.5 border-r border-slate-800/60 text-center">
									{#if col.is_primary_key}
										<span class="px-2 py-0.5 rounded bg-amber-500/20 text-amber-300 text-[10px] font-bold">PRIMARY KEY</span>
									{:else}
										<span class="text-slate-600">-</span>
									{/if}
								</td>
								<td class="px-4 py-2.5 text-slate-400">
									{col.default_value || '-'}
								</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		{/if}
	</div>
</div>
