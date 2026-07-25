<script lang="ts">
	import { Dialog as DialogPrimitive } from "bits-ui";

	let { ...restProps }: DialogPrimitive.PortalProps = $props();
</script>

<DialogPrimitive.Portal data-slot="dialog-portal" {...restProps} />
