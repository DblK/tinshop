<script lang="ts">
    import { variables } from '$lib/variables';
    export const prerender = true;
    export let link : string = ""
    export let icon : string
    export let alt : string
    export let bottom: boolean = false

</script>

<style lang="scss">
    // div {
        // background-color: green;
    // }
    .icon {
        width: 48px;
        text-align: center;
        padding-top: 0.5rem;
    }
    .bottom {
        // background-color: hotpink;
        flex-grow: 1;
        display: flex;
        flex-direction: column-reverse;
    }
</style>

<div class="{bottom ? ' bottom': ''}">
    <a href="{variables.basePath}{link}">
        <img class="icon" src="{variables.basePath}{icon}" alt="{alt}" title="{alt}">
    </a>
</div>