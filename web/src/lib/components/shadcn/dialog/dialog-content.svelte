<script lang="ts">
	import { Dialog as DialogPrimitive } from "bits-ui";
	import DialogPortal from "./dialog-portal.svelte";
	import DialogOverlay from "./dialog-overlay.svelte";
	import { cn } from "$lib/utils.js";
	import { X } from "@lucide/svelte";

	let {
		ref = $bindable(null),
		class: className,
		portalProps,
		children,
		showClose = true,
		...restProps
	}: DialogPrimitive.ContentProps & {
		portalProps?: DialogPrimitive.PortalProps;
		showClose?: boolean;
	} = $props();
</script>

<DialogPortal {...portalProps}>
	<DialogOverlay />
	<DialogPrimitive.Content
		bind:ref
		data-slot="dialog-content"
		class={cn(
			"data-open:animate-in data-closed:animate-out data-closed:fade-out-0 data-open:fade-in-0 data-closed:zoom-out-95 data-open:zoom-in-95 bg-card text-card-foreground border-border fixed top-1/2 left-1/2 z-50 grid w-full max-w-lg -translate-x-1/2 -translate-y-1/2 gap-4 rounded-xl border p-6 shadow-2xl duration-200 outline-none",
			className
		)}
		{...restProps}
	>
		{@render children?.()}
		{#if showClose}
			<DialogPrimitive.Close
				class="hover:bg-muted text-muted-foreground hover:text-foreground absolute right-4 top-4 rounded-sm p-1 transition-opacity focus:outline-none cursor-pointer"
			>
				<X class="size-4" />
				<span class="sr-only">Close</span>
			</DialogPrimitive.Close>
		{/if}
	</DialogPrimitive.Content>
</DialogPortal>
