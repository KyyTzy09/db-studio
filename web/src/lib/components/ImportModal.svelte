<script lang="ts">
	import { batchInsertOrUpdate } from '../api';

	let { tableName, isOpen = false, onClose, onSuccess } = $props<{
		tableName: string;
		isOpen?: boolean;
		onClose: () => void;
		onSuccess: () => void;
	}>();

	let isDragging = $state(false);
	let fileName = $state('');
	let parsedRows = $state<Record<string, any>[]>([]);
	let previewRows = $state<Record<string, any>[]>([]);
	let importMode = $state<'insert' | 'upsert'>('insert');
	let isUploading = $state(false);
	let errorMessage = $state('');
	let statusMessage = $state('');

	$effect(() => {
		if (isOpen) {
			resetState();
		}
	});

	function resetState() {
		isDragging = false;
		fileName = '';
		parsedRows = [];
		previewRows = [];
		importMode = 'insert';
		isUploading = false;
		errorMessage = '';
		statusMessage = '';
	}

	function handleDragOver(e: DragEvent) {
		e.preventDefault();
		isDragging = true;
	}

	function handleDragLeave(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
	}

	function handleDrop(e: DragEvent) {
		e.preventDefault();
		isDragging = false;
		if (e.dataTransfer && e.dataTransfer.files.length > 0) {
			processFile(e.dataTransfer.files[0]);
		}
	}

	function handleFileSelect(e: Event) {
		const target = e.target as HTMLInputElement;
		if (target.files && target.files.length > 0) {
			processFile(target.files[0]);
		}
	}

	function processFile(file: File) {
		errorMessage = '';
		fileName = file.name;
		const reader = new FileReader();

		if (file.name.endsWith('.json')) {
			reader.onload = (evt) => {
				try {
					const json = JSON.parse(evt.target?.result as string);
					if (Array.isArray(json)) {
						parsedRows = json;
						previewRows = json.slice(0, 5);
					} else {
						errorMessage = 'JSON file must contain an array of objects';
					}
				} catch (err: any) {
					errorMessage = 'Failed to parse JSON file: ' + err.message;
				}
			};
			reader.readAsText(file);
		} else if (file.name.endsWith('.csv')) {
			reader.onload = (evt) => {
				try {
					const text = evt.target?.result as string;
					const rows = parseCSVText(text);
					parsedRows = rows;
					previewRows = rows.slice(0, 5);
				} catch (err: any) {
					errorMessage = 'Failed to parse CSV file: ' + err.message;
				}
			};
			reader.readAsText(file);
		} else {
			errorMessage = 'Only .csv and .json files are supported!';
		}
	}

	function parseCSVText(text: string): Record<string, any>[] {
		const lines = text.split(/\r?\n/).filter((l) => l.trim() !== '');
		if (lines.length < 2) return [];

		const headers = parseCSVLine(lines[0]);
		const results: Record<string, any>[] = [];

		for (let i = 1; i < lines.length; i++) {
			const values = parseCSVLine(lines[i]);
			if (values.length === headers.length) {
				const row: Record<string, any> = {};
				headers.forEach((h, idx) => {
					row[h.trim()] = values[idx].trim();
				});
				results.push(row);
			}
		}
		return results;
	}

	function parseCSVLine(line: string): string[] {
		const result: string[] = [];
		let current = '';
		let inQuotes = false;

		for (let i = 0; i < line.length; i++) {
			const char = line[i];
			if (char === '"') {
				inQuotes = !inQuotes;
			} else if (char === ',' && !inQuotes) {
				result.push(current);
				current = '';
			} else {
				current += char;
			}
		}
		result.push(current);
		return result;
	}

	async function handleExecuteImport() {
		if (parsedRows.length === 0) return;
		isUploading = true;
		errorMessage = '';
		statusMessage = `Importing ${parsedRows.length} rows...`;

		try {
			const res = await batchInsertOrUpdate(tableName, parsedRows, importMode);
			statusMessage = `✅ Success! Processed ${res.affected_rows} rows.`;
			setTimeout(() => {
				onSuccess();
				onClose();
			}, 1000);
		} catch (err: any) {
			errorMessage = err.message || 'Failed to import data';
			statusMessage = '';
		} finally {
			isUploading = false;
		}
	}
</script>

