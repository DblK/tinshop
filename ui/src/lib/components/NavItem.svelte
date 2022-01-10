<script lang="ts">
    import { variables } from '$lib/variables';
    import IconButton from '@smui/icon-button';

    export const prerender = true;
    export let link : string = ""
    export let icon : string = ""
    export let alt : string
    export let bottom: boolean = false
    export let materialIcon: string = ""
</script>

<style lang="scss">
    // div {
        // background-color: green;
    // }

    .icon {
        display: flex;
        align-items: center;
        &.img {
            width: 48px;
            padding-top: 0.5rem;
        }
    }
    .bottom {
        // background-color: hotpink;
        flex-grow: 1;
        display: flex;
        flex-direction: column-reverse;
        & > div.icon {
            color: var(--mdc-theme-primary, #F66709);
        }
    }
</style>

<div class="{bottom ? 'bottom': ''}">
    {#if materialIcon}
        <div class="icon">
            <IconButton class="material-icons tinshop-button">{materialIcon}</IconButton>
        </div>
    {:else}
        <a href="{variables.basePath}{link}">
            <img class="icon img" src="{variables.basePath}{icon}" alt="{alt}" title="{alt}">
        </a>
    {/if}
</div>