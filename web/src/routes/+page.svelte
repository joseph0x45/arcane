<script lang="ts">
	import { goto } from '$app/navigation';
  import { API_URL } from '$lib/config';
	async function login() {
		const response = await fetch(`${API_URL}/auth/github`);
		if (response.status == 200) {
			const { url } = (await response.json()) as { url: string };
			goto(url)
		}
	}
</script>

<button on:click={login}>
	<h1>Login</h1>
</button>
