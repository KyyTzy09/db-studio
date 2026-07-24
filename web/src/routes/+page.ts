import type { PageLoad } from './$types';

export const load: PageLoad = async ({ fetch }) => {
	try {
		const statusRes = await fetch('/api/connection/status');
		const status = await statusRes.json();

		let tables: any[] = [];
		if (status.connected) {
			const tablesRes = await fetch('/api/tables');
			if (tablesRes.ok) {
				const data = await tablesRes.json();
				tables = data.tables || [];
			}
		}

		return {
			status,
			tables
		};
	} catch (err: any) {
		return {
			status: { connected: false, error: err.message },
			tables: []
		};
	}
};
