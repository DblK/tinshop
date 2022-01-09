<script lang="ts">
	export const prerender = true;
</script>

<style>
    div.content {
        /* background-color: green; */
        height: 100vh;
        width: 100%;
    }
</style>

<div class='content'>
    <slot/>
</div>