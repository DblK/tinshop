<script lang="ts">
	import { variables } from '$lib/variables';
    export const prerender = true;
</script>

<style type="scss">
    div.navbar {
        height: 100vh;
        width: 64px;
        // background-color: yellow;
        display: flex;
        flex-direction: column;
        align-items: center;
        // gap: 2rem;

        .icon {
            width: 48px;
            text-align: center;

            padding-top: 0.5rem;
        }
    }
</style>

<div class='navbar'>
    <div>
        <img class="icon" src="{variables.basePath}favicon.png" alt="tinshop">
    </div>
    <slot/>
</div>