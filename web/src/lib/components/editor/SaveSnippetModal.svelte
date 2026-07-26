<script lang="ts">
	import { Button } from '../shadcn/button';
	import { Input } from '../shadcn/input';
	import { Label } from '../shadcn/label';
	import {
		Dialog,
		DialogContent,
		DialogHeader,
		DialogTitle,
		DialogFooter,
		DialogDescription
	} from '../shadcn/dialog';
	import { Bookmark, Loader2 } from '@lucide/svelte';
	import { addOrUpdateSnippet } from '$lib/stores/historyStore';

	let {
		isOpen = $bindable(false),
		initialQuery = '',
		onSuccess
	} = $props<{
		isOpen: boolean;
		initialQuery?: string;
		onSuccess?: () => void;
	}>();

	let title = $state('');
	let description = $state('');
	let query = $state('');
	let isSaving = $state(false);
	let error = $state<string | null>(null);

	$effect(() => {
		if (isOpen) {
			query = initialQuery;
			title = '';
			description = '';
			error = null;
		}
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		if (!title.trim() || !query.trim()) {
			error = 'Title and Query are required';
			return;
		}
		isSaving = true;
		error = null;
		try {
			await addOrUpdateSnippet({
				title: title.trim(),
				description: description.trim(),
				query: query.trim()
			});
			isOpen = false;
			onSuccess?.();
		} catch (err: any) {
			error = err.message || 'Failed to save snippet';
		} finally {
			isSaving = false;
		}
	}
</script>

<Dialog bind:open={isOpen}>
	<DialogContent class="max-w-lg">
		<form onsubmit={handleSubmit}>
			<DialogHeader>
				<DialogTitle class="flex items-center gap-1.5 text-base font-bold">
					<Bookmark class="size-4 text-primary" /> Save Query Snippet
				</DialogTitle>
				<DialogDescription>
					Save your SQL query as a reusable snippet for future access.
				</DialogDescription>
			</DialogHeader>

			{#if error}
				<div class="mt-3 rounded-md bg-destructive/10 border border-destructive/20 p-2.5 text-xs text-destructive">
					{error}
				</div>
			{/if}

			<div class="space-y-4 py-4">
				<div class="space-y-1.5">
					<Label for="snippet-title" class="text-xs font-medium">Snippet Title <span class="text-destructive">*</span></Label>
					<Input
						id="snippet-title"
						placeholder="e.g. Monthly Top Sales Query"
						bind:value={title}
						class="h-9 text-xs"
						required
					/>
				</div>

				<div class="space-y-1.5">
					<Label for="snippet-desc" class="text-xs font-medium">Description (Optional)</Label>
					<Input
						id="snippet-desc"
						placeholder="Short explanation of what this query does"
						bind:value={description}
						class="h-9 text-xs"
					/>
				</div>

				<div class="space-y-1.5">
					<Label for="snippet-query" class="text-xs font-medium">SQL Query <span class="text-destructive">*</span></Label>
					<textarea
						id="snippet-query"
						bind:value={query}
						rows="4"
						class="w-full rounded-md border border-input bg-muted/40 p-2.5 font-mono text-xs focus:outline-none focus:ring-1 focus:ring-primary"
						required
					></textarea>
				</div>
			</div>

			<DialogFooter>
				<Button type="button" variant="outline" size="sm" onclick={() => (isOpen = false)} disabled={isSaving}>
					Cancel
				</Button>
				<Button type="submit" size="sm" disabled={isSaving}>
					{#if isSaving}
						<Loader2 class="mr-1.5 size-3.5 animate-spin" /> Saving...
					{:else}
						Save Snippet
					{/if}
				</Button>
			</DialogFooter>
		</form>
	</DialogContent>
</Dialog>
