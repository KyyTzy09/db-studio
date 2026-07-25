<script lang="ts">
	import { cn } from "$lib/utils.js";
	import type { HTMLAttributes } from "svelte/elements";

	let {
		class: className,
		children,
		...restProps
	}: HTMLAttributes<HTMLDivElement> = $props();
</script>

<div
	data-slot="dialog-footer"
	class={cn("flex flex-col-reverse sm:flex-row sm:justify-end sm:gap-2", className)}
	{...restProps}
>
	{@render children?.()}
</div>
