<script lang="ts">
	import { fetchTableSchema, type TableSchema } from '$lib/api';

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
	<div class="px-6 py-4 border-b border-slate-800 bg-slate-900/50 flex items-center justify-between">
		<div>
			<h2 class="text-sm font-bold text-white flex items-center gap-2">
				<svg class="w-4 h-4 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
					<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
				</svg>
				Schema Definition for "{tableName}"
			</h2>
			<p class="text-xs text-slate-400">Column metadata, data types, and constraint definitions</p>
		</div>
	</div>

	<div class="flex-1 overflow-auto custom-scrollbar p-6">
		{#if loading}
			<div class="h-64 flex flex-col items-center justify-center text-slate-400 gap-2">
				<div class="w-6 h-6 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin"></div>
				<span class="text-xs">Loading schema details...</span>
			</div>
		{:else if errorMsg}
			<div class="p-4 rounded-xl bg-red-500/10 border border-red-500/20 text-red-400 text-xs">
				❌ Error: {errorMsg}
			</div>
		{:else if schema}
			<div class="border border-slate-800 rounded-xl overflow-hidden shadow-2xl bg-slate-900/60">
				<table class="w-full text-left text-xs border-collapse">
					<thead>
						<tr class="bg-slate-900/90 text-slate-400 font-semibold border-b border-slate-800 uppercase tracking-wider">
							<th class="px-4 py-3 border-r border-slate-800/80">Column Name</th>
							<th class="px-4 py-3 border-r border-slate-800/80">Data Type</th>
							<th class="px-4 py-3 border-r border-slate-800/80">Primary Key</th>
							<th class="px-4 py-3 border-r border-slate-800/80">Nullable</th>
							<th class="px-4 py-3">Default Value</th>
						</tr>
					</thead>
					<tbody class="divide-y divide-slate-800/60 font-mono">
						{#each schema.columns as col}
							<tr class="hover:bg-slate-800/40 transition-colors">
								<td class="px-4 py-3 font-semibold text-indigo-300 border-r border-slate-800/60">
									{col.name}
								</td>
								<td class="px-4 py-3 text-slate-300 border-r border-slate-800/60">
									<span class="px-2 py-0.5 rounded bg-slate-800 text-slate-300 border border-slate-700/60 text-[11px]">
										{col.data_type}
									</span>
								</td>
								<td class="px-4 py-3 border-r border-slate-800/60">
									{#if col.is_primary_key}
										<span class="px-2 py-0.5 rounded bg-amber-500/10 text-amber-400 border border-amber-500/20 text-[10px] font-sans font-semibold">
											🔑 PRIMARY KEY
										</span>
									{:else}
										<span class="text-slate-600 font-sans">-</span>
									{/if}
								</td>
								<td class="px-4 py-3 border-r border-slate-800/60 font-sans">
									{#if col.is_nullable}
										<span class="px-2 py-0.5 rounded bg-slate-800 text-slate-400 text-[10px]">
											YES
										</span>
									{:else}
										<span class="px-2 py-0.5 rounded bg-slate-800/40 text-slate-500 text-[10px]">
											NO
										</span>
									{/if}
								</td>
								<td class="px-4 py-3 text-slate-400">
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
