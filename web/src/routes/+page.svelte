<script lang="ts">
	import { goto } from '$app/navigation';
	async function login() {
		const response = await fetch('http://localhost:8080/auth/github');
		if (response.status == 200) {
			const { url } = (await response.json()) as { url: string };
			goto(url)
		}
	}
</script>

<button on:click={login}>
	<h1>Login</h1>
</button>
