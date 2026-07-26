<script lang="ts">
	import { onMount } from 'svelte';
	import {
		queryHistory,
		querySnippets,
		isDrawerOpen,
		activeDrawerTab,
		refreshHistory,
		clearHistory,
		refreshSnippets,
		removeSnippet
	} from '$lib/stores/historyStore';
	import type { QueryHistoryItem, QuerySnippet } from '$lib/api';
	import { Button } from '../shadcn/button';
	import { Input } from '../shadcn/input';
	import {
		History,
		Bookmark,
		Trash2,
		CornerDownLeft,
		ChevronDown,
		ChevronUp,
		Clock,
		CheckCircle2,
		XCircle,
		Search,
		Plus
	} from '@lucide/svelte';

	let {
		onInsertQuery,
		onSaveSnippet
	} = $props<{
		onInsertQuery: (sql: string) => void;
		onSaveSnippet?: (sql: string) => void;
	}>();

	let searchQuery = $state('');

	onMount(() => {
		refreshHistory();
		refreshSnippets();
	});

	let filteredHistory = $derived(
		$queryHistory.filter((h) =>
			h.query.toLowerCase().includes(searchQuery.toLowerCase()) ||
			(h.error_message && h.error_message.toLowerCase().includes(searchQuery.toLowerCase()))
		)
	);

	let filteredSnippets = $derived(
		$querySnippets.filter((s) =>
			s.title.toLowerCase().includes(searchQuery.toLowerCase()) ||
			s.query.toLowerCase().includes(searchQuery.toLowerCase()) ||
			(s.description && s.description.toLowerCase().includes(searchQuery.toLowerCase()))
		)
	);

	function formatTime(isoString: string) {
		try {
			const d = new Date(isoString);
			return d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', second: '2-digit' });
		} catch {
			return isoString;
		}
	}
</script>

