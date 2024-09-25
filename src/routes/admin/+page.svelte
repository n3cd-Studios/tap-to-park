
<script>
    import { authStore } from "$lib/auth";
    import Card from "../../components/Card.svelte";

    let cards = [
        {
            name: "Kent State University",
            description: "Manages all parking spaces on the Kent State University campuses and regional campuses."
        },
        {
            name: "City of Kent",
            description: "Parking meter management for the city of Kent."
        }
    ]

    $: email = $authStore.user?.email
</script>

<div class="m-10">
    <p>Hello, {email}! You have access to these organizations:</p>
    <div class="pt-3 flex flex-row gap-1">
        {#each cards as card}
            <Card {...card} />
        {/each}
    </div>
</div>
