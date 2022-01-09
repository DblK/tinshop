<script lang="ts">
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
    }
</style>

<div class='navbar'>
    <slot/>
</div>