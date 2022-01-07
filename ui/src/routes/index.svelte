<script>
	export const prerender = true;
    import { onMount } from 'svelte';

    let stats;
    onMount(async () => {
        const res = await fetch("http://localhost:3000/api/stats");
        const JSONStats = await res.json();
        stats = JSON.stringify(JSONStats)
    })
</script>

<style>
    .icon {
        width: 64px;
    }
</style>

<main>
    <img class="icon" src="/favicon.png" alt="tinshop">
    
    <p>tinshop</p>
    
    <p>{stats}</p>
</main>