{#if isOpen}
	<div class="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-xs p-4">
		<div class="w-full max-w-2xl rounded-xl border border-slate-700 bg-slate-900 p-6 shadow-2xl">
			<div class="flex items-center justify-between border-b border-slate-800 pb-4">
				<h3 class="text-lg font-bold text-slate-100">
					📥 Bulk Import Data to <span class="text-purple-400 font-mono">{tableName}</span>
				</h3>
				<button onclick={onClose} class="text-slate-400 hover:text-slate-200 text-xl font-bold">✕</button>
			</div>

			{#if errorMessage}
				<div class="mt-4 rounded-lg bg-rose-950/70 border border-rose-800/80 p-3 text-xs text-rose-300">
					⚠️ {errorMessage}
				</div>
			{/if}

			{#if statusMessage}
				<div class="mt-4 rounded-lg bg-emerald-950/70 border border-emerald-800/80 p-3 text-xs text-emerald-300">
					{statusMessage}
				</div>
			{/if}

			<!-- Drag & Drop Zone -->
			<div
				ondragover={handleDragOver}
				ondragleave={handleDragLeave}
				ondrop={handleDrop}
				class="mt-4 border-2 border-dashed rounded-xl p-6 text-center transition flex flex-col items-center justify-center cursor-pointer {
					isDragging
						? 'border-purple-500 bg-purple-500/10'
						: 'border-slate-700 bg-slate-950/60 hover:border-slate-500 hover:bg-slate-950'
				}"
			>
				<input type="file" accept=".csv,.json" onchange={handleFileSelect} class="hidden" id="fileInput" />
				<label for="fileInput" class="cursor-pointer w-full flex flex-col items-center">
					<div class="text-3xl mb-2">📁</div>
					{#if fileName}
						<span class="text-sm font-semibold text-purple-300 font-mono">{fileName}</span>
						<span class="text-xs text-slate-400 mt-1">({parsedRows.length} rows detected)</span>
					{:else}
						<span class="text-sm font-semibold text-slate-200">Drag & Drop your .csv or .json file here</span>
						<span class="text-xs text-slate-400 mt-1">or click to browse files from your computer</span>
					{/if}
				</label>
			</div>

			<!-- Live Preview & Settings -->
			{#if parsedRows.length > 0}
				<div class="mt-4 space-y-4">
					<!-- Import Mode Toggle -->
					<div class="flex items-center justify-between bg-slate-950 rounded-lg p-3 border border-slate-800">
						<span class="text-xs font-semibold text-slate-300">Import Mode Strategy:</span>
						<div class="flex gap-2">
							<button
								type="button"
								onclick={() => (importMode = 'insert')}
								class="px-3 py-1.5 rounded-md text-xs font-semibold transition {
									importMode === 'insert' ? 'bg-purple-600 text-white shadow' : 'bg-slate-800 text-slate-400 hover:text-slate-200'
								}"
							>
								Insert Only
							</button>
							<button
								type="button"
								onclick={() => (importMode = 'upsert')}
								class="px-3 py-1.5 rounded-md text-xs font-semibold transition {
									importMode === 'upsert' ? 'bg-amber-600 text-white shadow' : 'bg-slate-800 text-slate-400 hover:text-slate-200'
								}"
							>
								Upsert (Update if PK Exists)
							</button>
						</div>
					</div>

					<!-- Sample Table Preview -->
					<div>
						<div class="text-xs font-semibold text-slate-400 mb-2">
							Sample Data Preview (First 5 of {parsedRows.length} rows):
						</div>
						<div class="max-h-36 overflow-auto rounded-lg border border-slate-800 bg-slate-950 text-xs">
							<table class="w-full text-left font-mono">
								<thead class="bg-slate-900 text-slate-400 border-b border-slate-800">
									<tr>
										{#each Object.keys(previewRows[0] || {}) as key}
											<th class="p-2 border-r border-slate-800">{key}</th>
										{/each}
									</tr>
								</thead>
								<tbody>
									{#each previewRows as row}
										<tr class="border-b border-slate-900/60 hover:bg-slate-900/40">
											{#each Object.keys(previewRows[0] || {}) as key}
												<td class="p-2 border-r border-slate-800 text-slate-300">{row[key]}</td>
											{/each}
										</tr>
									{/each}
								</tbody>
							</table>
						</div>
					</div>
				</div>
			{/if}

			<div class="pt-6 flex justify-end gap-3 border-t border-slate-800 mt-6">
				<button
					type="button"
					onclick={onClose}
					class="rounded-lg px-4 py-2 text-xs font-medium text-slate-400 hover:bg-slate-800 hover:text-slate-200 transition"
				>
					Cancel
				</button>
				<button
					type="button"
					onclick={handleExecuteImport}
					disabled={parsedRows.length === 0 || isUploading}
					class="rounded-lg bg-purple-600 px-5 py-2 text-xs font-semibold text-white shadow-lg hover:bg-purple-500 disabled:opacity-50 transition"
				>
					{isUploading ? 'Processing Import...' : `Import ${parsedRows.length} Rows`}
				</button>
			</div>
		</div>
	</div>
{/if}