<div class="border-t border-border bg-card transition-all duration-200">
	<!-- Drawer Header Bar -->
	<div class="flex items-center justify-between border-b border-border/60 px-3 py-1.5 bg-muted/20">
		<div class="flex items-center gap-1">
			<button
				type="button"
				class="flex items-center gap-1.5 rounded-md px-2.5 py-1 text-xs font-semibold transition-colors {$activeDrawerTab === 'history' ? 'bg-background shadow-xs text-foreground' : 'text-muted-foreground hover:text-foreground'}"
				onclick={() => {
					$activeDrawerTab = 'history';
					if (!$isDrawerOpen) $isDrawerOpen = true;
				}}
			>
				<History class="size-3.5" />
				<span>History</span>
				<span class="ml-1 rounded-full bg-muted px-1.5 py-0.2 text-[10px] font-mono text-muted-foreground">{$queryHistory.length}</span>
			</button>

			<button
				type="button"
				class="flex items-center gap-1.5 rounded-md px-2.5 py-1 text-xs font-semibold transition-colors {$activeDrawerTab === 'snippets' ? 'bg-background shadow-xs text-foreground' : 'text-muted-foreground hover:text-foreground'}"
				onclick={() => {
					$activeDrawerTab = 'snippets';
					if (!$isDrawerOpen) $isDrawerOpen = true;
				}}
			>
				<Bookmark class="size-3.5" />
				<span>Snippets</span>
				<span class="ml-1 rounded-full bg-muted px-1.5 py-0.2 text-[10px] font-mono text-muted-foreground">{$querySnippets.length}</span>
			</button>
		</div>

		<div class="flex items-center gap-2">
			{#if $isDrawerOpen}
				<div class="relative w-44">
					<Search class="absolute left-2 top-2 size-3.5 text-muted-foreground" />
					<Input
						placeholder="Search queries..."
						bind:value={searchQuery}
						class="h-7 pl-7 text-[11px]"
					/>
				</div>

				{#if $activeDrawerTab === 'history' && $queryHistory.length > 0}
					<Button variant="ghost" size="sm" class="h-7 text-xs text-destructive hover:bg-destructive/10" onclick={() => clearHistory()}>
						<Trash2 class="mr-1 size-3" /> Clear
					</Button>
				{/if}

				{#if $activeDrawerTab === 'snippets'}
					<Button variant="outline" size="sm" class="h-7 text-xs" onclick={() => onSaveSnippet?.('')}>
						<Plus class="mr-1 size-3" /> New Snippet
					</Button>
				{/if}
			{/if}

			<Button
				variant="ghost"
				size="sm"
				class="h-7 px-2 text-xs text-muted-foreground hover:text-foreground"
				onclick={() => ($isDrawerOpen = !$isDrawerOpen)}
			>
				{#if $isDrawerOpen}
					<ChevronDown class="size-4" />
				{:else}
					<ChevronUp class="size-4" />
				{/if}
			</Button>
		</div>
	</div>

	<!-- Drawer Content Area -->
	{#if $isDrawerOpen}
		<div class="max-h-60 min-h-36 overflow-y-auto p-2">
			{#if $activeDrawerTab === 'history'}
				{#if filteredHistory.length === 0}
					<div class="flex flex-col items-center justify-center py-8 text-center text-xs text-muted-foreground">
						<Clock class="mb-1.5 size-5 text-muted-foreground/50" />
						<span>No query history found</span>
					</div>
				{:else}
					<div class="space-y-1.5">
						{#each filteredHistory as item (item.id)}
							<div class="group flex items-start justify-between rounded-md border border-border/50 bg-background/80 p-2 text-xs transition-colors hover:border-border">
								<div class="flex-1 space-y-1 pr-3">
									<div class="flex items-center gap-2">
										{#if item.status === 'success'}
											<span class="inline-flex items-center gap-1 rounded-sm bg-success/10 px-1.5 py-0.5 text-[10px] font-semibold text-success">
												<CheckCircle2 class="size-3" /> Success
											</span>
										{:else}
											<span class="inline-flex items-center gap-1 rounded-sm bg-destructive/10 px-1.5 py-0.5 text-[10px] font-semibold text-destructive">
												<XCircle class="size-3" /> Error
											</span>
										{/if}

										<span class="font-mono text-[10px] text-muted-foreground">{item.duration_ms}ms</span>
										{#if item.rows_affected >= 0}
											<span class="font-mono text-[10px] text-muted-foreground">• {item.rows_affected} rows</span>
										{/if}
										<span class="ml-auto font-mono text-[10px] text-muted-foreground">{formatTime(item.executed_at)}</span>
									</div>

									<pre class="font-mono text-[11px] text-foreground max-h-16 overflow-x-auto whitespace-pre-wrap rounded bg-muted/30 p-1.5">{item.query}</pre>

									{#if item.error_message}
										<div class="font-mono text-[10px] text-destructive">{item.error_message}</div>
									{/if}
								</div>

								<div class="flex items-center gap-1 opacity-80 group-hover:opacity-100">
									<Button
										variant="ghost"
										size="sm"
										class="h-7 px-2 text-[11px]"
										title="Insert into SQL Editor"
										onclick={() => onInsertQuery(item.query)}
									>
										<CornerDownLeft class="mr-1 size-3" /> Use
									</Button>

									<Button
										variant="ghost"
										size="sm"
										class="h-7 px-2 text-[11px]"
										title="Save as Snippet"
										onclick={() => onSaveSnippet?.(item.query)}
									>
										<Bookmark class="size-3 text-muted-foreground" />
									</Button>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			{:else if $activeDrawerTab === 'snippets'}
				{#if filteredSnippets.length === 0}
					<div class="flex flex-col items-center justify-center py-8 text-center text-xs text-muted-foreground">
						<Bookmark class="mb-1.5 size-5 text-muted-foreground/50" />
						<span>No saved query snippets yet</span>
					</div>
				{:else}
					<div class="grid grid-cols-1 gap-2 md:grid-cols-2">
						{#each filteredSnippets as snippet (snippet.id)}
							<div class="group flex flex-col justify-between rounded-md border border-border/50 bg-background/80 p-2.5 text-xs transition-colors hover:border-border">
								<div class="space-y-1">
									<div class="flex items-center justify-between">
										<span class="font-bold text-foreground">{snippet.title}</span>
										<Button
											variant="ghost"
											size="sm"
											class="h-6 w-6 p-0 text-destructive opacity-0 group-hover:opacity-100"
											onclick={() => snippet.id && removeSnippet(snippet.id)}
										>
											<Trash2 class="size-3" />
										</Button>
									</div>

									{#if snippet.description}
										<p class="text-[11px] text-muted-foreground">{snippet.description}</p>
									{/if}

									<pre class="font-mono text-[11px] text-foreground max-h-20 overflow-x-auto whitespace-pre-wrap rounded bg-muted/30 p-1.5">{snippet.query}</pre>
								</div>

								<div class="mt-2 flex justify-end">
									<Button
										variant="outline"
										size="sm"
										class="h-6 px-2 text-[11px]"
										onclick={() => onInsertQuery(snippet.query)}
									>
										<CornerDownLeft class="mr-1 size-3" /> Insert Snippet
									</Button>
								</div>
							</div>
						{/each}
					</div>
				{/if}
			{/if}
		</div>
	{/if}
</div>
