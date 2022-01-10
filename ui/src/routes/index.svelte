<script lang="ts">
	export const prerender = true;
    import { onMount } from 'svelte';
    import LayoutGrid, { Cell } from '@smui/layout-grid';
    import Nav from '$lib/components/nav.svelte';
    import Content from '$lib/components/Content.svelte';
    import NavItem from '$lib/components/NavItem.svelte';
import { variables } from '$lib/variables';

    let stats;
    let keyMetrics = []
    onMount(async () => {
        const res = await fetch(variables.api + "/stats");
        const JSONStats = await res.json();
        // console.log(JSONStats)
        stats = JSONStats
        keyMetrics = [...keyMetrics, {
            description: "Overall Visits",
            value: stats.visit
        }, {
            description: "Unique switch",
            value: stats.uniqueSwitch
        }, {
            description: "Downloads",
            value: stats.downloadAsked
        }]
    })
</script>

<style type="scss">
    :global(body) {
		margin: 0;
	}
    main {
        display: flex;
    }
    .primary-stats {
        padding: 0.5rem;
        border-radius: 0.5rem;
        display: flex;
        justify-content: center;
        align-items: center;
        text-align: center;
        background-color: var(--mdc-theme-secondary, #F66709);
        color: var(--mdc-theme-on-secondary, #fff);

        flex-direction: column;
        .big-number {
            font-size: 4rem;
        }
        .description {
            font-style: italic;
        }
    }
</style>

<main>
    <Nav>
        <NavItem alt="tinshop" icon="favicon.png"></NavItem>
        <!-- <NavItem alt="tinshop" icon="favicon.png"></NavItem> -->
        <!-- <NavItem alt="tinshop" icon="favicon.png"></NavItem> -->
        <!-- <NavItem alt="tinshop" icon="favicon.png"></NavItem> -->
        
        <NavItem alt="Settings" materialIcon="settings" bottom={true}></NavItem>
    </Nav>
    <Content>
        <div class="mdc-layout-grid">
            <LayoutGrid>
                {#each keyMetrics as metric}
                    <Cell>
                        <div class='primary-stats'>
                            <div class="big-number">{metric.value}</div>
                            <div class="description">{metric.description}</div>
                        </div>
                    </Cell>
                {/each}
            </LayoutGrid>
        </div>
    </Content>
</main>
